package eth

import (
	"testing"
	"platform/om"
)

func TestRPCNew(t *testing.T) {
	err := RPCNew("http:192.168.0.192:10606")
	om.Equal(t, nil, err, "连接到正常的geth客户端")
	err = RPCNew("http:what ever:port")
	om.NotEqual(t, nil, err, "链接到错误geth服务地址")
}

func TestRPCCall(t *testing.T) {
	//链接错误
	RPCNew("http:what ever:port")
	resultstr := ""
	RPCCall("eth_getBalance", &resultstr, "0x35A2C322791005b7942863f49d5277405ca4e00b", "latest")
	om.NotEqual(t, 0, len(resultstr), "查找钱包地址为:",
		"0x35A2C322791005b7942863f49d5277405ca4e00b", "用户以太币:", resultstr)
	//链接正确
	RPCNew("http:192.168.0.192:10606")
	err := RPCCall("wrong RPC name", &resultstr, "what ever!")
	//调用不存在的api
	om.NotEqual(t, (*error)(nil), err, "调用不存在的以太坊EVM API")
}

func TestGetEthToken(t *testing.T) {
	RPCNew("http:192.168.0.192:10606")
	//获取智能合约实例
	GetEthToken()
}


