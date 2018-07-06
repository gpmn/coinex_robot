package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/nntaoli-project/GoEx"
	"github.com/nntaoli-project/GoEx/bigone"
	"github.com/nntaoli-project/GoEx/coinex"
	"github.com/nntaoli-project/GoEx/fcoin"
)

type subList []goex.SubAccount

func (sl subList) Len() int { return len(sl) }
func (sl subList) Less(i, j int) bool {
	left := sl[i].Currency.String()
	right := sl[j].Currency.String()
	for idx := 0; idx < len(left) && idx < len(right); idx++ {
		if left[idx] > right[idx] {
			return false
		}
		if left[idx] < right[idx] {
			return true
		}
	}
	return len(left) < len(right)
}
func (sl subList) Swap(i, j int) { sl[i], sl[j] = sl[j], sl[i] }

type depSort goex.DepthRecords

func (ds depSort) Len() int { return len(ds) }
func (ds depSort) Less(i, j int) bool {
	return ds[i].Price > ds[j].Price // 降序排列，反着Less
}
func (ds depSort) Swap(i, j int) { ds[i], ds[j] = ds[j], ds[i] }

// CreateRobot :
func CreateRobot(exchange, id, key string, pair goex.CurrencyPair, volin, volout, volignore, unit, discount float64,
	interval int, uplimit, downlimit, diff, exit float64) *Robot {
	// var netTransport = &http.Transport{
	// 	Dial: (&net.Dialer{
	// 		Timeout:   10 * time.Second,
	// 		KeepAlive: 30 * time.Second,
	// 	}).Dial,
	// 	TLSHandshakeTimeout:   5 * time.Second,
	// 	ResponseHeaderTimeout: 10 * time.Second,
	// 	ExpectContinueTimeout: 1 * time.Second,
	// }
	// var client = &http.Client{
	// 	Timeout: time.Second * 10,
	// 	// Transport: netTransport,
	// }
	var client = &http.Client{}
	var exc goex.API
	switch exchange {
	case "coinex":
		exc = coinex.New(client, id, key)
	case "bigone":
		exc = bigone.New(client, id, key)
	case "fcoin":
		exc = fcoin.NewFCoin(client, id, key)
	default:
		log.Printf("错误！exchange param invalid, must be one of coinex/bigone/fcoin.")
		os.Exit(1)
		return nil
	}

	return &Robot{
		exchangeStr: exchange,
		exchange:    exc,
		unit:        unit,
		pair:        pair,
		volin:       volin,
		volout:      volout,
		volignore:   volignore,
		discount:    discount,
		interval:    interval,
		uplimit:     uplimit,
		downlimit:   downlimit,
		diff:        diff,
		exit:        exit,
	}
}

// Robot :
type Robot struct {
	running                  bool
	runTimes                 int64
	volin, volout, volignore float64
	discount                 float64
	unit                     float64
	pair                     goex.CurrencyPair
	exchangeStr              string
	exchange                 goex.API // *coinex.CoinEx
	interval                 int
	uplimit                  float64
	downlimit                float64
	diff                     float64
	exit                     float64
}

// Run :
func (robot *Robot) Run() {
	robot.running = true
	robot.runTimes++
	times := robot.runTimes
	for times == robot.runTimes && robot.running {
		robot.monitor()
		time.Sleep(time.Duration(robot.interval*10) * time.Millisecond)
	}
}

func floatToStr(val float64) string {
	return fmt.Sprintf("%f", val)
}

