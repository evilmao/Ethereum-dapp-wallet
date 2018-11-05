package web

import (
	"platform/blockchain/eth"
	"time"
	"platform/om/log"
	"testing"
	"platform/om"
)

func init() {
	log.New("debug", "")
	eth.RPCNew("http://192.168.0.192:10606")
	//初始化用户表本地缓存
	eth.EthInfoUpdate(20 * time.Second)
}

func TestCheck(t *testing.T) {
	//一般用户 - 用户不存在
	reqInfo1 := &AccountBody{
		Username: "mine10001",
		Usertype: usernormal,
		TransferCurrency: "0",
	}
	//企业用户 - 正常
	reqInfo2 := &AccountBody{
		Username: "mine1",
		Usertype: userENT,
		TransferCurrency: "0",
	}
	//一般用户 - 已经存在
	reqInfo3 := &AccountBody{
		Username: "mine",
		Usertype: userENT,
		TransferCurrency: "0",
	}
	//不支持的币种
	reqInfo4 := &AccountBody{
		Username: "what ever",
		Usertype: userENT,
		TransferCurrency: "888",
	}

	//正常
	tranInfo1 := &TransferBody{
		SourceUser: "mine",
		DestUser: "mintest",
		TransferAmount: "1",
		TransferCurrency: "0",
	}
	//转账源账户不存在
	tranInfo2 := &TransferBody{
		SourceUser: "what ever",
		DestUser: "mintest",
		TransferAmount: "1",
		TransferCurrency: "0",
	}
	//转账目标账户不存在
	tranInfo3 := &TransferBody{
		SourceUser: "mine",
		DestUser: "what ever",
		TransferAmount: "1",
		TransferCurrency: "0",
	}
	//货币类型不支持
	tranInfo4 := &TransferBody{
		SourceUser: "mine",
		DestUser: "mintest",
		TransferAmount: "1",
		TransferCurrency: "888",
	}
	//转账金额为负数
	tranInfo5 := &TransferBody{
		SourceUser: "mine",
		DestUser: "mintest",
		TransferAmount: "-1",
		TransferCurrency: "0",
	}
	//转账金额大于源账户总金额
	tranInfo6 := &TransferBody{
		SourceUser: "mine",
		DestUser: "mintest",
		TransferAmount: "9999999999999",
		TransferCurrency: "0",
	}
	//正常
	CancelInfo1 := &CancellationBody{
		Username: "mine",
		Password: "123456",
	}
	//不存在的账户
	CancelInfo2 := &CancellationBody{
		Username: "what ever",
		Password: "123456",
	}
	//正常
	withdrawInfo1 := &WithdrawBody{
		SourceUser: "mine",
		ExternalAccount: "0x0409BeA70bABf2c240d659497D01c4d7f410ae99",
		TransferAmount: "1",
		TransferCurrency: "0",
	}
	//提现源账户不存在
	withdrawInfo2 := &WithdrawBody{
		SourceUser: "what ever",
		ExternalAccount: "0x0409BeA70bABf2c240d659497D01c4d7f410ae92",
		TransferAmount: "1",
		TransferCurrency: "0",
	}
	//提现目标账户错误
	withdrawInfo3 := &WithdrawBody{
		SourceUser: "mine",
		ExternalAccount: "what ever",
		TransferAmount: "1",
		TransferCurrency: "0",
	}
	//提现金额为负数
	withdrawInfo4 := &WithdrawBody{
		SourceUser: "mine",
		ExternalAccount: "0x0409BeA70bABf2c240d659497D01c4d7f410ae93",
		TransferAmount: "-1",
		TransferCurrency: "0",
	}
	//提现金额大于提现账户的总金额
	withdrawInfo5 := &WithdrawBody{
		SourceUser: "mine",
		ExternalAccount: "0x0409BeA70bABf2c240d659497D01c4d7f410ae94",
		TransferAmount: "8000000000",
		TransferCurrency: "0",
	}
	//提现为负数
	//开户校验
	om.Equal(t, MSG_SUCCESS, accountcheck(reqInfo1, method_create),"创建一般账户不存在的用户")
	om.Equal(t, MSG_SUCCESS, accountcheck(reqInfo2, method_create), "创建企业账户")
	om.Equal(t, MSG_EXIST_ACCOUNT, accountcheck(reqInfo3, method_create), "创建已经存在的账户")
	om.Equal(t, MSG_UNSUPPORTED_CURRENCY, accountcheck(reqInfo4, method_create), "币种不支持")
	//查询校验
	om.Equal(t, MSG_INVALID_ACCOUNT, accountcheck(reqInfo1, method_query), "查询不存在的账户的余额")
	om.Equal(t, MSG_SUCCESS, accountcheck(reqInfo2, method_query), "查询企业账户余额")
	om.Equal(t, MSG_SUCCESS, accountcheck(reqInfo3, method_query), "查询一般已经存在的账户的余额")
	om.Equal(t, MSG_UNSUPPORTED_CURRENCY, accountcheck(reqInfo4, method_query), "查询的币种不支持")
	//转账校验
	om.Equal(t, MSG_SUCCESS, transfercheck(tranInfo1), "正常转账")
	om.Equal(t, MSG_NOTEXIST_FROMACCOUNT, transfercheck(tranInfo2), "源账户存在")
	om.Equal(t, MSG_NOTEXIST_TOACCOUNT, transfercheck(tranInfo3), "目标账户不存在")
	om.Equal(t, MSG_UNSUPPORTED_CURRENCY, transfercheck(tranInfo4), "不支持的货币类型")
	om.Equal(t, MSG_NEGATIVE_AMOUNT, transfercheck(tranInfo5), "转账金额为负数")
	om.Equal(t, SUBMSG_INSUFFICIENT, transfercheck(tranInfo6), "源账户余额不足")
	//销户校验
	om.Equal(t, MSG_SUCCESS, cancellationcheck(CancelInfo1), "正常销户")
	om.Equal(t, MSG_INVALID_ACCOUNT, cancellationcheck(CancelInfo2), "消不存在的账户")
	//提现校验
	om.Equal(t, MSG_SUCCESS, withdrawcheck(withdrawInfo1), "正常提现")
	om.Equal(t, MSG_NOTEXIST_FROMACCOUNT, withdrawcheck(withdrawInfo2), "源账户不存在")
	om.Equal(t, MSG_INVALID_ACCOUNT, withdrawcheck(withdrawInfo3), "目标账户不存在")
	om.Equal(t, MSG_NEGATIVE_AMOUNT, withdrawcheck(withdrawInfo4), "提现金额为负数")
	om.Equal(t, SUBMSG_INSUFFICIENT, withdrawcheck(withdrawInfo5), "原账户金额不足")
}




