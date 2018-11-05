/**************************************************************
【提交类型】:BUG/新功能/需求修改/版本制作/代码整理/解决编译不过
【问题描述】:db模块代码
【修改内容】:初次提交
【提交人】:failymao
【评审人】:
***************************************************************/
package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"regexp"
	"strconv"
	"strings"
	"time"
	"platform/config"
	"platform/om/log"
	"github.com/gin-gonic/gin"
	"crypto/sha256"
	"encoding/hex"
	glog "log"
)

//const(
//	MSG_DEFICIENT_BALANCE = "deficient-balance"         //转账账户余额不足
//	MSG_TXFAILED = "tx-failed"                          //转账交易失败
//)

const (
	FormatNormalTime = "2006-01-02 15:04:05"
	UserInfo = "UserInfo"
	AppID = "AppID"
	UserName = "Username"
	PassWord = "Password"
	AccountStatus = "Status"
	WalletAddress = "WalletAddress"
	TimeStamp = "Timestamp"
	StatusIsZero = "' and Status = 0"
	UsernameEqual = "Username = '"
	Keypath = "Keypath"
	Nounce = "nounce"
	TransactionDetail = "transaction_detail"
	TXStatus = "TXstatus"
)

const(
	MSG_SUCCESS = "success"                             //操作成功
	MSG_INVALID_INPUTINFO = "invalid-inputinfo"         //无效的输入信息（用户名、密码错误或发出请求的是已被注销的账户）
	MSG_BALANCE_NOTZERO = "balance-notzero"             //发出销户请求的账户余额不为0，销户失败
	MSG_DBERR = "db-error"                              //数据库错误
	//MSG_RMVKEYFILE_FAILED = "remove-keyfile-failed"   //删除key文件失败
)

const(
	//MSG_TX_LIFE_CYCLE = 15                            //转账等待超时时间（minute）
	MSG_AWAIT_CONFIRM = 0                               //转账交易待确认
	MSG_CONFIRMED = 1                                   //转账交易已确认
	MSG_OVERTIME = 2                                    //转账交易已被超时取消
)

const(
	MSG_ACTIVE = 0                                      //有效活跃账户
	MSG_DEACTIVED = 1                                   //已注销账户
)

const(
	MSG_LICENSE_TIME = 365                              //用户license有效期（day）
	MSG_EXPIRED_USER = "expired-user"                   //license过期用户
	MSG_VERIFY_FAILED = "verify-failed"                 //公共请求消息签名串与实际不符，校验失败
)

const(
	maxOpen = 200                       //最大打开连接数，避免并发太高时连接mysql出现too many connections错误
	maxIdle = 50                        //最大闲置连接数，已开启的连接使用完成后归还连接池，等待下一次使用
	maxConnLifetime = 300*time.Second   //数据库连接超时时间，超时后数据库会单方面断掉连接
)

const(
	maxByte = 6                         //同秒级转账交易顺序号不足值补零位数
)

const(
	ethernet  = "0"                    //以太币标志：  "0"
	bit = "1"                          //比特币标志：  "1"
	ETHERNETCOIN = "ethernetcoin"
	BITCOIN = "bitcoin"

	//新币种加入新增币种名即可，目前数据库支持5个币种
)



type Model struct {
	db        *sql.DB
	tablename string
	param     []string
	columnstr string
	where     string
	pk        string
	orderby   string
	limit     string
	join      string
}

var Connect *Model = new(Model)
var trans_order string
var counter int
var Currencies = make(map[string]string)

//初始化，连接数据库并完成连接池设置
func init() {
	//db, _ = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/test?charset=utf8")
	var _ error
	Connect, _ = SetConfig()
	Connect.db.SetMaxOpenConns(maxOpen)
	Connect.db.SetMaxIdleConns(maxIdle)
	Connect.db.SetConnMaxLifetime(maxConnLifetime)
	Connect.db.Ping()
	Currency()
}


//初始化币种map
func Currency() (map[string]string) {
	Currencies ["0"] = ETHERNETCOIN
	Currencies ["1"] = BITCOIN
	return Currencies
}


