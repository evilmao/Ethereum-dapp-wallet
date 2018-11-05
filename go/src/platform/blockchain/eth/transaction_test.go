/*****************************************************************************
File name: transaction.go
Description: 转账测试
Author: 尚文静
Version: V1.0
Date: 2018/07/23
History:
*****************************************************************************/
package eth

import (
	"math/big"
	"github.com/ethereum/go-ethereum/common"
	"platform/tools"
	"platform/om/log"
	"strconv"
	"context"
	"github.com/ethereum/go-ethereum"
	"fmt"

	"testing"
	"platform/db"
)

func init() {
	log.New("debug", "")
	RPCNew("http://192.168.0.192:10606")
	//初始化用户表本地缓存
	dbMode, _ := db.SetConfig()
	getAccountNamePare(dbMode)
	ethBalRecover()
}
/*************************************************
Function: TransferAccountsTest
Description:转账测试
Author:
Date: 2018/07/14
History:
*************************************************/

func TestTransferAccounts(t *testing.T){
	toAccAccount:="0x9ecfb01035e958a41ff0e542b309937c0917d503"  //转入账户地址
	keypath:="/keystore/UTC--2018-06-28T03-33-18.897438811Z--35a2c322791005b7942863f49d5277405ca4e00b" //转出账户keystore路径
	var  transferAmount  int64=1000000000//转账金额
	password:="123456"//转出账户密码
	serialno :="1234560001"//交易流水号
	fmt.Println("rpc.Dial transferAmount:", transferAmount)
	TransferAccounts(&serialno,&keypath,&password,&toAccAccount, &transferAmount)
}



/*************************************************
Function: GetBalance
Description:根据账户地址获取代币数量
Author:
Date: 2018/06/19
History:
*************************************************/
func  TestGetBalance(t *testing.T) {
	//start := time.Now()
	toAccAccount:="0x9ecfb01035e958a41ff0e542b309937c0917d503"  //账户地址
	//实例化一个智能合约
	token := GetEthToken()
	//该账户代币数量
	val, err :=token.BalanceOf(nil,  common.HexToAddress(toAccAccount))
	if err !=nil {
		log.Error("Failed to get the balance :", err)
		return
	}
	fmt.Printf("balance is %d\n",val.Div(val, big.NewInt(GWei)).Int64())
}
/*************************************************
Function: GetEthBalance
Description:根据账户地址获取以太币数量
Author:
Date: 2018/07/02
History:
*************************************************/
func TestGetEthBalance(t *testing.T) {

	resultstr := ""
	toAccAccount:="0x9ecfb01035e958a41ff0e542b309937c0917d503"  //账户地址
	err := RPCCall("eth_getBalance", &resultstr, toAccAccount, "latest")
	if err != nil {
		log.Error("Failed to get the eth_balance :", err)
	}
	result:=tools.HexDec(resultstr)//将16进制转换为10进制
	fmt.Printf("eth_balance is %d\n",result)
}
/*************************************************
Function: EstimateTranferFee
Description:计算某笔交易的手续费
Author:
Date: 2018/07/02
History:
*************************************************/
func TestEstimateTranferFee(t *testing.T) {
	toAccAccount:="0x9ecfb01035e958a41ff0e542b309937c0917d503"  //转入账户地址
	fromAccount:="0x35a2c322791005b7942863f49d5277405ca4e00b" //转出账户地址
	var  transferAmount  int64=100000//转账金额
	ctx := context.Background()
	toaddr := common.HexToAddress(toAccAccount)

	call := ethereum.CallMsg{
		From: common.HexToAddress(fromAccount),
		To: &toaddr,
	//	Gas: ethBal/uint64(config.Gconfig.Ethcfg.GasPrice), //todo:不确定
	//	GasPrice: big.NewInt(config.Gconfig.Ethcfg.GasPrice),
		Value: big.NewInt(transferAmount),
	//	Data: []byte(""), //todo : 这里填什么
	}
	gas, err := GetEthClient().EstimateGas(ctx, call)
	if err != nil {
		log.Debug("EstimateGas:%s", err.Error())
	}
	fmt.Printf("gas is %d\n",gas)
	gasprice, err := GetEthClient().SuggestGasPrice(context.Background())
	if err != nil {
		return
	}
	fmt.Printf("gasprice is %d\n",gasprice)
//	fmt.Println("gasprice:", gasprice)
	egas :=gas * gasprice.Uint64()
	fmt.Printf("egas is %d\n",egas)

}

/*************************************************
Function: GetCurrentBlock
Description:获取当前区块链的区块高度
Author:
Date: 2018/07/02
History:
*************************************************/
func TestGetCurrentBlock(t *testing.T) {
	var blocknumstr string
	//调用RPCCall 获取当前最新区块高度
	err := RPCCall("eth_blockNumber", &blocknumstr)
	if err != nil {
		fmt.Println("new account  err", err)
		return
	}
	blocknumstr = blocknumstr[2:]
	currentblocknumber,err :=strconv.ParseInt(blocknumstr, 16, 32)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("当前区块链的区块高度： %d\n",currentblocknumber)

}