func (robot *Robot) monitor() error {
	account, err := robot.exchange.GetAccount()
	if nil != err {
		log.Printf("Robot.monitor - 错误！GetAccount failed : %v", err)
		if strings.Contains(err.Error(), "too quick, slow down") {
			log.Printf("服务器抱怨太快，休息5秒重试")
			time.Sleep(5 * time.Second)
		}
		return err
	}
	fmt.Printf("\n=========================当前资产 @ %s============================\n", time.Now().Format("2006-01-02 15:04:05"))

	var subs subList
	for _, sub := range account.SubAccounts {
		if sub.Amount > 0.0001 {
			subs = append(subs, sub)
		}
	}
	sort.Sort(subs)

	maxCol := 8
	for row := 0; row < (len(subs)+maxCol-1)/maxCol; row++ {
		for col := 0; col < maxCol && (col+maxCol*row) < len(subs); col++ {
			fmt.Printf(" %-9s", subs[row*maxCol+col].Currency.String())
		}
		fmt.Println()
		for col := 0; col < maxCol && (col+maxCol*row) < len(subs); col++ {
			fmt.Printf(" %-9.5f", subs[row*maxCol+col].Amount)
		}
		fmt.Println()
	}
	fmt.Printf("---------------------------------------------------------------------------------\n")

	orders, err := robot.exchange.GetUnfinishOrders(robot.pair)
	if nil != err {
		log.Printf("Robot.monitor - 错误！exchange.GetUnfinishOrders() failed : %v", err)
		return err
	}

	log.Printf("    挂单号    品种     方向  挂单价     挂单量\n")
	for _, order := range orders {
		log.Printf("    %s  %-7s %-4s  %-10.4f %-10.4f\n", order.OrderID2, order.Currency.String(), order.Side, order.Price, order.Amount)
	}
	fmt.Printf("=================================================================================\n\n")
	depth, err := robot.exchange.GetDepth(20, robot.pair)
	if nil != err {
		log.Printf("Robot.monitor - 错误！exchange.GetUnfinishOrders() failed : %v", err)
		return err
	}
	sort.Sort(depSort(depth.AskList)) // 统一设置成降序排列
	sort.Sort(depSort(depth.BidList)) // 统一设置成降序排列

	// 检查挂单，排位靠后就撤销; 买单和卖单都是价格低减排序
	closed := false
	pendingVol := 0.0
	sellingVol := 0.0
	highestBuy := 0.0
	lowestSell := 1000000000.0
	for idx := range orders {
		order := &orders[idx]
		accVol := 0.0
		if order.Side == goex.BUY {
			if order.Price >= highestBuy {
				highestBuy = order.Price
			}
			pendingVol += order.Amount
			for idx := range depth.BidList {
				buy := &depth.BidList[idx]
				if buy.Price < order.Price {
					break
				}
				accVol += buy.Amount

				if accVol > robot.volout || idx >= len(depth.BidList)-1 { //关单
					log.Printf("关闭买单 %s - %s", order.OrderID2, robot.pair)
					closed = true
					_, err = robot.exchange.CancelOrder(order.OrderID2, robot.pair)
					if nil != err {
						log.Printf("Robot.monitor - 错误！robot.exchange.CancelOrder(%s, %s) failed : %v",
							order.OrderID2, robot.pair, err)
					}
					break
				}
			}
		} else if order.Side == goex.SELL {
			if order.Price <= lowestSell {
				lowestSell = order.Price
			}
			sellingVol += order.Amount
			for idx := len(depth.AskList) - 1; idx >= 0; idx-- {
				sell := &depth.AskList[idx]
				if sell.Price > order.Price {
					break
				}
				accVol += sell.Amount

				if accVol > robot.volout || idx == 0 { // 关单
					log.Printf("关闭卖单 %s - %s", order.OrderID2, robot.pair)
					closed = true
					_, err = robot.exchange.CancelOrder(order.OrderID2, robot.pair)
					if nil != err {
						log.Printf("Robot.monitor - 错误！robot.exchange.CancelOrder(%s, %s) failed : %v",
							order.OrderID2, robot.pair, err)
					}
					break
				}
			}
		} else {
			log.Printf("Robot.monitor - 错误！Unknown order.Side : %d", order.Side)
		}
		if closed { // 有过关单，下次再试
			continue
		}
	}

	// coinex 控制难度
	if robot.exchangeStr == "coinex" {
		exc := robot.exchange.(*coinex.CoinEx)
		limit, cur, err := exc.GetDifficulty()

		if nil != err && cur > limit*robot.exit {
			log.Printf("Robot.Monitor - 提示：CoinEx暂停开仓 : 当前难度%.3f > 限制难度%.3f * 退出比例%.2f", cur, limit, robot.exit)
			return nil
		}
		log.Printf("CoinEx 当前难度%.3f < 限制难度%.3f * 退出比例%.2f", cur, limit, robot.exit)
	}

	// 尝试买
	accVol := 0.0 // 按量找到买价
	price := 0.0000000001

	for idx := range depth.BidList {
		buy := &depth.BidList[idx]
		if accVol+buy.Amount > robot.volin {
			price = buy.Price + robot.discount
			break
		}
		accVol += buy.Amount
	}

	if price < 0.00000000011 { // 用最后一个挂单
		if len(depth.BidList) > 0 {
			price = depth.AskList[len(depth.BidList)-1].Price + robot.discount
		}
	}

	if highestBuy < 0.000001 || // 没有挂买单
		price > highestBuy+robot.diff { // 或者价格明显有差距才开新单
		vol := 0.0
		if base, ok := account.SubAccounts[robot.pair.CurrencyB]; ok {
			vol = base.Amount / price * 0.99 // 避免无法开仓
		}
		if vol > robot.unit {
			vol = robot.unit
		}

		if vol < robot.volignore {
			log.Printf("剩余%s太少，不买 %f < %f",
				robot.pair.String(), vol, robot.volignore)
		} else if price > 0.0000000001 && price < robot.uplimit { // 用最后一个挂单
			log.Printf("尝试开仓买入 量:%f 价:%f %s", vol, price, robot.pair.String())
			_, err = robot.exchange.LimitBuy(floatToStr(vol), floatToStr(price), robot.pair)
			if nil != err {
				log.Printf("Robot.monitor - 错误！LimitBuy(%f, %f, %s) failed : %v",
					vol, price, robot.pair.String(), err)
			}
		} else {
			log.Printf("没有买单参考，不开买单")
		}
	}

	// 尝试卖
	accVol = 0.0 // 按量找到卖价
	price = 1000000000.0

	for idx := len(depth.AskList) - 1; idx >= 0; idx-- {
		sell := &depth.AskList[idx]
		if accVol+sell.Amount > robot.volin {
			price = sell.Price - robot.discount
			break
		}
		accVol += sell.Amount
	}
	if price >= 1000000000.0-0.1 { // 用最后一个挂单
		if len(depth.AskList) > 0 {
			price = depth.AskList[len(depth.AskList)-1].Price - robot.discount
		}
	}

	if lowestSell < 1000000000.0-1.0 && // 已经有卖单
		price < lowestSell-robot.diff { //价格差距过近
		return nil // 不新卖
	}
	vol := 0.0
	if refer, ok := account.SubAccounts[robot.pair.CurrencyA]; ok {
		vol = refer.Amount * 0.99 // * 0.99 是因为有时候因为极小的差值不开仓
	}
	if vol > robot.unit {
		vol = robot.unit
	}

	if vol < robot.volignore {
		log.Printf("剩余%s太少，不卖 %f < %f",
			robot.pair.String(), vol, robot.volignore)
		return nil
	}
	if price < 1000000000.0-0.1 && price > robot.downlimit {
		log.Printf("尝试开仓卖出 量:%f 价:%f %s", vol, price, robot.pair.String())
		_, err = robot.exchange.LimitSell(floatToStr(vol), floatToStr(price), robot.pair)
		if nil != err {
			log.Printf("Robot.monitor - 错误！LimitSell(%f, %f, %s) failed : %v",
				vol, price, robot.pair.String(), err)
		}
	} else {
		log.Printf("没有挂单作参考，不开卖单")
	}
	return err
}