//读取配置文件信息，连接数据库
func SetConfig() (*Model, error) {
	c := new(Model)
	charset := config.Gconfig.Mysqlcfg.Charset
	username := config.Gconfig.Mysqlcfg.Username
	password := config.Gconfig.Mysqlcfg.Password
	hostname := config.Gconfig.Mysqlcfg.Hostname
	database := config.Gconfig.Mysqlcfg.Database
	port := config.Gconfig.Mysqlcfg.Port
	
	db, err := sql.Open("mysql", username+":"+password+"@tcp("+hostname+":"+fmt.Sprintf("%d",port)+")/"+database+"?charset="+charset)
	err = db.Ping()
	if err != nil {
		//if connect error then return the error message
		glog.Fatal("connect DB err", err.Error())
		return c, err
	}
	c.db = db
	return c, err
}


//查找所有匹配结果
func (m *Model) FindAll() map[int]map[string]string {

	result := make(map[int]map[string]string)
	if m.db == nil {
		fmt.Printf("Findall db notconnect")
		return result
	}
	//fmt.Printf("=========len(m.param): %d=========\n", len(m.param))
	if len(m.param) == 0 {
		m.columnstr = "*"
	} else {
		if len(m.param) == 1 {
			m.columnstr = m.param[0]
		} else {
			m.columnstr = strings.Join(m.param, ",")
		}

	}

	query := fmt.Sprintf("Select %v from %v %v %v %v %v;", m.columnstr, m.tablename, m.join, m.where, m.orderby, m.limit)
	fmt.Println("SQL statement is :",query)
	rows, err := m.db.Query(query)
	defer rows.Close()
	if err != nil {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("SQL syntax errors :%s\n", err)
				log.Error("SQL syntax errors")
			}
		}()
		err = errors.New("select sql failure")
	}
	result = QueryResult(rows)
	return result
}


//查找一个匹配结果
func (m *Model) FindOne() map[int]map[string]string {
	empty := make(map[int]map[string]string)
	if m.db != nil {
		data := m.Limit(1).FindAll()
		return data
	}
	log.Error("mysql not connect\r\n")
	return empty
}


//插入数据
func (m *Model) Insert(param map[string]interface{}) (num int, err error) {
	if m.db == nil {
	    log.Error("mysql not connect\r\n")
		return 0, errors.New("IN Insert, mysql not connect")
	}
	var keys []string
	var values []string
	if len(m.pk) != 0 {
		delete(param, m.pk)
	}

	for key, value := range param {
		keys = append(keys, key)
		switch value.(type) {
		case int, int64, int32:
			values = append(values, strconv.FormatInt(int64(value.(int)), 10))
		case uint64, uint32:
			values = append(values, strconv.FormatUint(value.(uint64), 10))
		case string:
			values = append(values, "'" + value.(string) + "'")
		//case float32, float64:
		//	values = append(values, strconv.FormatFloat(value.(float64), 'f', -1, 64))
		}
	}
	fileValue := strings.Join(values, ",")
	fileds := "`" + strings.Join(keys, "`,`") + "`"
	sql := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v);", m.tablename, fileds, fileValue)
	var query = strings.TrimSpace(sql)
	fmt.Printf("insert sql :%s\n", query)
	//result, err := m.db.Exec(sql)
	result, err := m.db.Exec(query)
	if err != nil {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("SQL syntax errors ")
			}
		}()
		err = errors.New("inster sql failure")
		log.Error("inster sql failure.error :%s", err)
		return 0, err
	}
	//i, err := result.LastInsertId()
	i, err := result.RowsAffected()
	s, _ := strconv.Atoi(strconv.FormatInt(i, 10))
	if err != nil {
		err = errors.New("insert failure")
	}
	return s, err

}


//指定字段
func (m *Model) Fileds(param ...string) *Model {
	m.param = param
	return m
}


