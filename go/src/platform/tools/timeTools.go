package tools

import "time"

//计算分钟时间差
func GetMinDiffer(startTime string, endTime string) int64 {
	var minute int64
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", startTime, time.Local)
	t2, err := time.ParseInLocation("2006-01-02 15:04:05", endTime, time.Local)
	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix()
		minute = diff / 60
		return minute
	} else {
		return minute
	}
}