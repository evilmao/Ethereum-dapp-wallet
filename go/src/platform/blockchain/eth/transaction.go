/*****************************************************************************
File name: transaction.go
Description: 转账
Author: 尚文静
Version: V1.0
Date: 2018/06/19
History:
*****************************************************************************/
package eth

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"strings"
	"math/big"
	"github.com/ethereum/go-ethereum/common"
	"platform/tools"
	"platform/om/log"
	"strconv"
	"platform/config"
	"context"
	"github.com/ethereum/go-ethereum"
	"fmt"
	"time"
	"platform/db"
)

type Traninfo struct {
	Serialno string
	keypath string
	password string
	transferCurrency string
	toAccAccount string
	fromAccAccount string
	Amount	int64
	TxHash  string
	TxStatus int
	Ntime  int64  //交易再次打包执行次数
}
var tranokChann chan *Traninfo = make(chan *Traninfo, 100)
var tranfailChann chan *Traninfo = make(chan *Traninfo, 100)
/*************************************************
Function: TransferAccounts
Description:转账 根据转出账户key文件路径 以及 密码 ，转入账户，转账金额 进行转账操作
Author:
Date: 2018/06/19
History: 	以太坊公链转账时，需要返回交易的hash值  add 2018-7-4  shangwj
*************************************************/
func TransferAccounts(serial_num,keypath, password, toAccAccount *string, transferAmount *int64)  (string, error) {
	traninfo := &Traninfo{
		Serialno : *serial_num,
		toAccAccount: *toAccAccount,
		fromAccAccount:"",
		keypath:  *keypath,
		password: *password,
		Amount:*transferAmount,
		TxHash:"",
		TxStatus:0,
	}
	if len(*keypath)>42{
		traninfo.fromAccAccount="0x"+(*keypath)[len(*keypath)-40:len(*keypath)] //截取key文件名称倒数41位的账户地址
	}
	//根据key文件的路径，获取key文件内容
	log.Debug("key dir:%s", *keypath)
	log.Debug("traninfo.fromAccAccount:", traninfo.fromAccAccount)
	keycontent, err := tools.ReadFile(*keypath)
	if err!=nil {
		log.Error("fail to read file: ", err)
		tranfailChann <- traninfo
		return  "", err
	}
	log.Debug("keyStore:%s",keycontent)

	// 创建授权账户(合约账户所绑定的外部账户，必须解锁后才能进行交易)
	auth, err := bind.NewTransactor(strings.NewReader(keycontent), *password)
	if err != nil {
		log.Error("Failed to create authorized transactor: %v", err)
		tranfailChann <- traninfo
		return "", err
	}

   //估算下以太币够不够手续费
	//获取转出账户的以太币
	fromethc,err :=GetEthBalance(&traninfo.fromAccAccount)
	log.Debug("转出账户以太坊币：%d",fromethc)
	//本次交易需要的费用最大为：
	maxGas,_:=EstimateTranferFee(traninfo.fromAccAccount,*toAccAccount,  *transferAmount)
	log.Debug("以太坊转账手续费=:%d",maxGas)

	//if(maxGas>fromethc){
	//	fmt.Println("Failed to gas is out")
	//	return  "",nil
	//}
	//调用合约Transfer方法
	go AsyncPost(auth, traninfo)
	return traninfo.TxHash, nil
}



