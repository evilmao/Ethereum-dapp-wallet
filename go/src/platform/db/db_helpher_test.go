package db

import (
	"fmt"
	"testing"
	"time"
	"strings"
	"math/rand"
	"platform/om/log"
)

func init() {
	log.New("debug", "")
}

func TestSetConfig(*testing.T) {
	_, err := SetConfig()
	if err != nil {
		fmt.Printf("SetConfig failed, error :%s\n", err)
	} else {
		fmt.Printf("SetConfig success.\n")
	}
	Currency()
}

func TestFind(*testing.T) {
	data := Connect.SetTable("UserInfo").Limit(1,20).FindAll()
	data = Connect.SetTable("UserInfo").Where("Balance = 0").FindAll()
	data = Connect.SetTable("UserInfo").Where("Balance = 0").Fileds("Ethernet_Current_Balance").OrderBy("Username ASC").FindAll()
	data = Connect.SetTable("UserInfo").Where("Balance = 0").Fileds("Ethernet_Current_Balance").OrderBy("Username ASC").Limit(1).FindAll()
	data = Connect.SetTable("UserInfo").Where("Username = 'Haxima001'").FindOne()
	Print(data)
}

func TestInsertandDelete(*testing.T) {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 5; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	res := string(result)

	tm := time.Now().Format(FormatNormalTime)
	var newuser = make(map[string]interface{})
	newuser["Username"] = res
	newuser["Password"] = "myDBtest"
	newuser["Status"] = MSG_ACTIVE
	newuser["WalletAddress"] = "myDBtest"
	newuser["Timestamp"] = tm
	newuser["Keypath"] = "myDBtest"
	newuser["nounce"] = 0
	data := Connect.SetTable("UserInfo").Fileds("Username", "Password", "Status", "WalletAddress", "Timestamp", "Keypath", "nounce").Where(strings.Join([]string{"Username = '" + res + "'"}, "")).FindOne()
	fmt.Println("TestFind data :", data)
	if 0 == len(data) {
		n,err := Connect.SetTable("UserInfo").Insert(newuser)
		if err != nil {
			fmt.Printf("TestInsert failed, error :%s\n", err)
		} else {
			fmt.Printf("TestInsert rows effected, n :%d\n", n)
			fmt.Printf("TestInsert success.\n")
		}
	}

	n,err := Connect.SetTable("UserInfo").Delete(strings.Join([]string{"Username = '" + res + "'"}, ""))
	if nil != err {
		fmt.Printf("TestDelete fail, error :%v\n", err)
	} else {
		fmt.Printf("TestDelete rows effected, n :%d\n", n)
		fmt.Printf("TestDelete success.\n")
	}

	sqlstr := "INSERT INTO UserInfo (`WalletAddress`,`Timestamp`,`Keypath`,`nounce`,`Username`,`Password`,`Status`) VALUES ('myDBtest','2018-07-31 11:15:36','myDBtest',0,'3rgks','myDBtest',0)"
	Connect.Query(sqlstr)

	n,err = Connect.SetTable("UserInfo").Delete(strings.Join([]string{"Username = '" + res + "'"}, ""))
	if nil != err {
		fmt.Printf("TestDelete fail, error :%v\n", err)
	} else {
		fmt.Printf("TestDelete rows effected, n :%d\n", n)
		fmt.Printf("TestDelete success.\n")
	}

	NewUserDBInsert(res, "myDBtest", "myDBtest", "myDBtest")
	msg := AccountCancellation("UserInfo", res, "myDBtest")
	if MSG_SUCCESS == msg {
		fmt.Printf("TestAccountCancellation success.\n")
	} else {
		fmt.Printf("TestAccountCancellation failed :%s\n",msg)
	}

	msg = AccountCancellation("UserInfo", "Nonexistent username", "myDBtest")
	if MSG_SUCCESS == msg {
		fmt.Printf("TestAccountCancellation Nonexistent username failed :%s\n",msg)
	} else {
		fmt.Printf("TestAccountCancellation Nonexistent username success :%s\n",msg)
	}

	msg = AccountCancellation("UserInfo", "test02", "myDBtest")
	if MSG_SUCCESS == msg {
		fmt.Printf("TestAccountCancellation test02 failed :%s\n",msg)
	} else {
		fmt.Printf("TestAccountCancellation test02 success :%s\n",msg)
	}
}

