/*****************************************************************************
File name: adapter.go
Description: web模块
Author: failymao
Version: V1.0
Date: 2018/06/14
History:
*****************************************************************************/
package web

import (
	"platform/blockchain/eth"
	"platform/om/log"
	"strconv"
	"platform/config"
	"fmt"
)

type transfercode struct {
	sub_code string
	sub_msg  string
}

type withdrawcode struct {
	sub_code string
	sub_msg string
}

func openaccount (reqInfo *AccountBody, account *string, password *string, appid, ctime string)  string{
	//调用接口
	//*account = "0x99A3417b2FCa03aaF1F36f4bC7e6fdD8639B600f"
	//db.NewUserDBInsert(reqInfo.Username, *password, *account)
	if reqInfo.Usertype == usernormal {
		*account, *password , _= eth.GetCreatedAccount(reqInfo.Username, appid, ctime)
		if "" == *account || "" == *password {
			return MSG_NOPREPARED_USER
		}
	} else {
		*account = config.Gconfig.Ethcfg.EnterpriseAccount
		*password = config.Gconfig.Ethcfg.EnterPrisePasswd
	}
	//log.Debug("openaccount account %s\n", *account)
	return MSG_SUCCESS
}

func starttransfer(reqInfo *TransferBody, serial_num *string) (string, string) {
	//根据用户名获取账户的key路径以及密码
	fmt.Println("SourceUser:",reqInfo.SourceUser)
	fmt.Println("DestUser:",reqInfo.DestUser)
	susr := eth.GetUser(reqInfo.SourceUser)
	tusr := eth.GetUser(reqInfo.DestUser)
	if susr == nil || tusr == nil {
		return "", MSG_INVALID_NAME
	}

	//根据用户名获取账户地址
	password, fromAccKeyPath, toAccAccount := susr.Password, susr.KeyPath, tusr.WalletAddr
	 //将转账金额由string型转换为int64
	transferAmount, _ :=strconv.ParseInt(reqInfo.TransferAmount, 10, 64)
	//调用转账函数，返回转入账户转账成功后余额
	//以太坊公链转账时，需要返回交易的hash值  add 2018-7-4  shangwj
	txhansh, err := eth.TransferAccounts(serial_num,&fromAccKeyPath,&password,&toAccAccount,&transferAmount)
	if err != nil {
		log.Error("Transfer : %s", err.Error())
		return "",MSG_UNKNOW
	}

	return txhansh,MSG_SUCCESS
}

func querybalance (reqInfo *AccountBody, balance *int64) string {
	var err error
	usr := eth.GetUser(reqInfo.Username)
	if nil == usr {
		return MSG_INVALID_NAME
	}

	switch reqInfo.TransferCurrency {
	case "0":
		if reqInfo.Usertype == userENT {
			*balance, err = eth.GetBalance(&config.Gconfig.Ethcfg.EnterpriseAccount)
		}
		if reqInfo.Usertype == usernormal {
			*balance, err = eth.GetBalance(&usr.WalletAddr)
		}
	    if err != nil {
	    	log.Error("query balance:%s", err.Error())
	    	return MSG_INVALID_ACCOUNT
        }
	case "1":
		//todo:添加获取以太币utxo的api
		fmt.Println("Bitcoin transfer operations is to be supported.")
		return MSG_UNSUPPORTED_CURRENCY
	}
	return MSG_SUCCESS
}

//提现操作
func withdrawtransfer(reqInfo *WithdrawBody, serial_num *string) (string, string) {
	inusr := eth.GetUser(reqInfo.SourceUser) //获取钱包账户
	exusr := eth.GetUser(reqInfo.ExternalAccount)

	fmt.Println(1)
	fmt.Println(exusr)
	if inusr == nil || exusr != nil {
		return "", MSG_INVALID_ACCOUNT
	}

	//to_addr:直接解析客户端发送的钱包地址，该账户地址没有在本地数据库上
	password, fromAccKeyPath, toAccAccount := inusr.Password, inusr.KeyPath, reqInfo.ExternalAccount //from_addr:数据库查找当前用户地址

	transferAmount, _ := strconv.ParseInt(reqInfo.TransferAmount, //value:  提现金额
		10, 64)

	wd_hash, err := eth.Withdraw(serial_num, &fromAccKeyPath, &password, &toAccAccount, //返回交易hash值
		&transferAmount)

	if err != nil {
		log.Error("Withdraw : %s", err.Error())
		return "", MSG_UNKNOW
	}
	return wd_hash, MSG_SUCCESS
}
