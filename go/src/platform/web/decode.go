/*****************************************************************************
File name: decode.go
Description: web模块消息解析
Author: failymao
Version: V1.0
Date: 2018/06/14
History:
*****************************************************************************/
package web

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"unicode/utf8"
	"platform/blockchain/eth"
	"strconv"
	"platform/db"
	"regexp"
)

const (
	method_create = "CreateAccount"
	method_transfer = "TransferAccounts"
	method_query = "QueryAccount"
	method_cancellation = "CancellationAccount"
	method_withdraw     = "WithdrawAccount"         //提现接口常量

	usernamemax = 32
	usernormal = "NORMAL"
	userENT = "ENTERPRISE"
)

type AccountBody struct {
	Username string `json:"username"`
	Usertype string `json:"usertype"`
	TransferCurrency string `json:"TransferCurrency"`      //转账币种：以太币("0") 比特币("1") （"2" "3" "4"是预留币种，尚未加入），开户时可不填入
}

type TransferBody struct {
	SourceUser string `json:"SourceUser"`
	DestUser string `json:"DestUser"`
	TransferAmount string `json:"TransferAmount"`           //单位为GWei
	TransferCurrency string `json:"TransferCurrency"`       //转账币种：以太币("0") 比特币("1") （"2" "3" "4"是预留币种，尚未加入）
}

type CancellationBody struct {
	Username string `json:"username"`
	Password string `json:"Password"`
}

// 定义提现请求参数
type WithdrawBody struct {
	SourceUser       string `json:"SourceUser"`
	ExternalAccount  string `json:"ExternalAccount"`
	TransferAmount   string `json:"TransferAmount"`
	TransferCurrency string `json:"TransferCurrency"`      //提现币种：以太币("0") 比特币("1") （"2" "3" "4"是预留币种，尚未加入）

}

func accountdecode (reqInfo *AccountBody, c *gin.Context)  string{
	err := c.BindJSON(reqInfo)
	if err != nil {
		return MSG_PARATYPEERR
	}
	return MSG_SUCCESS
}

func transferdecode (reqInfo *TransferBody, c *gin.Context)  string{
	err := c.BindJSON(reqInfo)
	if err != nil {
		return MSG_PARATYPEERR
	}
	return MSG_SUCCESS
}

//提现客户端参数解析函数
func withdrawdecode(reqInfo *WithdrawBody, c *gin.Context) string {
	err := c.BindJSON(reqInfo)
	if err != nil {
		return MSG_PARATYPEERR
	}
	return MSG_SUCCESS
}

func cancellationdecode(reqInfo *CancellationBody, c *gin.Context) string {
	err := c.BindJSON(reqInfo)
	if err != nil {
		return MSG_PARATYPEERR
	}
	return MSG_SUCCESS
}


func accountcheck (reqInfo *AccountBody, method string)  string{
	//Username校验
	if "" == reqInfo.Username || usernamemax < utf8.RuneCountInString(reqInfo.Username) {
		fmt.Printf("accountcheck failed in check Username:%s\n", reqInfo.Username)
		return MSG_INVALID_NAME
	}
	//Usertype校验
	if usernormal != reqInfo.Usertype && userENT != reqInfo.Usertype{
		fmt.Printf("accountcheck failed in Usertype:%s\n", reqInfo.Usertype)
		return MSG_INVALID_TYPE
	}
	//校验转账币种是否支持
	if "" == reqInfo.TransferCurrency || "" == db.Currencies[reqInfo.TransferCurrency] {
		fmt.Printf("transfercheck failed in check TransferCurrency:%s\n", reqInfo.TransferCurrency)
		return MSG_UNSUPPORTED_CURRENCY
	}
	//Username是否已经存在
	acc := eth.GetUser(reqInfo.Username)
	if nil != acc && method_create == method {
		fmt.Printf("The account has already existed:%s\n", reqInfo.Username)
		return MSG_EXIST_ACCOUNT
	}
	if nil == acc && method_query == method {
		fmt.Printf("The account doesn't existed:%s\n", reqInfo.Username)
		return MSG_INVALID_ACCOUNT
	}
	return MSG_SUCCESS
}

