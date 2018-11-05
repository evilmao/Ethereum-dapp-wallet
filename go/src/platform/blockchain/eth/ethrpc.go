/*****************************************************************************
File name: ethrpc.go
Description: rpc连接geth的session
Author:刘刚科
Version: V1.0
Date: 2018/06/25
History:
*****************************************************************************/
package eth

import (
	"errors"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"platform/blockchain/mytoken"
	"platform/config"
	"platform/om/log"
)

const GWei = int64(1000000000)
/*
type IReConn interface {
	Reconn() error
}

var ReconnChan chan IReConn= make(chan IReConn)



func (c *RPCAgent) Reconn() (error) {

	c.Cli.Close()

	newcli, err := rpc.Dial(c.ConnStr)
	if err != nil {
		log.Error("Dial[%s] fail:%s", c.ConnStr, err)
		return err
	}

	c.Cli = newcli
	return nil
}
*/

type RPCAgent struct {
	ConnStr string
	Cli     *rpc.Client
}

var rpcAgent *RPCAgent
var keyStore *keystore.KeyStore

func init() {
	keyStore = keystore.NewKeyStore(config.Gconfig.Ethcfg.Keystoredir, keystore.StandardScryptN, keystore.StandardScryptP)
}

/*************************************************
Function: GetEthClient
Description: 转换geth的rpc链接->geth的api的ethclient
Author:
Date: 2018/07/06
History:
*************************************************/
func GetEthClient() *ethclient.Client {
	return ethclient.NewClient(rpcAgent.Cli)
}

/*************************************************
Function: GetEthToken
Description: 获取智能合约实例
Author:
Date: 2018/07/06
History:
*************************************************/
func GetEthToken() *mytoken.Token {
	conn := GetEthClient()
	token, err := mytoken.NewToken(common.HexToAddress(config.Gconfig.Ethcfg.ContractAccount), conn)

	if err != nil {
		log.Error("Failed to instantiate a Token contract: %v", err)
	}
	return token
}
/*************************************************
Function: RPCNew
Description: 链接到geth,并保存session
Author:
Date: 2018/07/06
History:
*************************************************/
func RPCNew(addr string) error {

	client, err := rpc.Dial(addr)
	if err != nil {
		log.Error("Dial[%s] fail:%s", addr, err)
		return err
	}

	rpcAgent = &RPCAgent{
		ConnStr: addr,
		Cli:     client,
	}
	return nil
}
/*************************************************
Function: RPCCall
Description: 链接到geth,并保存session
Author:
Date: 2018/07/06
History:
*************************************************/
func RPCCall(method string, reply interface{}, rpcArgs ...interface{}) error {

	if rpcAgent == nil {
		errstr := "rpc Agent must dial first"
		log.Error(errstr)
		return errors.New(errstr)
	}

	err := rpcAgent.Cli.Call(reply, method, rpcArgs...)
	if err != nil {
		log.Error("RPC Call Err:%s", err.Error())
		return err
	}
	return nil
}
/*************************************************
Function: GetKeyStore
Description: 返回已经创建好的Keystore(一次初始化)
Author:
Date: 2018/07/06
History:
*************************************************/
func GetKeyStore() *keystore.KeyStore {
	return keyStore
}