func TestUpdate(*testing.T) {
	str := "123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 5; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	res := string(result)

	var value = make(map[string]interface{})
	value["Bitcoin_Current_Balance"] = res
	_, err := Connect.SetTable("UserInfo").Where("Username = 'Haxima006'").Update(value)
	if err != nil {
		fmt.Printf("TestUpdate failed.\n")
	} else {
		fmt.Printf("TestUpdate success.\n")
	}
}

func TestQuery(*testing.T) {
	sqlstr := "SHOW DATABASES"
	result := Connect.Query(sqlstr)
	fmt.Printf("DATABASES :%s\n", result)
	sqlstr = "SELECT * FROM UserInfo"
	result = Connect.Query(sqlstr)
	switch result.(type) {
	case map[int]map[string]string:
		fmt.Printf("TestQuery success.\n")
	default:
		fmt.Printf("TestQuery failed.\n")
	}
}

func TestTimeDiff(*testing.T) {
	start_time := "2018-07-25 00:00:00"
	end_time := time.Now().Format(FormatNormalTime)

	diff := getSecondDiffer(start_time, end_time)
	fmt.Printf("start_time before end_time getSecondDiffer :%d\n", diff)
	diff = getSecondDiffer(end_time, start_time)
	fmt.Printf("start_time after end_time getSecondDiffer :%d\n", diff)

	diff = getMinDiffer(start_time, end_time)
	fmt.Printf("start_time before end_time getMinDiffer :%d\n", diff)
	diff = getMinDiffer(end_time, start_time)
	fmt.Printf("start_time after end_time getMinDiffer :%d\n", diff)

	diff = getHourDiffer(start_time, end_time)
	fmt.Printf("start_time before end_time getHourDiffer :%d\n", diff)
	diff = getHourDiffer(end_time, start_time)
	fmt.Printf("start_time after end_time getHourDiffer :%d\n", diff)

	diff = getDayDiffer(start_time, end_time)
	fmt.Printf("start_time before end_time getDayDiffer :%d\n", diff)
	diff = getDayDiffer(end_time, start_time)
	fmt.Printf("start_time after end_time getDayDiffer :%d\n", diff)
}

func TestTransBetweenUsers(*testing.T) {
	serial_num := GetSerialNum()
	msg := TransBetweenUsers (serial_num, "Haxima006", "Haxima001", "1", "0", "1")
	if MSG_SUCCESS == msg {
		fmt.Printf("TransBetweenUsers success.\n")
	} else {
		fmt.Printf("TransBetweenUsers failed :%s\n", msg)
	}

	result := UpdateBalanceTXstatus(serial_num, "0", "Haxima001", "Haxima006", 1, 1)
	if result {
		fmt.Printf("TestUpdateBalanceTXstatus success.\n")
	} else {
		fmt.Printf("TestUpdateBalanceTXstatus failed.\n")
	}

	serial_num = GetSerialNum()
	msg = TransBetweenUsers (serial_num, "Haxima006", "Haxima001", "1", "1", "1")
	if MSG_SUCCESS == msg {
		fmt.Printf("TransBetweenUsers success.\n")
	} else {
		fmt.Printf("TransBetweenUsers failed :%s\n", msg)
	}
	result = UpdateBalanceTXstatus(serial_num, "1", "Haxima001", "Haxima006", 0, 1)
	if result {
		fmt.Printf("TestUpdateBalanceTXstatus success.\n")
	} else {
		fmt.Printf("TestUpdateBalanceTXstatus failed.\n")
	}

	serial_num = GetSerialNum()
	result = UpdateBalanceTXstatus(serial_num, "1", "Haxima001", "Haxima006", 3, 1)
	if result {
		fmt.Printf("TestUpdateBalanceTXstatus invalid TXstatus failed.\n")
	} else {
		fmt.Printf("TestUpdateBalanceTXstatus invalid TXstatus success.\n")
	}

	balance := QueryBalance("UserInfo", "Haxima001")
	fmt.Printf("TestQueryBalance balance before :%d\n", balance)
}


func TestTransWithdraw(*testing.T) {
	serial_num := GetSerialNum()

	msg := TransWithdraw(serial_num, "Haxima006", "mine", "1","0", "1")
	if MSG_SUCCESS == msg {
		fmt.Printf("TestTransWithdraw success.\n")
	} else {
		fmt.Printf("TestTransWithdraw failed :%s\n", msg)
	}
}