func transfercheck (reqInfo *TransferBody)  string{

	//校验转出账户是否已经开户
	susr := eth.GetUser(reqInfo.SourceUser)
	tusr := eth.GetUser(reqInfo.DestUser)
	if nil == susr {
		fmt.Printf("transfercheck failed in check SourceUser:%s\n", reqInfo.SourceUser)
		return MSG_NOTEXIST_FROMACCOUNT
	}
	if nil == tusr {
		fmt.Printf("transfercheck failed in check DestUser:%s\n", reqInfo.SourceUser)
		return MSG_NOTEXIST_TOACCOUNT
	}
	//校验转账币种是否支持
	if "" == reqInfo.TransferCurrency|| "" == db.Currencies[reqInfo.TransferCurrency]{
		fmt.Printf("transfercheck failed in check TransferCurrency:%s\n", reqInfo.TransferCurrency)
		return MSG_UNSUPPORTED_CURRENCY
	}
	//TransferAmount校验
	var  transferAmountStr string= reqInfo.TransferAmount //转账金额
	transferAmount, _ :=strconv.ParseInt(transferAmountStr, 10, 64)
	//校验提现的金额是不是为负数
	if transferAmount < 0 {
		fmt.Printf("transfercheck failed in TransferAmount:%s < 0\n", reqInfo.TransferAmount)
		return MSG_NEGATIVE_AMOUNT
	}
	//不上链方案只查数据库，不管链上
	if susr.Balance < transferAmount {
		fmt.Printf("TransferAmount :%d\n", transferAmount)
		return SUBMSG_INSUFFICIENT
	}
	//上链方案还需要预估手续费
	return MSG_SUCCESS
}

func cancellationcheck (reqInfo *CancellationBody)  string{
	//Username校验
	if "" == reqInfo.Username || usernamemax < utf8.RuneCountInString(reqInfo.Username) {
		fmt.Printf("cancellationcheck failed in check Username:%s\n", reqInfo.Username)
		return MSG_INVALID_NAME
	}
	if acc := eth.GetUser(reqInfo.Username) ; acc == nil {
		fmt.Printf("The account doesn't existed:%s\n", reqInfo.Username)
		return MSG_INVALID_ACCOUNT
	}
	return MSG_SUCCESS
}

/*************************************************
Function: withdrawcheck
Description: 提现接口客户端请求参数校验：from_addr,to_addr,transfer_amount
Author:failymao
Date: 2018/7/11
History:
*************************************************/
func withdrawcheck(reqInfo *WithdrawBody) string {

	//from_addr校验：账户是否在钱包数据库中，是否激活。
	inusr := eth.GetUser(reqInfo.SourceUser)
	if nil == inusr {
		fmt.Printf("withdrawcheck failed in check SourceUser:%s\n", reqInfo.SourceUser)
		return MSG_NOTEXIST_FROMACCOUNT
	}

	//to_addr校验1：外部账户需不在本地账户
	exusr := eth.GetUserName(reqInfo.ExternalAccount)
	if exusr != "" {
		//fmt.Printf("withdrawcheck failed in check ExternalAccount:%s\n", reqInfo.ExternalAccount)
		return MSG_INTERNAL_ACCOUNT
	}

	//to_addr校验2：地址是否为空
	if "" == reqInfo.ExternalAccount {
		fmt.Printf("withdrawcheck failed in check ExternalAccount:%s\n", reqInfo.ExternalAccount)
		return MSG_INVALID_ACCOUNT
	}

	// to_addr校验3：是否为ethereum地址
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if re.MatchString(reqInfo.ExternalAccount) != true {
		return MSG_INVALID_ACCOUNT
	}

	//提现币种校验
	if "" == reqInfo.TransferCurrency || "" == db.Currencies[reqInfo.TransferCurrency] {
		fmt.Printf("withdrawcheck failed in check TransferCurrency:%s\n", reqInfo.TransferCurrency)
		return MSG_UNSUPPORTED_CURRENCY
	}

	//value:类型转换string-->int
	var transferAmountStr string = reqInfo.TransferAmount //提现金额
	transferAmount, _ := strconv.ParseInt(transferAmountStr, 10, 64)
	//校验提现的金额是不是为负数
	if transferAmount < 0 {
		fmt.Printf("transfercheck failed in TransferAmount:%s < 0\n", reqInfo.TransferAmount)
		return MSG_NEGATIVE_AMOUNT
	}
	//value校验： 转账金额是否大于可用金额(数据库查询)
	if inusr.Balance < transferAmount {
		fmt.Printf("TransferAmount :%d\n", transferAmount)
		return SUBMSG_INSUFFICIENT
	}
	return MSG_SUCCESS
}