/*************************************************
Function: GetBalance
Description:根据账户地址获取代币数量
Author:
Date: 2018/06/19
History:
*************************************************/
func  GetBalance(AccAccount *string,) (int64, error){
	//start := time.Now()

	//实例化一个智能合约
	token := GetEthToken()
	//该账户代币数量
	val, err :=token.BalanceOf(nil,  common.HexToAddress(*AccAccount))
	if err !=nil {
		log.Error("Failed to get the balance :", err)
		return 0, err
	}
	//defer log.Debug(">>>>>GetBalance:%v", time.Since(start))
	return val.Div(val, big.NewInt(GWei)).Int64(), nil
}
/*************************************************
Function: GetEthBalance
Description:根据账户地址获取以太币数量
Author:
Date: 2018/07/02
History:
*************************************************/
func GetEthBalance(AccAccount *string) (int64, error) {

	resultstr := ""

	err := RPCCall("eth_getBalance", &resultstr, AccAccount, "latest")
	if err != nil {
		return 0, err
	}
	result:=tools.HexDec(resultstr)//将16进制转换为10进制
	return result, nil
}
/*************************************************
Function: EstimateTranferFee
Description:计算某笔交易的手续费
Author:
Date: 2018/07/02
History:
*************************************************/
func EstimateTranferFee(fromAccount, toAccount string, amount int64) (uint64, error) {

	ctx := context.Background()
	toaddr := common.HexToAddress(toAccount)

	call := ethereum.CallMsg{
		From: common.HexToAddress(fromAccount),
		To: &toaddr,
	//	Gas: ethBal/uint64(config.Gconfig.Ethcfg.GasPrice), //todo:不确定
	//	GasPrice: big.NewInt(config.Gconfig.Ethcfg.GasPrice),
		Value: big.NewInt(amount),
	//	Data: []byte(""), //todo : 这里填什么
	}
	gas, err := GetEthClient().EstimateGas(ctx, call)
	if err != nil {
		log.Debug("EstimateGas:%s", err.Error())
		return 0, err
	}
	gasprice,_:=	GetEthClient().SuggestGasPrice(context.Background())
//	fmt.Println("gasprice:", gasprice)
	egas :=gas * gasprice.Uint64()
	return egas, nil

}
/*************************************************
Function: AsyncPost
Description:异步投递交易任务，每个异步交易一个协程
Author:
Date: 2018/07/02
History:
*************************************************/
func AsyncPost(auth *bind.TransactOpts, traninfo *Traninfo) {
	conn:=GetEthClient()
	token := GetEthToken()
	//实例化转入账户
	toAddress := common.HexToAddress(traninfo.toAccAccount)
	//转账金额
	convertAmount :=   big.NewInt(traninfo.Amount*GWei) //转换Gwei->wei
	//获取当前交易的区块高度
	oldBlock:=  GetCurrentBlock()
	//获取当前时间
	oldtime:=time.Now().Format("2006-01-02 15:04:05")
	//找交易所在连
	tranBlock := int64(0)
	j := int64(0)
	//当前时间
	var newtime string
	tx, err := token.Transfer(auth, toAddress, convertAmount)
	if err != nil {
		log.Error("Failed to request token transfer: %v", err)
		//再投递一次本次交易
		traninfo.TxStatus=2  //交易打包失败
		tranfailChann <- traninfo
		return
	}
	//获取以太坊转账交易的hash值
	txHash := tx.Hash().Hex()
	traninfo.TxHash=txHash
	for {
		time.Sleep(10*time.Second)
		newBlock := GetCurrentBlock()
		for i := oldBlock ; i < newBlock ; i++ {
			//根据区块高度获取该区块的区块信息
			curBlock, _ := conn.BlockByNumber(context.Background(),big.NewInt(int64(i)))
			for _, tran :=  range curBlock.Transactions() {//遍历区块上的交易
				newtime =time.Now().Format("2006-01-02 15:04:05")
				minute := tools.GetMinDiffer(oldtime, newtime) //计算分钟时间差
				if txHash == tran.Hash().Hex() { //对比当前交易的hash值是否在查询的区块高度上，如果在返回当前查询的区块高度
					tranBlock = i
					break
				}
				if minute > config.Gconfig.Ethcfg.TimeOut { //检查是否超时
					//oldBlock=  GetCurrentBlock()
					traninfo.TxStatus=3  //查询交易打包块超时
					tranfailChann <- traninfo
					if j >= config.Gconfig.Ethcfg.TimeoutTimes {
						log.Error("transaction fails times")
						tranfailChann <- traninfo
					}
					j++
					break
				}
			}
		}
		oldBlock = newBlock
		//获取到当前交易
		if tranBlock>0 || j>=config.Gconfig.Ethcfg.TimeoutTimes{
			break
		}

	}
	//	fmt.Printf("tranBlock :%d\n", tranBlock)
	//检验高度是否超过6
	if(tranBlock>0){
		for {
			time.Sleep(2*time.Second)
			newBlock:=  GetCurrentBlock()
			newtime =time.Now().Format("2006-01-02 15:04:05")
			minute2 := tools.GetMinDiffer(oldtime, newtime)

			if newBlock - tranBlock >= config.Gconfig.Ethcfg.ConfirmeTimes {
				traninfo.TxStatus=1 //交易成功
				tranokChann <- traninfo
				break
			}else{
				if (minute2 > config.Gconfig.Ethcfg.TimeOut) {
					traninfo.TxStatus=4 //已打包，超时未确认
					tranfailChann <- traninfo
					break
				}
			}
		}
	}
}
/*************************************************
Function: GetCurrentBlock
Description:获取当前区块链的区块高度
Author:
Date: 2018/07/02
History:
*************************************************/
func GetCurrentBlock() int64{
	var blocknumstr string
	//调用RPCCall 获取当前最新区块高度
	err := RPCCall("eth_blockNumber", &blocknumstr)
	if err != nil {
		fmt.Println("new account  err", err)
		return 0
	}
	blocknumstr = blocknumstr[2:]
	currentblocknumber,err :=strconv.ParseInt(blocknumstr, 16, 32)
	if err != nil {
		fmt.Println(err)
	}
	return currentblocknumber
}

//异步操作channel
func AsyncTranfer() {
	go func() {
		for traninfo := range tranfailChann {
			//失败再投递处理：失败方式超时||超次数
			traninfo.Ntime++
			if traninfo.Ntime > config.Gconfig.Ethcfg.TimeoutTimes {
				//彻底失败，写数据库彻底失败，不再重试
				if traninfo.Serialno != "" {
					//更新数据库中交易状态
					db.UpdateBalanceTXstatus(traninfo.Serialno, traninfo.transferCurrency, GetUserName(traninfo.toAccAccount), GetUserName(traninfo.toAccAccount), traninfo.TxStatus, traninfo.Amount)
				}
			} else {
				//失败再投递
				TransferAccounts(&traninfo.Serialno, &traninfo.keypath, &traninfo.password, &traninfo.toAccAccount, &traninfo.Amount)
			}
		}
	}()
	go func() {
		for traninfo := range tranokChann {
		    log.Debug("traninfo.Serialno:%s", traninfo.Serialno)
			//成功写数据库
			if traninfo.Serialno != "" {
				//更新数据库中交易状态
			    db.UpdateBalanceTXstatus(traninfo.Serialno, traninfo.transferCurrency, GetUserName(traninfo.toAccAccount), GetUserName(traninfo.toAccAccount), traninfo.TxStatus, traninfo.Amount)
		    }
		}
	}()
}