// Stop :
func (robot *Robot) Stop() {
	robot.running = false
	robot.runTimes++
}

func welcome() {
	fmt.Printf(`*自我介绍
    大家好！我是微信公众号 “呆萌小机器人”和“投资呆萌小机器人”的作者，
    也是eosforce节点jiqix的维护者，也是这个机器人的作者。
    这次因为时间紧张，来不及挂到微信后台，以后类似的工具我会直接通过微
    信提供，所以希望您能关注本公众号，谢谢关注我的公众号！

    如果您在EOS创世名单中，请访问eosforce.io下载钱包，灌入私钥，即可领取
    EOSForce链上的EOS。eosforce投票有分红，这是您投票的最大原因和动力！
    如果担心安全隐患，可以在EOS主网上新建帐号，把账户余额转给信号，再领
    取EosForce。
    jiqix目前为投票者分红95%%，所以必须为jiqix投票，谢谢支持！

    如果有其他疑问，请在github或者微信后台和我联系。

*安全申明
    本人保证github上发布的版本没有恶意代码。但是为了安全起见，请用虚拟机，
    或者完全不带敏感信息的电脑运行此程序。如果出现任何安全事件，本程序概
    不负责。
    另:可以考虑买丐版云服务器运行本程序，在腾讯官方价格基础上，联系微信号
    wshinewmm还有进一步优惠。

*使用说明
    不带参数运行可以看到帮助，我现在测试使用的参数如下，策略不同、品种不同的话需要调整。
    如果纯粹以刷量为目的，volin,volout调小，discount调大。
    如果试图主动赢利，把volin,volout调大，discount调小。
	如下是bigone上交易ONEUSDT的一个例子：
 robot.go -id xxxxxxxxxxxxxxxxxxxxxxxx\
   -key yyyyyyyyyyyyyyyyyyyyyyyyyyyyy \
   -discount 0.001 \
   -lsym ONE \
   -rsym USDT \
   -unit 30 \
   -volignore 10\
   -volin 3000.0\
   -volout 50000.0 \
   -interval 100 \
   -diff 0.1 \
   -exchange bigone

*归档
    目前只上传程序在此，代码待挖矿结束后提供。
    https://github.com/gpmn/coinex_robot
`)
}

