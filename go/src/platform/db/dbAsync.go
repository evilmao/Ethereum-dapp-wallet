package db

import (
	"strconv"
	"platform/om/log"
	"time"
)

const (
	Get_prepared_user = iota
)

type AsyncQuery struct {
	QueryStr	string	//失败语句
	Count		int 	//失败次数
	Type		int		//query类型
}

var AsyncQuerysCh chan *AsyncQuery = make(chan *AsyncQuery, 1000)
var FailedQeurysCh chan *AsyncQuery = make(chan *AsyncQuery, 10)

func init() {
	dbMod, _ := SetConfig()
	go Asyncquery(dbMod)
	go Failquery(dbMod)
}

func Asyncquery(mod *Model) {
	for fquery := range AsyncQuerysCh {
		switch fquery.Type  {
		case Get_prepared_user :
			log.Debug(">>>>>>>>>>>>>%s", fquery.QueryStr)
			//数据库置未分配状态为已分配
			lidintr := mod.Query(fquery.QueryStr)
			switch lidintr.(type) {
			case string :
				lidstr := lidintr.(string)
				lid, err := strconv.Atoi(lidstr)
				if err != nil {
					log.Error("convert update line id err:%s", err.Error())
					fquery.Count++
					FailedQeurysCh <- fquery
				}
				if lid <= 0 {
					log.Error("update Userid[%d] fail", lid)
					fquery.Count++
					FailedQeurysCh <- fquery
				}
			default:
				log.Error("GetCreatedAccount:unexpected return %+v", lidintr)
				fquery.Count++
				FailedQeurysCh <- fquery
			}
		}
	}
}

func Failquery(mod *Model) {

	for fquery := range FailedQeurysCh {
		log.Debug("DB Query err:%s", fquery.QueryStr)
		if fquery.Count > 5 {
			continue
		} else {
			switch fquery.Type {
			case Get_prepared_user:

				time.Sleep(10 * time.Second)
				lidintr := mod.Query(fquery.QueryStr)
				switch lidintr.(type) {
				case string:

				default:
					fquery.Count++
					FailedQeurysCh <- fquery
					log.Error("GetCreatedAccount:unexpected return %+v", lidintr)
				}
			}
		}
	}
}