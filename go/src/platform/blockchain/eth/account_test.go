
package eth

import (
	"platform/om/log"
	"testing"
	"fmt"
	"time"
	"platform/db"
	"platform/om"
	"platform/tools"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"strings"
	"math/big"
	"github.com/ethereum/go-ethereum/common"
	"context"
)

/*****************************************************************************
Function: TestUpdateBanlanceoff 和 TransferAccountsV1
Description: 用于链上充值测试，想充值修改相应的参数，就可以完成充值
Author: 刘刚科
Version: V1.0
Date: 2018/07/10
History:
*****************************************************************************/
var dbMode *db.Model = nil
func init() {
	log.New("debug", "")
	RPCNew("http://192.168.0.192:10606")

	dbMode, _ = db.SetConfig()
}

/*
func TestUpdateBanlanceoff(t *testing.T) {
	//dbMode, _ := db.SetConfig()
	//getAccountNamePare(dbMode)
	//updateBanlanceoff(dbMode)
	RPCNew("http://192.168.0.192:10606")
	log.New("debug", "")

	key := "D:\\keystore\\UTC--2018-06-28T03-33-18.897438811Z--35a2c322791005b7942863f49d5277405ca4e00b"
	passwd := "123456"
	toacc := "0x3368fFE997fF5980b3fb806AA7b17ECa127b4665"
	amount := int64(4000)
	TransferAccountsV1(&key, &passwd, &toacc, &amount)
}
*/
func TransferAccountsV1(keypath, password, toAccAccount *string, transferAmount *int64) int64{

	conn := GetEthClient()
	token := GetEthToken()
	//实例化转入账户
	toAddress := common.HexToAddress(*toAccAccount)
	//获取转账前的转入账户代币数量
	val, _ := token.BalanceOf(nil, toAddress)

	//根据key文件的路径，获取key文件内容
	log.Debug("key dir:%s", *keypath)
	keycontent,err:=tools.ReadFile(*keypath)

	if err!=nil {
		log.Error("fail to read file: ", err)
	}
	// 创建授权账户(合约账户所绑定的外部账户，必须解锁后才能进行交易)
	auth, err := bind.NewTransactor(strings.NewReader(keycontent), *password)
	if err != nil {
		log.Error("Failed to create authorized transactor: %v", err)
	}
	//转账金额
	convertAmount :=   big.NewInt(*transferAmount)
	// 调用合约Transfer方法
	tx, err := token.Transfer(auth,toAddress,convertAmount)
	if err != nil {
		log.Error("Failed to request token transfer: %v", err)
	}

	//通过bind.WaitMined来等待事务真正被矿工处理完毕后，才会进行下一步操作
	ctx := context.Background()
	_,err = bind.WaitMined(ctx, conn, tx)

	if err != nil {
		log.Error("tx mining error:%v\n", err)
	}

	//获取转账后的代币数量
	val, _ = token.BalanceOf(nil, toAddress)
	return  val.Int64()
}

// account.go 测试
func TestCreateAccount(t *testing.T) {
	_, passwd, _ := CreateAccount()
	om.Equal(t, password, passwd, "创建以太坊新用户")
}

func TestGetCreatedAccount(t *testing.T) {
	ctime := fmt.Sprintf(time.Now().Format("2006-01-02 15:04:05"))
	name := fmt.Sprintf("%7d", time.Now().UnixNano())
	_, passwd, _ := GetCreatedAccount(name, "whatever", ctime)

	om.Equal(t, password, passwd, "分配已经新建好的以太坊账户成功")
}

func TestDelUserInfo(t *testing.T) {
	//删除不存在的账户
	DelUserInfo("what ever")
	//删除正常的账户
	DelUserInfo("testDel")
}



