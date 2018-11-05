package eth

import (
	"strconv"
	"platform/om/log"
	"fmt"
	"platform/db"
)
/*************************************************
Function: PrepareUser
Description: 准备一定数量的用户，以备新用户注册时分配，
	入参是需要准备的数量，返回准备好的用户数量
Author:
Date: 2018/06/15
History:
*************************************************/
func PrepareUser(count int) (prepared int) {

	result := db.Connect.Query("SELECT COUNT(1) AS COUNT FROM UserInfo WHERE Status = 1")
	switch result.(type) {
	case map[int]map[string]string:
		tabRes := result.(map[int]map[string]string)
		if len(tabRes) == 1 && len(tabRes[1]["COUNT"]) > 0 {
			count := tabRes[1]["COUNT"]
			counti, err := strconv.Atoi(count)
			if err == nil && counti >= 0{
				prepared = counti
				log.Debug("already prepared %d user", prepared)
			} else {
				return
			}
		}
	default:
		log.Error("unspected result:%+v", result)
		return
	}
	if prepared >= count {
		return
	}

	for i := count - prepared ; i > 0 ; i-- {
		acc, passwd, keystorePath := CreateAccount()
		insstr := fmt.Sprintf(`INSERT INTO UserInfo(Password, WalletAddress, Balance, Status, Timestamp, Keypath) VALUES('%s', '%s', '0', 1, NOW(), '%s')`,
									passwd, acc, keystorePath)
		result = db.Connect.Query(insstr)
		switch result.(type) {
		case string:
			lidstr := result.(string)
			lid , err := strconv.Atoi(lidstr)
			if err == nil && lid > 0 {
				prepared++
			} else {
				log.Error("unspected insert id:%+v", result)
			}
		default:
			log.Error("unspected result:%+v", result)
		}
	}

	return
}