//更新表数据
func (m *Model) Update(param map[string]interface{}) (num int, err error) {
	if m.db == nil {
		return 0, errors.New("mysql not connect")
	}
	var setValue []string
	for key, value := range param {
		switch value.(type) {
		case int, int64, int32:
			set := fmt.Sprintf("%v = %v", key, value.(int))
			setValue = append(setValue, set)
		case string:
			set := fmt.Sprintf("%v = '%v'", key, value.(string))
			setValue = append(setValue, set)
		//case float32, float64:
		//	set := fmt.Sprintf("%v = '%v'", key, strconv.FormatFloat(value.(float64), 'f', -1, 64))
		//	setValue = append(setValue, set)
		}

	}
	setData := strings.Join(setValue, ",")
	sql := fmt.Sprintf("UPDATE %v SET %v %v", m.tablename, setData, m.where)
	fmt.Printf("update_sql :%s\n", sql)
	result, err := m.db.Exec(sql)
	if err != nil {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("SQL syntax errors ")
			}
		}()
		err = errors.New("update sql failure")
		return 0, err
	}
	i, err := result.RowsAffected()
	if err != nil {
		err = errors.New("update failure")
		log.Error("update tabledata error:%s", err)
		return 0, err
	}
	s, _ := strconv.Atoi(strconv.FormatInt(i, 10))

	return s, err
}


//删除数据
func (m *Model) Delete(param string) (num int, err error) {
	if m.db == nil {
		return 0, errors.New("mysql not connect")
	}
	h := m.Where(param).FindOne()
	if len(h) == 0 {
		return 0, errors.New("no Value")
	}
	sql := fmt.Sprintf("DELETE FROM %v WHERE %v", m.tablename, param)
	result, err := m.db.Exec(sql)
	if err != nil {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("SQL syntax errors: %+v", err)
				log.Error("SQL syntax errors:%+v", err)
			}
		}()
		err = errors.New("delete sql failure")
		return 0, err
	}
	i, err := result.RowsAffected()
	s, _ := strconv.Atoi(strconv.FormatInt(i, 10))
	if i == 0 {
		err = errors.New("delete failure")
	}

	return s, err
}


//执行自定义sql语句
func (m *Model) Query(sql string) interface{} {
	if m.db == nil {
		return errors.New("mysql not connect")
	}
	var query = strings.TrimSpace(sql)
	s, err := regexp.MatchString(`(?i)^(select|call)`, query)
	if nil == err && s {
		result, _ := m.db.Query(sql)
		defer result.Close()
		c := QueryResult(result)
		return c
	}
	exec, err := regexp.MatchString(`(?i)^(update|delete)`, query)
	if nil == err && exec {
		m_exec, err := m.db.Exec(query)
		if err != nil {
			return err
		}
		num, _ := m_exec.RowsAffected()
		id := strconv.FormatInt(num, 10)
		return id
	}

	insert, err := regexp.MatchString(`(?i)^insert`, query)
	if nil == err && insert {
		m_exec, err := m.db.Exec(query)
		if err != nil {
			return err
		}
		num, _ := m_exec.LastInsertId()
		id := strconv.FormatInt(num, 10)
		return id
	}
	result, _ := m.db.Exec(query)

	return result

}


//返回sql语句执行结果
func QueryResult(rows *sql.Rows) map[int]map[string]string {
	var result = make(map[int]map[string]string)
	columns, _ := rows.Columns()
	values := make([]sql.RawBytes, len(columns))
	scanargs := make([]interface{}, len(values))
	for i := range values {
		scanargs[i] = &values[i]
	}

	var n = 1
	for rows.Next() {
		result[n] = make(map[string]string)
		err := rows.Scan(scanargs...)

		if err != nil {
			fmt.Println(err)
		}

		for i, v := range values {
			result[n][columns[i]] = string(v)
		}
		n++
	}

	return result
}


//指定待查询表名
func (m *Model) SetTable(tablename string) *Model {
	m.tablename = tablename
	return m
}


//设置where条件
func (m *Model) Where(param string) *Model {
	m.where = fmt.Sprintf(" where %v", param)
	return m
}

/*
//设置自增主键字段
func (m *Model) SetPk(pk string) *Model {
	m.pk = pk
	return m
}*/


//设置排序方式
func (m *Model) OrderBy(param string) *Model {
	m.orderby = fmt.Sprintf("ORDER BY %v", param)
	return m
}


