/*****************************************************************************
File name: server.go
Description: web server
Author:
Version: V1.0
Date: 2018/06/14
History:
*****************************************************************************/
package web
import (
	"github.com/gin-gonic/gin"
	"fmt"
	"platform/om/log"
	"platform/db"
	"platform/config"
	"github.com/DeanThompson/ginpprof"
)
import (
	_ "net/http/pprof"
	"platform/bill"
	"platform/blockchain/eth"
)


func ServerInit() {
	//设置server参数
	//key_path := config.Gconfig.Httpscfg.Keypath 				   //https证书加载路径
	//cert_path := config.Gconfig.Httpscfg.Certpath
	addr := fmt.Sprintf(":%d", config.Gconfig.Servercfg.Port)

	router := gin.Default()            							  // 注册一个默认的路由器
	router.POST("/wallet", postInvoke) 							  // 注册POST回调
	router.Use(LimitMiddleware(config.Gconfig.Servercfg.MaxConn))
	ginpprof.Wrapper(router)									  //pprof注册，
	//err := router.RunTLS(addr, cert_path, key_path) 			  //https启动
	err := router.Run(addr) 									  //http启动
	if err != nil {
		log.Error("Authentication failed : %s", err.Error())
	}
}


/*************************************************
Function: postInvoke
Description: post 路由器，context在其他goroutine中
不能修改，只能传递只读的context！
Author:
Date: 2018/06/14
History:
*************************************************/
func postInvoke(c *gin.Context){
	//common.Dispather.AddJob(c.Copy())
	//fmt.Printf("postInvoke IN !!!\n")
	manage(c)
	//c.String(http.StatusOK, retmsg)


}

/*************************************************
Function: LimitMiddleware
Description: 限制gin最大链接个数
Author:
Date: 2018/07/12
History:
*************************************************/
func LimitMiddleware(limit int) gin.HandlerFunc {
	// create a buffered channel with 1000 spaces
	semaphore := make(chan struct{}, limit)
	return func (c *gin.Context) {
		select {
		case semaphore  <- struct{}{}: // Try putting a new val into our semaphore
			// Ok, managed to get a space in queue. execute the handler
			c.Next()

			// Don't forget to release a handle
			<-semaphore
		default:
			// Buffer full, so drop the connection. Return whatever status you want here
			return
		}
	}
}

