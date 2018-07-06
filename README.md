目前只提供二进制程序，秉承开源精神，会在挖矿活动结束后时间以后上传源码。

目前提供了windows/linux的64位版本，以及mac的32/64版本。解压后请注意核对md5校验码，防止被篡改。  

目前票数已经够了，谢谢大家支持！我也不太小气了，v1.5开始随便下载吧，大家随便用。其实这个刷量机器人也没啥技术含量。只是希望大家在启动的时候多看几眼我的广告，所以还是等等再放源码上来吧 (^_^)

大家数到大钱的话，可以给我EOS小费打赏，帐号：testgambling，谢谢大家！

## 終章  
	挖矿红利逐渐结束了，现在的coinex加上返佣、分红勉强打平。参数调整的好的话，可能会小有收益。   
	我这个程序也就不用再保密了，大家可以自行编译。大家以后有什么意见、建议、新策略，欢迎和我交流。  
	
![太平一犬](https://github.com/gpmn/coinex_robot/raw/master/webwxgetmsgimg.jpg)
	
## 修改历史  

	v1.1   增加了一个interval参数，默认是3，即3秒尝试一次交易。有需要的话可以改成1，甚至0。  
	v1.2   修改了一个bug，在给定的volin过大的时候，老版本会尝试在100000000卖，0.000000001买，而且不会自动关单，修改此错误。 
	v1.2.1 压缩包密码错了，我自己都解不开，重新传一个，内容不变，md5校验值不变。
	v1.3   1).interval参数从1秒为单位，改为1tick即百分之一秒为单位。
	       2).新增一个exchange参数，支持fcoin、bigone、coinex，默认coinex。  
	v1.4   解决bug -- 会连续关单多次，导致被服务器流控;同时在被服务器流控后，主动延时5秒钟后再重新工作。
	v1.5   增加了开单上限、下限控制，避免极端行情入市交易被割韭菜。  
	v1.6   1).由于bigone接口变更，之前的版本不能交易bigone，重新适配bigone接口。  
	       2).新增diff参数，用以控制开仓间隔，默认三倍discount。  
	       3).美化了一下打印  
	       4).之前连续关单导致被流控的问题没有彻底解决，这次一起处理掉。  
	v1.7   CoinEx增加难度控制，达到难度后就停止开仓，直到难度更新。
	v1.8   加了一个exit参数，用以控制暂停挖矿的时间。因为据矿工反映，coinex的难度api有明显延迟，为了避免超挖，所以得在api返回的当前难度/总难度的基础上，提前退出。
	v1.8.1 1.8的exit没有生效，这个版本用上了。

## 自我介绍  

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
	
## 安全申明  

	本人保证github上发布的版本没有恶意代码。但是为了安全起见，请用虚拟机，
	或者完全不带敏感信息的电脑运行此程序。如果出现任何安全事件，本程序概
	不负责。  
	另:可以考虑买丐版云服务器运行本程序，在腾讯官方价格基础上，联系微信号
	wshinewmm还有进一步优惠。  
	
## 使用说明  

    不带参数运行可以看到帮助，我现在测试使用的参数如下，策略不同、品种不同的话需要调整。  
    如果纯粹以刷量为目的，volin,volout调小，discount调大。  
    如果试图主动赢利，把volin,volout调大，discount调小。  
    
    ./robot -id XXXXXXXXXXXXX \
        -key YYYYYYYYYYYYYYY \
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

   
## 归档  

    目前只上传程序在此，代码待挖矿结束后提供。  
    https://github.com/gpmn/coinex_robot  
    
## 命令行参数  

    -discount float    
        在找到的价格点基础上调整多少入场，默认0.0001，按rsym计。 (default 0.0001)  
    -exchange string  
        one of coinex/bigone/fcoin, will support more later (default "coinex")  
    -id string  
        API ID。  
    -interval int  
        每一轮检查间隔tick数，1tick==10毫秒,即百分之1秒。 (default 100)    
    -key string  
        API Secret。  
    -lsym string  
        交易对中的左值，比如BTCUSDT中的BTC。  
    -rsym string  
        默认USDT，交易对中的右值，即BTCUSDT中的USDT。 (default "USDT")  
    -unit float  
        单次开仓的最大量，lsym计。  
    -volignore float  
        如果小于这个量，就不开仓，按lsym计算  
    -volin float  
        在累计挂单量到这个阈值之前挂单，按lsym计。  
    -volout float  
        之前挂单量超过这个阈值就撤销重新挂，建议三倍volin，按lsym计。  
    -uplimit float
        大于这个极端值不开仓，lsym/rsym计算，即价格，请酌情配置。 (default 1e+13)
    -downlimit float
        小于这个极端值不开仓，lsym/rsym计算，即价格，请酌情配置。 (default 1e-13)
    -diff float
        挂单最小间隔，默认等于3倍discount。  
    -exit float
       只对coinex有效。按照coinex公布的当前进度/难度，如果达到这个比例就停止挖矿。由于coinex这个数据有延迟，所以用exit阈值控制，提前停工，减少无效挖矿。 (default 0.9)


## 使用疑问  

#### 接入coinex，报告tonce超过服务器时间1分钟以上  
    通常是时间和服务器不同步导致的。linux打开ntp服务，windows点击右下时钟小窗口，调整为自动设置时间。  
    
 #### 接入bigone报告connection error  
    bigone被封，需要翻墙。windows下用set如下设置，代理地址酌情修改，可以加到批处理开头：  
    set http_proxy=http://127.0.0.1:1189  
    set https_proxy=http://127.0.0.1:1189  
    Linux下把set改成export，加到脚本开头。  
    
 #### 刷亏了？？  
    首先请读一下策略说明。  
    我们这个策略主要目的不是赢利，而是争取在亏得不多的前提下刷量。  
    如果要刷量，把volin，volout都调小。如果要减小交易亏损，把volin，volout调大。  
    而且，两种情况下，unit都不能太大。  
    volin,volout太小，会导致我们这个策略从被动等着被大单扫，变成主动做市的一方，不符合本来的意图。  
    随着coinex红利减小，建议稍微调大volin和voluot，增加安全垫。  
[策略说明](http://8btc.com/thread-93841-1-1.html)  