//设置返回结果个数
func (m *Model) Limit(size ...int) *Model {
	var end int
	start := size[0]
	//fmt.Printf("=========len(size): %d=========\n", len(size))
	if len(size) > 1 {
		end = size[1]
		m.limit = fmt.Sprintf("Limit %d,%d", start, end)
		return m
	}
	m.limit = fmt.Sprintf("Limit %d", start)
	return m
}

/*
//左连接
func (m *Model) LeftJoin(table, condition string) *Model {
	m.join = fmt.Sprintf("LEFT JOIN %v ON %v", table, condition)
	return m
}


//右连接
func (m *Model) RightJoin(table, condition string) *Model {
	m.join = fmt.Sprintf("RIGHT JOIN %v ON %v", table, condition)
	return m
}


//内连接
func (m *Model) Join(table, condition string) *Model {
	m.join = fmt.Sprintf("INNER JOIN %v ON %v", table, condition)
	return m
}


//外连接
func (m *Model) FullJoin(table, condition string) *Model {
	m.join = fmt.Sprintf("FULL JOIN %v ON %v", table, condition)
	return m
}
*/

//将结果输出到屏幕
func Print(slice map[int]map[string]string) {
	for _, v := range slice {
		for key, value := range v {
			fmt.Println(key, value)
		}
		fmt.Println("---------------")
	}
}


//关闭数据库
//func (m *Model) DbClose() {
//	m.db.Close()
//}


//计算秒数时间差
func getSecondDiffer(start_time string, end_time string) int64 {
	var second int64
	t1, err := time.ParseInLocation(FormatNormalTime, start_time, time.Local)
	t2, err := time.ParseInLocation(FormatNormalTime, end_time, time.Local)
	if err == nil && t1.Before(t2) {
		second = t2.Unix() - t1.Unix()
		return second
	} else {
		return second
	}
}


//计算分钟时间差
func getMinDiffer(start_time string, end_time string) int64 {
	var minute int64
	t1, err := time.ParseInLocation(FormatNormalTime, start_time, time.Local)
	t2, err := time.ParseInLocation(FormatNormalTime, end_time, time.Local)
	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix()
		minute = diff / 60
		return minute
	} else {
		return minute
	}
}


//计算小时时间差
func getHourDiffer(start_time string, end_time string) int64 {
	var hour int64
	t1, err := time.ParseInLocation(FormatNormalTime, start_time, time.Local)
	t2, err := time.ParseInLocation(FormatNormalTime, end_time, time.Local)
	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix()
		hour = diff / 3600
		return hour
	} else {
		return hour
	}
}


//计算天数时间差
func getDayDiffer(start_time string, end_time string) int64 {
	var day int64
	t1, err := time.ParseInLocation(FormatNormalTime, start_time, time.Local)
	t2, err := time.ParseInLocation(FormatNormalTime, end_time, time.Local)
	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix()
		day = diff / 86400
		return day
	} else {
		return day
	}
}


//用户鉴权
func Authentication (c *gin.Context) (msg string){
	appid := c.Query("app_id")
	sign := c.Query("sign")
	timestamp := c.Query("timestamp")

    current := time.Now().Format(FormatNormalTime)
	day := getDayDiffer(timestamp, current)
	if day >= MSG_LICENSE_TIME {
		return MSG_EXPIRED_USER
	}

	data := Connect.SetTable(UserInfo).Fileds(AppID, PassWord, TimeStamp).Where(strings.Join([]string{"AppID = '" + appid + StatusIsZero}, "")).FindOne()
	info := data[1]
	if nil == info {
		fmt.Println("Invalid Authentication input information")
		return MSG_INVALID_INPUTINFO
	}
	//不能从缓存中取，否则造成eth和db的循环引用
	//从数据库中读取用户信息并加密，判断解析后的公共参数是否与从数据库取出并加密后的哈希值相等，相等则鉴权成功
	plaintext := strings.Join([]string{info[AppID] + ":" + info[PassWord]+ ":" + info[TimeStamp]}, "")
	hash := sha256.New()
	hash.Write([]byte(plaintext))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)

	if sign != mdStr {
		fmt.Println("Verification failed!")
		return MSG_VERIFY_FAILED
	}
	fmt.Println("Verification success!")
	return MSG_SUCCESS
}


