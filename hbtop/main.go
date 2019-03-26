package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gpmn/sheep/huobi"
	"github.com/gpmn/sheep/proto"
)

func main() {
	fmt.Printf(`
应朋友要求写的，我没测试过，也不保证能用。
用于帮着大家抢TOP，每秒抢10次，每次抢100$。不要多开，因为限流10s 100次请求。
使用前需要到火币开通api/secret，api/secret只给交易权限即可。
最好自己编译，最好假设我提供的exe有安全问题。

必须传入api secret参数。默认是市价单用ht买入100ht的top。

用法如下：
	./hbtop -api xxxxx -sec xxxxxx 

此外还有quote,base,amt三个参数,比如上述命令等于：
	./hbtop -api xxxxx -sec xxxxxx -base top -quote ht -amt 100

测试的话，钱别放多了，要不然很快买满：
	./hbtop -api xxxxx -sec xxxxxx -base btc -quote usdt -usd 10

通常来说测试买入btcOk，也就能买top了。
`)

	api := flag.String("api", "", "你的交易api")
	secret := flag.String("sec", "", "你的交易secret")
	amt := flag.Float64("amt", 100.0, "按quote单次买入数量,默认100$")
	base := flag.String("base", "top", "想买入的品种，默认top")
	quote := flag.String("quote", "ht", "计价币种，默认ht")

	flag.Parse()
	if *api == "" || *secret == "" {
		log.Printf("must give api and secret param")
		flag.Usage()
		os.Exit(1)
	}
	hbex, err := huobi.NewHuobi(*api, *secret)
	if nil != err {
		log.Printf("Trader.tradeRoutine - huobi.NewHuobi() failed : %v", err)
		return
	}
	for {
		time.Sleep(time.Millisecond * 100)
		ret, err := hbex.OrderPlace(&proto.OrderPlaceParams{
			Price:           0,
			Amount:          *amt,
			BaseCurrencyID:  *base,  // E.g. btc
			QuoteCurrencyID: *quote, // e.g. usdt
			Type:            "buy-market",
		})

		if nil != err {
			log.Printf("hbex.OrderPlace - failed : %v", err)
			continue
		}
		log.Printf("hbex.OrderPlace - OK, order id : %v", ret)
	}
}
