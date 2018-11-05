/*****************************************************************************
File name: bill.go
Description: 账单
Author:刘刚科
Version: V1.0
Date: 2018/06/29
History:
*****************************************************************************/
package bill

import (
	"os"
	"time"
	"fmt"
	"path"
	"log"
	"platform/config"
	"platform/blockchain/eth"
)

const (
	actionCreate = 1
	actionTransfer = 2
	actionDelete = 3

	actionCreateStr = 		"[Create Account] user[%50s] account[%50s]"
	actionTransferstr = 	"[Transfter     ] from[%50s:%50s]->[%50s:%50s] amount[%32s]"
	actionDeletestr = 		"[Delete Account] user[%50s] account[%50s]"
)

var actions map[int]string = map[int]string {
	actionCreate : actionCreateStr,
	actionTransfer : actionTransferstr,
	actionDelete : actionDeletestr,
}

var bill *Bill = nil

type Bill struct {
	base	 	*log.Logger
	baseFile   	*os.File
	basePath	string
	size	   	int64
}

func New(pathname string)  error {

	now := time.Now()
	filename := fmt.Sprintf("bill%d%02d%02d_%02d_%02d_%02d.txt",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second())

	file, err := os.Create(path.Join(pathname, filename))
	if err != nil {
		return err
	}

	bill = new(Bill)
	bill.base = log.New(file, "", log.Lshortfile)
	bill.baseFile = file
	bill.basePath = pathname
	return nil
}

func (b *Bill) doPrintf(act int, a ...interface{}) {
	if b.base == nil {
		panic("logger closed")
	}
	if b.size > int64(config.Gconfig.Logcfg.LogMax * 1024 * 1024) {
		b.baseFile.Close()
		b.rotate()
	}
	format := actions[act]
	outstr := fmt.Sprintf(format, a...)
	b.size += int64(len(outstr))
	b.base.Output(3, outstr)
}
// 大于规定大小之后，重新新建文件
func (b *Bill) rotate() error {

	now := time.Now()
	filename := fmt.Sprintf("bill%d%02d%02d_%02d_%02d_%02d.txt",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second())

	file, err := os.Create(path.Join(b.basePath, filename))
	if err != nil {
		fmt.Printf("log create new err:%s", err.Error())
		return err
	}

	b.base = log.New(file, "", log.LstdFlags)
	b.baseFile = file
	b.size = 0

	return nil
}

func CreateAcct(user string) {
	bill.doPrintf(actionCreate, user, eth.GetUser(user).WalletAddr)
}

func CreateEntAcct(user string) {
	bill.doPrintf(actionCreate, user, config.Gconfig.Ethcfg.EnterpriseAccount)
}

func TransferAcct(from, to, amount string) {
	bill.doPrintf(actionTransfer, from, eth.GetUser(from).WalletAddr, to, eth.GetUser(to).WalletAddr, amount)
}

func DeleteAcct(user string) {
	bill.doPrintf(actionDelete, user, eth.GetUser(user).WalletAddr)
}


