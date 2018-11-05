package eth

import (
	"testing"
	"platform/om"
	"platform/db"
	"platform/om/log"
	"time"
)

func init() {
	log.New("debug", "")
	RPCNew("http://192.168.0.192:10606")
	//初始化用户表本地缓存
	dbMode, _ := db.SetConfig()
	getAccountNamePare(dbMode)
	ethBalRecover()
	time.Sleep( time.Second / 2)
}

func TestGetUser(t *testing.T) {
	usr1 := GetUser("")
	om.Equal(t, (*UserInfo)(nil), usr1, "查找为空的用户:", usr1)
	usr2 := GetUser("Haxima006")
	om.NotEqual(t, nil, usr2, "查找名字叫Haxima006的用户:", *usr2)
}


func TestGetUserName(t *testing.T) {
	usr1 := GetUserName("")
	om.Equal(t, "", usr1, "查找钱包地址为空的用户的用户名：", usr1)
	usr2 := GetUserName("0x35A2C322791005b7942863f49d5277405ca4e00b")
	om.Equal(t, "Haxima006", usr2, "查找钱包地址:",
		"0x35A2C322791005b7942863f49d5277405ca4e00b", "的用户名:", usr2)
}

func TestEth_comm(t *testing.T) {
	ethBalRecover()
	getAccountNamePare(dbMode)
	updateBanlanceoff(dbMode)
	time.Sleep( time.Second * 2)
	updateBanlance(dbMode)
}

func TestEthInfoUpdate(t *testing.T) {

	EthInfoUpdate(3 * time.Second)
	//保证EthInfoUpdate里面的写成至少能执行一次。
	time.Sleep(5 * time.Second)
}
