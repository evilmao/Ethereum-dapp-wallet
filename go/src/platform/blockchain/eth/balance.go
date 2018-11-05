/*****************************************************************************
File name: balance.go
Description: 账户资金查询
Author:刘刚科
Version: V1.0
Date: 2018/06/20
History:
*****************************************************************************/
package eth

import (
	"time"
	"sync"
	"platform/om/log"
	"platform/db"
	"reflect"
	"strconv"
	"fmt"
	"platform/config"
	"sync/atomic"
	"math/big"
)

type UserInfo struct {
	Username	string
	Password	string
	WalletAddr	string
	Balance		int64
	ChBalance   int64
	KeyPath		string
	Appid		string
	Timestamp 	string
}

var balanceList map[string]int64= make(map[string]int64, 0)    //代币数目
var balanceEthList map[string]int64 = make(map[string]int64, 0) //以太币数目
var usrmap map[string]*UserInfo = make(map[string]*UserInfo, 0)				//用户表本地缓存 --已分配
var usrmapunused map[string]*UserInfo = make(map[string]*UserInfo, 0)		//用户表本地缓存 --未分配
var usrmapaddr  map[string]*UserInfo = make(map[string]*UserInfo, 0)		//用户表本地缓存 --key为walletaddr
var usrmapRWLoker *sync.RWMutex = &sync.RWMutex{}
/*************************************************
Function: updateBanlanceoff
Description: 不上链方案的第三方充值更新余额
Author:
Date: 2018/07/06
History:
*************************************************/
func updateBanlanceoff(mod *db.Model) error {
	ks := GetKeyStore()
	accounts := ks.Accounts()
	token := GetEthToken()

	//从链上查询余额
	baldiff := make(map[string]int64)
	for _, account := range accounts {

		balance, err := token.BalanceOf(nil, account.Address)
		if err != nil {
			//log.Error("Query Accout[%v] err:%s",account, err.Error())
			continue
		}
		gbal := balance.Div(balance, big.NewInt(GWei)).Int64()
		//log.Debug("%s -> ", account.Address.Hex(), balance.Uint64())
		//计算充值金额
		if gbal > balanceList[account.Address.Hex()] {
			baldiff[account.Address.Hex()] = gbal - balanceList[account.Address.Hex()]
			balanceList[account.Address.Hex()] = gbal
		}
	}
	//余额有变化的写入数据库
	//log.Debug("链上金额：%+v", balanceList)
	//for _, v := range usermap {
	//	log.Debug("数据库：%+v", *v)
	//}

	sqlstr := ""
	for k, v := range baldiff {
		if v == 0 {
			continue
		}
		if _, ok := usrmapaddr[k] ; !ok {
			continue
		}

		sqlstr += fmt.Sprintf("%s,%d,", k, v)
		tid := ""
		amount := int64(v)
		if _, ok := usrmapaddr[k] ; ok {
			TransferAccounts(&tid, &usrmapaddr[k].KeyPath, &usrmapaddr[k].Password,
				&config.Gconfig.Ethcfg.EnterpriseAccount, &amount)
		}

	}
	if "" == sqlstr {
		return nil
	}
	sqlstr = sqlstr[:len(sqlstr) - 1]
	sqlstr = fmt.Sprintf(`CALL UpdateBalanceoff('%s')`, sqlstr)
	log.Debug("UpdateBalance str:%s", sqlstr)
	result := mod.Query(sqlstr)
	switch result.(type) {
	case map[int]map[string]string:
		tabRes := result.(map[int]map[string]string)
		if len(tabRes) != 1 && tabRes[0]["result"] != "0" {
			return fmt.Errorf("CALL UpdateBalance: result[%+v]",tabRes)
		} else {
			//写入成功，还需要更新内存中的余额
			for acc, bal := range baldiff {
				for _, usr := range usrmap {
					if usr.WalletAddr == acc {
						atomic.AddInt64(&usr.Balance, bal)
						atomic.AddInt64(&usr.ChBalance, bal)
						log.Debug(">>>>>>>>>>>user[%s] balance diff[%d]<<<<<<<<<<", usr.Username, bal)
					}
				}
			}
		}
	default:
		return fmt.Errorf("CALL UpdateBalance: result[%+v]",result)
	}
	return nil
}