/*************************************************
Function: manage
Description: web模块POST消息处理
Author: failymao
Date: 2018/06/14
History:
*************************************************/
func manage (c *gin.Context) {

	//用户鉴权
	//示例：http://192.168.0.191:10086/wallet?app_id=0147&sign=38592ef2a9256617e7b7110480df2510c8be75d3903c8053c464ac61708c04cb&timestamp=2018-06-20 14:07:33&method=TransferAccounts
	result := db.Authentication(c)
	if MSG_SUCCESS != result {
		fmt.Printf("authentication failure:%s\n", result)
		errencode(c, result)
		return
	}

	//公共参数解析
	method := c.Query("method")
	//消息类型判断
	if method_create == method {
		//开户消息处理
		var reqInfo AccountBody
		retmsg := accountdecode(&reqInfo, c)
		if MSG_SUCCESS != retmsg {
			fmt.Printf("accountdecode return Msgcode:%s\n", retmsg)
			errencode(c, retmsg)
			return
		}

		//参数校验
		retmsg = accountcheck(&reqInfo, method)
		if MSG_SUCCESS != retmsg {
			fmt.Printf("accountcheck failed with Msgcode:%s\n", retmsg)
			errencode(c, retmsg)
			return
		}

		//开户
		var account, password string
		//start := time.Now()
		retmsg = openaccount(&reqInfo, &account, &password, c.Query("app_id"), c.Query("timestamp"))
		//defer log.Debug(">>>>开户总:%v", time.Since(start))
		if MSG_SUCCESS != retmsg {
			fmt.Printf("openaccount failed with Msgcode:%s\n", retmsg)
			errencode(c, retmsg)
			return
		}
		if usernormal == reqInfo.Usertype {
			bill.CreateAcct(reqInfo.Username)
		}
		if userENT == reqInfo.Usertype {
			bill.CreateEntAcct(reqInfo.Username)
		}
		//rechargetest()
		//响应构造
		accountencode(c, reqInfo, retmsg, account, password)

	} else if method_transfer == method {
		//转账消息处理
		var reqTransInfo TransferBody
		retmsg := transferdecode(&reqTransInfo, c)
		if MSG_SUCCESS != retmsg {
			fmt.Printf("Transferdecode fail:%s\n", retmsg)
			errencode(c, retmsg)
			return
		}

		//参数校验
		retmsg = transfercheck(&reqTransInfo)
		if MSG_SUCCESS != retmsg {
			fmt.Printf("transfercheck failed with Msgcode:%s\n", retmsg)
			errencode(c, retmsg)
			return
		}

		//获取本次交易的流水号
		serial_num := db.GetSerialNum()
		//log.Debug("交易序列号：%s", serial_num)
		//转账
		var code transfercode
		var txhash string="0000" //交易的hash值
		switch reqTransInfo.TransferCurrency {
		case "0":   //查找以太币账户余额
			//当transflag为true时，需要调用以太坊智能合约实现转账功能
			if  true == config.Gconfig.Ethcfg.Transflag {
				//转账 调用以太坊智能合约，实现转账功能
				//start := time.Now()
				//以太坊公链转账时，需要返回交易的hash值  add 2018-7-4  shangwj
				txhash, retmsg = starttransfer(&reqTransInfo, &serial_num)
				//defer log.Debug(">>>>转账总:%v", time.Since(start))
				if MSG_SUCCESS != retmsg {
					fmt.Printf("starttransfer failed with Msgcode:%s\n", retmsg)
					errencode(c, retmsg)
					return
				}
			}
		case "1":   //查找比特币账户余额
			//todo:比特币币种查询，现在只支持以太币，参数检查时已将比特币过滤，此处留做比特币支持时填补
			fmt.Println("Bitcoin transfer operations is to be supported.")
			return
		}

		//转账 数据库记账操作
		//add 2018-7-4  shangwj  交易明细表增加 以太坊交易hash值  txhash
		//1.转账需要上以太坊公链时，需要更新交易的hash值 到数据库的交易明细表中
		//2.转账不上以太坊公链时，txhash为空
		retmsg = db.TransBetweenUsers(serial_num, reqTransInfo.SourceUser, reqTransInfo.DestUser, reqTransInfo.TransferAmount, reqTransInfo.TransferCurrency, txhash)
		if MSG_SUCCESS != retmsg {
			fmt.Printf("transfercheck failed with Msgcode:%s\n", retmsg)
			errencode(c, retmsg)
			return
		}
		bill.TransferAcct(reqTransInfo.SourceUser, reqTransInfo.DestUser, reqTransInfo.TransferAmount)
		//响应构造
		transferencode(c, code, retmsg)

	} else if method_query == method {
		//查询消息处理
		var reqInfo AccountBody
		err := c.BindJSON(&reqInfo)
		if err != nil {
			log.Error("handle http json:%s", err.Error())
			errencode(c, MSG_PARATYPEERR)
			return
		}

		//参数校验
		retmsg := accountcheck(&reqInfo, method)
		if MSG_SUCCESS != retmsg {
			fmt.Printf("querybalance failed with Msgcode:%s\n", retmsg)
			errencode(c, retmsg)
			return
		}

		//查询
		var balance int64
		//start := time.Now()
		retmsg = querybalance(&reqInfo, &balance)
		//defer log.Debug(">>>>查询总:%v", time.Since(start))
		if MSG_SUCCESS != retmsg {
			fmt.Printf("querybalance failed with Msgcode:%s\n", retmsg)
			errencode(c, retmsg)
			return
		}

		//响应构造
		queryencode(c, retmsg, balance)
	} else if method_cancellation == method {
		//销户消息处理
		var reqCancellation CancellationBody
		retmsg := cancellationdecode(&reqCancellation, c)
		if MSG_SUCCESS != retmsg {
			fmt.Printf("Cancellationdecode failed with Msgcode:%s\n", retmsg)
			log.Error("Cancellationdecode failed with Msgcode:%s\n", retmsg)
			errencode(c, retmsg)
			return
		}

		//参数校验
		retmsg = cancellationcheck(&reqCancellation)
		if MSG_SUCCESS != retmsg {
			fmt.Printf("cancellationdecode failed with Msgcode:%s\n", retmsg)
			errencode(c, retmsg)
			return
		}

		//销户
		retmsg = db.AccountCancellation("UserInfo", reqCancellation.Username, reqCancellation.Password)
		if MSG_SUCCESS != retmsg {
			fmt.Printf("cancellation failed")
			errencode(c, retmsg)
			return
		}
		//清UserInfo本地缓存
		eth.DelUserInfo(reqCancellation.Username)
		//账单记录
		bill.DeleteAcct(reqCancellation.Username)
		//响应构造
		cancellationcode(c, retmsg, reqCancellation.Username)

	} else if method_withdraw == method {
		var WithdrawInfo WithdrawBody              //提现客户端发送参数实例化
		retmsg := withdrawdecode(&WithdrawInfo, c) //解析参数
		if MSG_SUCCESS != retmsg {
			fmt.Printf("withdrawdecode fail:%s\n", retmsg)
			errencode(c, retmsg)
			return
		}

		retmsg = withdrawcheck(&WithdrawInfo) //请求参数校验
		if MSG_SUCCESS != retmsg {
			fmt.Printf("withdrawcheck failed with Msgcode:%s\n", retmsg)
			errencode(c, retmsg)
			return
		}

		var code withdrawcode           //返回客户前消息处理
		serial_num := db.GetSerialNum() //生成交易流水帐号
		var wd_hash string = "0000"     //默认不上链hash:0000
		switch WithdrawInfo.TransferCurrency {
		case "0":
			if true == config.Gconfig.Ethcfg.Transflag { // flag为true时，上链操作
				wd_hash, retmsg = withdrawtransfer(&WithdrawInfo, &serial_num)
				if MSG_SUCCESS != retmsg {
					log.Error("starttransfer failed with Msgcode:%s\n", retmsg)
					errencode(c, retmsg)
					return
				}
			}
		case "1":
			fmt.Println("Bitcoin transfer operations is to be supported.")
			return

		}
		retmsg = db.TransWithdraw(serial_num, WithdrawInfo.SourceUser,
			WithdrawInfo.ExternalAccount, WithdrawInfo.TransferAmount, WithdrawInfo.TransferCurrency, wd_hash)

		if MSG_SUCCESS != retmsg {
			log.Error("withdrawcheck failed with Msgcode:%s\n", retmsg)
			errencode(c, retmsg)
			return
		}

		//bill.TransferAcct(WithdrawInfo.SourceUser, WithdrawInfo.ExternalAccount, WithdrawInfo.TransferAmount)
		withdrawencode(c, code, retmsg) //返回客户端信息

	} else {
		fmt.Printf("unsupport method:%s\n", method)
		errencode(c, MSG_UNSUPPORT)
		return
	}
}