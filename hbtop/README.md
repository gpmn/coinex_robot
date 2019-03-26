应朋友要求写的，我没测试过，也不保证能用。  

用于帮着大家抢TOP，每秒抢10次，每次抢100$。不要多开，因为限流10s 100次请求。  

必须传入api secret参数。默认是市价单用ht买入100ht的top。   

使用前需要到火币开通api/secret，api/secret只给交易权限即可。  

**最好自己编译，最好假设我提供的exe有安全问题。**  

用法如下：  

	./hbtop -api xxxxx -sec xxxxxx   

此外还有quote,base,amt三个参数,比如上述命令等于：  

	./hbtop -api xxxxx -sec xxxxxx -base top -quote ht -amt 100  

测试的话，钱别放多了，要不然很快买满：  

	./hbtop -api xxxxx -sec xxxxxx -base btc -quote usdt -usd 10  

通常来说测试买入btcOk，也就能买top了。


