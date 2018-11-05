/*****************************************************************************
File name: server.go
Description: web server,端到端测试，不涉及覆盖率
Author: failymao
Version: V1.0
Date: 2018/06/14
History:
*****************************************************************************/
package web

import (
	"net/http"
	"fmt"
	"strings"
	"io/ioutil"
	"testing"
)

func Test_Post_1(t *testing.T)  {

	Registerlist := []string{
		`{"Username":"test001","Usertype":"NORMAL"}`,     //注册一般用户
		`{"Username":"test001","Usertype":"ENTERPRISE"}`, //注册企业账户
	}
	Transferlist := []string{
		`{"SourceUser":"mine", "DestUser":"minetest", "TransferAmount":"1"}`, //转账
		`{"SourceUser":"minetest", "DestUser":"mine", "TransferAmount":"1"}`, //转账
	}

	Withdraw_list :=[]string{
		`{"SourceUser":"Haxima006","ExternalAccount":"0x323b5d4c32345ced77393b3530b1eed0f346429d","TransferAmount":"10","Password":"123456","TransferCurrency":"0"}`, //提现正常,以太坊
		`{"SourceUser":"Haxima006","ExternalAccount":"0x323b5d4c32345ced77393b3530b1eed0f346429d","TransferAmount":"10","Password":"123456","TransferCurrency":"1"}`, //提现正常,其他币种
		`{"SourceUser":"Haxima006","ExternalAccount":"","TransferAmount":"10","Password":"123456","TransferCurrency":"0"}`,                                           //提现外部地址为空
		`{"SourceUser":"Haxima006","ExternalAccount":"0x323b5d4c32345ced77393b3530b1eed0f34642","TransferAmount":"10","Password":"123456","TransferCurrency":"0"}`,   //提现外部地址为不正确
		`{"SourceUser":"Haxima006","ExternalAccount":"0xC709c71884F1bd4F90C85b1F2229e5094EE8bEF1","TransferAmount":"10","Password":"123456","TransferCurrency":"0"}`, //提现：外部地址为数据库地址
		`{"SourceUser":"5fe0a050caf7af9","ExternalAccount":"0x323b5d4c32345ced77393b3530b1eed0f34642","TransferAmount":"10","Password":"123456","TransferCurrency":"0"}`, //钱包地址未注册
		`{"SourceUser":"mine","ExternalAccount":"0x323b5d4c32345ced77393b3530b1eed0f34642","TransferAmount":"10","Password":"123456","TransferCurrency":"0"}`,           //钱包地址余额不足
		`{"SourceUser":"mine","ExternalAccount":"0x323b5d4c32345ced77393b3530b1eed0f34642","TransferAmount":"10","Password":"123456"}`,                                 //提现参数缺失
		`"SourceUser":"Haxima006","ExternalAccount":"0x323b5d4c32345ced77393b3530b1eed0f346429d","TransferAmount":"10","Password":"123456","TransferCurrency":"0"`,     //提现发送参数格式错误--string
	}

	Cancellist := []string {
		`{"username":"test001", "Password":"123456"}`, //销户
	}

	req := make([]*http.Request, 0)
	for _, data := range Registerlist {
		request, _ := http.NewRequest("POST", "http://localhost:12306/wallet?method=CreateAccount",
			strings.NewReader(data))
		req = append(req, request)
	}
	for _, data := range Transferlist {
		request, _ := http.NewRequest("POST", "http://localhost:12306/wallet?method=TransferAccounts",
			strings.NewReader(data))
		req = append(req, request)
	}
	for _, data := range Cancellist {
		request, _ := http.NewRequest("POST", "http://localhost:12306/wallet?method=CancellationAccount",
			strings.NewReader(data))
		req = append(req, request)
	}


	for _,data :=range Withdraw_list{
		request, _ := http.NewRequest("POST", "http://localhost:12306/wallet?method=WithdrawAccount",
			strings.NewReader(data))
		req = append(req, request)

	}
	//post数据并接收http响应
	for _, reqi := range req {
		resp,err :=http.DefaultClient.Do(reqi)
		if err!=nil{
			t.Error("Test_Post_1 return error\n", err)
		}else {
			fmt.Println("Test_Post_1 post a data successful.")
			respBody,_ :=ioutil.ReadAll(resp.Body)
			t.Log("Test_Post_1 return:",string(respBody))
		}
	}
}
