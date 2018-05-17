package main

import (
	"flag"
	"neo-dev/cmd/coinsrv/global"
	"neo-dev/cmd/coinsrv/neoanderls"
	"neo-dev/configure"
)

func init() {
	flag.StringVar(&configure.RPCAddr, "rpcaddr", "http://127.0.0.1:20332", "rpc地址")
	flag.StringVar(&configure.Port, "port", "6000", "服务端监听端口")
	flag.BoolVar(&configure.Debug, "debug", false, "开启debug模式")
	flag.StringVar(&configure.NotifyCoin, "notifyCoin", "neo", "定时器支持币种")
	flag.BoolVar(&configure.StartTxCallBack, "startTxCallBack", configure.StartTxCallBack, "开启定时器")
	// 数据库
	flag.StringVar(&configure.DBAddr, "dbaddr", configure.DBAddr, "数据库ip:port")
	flag.StringVar(&configure.DBName, "dbname", configure.DBName, "数据库名称")
	flag.StringVar(&configure.DBUser, "dbuser", configure.DBUser, "数据库用户名")
	flag.StringVar(&configure.DBPasswd, "dbpasswd", configure.DBPasswd, "数据库密码")
	// 交易通知配置
	flag.StringVar(&configure.NotifyAddr, "notifyaddr", configure.NotifyAddr, "交易通知地址")
	flag.StringVar(&configure.NotifyCoin, "notifycoin", configure.NotifyCoin, "那些币种启动交易通知")
	flag.BoolVar(&configure.StartNotify, "startnotify", configure.StartNotify, "是否启动交易通知")

	flag.BoolVar(&configure.IsProduce, "produce", configure.IsProduce, "是否是生产环境")

	flag.Parse()
}

func main() {
	webapp, err := global.NewWebApp()
	if err != nil {
		return
	}
	neoRouter := webapp.Group("/neo")
	{
		// neoRouter.Use(midware.VerfSign)
		neoRouter.POST("/account/query", neoanderls.Balance)
		neoRouter.POST("/account/create", neoanderls.Create)
		neoRouter.POST("/trade/send", neoanderls.SendTransaction)
		neoRouter.POST("/trade/query", neoanderls.TransactionInfoByHash)
		// neoRouter.GET("/account/list", neoanderls.AccountList)
	}

	webapp.Logger.Fatal(webapp.Start(":" + configure.Port))
}