func main() {
	welcome()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	id := flag.String("id", "", "API ID。")
	key := flag.String("key", "", "API Secret。")
	lsym := flag.String("lsym", "", "交易对中的左值，比如BTCUSDT中的BTC。")
	rsym := flag.String("rsym", "USDT", "默认USDT，交易对中的右值，即BTCUSDT中的USDT。")
	volin := flag.Float64("volin", 0.0, "在累计挂单量到这个阈值之前挂单，按lsym计。")
	volout := flag.Float64("volout", 0.0, "之前挂单量超过这个阈值就撤销重新挂，建议三倍volin，按lsym计。")
	volignore := flag.Float64("volignore", 0.0, "如果小于这个量，就不开仓，按lsym计")
	unit := flag.Float64("unit", 0.0, "单次开仓的最大量，lsym计。")
	discount := flag.Float64("discount", 0.0001, "在找到的价格点基础上调整多少入场，默认0.0001，按rsym计。")
	interval := flag.Int("interval", 100, "每一轮检查间隔tick数，1tick==10毫秒,即百分之1秒。")
	exchange := flag.String("exchange", "coinex", "默认coinex，支持coinex/bigone/fcoin。需要更多交易所可以向作者反馈。")
	uplimit := flag.Float64("uplimit", 10000000000000.0, "大于这个极端值不开仓，lsym/rsym计算，即价格，请酌情配置。")
	downlimit := flag.Float64("downlimit", 1.0/10000000000000.0, "小于这个极端值不开仓，lsym/rsym计算，即价格，请酌情配置。")
	diff := flag.Float64("diff", 0.0, "挂单最小间隔，默认等于3倍discount")
	exit := flag.Float64("exit", 0.90, "只对coinex有效。按照coinex公布的当前进度/难度，如果达到这个比例就停止挖矿。由于coinex这个数据有延迟，所以用exit阈值控制，提前停工，减少无效挖矿。")
	flag.Parse()

	if *id == "" {
		flag.Usage()
		log.Printf("error - must supply id of your account")
		os.Exit(1)
		return
	}

	if *key == "" {
		flag.Usage()
		log.Printf("error - must supply key of your account")
		os.Exit(1)
		return
	}

	if *lsym == "" {
		flag.Usage()
		log.Printf("error - must supply lsym for trade")
		os.Exit(1)
		return
	}

	if *rsym == "" {
		flag.Usage()
		log.Printf("error - must supply rsym for trade")
		os.Exit(1)
		return
	}

	if *discount < 0 {
		flag.Usage()
		log.Printf("error - must supply positive discount")
		os.Exit(1)
		return
	}

	if *unit <= 0 {
		flag.Usage()
		log.Printf("error - must supply positive unit")
		os.Exit(1)
		return
	}

	if *volignore < 0 || *volignore >= *unit {
		flag.Usage()
		log.Printf("error - must supply positive volignore to trade, and must lower than unit")
		os.Exit(1)
		return
	}

	if *volin <= 0 {
		flag.Usage()
		log.Printf("error - must supply positive volin to trade")
		os.Exit(1)
		return
	}

	if *volout <= 0 || *volout < *volin*2 {
		flag.Usage()
		log.Printf("error - must supply positive volout to trade, and must volout greater than 2*volin")
		os.Exit(1)
		return
	}

	if *diff <= 0 {
		*diff = math.Abs(3 * (*discount))
	}

	if *exit <= 0 {
		*exit = 0.9
	}

	pair := goex.CurrencyPair{
		CurrencyA: goex.Currency{Symbol: strings.ToUpper(*lsym), Desc: ""},
		CurrencyB: goex.Currency{Symbol: strings.ToUpper(*rsym), Desc: ""},
	}

	robot := CreateRobot(*exchange, *id, *key, pair, *volin, *volout, *volignore, *unit, *discount, *interval, *uplimit, *downlimit, *diff, *exit)

	for idx := 0; idx < 100; idx++ {
		time.Sleep(time.Millisecond * 10)
		fmt.Printf(".")
	}
	fmt.Println(".")

	robot.Run()
}