//创建新用户
func NewUserDBInsert (username string, password string, account string, keypath string) {

	tm := time.Now().Format(FormatNormalTime)

	var newuser = make(map[string]interface{})
	newuser[UserName] = username
	newuser[PassWord] = password

	newuser[AccountStatus] = MSG_ACTIVE
	newuser[WalletAddress] = account
	newuser[TimeStamp] = tm
	newuser[Keypath] = keypath
	newuser[Nounce] = 0
	t := Connect.SetTable(UserInfo)
	data := t.Fileds(UserName, PassWord, AccountStatus, WalletAddress, TimeStamp, Keypath, Nounce).Where(strings.Join([]string{UsernameEqual + username + "'"}, "")).FindOne()
	if len(data) == 0{
		_,err := t.Insert(newuser)
		if err != nil {
			log.Error("openaccount insert value to UserInfo failed, error :%s", err)
		}
	}
}


//更新账户余额和交易状态
func UpdateBalanceTXstatus(serial_num string, transferCurrency string, transferor string, receiptor string, TXstatus int, amount int64) (bool) {

	var value = make(map[string]interface{})
	if MSG_CONFIRMED == TXstatus {
		value[TXStatus] = MSG_CONFIRMED
		n, err := Connect.SetTable(TransactionDetail).Where(strings.Join([]string{"serial_number = '" + serial_num + "'"}, "")).Update(value)
		fmt.Printf("n :%d\n", n)
		if err != nil {
			fmt.Printf("UpdateBalance TXstatus(1) MSG_CONFIRMED to transaction_detail failed: %s\n", err)
			return false
		}
		UpdateBalance(transferCurrency, transferor, -1*amount)
		UpdateBalance(transferCurrency, receiptor, amount)
	} else if MSG_AWAIT_CONFIRM == TXstatus{
		value[TXStatus] = MSG_OVERTIME
		_, err := Connect.SetTable(TransactionDetail).Where(strings.Join([]string{"serial_num = '" + serial_num + "'"}, "")).Update(value)
		if err != nil {
			fmt.Println("UpdateBalance TXstatus(0) MSG_AWAIT_CONFIRM to transaction_detail failed:")
		}
		fmt.Println("The transaction has been automatically cancelled for timeout.")
		return false
	}else {
		fmt.Println("An unknown error occurred and the transaction has been cancelled.")
		return false
	}
	return true
}


//查找操作前账户余额
func QueryBalance(tablename string, username string) (balance_before int){

	var data map[int]map[string]string
	var result string
	for key := range Currencies {
		switch key {
		case ethernet: //查找以太币账户余额
			data = Connect.SetTable(tablename).Fileds("Ethernet_Current_Balance").Where(strings.Join([]string{UsernameEqual + username + StatusIsZero}, "")).FindOne()
			if len(data) != 0 {
				result = data[1]["Ethernet_Current_Balance"]
				ethBalance, err := strconv.Atoi(result)
				fmt.Printf("ethBalance :%d\n", ethBalance)
				if (nil == err) && (0 != ethBalance){
					balance_before = ethBalance
					fmt.Println("Ethernet_Current_Balance :", ethBalance)
					return balance_before
				}
			}
		case bit: //查找比特币账户余额
			data = Connect.SetTable(tablename).Fileds("Bitcoin_Current_Balance").Where(strings.Join([]string{UsernameEqual + username + StatusIsZero}, "")).FindOne()
			if len(data) != 0 {
				result = data[1]["Bitcoin_Current_Balance"]
				bitBalance, err := strconv.Atoi(result)
				fmt.Printf("bitBalance :%d\n", bitBalance)
				if (err == nil) && (0 != bitBalance){
					balance_before = bitBalance
					fmt.Println("Bitcoin_Current_Balance :", bitBalance)
					return balance_before
				}
			}
		default:
			return 0
		}
	}
	return balance_before
}


