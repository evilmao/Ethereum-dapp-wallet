/*****************************************************************************
File name: account.go
Description: 创建账户
Author:尚文静
Version: V1.0
Date: 2018/06/15
History:
*****************************************************************************/
package eth

import (
	"platform/om/log"
	"fmt"
	"platform/db"
)

const password string="123456"   						//账户密码
var allocator chan struct{} = make(chan struct{}, 1) 	//分配锁


func init() {
	allocator <- struct{}{} //分配锁初始化
}

/*************************************************
Function: createAccount
Description: 创建以太坊账户
Author:
Date: 2018/06/15
History:
*************************************************/
func CreateAccount() (string ,string,string) {

	//start := time.Now()
	//创建账户
	ks := GetKeyStore()

	addressContent, _ := ks.NewAccount(password)
	//defer log.Debug(">>>>>ks.NewAccount:%v", time.Since(start))
	_, err := ks.Export(addressContent, password, password)
	if err != nil {
		log.Error("new account fail :",err)
	}
	//defer log.Debug(">>>>>ks.Export:%v", time.Since(start))
	address :=addressContent.Address.Hex()
    keystorePath :=addressContent.URL.Path
	return address,password,keystorePath
}

/*************************************************
Function: GetCreatedAccount
Description:从数据库中获取已经准备好的账户给用户
Author:
Date: 2018/06/15
History:
*************************************************/
func GetCreatedAccount(usrname, appid, ctime string) (string, string, string) {
	<-allocator
	for name, usr := range usrmapunused {
		usr.Username = usrname
		delete(usrmapunused, name)
		usrmapRWLoker.Lock()
		usrmap[usrname] = usr
		usrmapRWLoker.Unlock()
		db.AsyncQuerysCh <- &db.AsyncQuery{
			QueryStr: fmt.Sprintf(`UPDATE UserInfo SET Username = '%s', 
Status = 0 , timestamp = '%s', AppID = '%s' WHERE Username = '%s'`, usrname, ctime, appid, name),
			Type: db.Get_prepared_user,
		}
		allocator<-struct{}{}
		//修改本地缓存
		return usr.WalletAddr, usr.Password, usr.KeyPath
	}

	log.Error("no user prepared")
	return "", "", ""
}
/*************************************************
Function: DelUserInfo
Description:删除本地缓存中的指定用户信息
Author:
Date: 2018/07/05
History:
*************************************************/
func DelUserInfo(usr string) {
	usrmapRWLoker.Lock()
	defer usrmapRWLoker.Unlock()
	if _, ok := usrmap[usr] ; ok {
		delete(usrmap, usr)
	}
}