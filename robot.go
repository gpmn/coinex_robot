package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/nntaoli-project/GoEx"
	"github.com/nntaoli-project/GoEx/coinex"
)

// CreateRobot :
func CreateRobot(id, key string, pair goex.CurrencyPair, volin, volout, volignore, unit, discount float64) *Robot {
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	var client = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}

	return &Robot{
		exchange:  coinex.New(client, id, key),
		unit:      unit,
		pair:      pair,
		volin:     volin,
		volout:    volout,
		volignore: volignore,
		discount:  discount,
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
	exchange                 *coinex.CoinEx
}

// Run :
func (robot *Robot) Run() {
	robot.running = true
	robot.runTimes++
	times := robot.runTimes
	for times == robot.runTimes && robot.running {
		robot.monitor()
		time.Sleep(5 * time.Second)
	}
}

func floatToStr(val float64) string {
	return fmt.Sprintf("%f", val)
}

func (robot *Robot) monitor() error {
//挖矿结束后补充，^_^
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
	jiqix目前为投票者分红95%，所以必须为jiqix投票，谢谢支持！

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

robot -id xxxxxxxxxxxxxxx \
   -key xxxxxxxxxxxxxxxx \
   -discount 0.01 \
   -lsym ETH \
   -rsym USDT \
   -unit 0.01 \
   -volignore 0.005\
   -volin 0.1 \
   -volout 0.3

*归档
    目前只上传程序在此，代码待挖矿结束后提供。
    https://github.com/gpmn/coinex_robot
`)
}

func main() {
	welcome()
	for idx := 0; idx < 100; idx++ {
		time.Sleep(time.Millisecond * 100)
		fmt.Printf(".")
	}
	fmt.Println(".")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	id := flag.String("id", "", "API ID。")
	key := flag.String("key", "", "API Secret。")
	lsym := flag.String("lsym", "", "交易对中的左值，比如BTCUSDT中的BTC。")
	rsym := flag.String("rsym", "USDT", "默认USDT，交易对中的右值，即BTCUSDT中的USDT。")
	volin := flag.Float64("volin", 0.0, "在累计挂单量到这个阈值之前挂单，按lsym计。")
	volout := flag.Float64("volout", 0.0, "之前挂单量超过这个阈值就撤销重新挂，建议三倍volin，按lsym计。")
	volignore := flag.Float64("volignore", 0.0, "如果小于这个量，就不开仓，按lsym计算")
	unit := flag.Float64("unit", 0.0, "单次开仓的最大量，lsym计。")
	discount := flag.Float64("discount", 0.0001, "在找到的价格点基础上调整多少入场，默认0.0001，按rsym计。")

	flag.Parse()

	if *id == "" {
		log.Printf("error - must supply id of your account")
		flag.Usage()
		os.Exit(1)
		return
	}

	if *key == "" {
		log.Printf("error - must supply key of your account")
		flag.Usage()
		os.Exit(1)
		return
	}

	if *lsym == "" {
		log.Printf("error - must supply lsym for trade")
		flag.Usage()
		os.Exit(1)
		return
	}

	if *rsym == "" {
		log.Printf("error - must supply rsym for trade")
		flag.Usage()
		os.Exit(1)
		return
	}

	if *discount < 0 {
		log.Printf("error - must supply positive discount")
		flag.Usage()
		os.Exit(1)
		return
	}

	if *unit <= 0 {
		log.Printf("error - must supply positive unit")
		flag.Usage()
		os.Exit(1)
		return
	}

	if *volignore < 0 || *volignore >= *unit {
		log.Printf("error - must supply positive volignore to trade, and must lower than unit")
		flag.Usage()
		os.Exit(1)
		return
	}

	if *volin <= 0 {
		log.Printf("error - must supply positive volin to trade")
		flag.Usage()
		os.Exit(1)
		return
	}

	if *volout <= 0 || *volout < *volin*2 {
		log.Printf("error - must supply positive volout to trade, and must volout greater than 2*volin")
		flag.Usage()
		os.Exit(1)
		return
	}

	pair := goex.CurrencyPair{
		CurrencyA: goex.Currency{Symbol: strings.ToUpper(*lsym), Desc: ""},
		CurrencyB: goex.Currency{Symbol: strings.ToUpper(*rsym), Desc: ""},
	}

	robot := CreateRobot(*id, *key, pair, *volin, *volout, *volignore, *unit, *discount)
	robot.Run()
}