//更新交易后账户余额 正数增加，负数减少
func UpdateBalance(transferCurrency string, username string, amount int64){

	var sqlstr = ""

	switch transferCurrency {
	case ethernet: //查找以太币账户余额
		if amount >= 0 {
			sqlstr = fmt.Sprintf("UPDATE UserInfo SET Ethernet_Current_Balance = CAST( (CAST(Ethernet_Current_Balance AS UNSIGNED) + %d) AS CHAR ) WHERE Username = '%s'; ", amount, username)
			fmt.Printf("sqlstr :%s\n", sqlstr)
		} else {
			amount = amount * -1
			sqlstr = fmt.Sprintf("UPDATE UserInfo SET Ethernet_Current_Balance = CAST( (CAST(Ethernet_Current_Balance AS UNSIGNED) - %d) AS CHAR ) WHERE Username = '%s' AND CAST(Ethernet_Current_Balance AS UNSIGNED) > %d;",
				amount, username, amount)
			fmt.Printf("sqlstr :%s\n", sqlstr)
		}
	case bit: //查找比特币账户余额
		if amount >= 0 {
			sqlstr = fmt.Sprintf("UPDATE UserInfo SET Bitcoin_Current_Balance = CAST( (CAST(Bitcoin_Current_Balance AS UNSIGNED) + %d) AS CHAR ) WHERE Username = '%s'; ", amount, username)
			fmt.Printf("sqlstr :%s\n", sqlstr)
		} else {
			amount = amount * -1
			sqlstr = fmt.Sprintf("UPDATE UserInfo SET Bitcoin_Current_Balance = CAST( (CAST(Bitcoin_Current_Balance AS UNSIGNED) - %d) AS CHAR ) WHERE Username = '%s' AND CAST(Bitcoin_Current_Balance AS UNSIGNED) > %d;",
				amount, username, amount)
			fmt.Printf("sqlstr :%s\n", sqlstr)
		}
	}
	result := Connect.Query(sqlstr)
	switch result.(type) {
	case string :
		afrowstr := result.(string)
		lid, err := strconv.Atoi(afrowstr)
		if err != nil {
			log.Error("convert update line id err:%s", err.Error())
		}
		if lid < 0 {
			log.Error("update balance[%d] fail", lid)
		}
	default:
		log.Error("UpdateBalance:unexpected return %+v", result)
	}
}


//判断该交易是同一秒内的第几个交易
func TransaOrder(tm string) (int){
	if trans_order == tm{
		counter = counter + 1
	}else{
		trans_order = tm
		counter = 1
	}
	return counter
}


//生成交易流水号
func GetSerialNum()(string) {
	//获取当前时间
	current := time.Now().Format(FormatNormalTime)
	tm := strings.Replace(current, " ", "", -1)
	tm = strings.Replace(tm, "-", "", -1)
	tm = strings.Replace(tm, ":", "", -1)
	//获取当前秒内交易顺序号
	num := TransaOrder(tm)
	subNum := strconv.Itoa(num)
	//获取交易流水号（一秒内交易次数大于等于十万时，直接拼接不再补零）
	length := maxByte - len([]rune(subNum))
	if length > 0 {
		for i := 0; i < length; i++ {
			subNum = strings.Join([]string{ "0" + subNum}, "")
		}
	}

	serial_number := tm + subNum
	fmt.Printf("serial_number: %s\n", serial_number)
	return serial_number
}



//将交易记录插入交易明表
// add 2018-7-4  shangwj  交易明细表增加 以太坊交易hash值  txhash
func InsertDetail(Currency string, SourceUser string, DestUser string, transferAmount uint64, serialNumber string, txhash string) (err error) {

	switch Currency {
	case ethernet:   //查找以太币账户余额
		Currency = ETHERNETCOIN
		fmt.Printf("Currency: %s\n", Currency)
	case bit:   //查找比特币账户余额
		Currency = BITCOIN
		fmt.Printf("Currency: %s\n", Currency)
		//转账参数校验已经校验过币种，此处不用再校验
	}

	var value= make(map[string]interface{})
	value["serial_number"] = serialNumber
	value["currency"] = Currency
	value["transferor"] = SourceUser
	value["receiptor"] = DestUser
	value["transfer_amount"] = transferAmount
	value[TXStatus] = MSG_AWAIT_CONFIRM
	value["Txhash"] = txhash                          //add 2018-7-4  shangwj  交易明细表增加以太坊交易hash值
	n,err := Connect.SetTable(TransactionDetail).Insert(value)
	if err != nil {
		fmt.Println("UpdateBalance insert value to transaction_detail failed:", err)
		return err
	}
	fmt.Println("UpdateBalance effected rows quantity is :", n)
	return err
}


