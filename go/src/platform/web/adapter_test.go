package web

import (
	"testing"
	"fmt"
	"time"
	"platform/om"
	"platform/om/log"
	"platform/blockchain/eth"
)

func init() {
	log.New("debug", "")
	eth.RPCNew("http://192.168.0.192:10606")
	//初始化用户表本地缓存
	eth.EthInfoUpdate(20 * time.Second)
}

func TestOpenaccount(t *testing.T) {

	reqInfo := &AccountBody{
		Username: fmt.Sprintf("%10d", time.Now().Nanosecond()),
		Usertype: usernormal,
		TransferCurrency: "0",
	}
	account, passwd, appid := "", "", "what ever"
	ctime := time.Now().Format("2006-01-02 15:04:05")
	ret := openaccount(reqInfo, &account, &passwd, appid, ctime)
	om.Equal(t, MSG_SUCCESS, ret, "新建用户名为:", reqInfo.Username)
}

func TestStarttransfer(t *testing.T) {
	tranInfo := &TransferBody {
		SourceUser:"mintest",
		DestUser: "mine",
		TransferAmount: "1",
		TransferCurrency: "0",
	}
	seriID := ""
	_, ret := starttransfer(tranInfo, &seriID)
	om.Equal(t, MSG_SUCCESS, ret, "从", tranInfo.SourceUser, "转账到",
		tranInfo.DestUser, "金额为：", tranInfo.TransferAmount)
}

func TestQuerybalance(t *testing.T) {

	reqInfo := &AccountBody{
		Username: "mintest",
		Usertype: usernormal,
		TransferCurrency: "0",
	}
	bal := int64(0)
	ret := querybalance(reqInfo, &bal)
	om.Equal(t, MSG_SUCCESS, ret, "查询用户", reqInfo.Username, "余额为：",
				bal)
}

func TestWithdrawtransfer(t *testing.T) {
	withdInfo := &WithdrawBody {
		SourceUser: "teston193",
		ExternalAccount: "0x0409BeA70bABf2c240d659497D01c4d7f410ae92",
		TransferAmount: "1",
		TransferCurrency: "0",
	}
	seriID := ""
	_, ret := withdrawtransfer(withdInfo, &seriID)
	om.Equal(t, MSG_SUCCESS, ret, "用户提现", withdInfo.SourceUser, "到外部账户:",
				withdInfo.ExternalAccount, "金额:", withdInfo.TransferAmount)
}

