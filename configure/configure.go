package configure

// 全局配置信息

var (
	// RPCAddr rpc地址
	RPCAddr string
	// Port 服务端口号
	Port string
	// Debug 是否开启debug模式
	Debug bool
	//NotifyAddr 交易通知地址
	NotifyAddr = "http://www.tigerft.com/fag/recv?service=deposit"
	// StartNotify 是否启动通知
	StartNotify bool 
	// StartTxCallBack 是否启动通知
	StartTxCallBack bool
	// NotifyCoin 那些币种启动通知用逗号隔开 btc,eth,
	NotifyCoin string
	// EthSuccessNumber eth验证成功的区块数量
	EthSuccessNumber int64 = 6
	// IsProduce 是否是生产环境
	IsProduce bool 

)

// NeoAsset NEO 资产类型
var NeoAsset = "c56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b"

// GasAsset GAS 资产
var GasAsset = "602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7"