/*************************************************
Function: updateBanlance
Description: 上链方案的第三方充值更新余额
Author:
Date: 2018/07/06
History:
*************************************************/
func updateBanlance(mod *db.Model) error {
	ks := GetKeyStore()
	accounts := ks.Accounts()
	token := GetEthToken()

	//从数据库更新链上余额
	for _, acc := range accounts {
		for _, usr := range usrmap {
			//log.Debug("%s----------%s", strings.ToLower(usr.WalletAddr), strings.ToLower(acc.Address.Hex()))
			if usr.WalletAddr == acc.Address.Hex() {
				balanceList[usr.WalletAddr] = usr.ChBalance
			}
		}
	}
	//从链上查询余额
	for _, account := range accounts {
		//log.Debug("%+v", account)
		balance, err := token.BalanceOf(nil, account.Address)
		if err != nil {
			log.Error("Query Accout[%v] err:%s",account, err.Error())
			continue
		}

		if _, ok := balanceList[account.Address.Hex()] ; !ok {
			log.Error("Database has account %s, but don't have keystore", balanceList[account.Address.Hex()])
			continue
		}

		//计算充值金额
		balanceList[account.Address.Hex()] = balance.Div(balance, big.NewInt(GWei)).Int64()
	}
	//余额有变化的写入数据库
	//log.Debug("链上金额：%+v", balanceList)
	//for _, v := range usermap {
	//	log.Debug("数据库：%+v", *v)
	//}
	sqlstr := ""
	var changed map[string]int64 = make(map[string]int64, 0)
	for k, v := range balanceList {
		for _, usr := range usrmap {
			if k == usr.WalletAddr && v != usr.Balance {
				sqlstr += fmt.Sprintf("%s,%d,",usr.Username, v)
				changed[usr.Username] = v
			}
		}
	}
	if "" == sqlstr {
		return nil
	}
	sqlstr = sqlstr[:len(sqlstr) - 1]

	sqlstr = fmt.Sprintf(`CALL UpdateBalance('%s')`, sqlstr)

	log.Debug("UpdateBalance str:%s", sqlstr)
	result := mod.Query(sqlstr)
	switch result.(type) {
	case map[int]map[string]string:
		tabRes := result.(map[int]map[string]string)
		if len(tabRes) != 1 && tabRes[0]["result"] != "0" {
			return fmt.Errorf("CALL UpdateBalance: result[%+v]",tabRes)
		} else {
			//写入成功，还需要更新内存中的余额
			for name, bal := range changed {
				atomic.StoreInt64(&usrmap[name].Balance, bal)
				atomic.StoreInt64(&usrmap[name].ChBalance, bal)
				log.Debug(">>>>>>>>>>>user[%s] balance[%d]<<<<<<<<<<", name,bal)
			}
		}
	default:
		return fmt.Errorf("CALL UpdateBalance: result[%+v]",result)
	}
	return nil
}
/*************************************************
Function: GetUser
Description: 由用户名获取用户所有信息
Author:
Date: 2018/07/06
History:
*************************************************/
func GetUser(name string) *UserInfo {
	usrmapRWLoker.RLock()
	defer usrmapRWLoker.RUnlock()
	if usr, ok := usrmap[name] ; ok {
		return usr
	} else {
		return nil
	}
}

/*************************************************
Function: GetUserName
Description: 由用户名钱包地址获取用户名
Author:
Date: 2018/07/06
History:
*************************************************/
func GetUserName(account string) string{
	if usr, ok := usrmapaddr[account] ; ok {
		return usr.Username
	} else {
		return ""
	}
}
/*************************************************
Function: getAccountNamePare
Description: 从数据库获缓存用户信息到内存
Author:
Date: 2018/07/06
History:
*************************************************/
func getAccountNamePare(mod *db.Model) {
	result := mod.Query("SELECT Username, Password, Balance, ChBalance, Keypath, WalletAddress, AppID, Timestamp, Status FROM UserInfo")

	switch result.(type) {
	case map[int]map[string]string :
		tabRes := result.(map[int]map[string]string)
		for _, row := range tabRes {
			usr := &UserInfo{}
			usr.Username = row["Username"]
			usr.Password = row["Password"]
			balance, err := strconv.ParseInt(row["Balance"], 10, 64)
			if err != nil {
				continue
			}
			chbalance, err := strconv.ParseInt(row["ChBalance"], 10, 64)
			if err != nil {
				continue
			}
			usr.Balance = balance
			usr.ChBalance = chbalance
			usr.KeyPath = row["Keypath"]
			usr.WalletAddr = row["WalletAddress"]
			usr.Timestamp = row["Timestamp"]
			usr.Appid = row["AppID"]
			if "0" == row["Status"] {
				usrmap[usr.Username] = usr
			}
			if "1" == row["Status"] {
				usrmapunused[usr.Username] = usr
			}
			usrmapaddr[usr.WalletAddr] = usr
		}
	default:
		log.Error("unexpected result type:%+v", reflect.TypeOf(result))
	}
	log.Debug("registed user [%d], prepared user [%d]", len(usrmap), len(usrmapunused))
	//log.Debug("%+v", usernameAccountPares)
}
/*************************************************
Function: ethBalRecover
Description: 启动时获取到用户此时链上的余额
Author:
Date: 2018/07/06
History:
*************************************************/
func ethBalRecover() {

	ks := GetKeyStore()
	accounts := ks.Accounts()

	for _, acc := range accounts {
		//log.Debug("%s\n", acc.Address.Hex())
		balanceList[acc.Address.Hex()] = 0
	}
	for _, usr := range usrmap {
		//log.Debug("%s----------%s", strings.ToLower(usr.WalletAddr), strings.ToLower(acc.Address.Hex()))
		if _, ok := balanceList[usr.WalletAddr] ; !ok {
			continue
		} else {
			balanceList[usr.WalletAddr] = usr.ChBalance
		}
	}
}
/*************************************************
Function: EthInfoUpdate
Description: 不上链时：定时更新用户链上的充值行为，
			 上链时：定时更新用户转账，充值，提现等行为的账户余额变动
Author:
Date: 2018/07/06
History:
*************************************************/
func EthInfoUpdate(dura time.Duration) {
	//读取username和account的对应关系
	dbMode, _ := db.SetConfig()
	getAccountNamePare(dbMode)
	ethBalRecover()
	ticker := time.NewTicker(dura)
	go func() {
		for {
			select {
			case <-ticker.C:
				usrmapRWLoker.Lock()
				if !config.Gconfig.Ethcfg.Transflag {
					updateBanlanceoff(dbMode)
				} else {
					updateBanlance(dbMode)
				}
				usrmapRWLoker.Unlock()
			}
		}
	}()
}