//用户间转账
func TransBetweenUsers (serialNum string, SourceUser string, DestUser string, TransferAmount string, TransferCurrency string, Txhash string)  (msg string){

	transferAmount, err := strconv.Atoi(TransferAmount)
	if err != nil{
		log.Error("change amount string to int failed, error:%s", err)
	}
	//将最新的交易明细插入记账表
	//  add 2018-7-4  shangwj  交易明细表增加 以太坊交易hash值  Txhash
	err = InsertDetail(TransferCurrency, SourceUser, DestUser, uint64(transferAmount), serialNum, Txhash)
	if err != nil {
		log.Error("UpdateBalance failed, error:%s", err)
		fmt.Println("UpdateBalance failed, error :", err)
		return MSG_DBERR
	}

	UpdateBalance(TransferCurrency, SourceUser, -1*int64(transferAmount))
	UpdateBalance(TransferCurrency, DestUser, int64(transferAmount))
	return MSG_SUCCESS
}


//销户
func AccountCancellation(tablename string, Username string, Password string) (msg string) {

	//核对用户名、密码以及账户活跃状态
	n := Connect.SetTable(tablename).Where(strings.Join([]string{UsernameEqual + Username + "' and  Password = '"+ Password + StatusIsZero}, "")).FindOne()
	if nil == n[1] {
		fmt.Println("invalid AccountCancellation username or password")
		return MSG_INVALID_INPUTINFO
	}
    //余额为0时方可销户
    balance := QueryBalance(tablename, Username)
	if 0 != balance {
		fmt.Printf("cancellation can only succeed if the balance of account is zero,now your balance is:%d\n", balance)
		return MSG_BALANCE_NOTZERO
	}
	//销户
	/*
	keypath := Connect.SetTable(tablename).Fileds(Keypath).Where(strings.Join([]string{UsernameEqual + Username + "'"}, "")).FindOne()
	//删除key文件
	err := os.Remove(keypath[1][Keypath])
	if err != nil {
		fmt.Printf("keypath file remove failed! err :%s\n", err)
		return MSG_RMVKEYFILE_FAILED
	}*/
    //更新销户账号用户状态和地址
	var value = make(map[string]interface{})
	value[AccountStatus] = MSG_DEACTIVED
	//value[WalletAddress] = "0"
	_, err := Connect.SetTable(tablename).Where(strings.Join([]string{UsernameEqual + Username + "'"}, "")).Update(value)
	if err != nil {
		fmt.Printf("Update insert value to %s failed, error : %s\n", tablename, err)
		return MSG_DBERR
	}
	return MSG_SUCCESS
}


/*************************************************
Function: TransWithdraw
Description:提现数据库操作
Author:failymao
Date: 2018/7/11
History:
*************************************************/
// 提现：钱包地址向外部账户转账
func TransWithdraw(serialNum string, SourceUser string, ExternalAccount string,
	TransferAmount string,TransferCurrency string, Txhash string) (msg string) {

	transferAmount, err1 := strconv.Atoi(TransferAmount) //提现金额类型转换，string--->int
	if err1 != nil {
		log.Error("change amount string to int failed, error:%s", err1)
	}

	err := InsertDetail(TransferCurrency,SourceUser, ExternalAccount, uint64(transferAmount), serialNum, Txhash) //交易明细插入记账表
	if err != nil {
		log.Error("UpdateBalance failed, error:%s", err)
		return MSG_DBERR
	}
	UpdateBalance(TransferCurrency, SourceUser, -1*int64(transferAmount))
	UpdateBalance(TransferCurrency, ExternalAccount, int64(transferAmount))
	return MSG_SUCCESS
}
