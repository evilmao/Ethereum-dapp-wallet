/*****************************************************************************
File name: main.go
Description: web模块初始化及启动
Author: failymao
Version: V1.0
Date: 2018/06/14
History:
*****************************************************************************/
package main

import (
	"fmt"
	"os"
	"platform/bill"
	"platform/blockchain/eth"
	"platform/config"
	"platform/om/log"
	"platform/web"
	"time"
)

/*************************************************
Function: main
Description: web模块初始化及gin启动
Author: failymao
Date: 2018/06/14
History:
*************************************************/
func main() {
	//new logger
	err := log.New(config.Gconfig.Logcfg.LogLev,
		config.Gconfig.Logcfg.LogPath)
	if err != nil {
		fmt.Sprintf("create logger:%s", err.Error())
		os.Exit(0)
	}
	//new bill
	err = bill.New(config.Gconfig.Billcfg.BillPath)
	if err != nil {
		log.Fatal("create bill:%s", err.Error())
	}
	//ethrpc client create
	err = eth.RPCNew(config.Gconfig.Ethcfg.RPCUrl)
	if err != nil {
		log.Fatal("create eth rpc client:%s", err.Error())
	}
	//准备好已经注册的用户
	count := eth.PrepareUser(config.Gconfig.UsrPrepared)
	log.Debug("need %d user, prepared %d user", config.Gconfig.UsrPrepared, count)

	config.ConfigRuntime()

	eth.EthInfoUpdate(time.Duration(config.Gconfig.UpdateInt) * time.Second)
	eth.AsyncTranfer() //异步 上链交易确认  shangwj  2018-7-9
	web.ServerInit()

}
