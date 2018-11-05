/*****************************************************************************
File name: radixConvert.go
Description:进制转换
Author: failymao
Version: V1.0
Date: 2018/07/16
History:
*****************************************************************************/
package tools

import (
	"strings"
	"strconv"
	"log"
	"math"
)

/*************************************************
Function: HexDec
Description:十六进制转换为十进制
Author:failymao
Date: 2018/07/16
History:
*************************************************/
func HexDec(h string) (n int64) {
	s := strings.Split(strings.ToUpper(h), "")
	l := len(s)
	i := 0
	d := float64(0)
	hex := map[string]string{"A": "10", "B": "11", "C": "12", "D": "13", "E": "14", "F": "15"}
	for i = 0; i < l; i++ {
		c := s[i]
		if v, ok := hex[c]; ok {
			c = v
		}
		f, err := strconv.ParseFloat(c, 10)
		if err != nil {
			log.Println("Hexadecimal to decimal error:", err.Error())
			return -1
		}
		d += f * math.Pow(16, float64(l-i-1))
	}
	return int64(d)
}


/*************************************************
Function: HexDec
Description:二进制转换为十进制
Author:failymao
Date: 2018/07/16
History:
*************************************************/
func BinDec(b string) (n int64) {
	s := strings.Split(b, "")
	l := len(s)
	i := 0
	d := float64(0)
	for i = 0; i < l; i++ {
		f, err := strconv.ParseFloat(s[i], 10)
		if err != nil {
			log.Println("Binary to decimal error:", err.Error())
			return -1
		}
		d += f * math.Pow(2, float64(l-i-1))
	}
	return int64(d)
}

/*************************************************
Function: HexDec
Description:八进制转换为十进制
Author:failymao
Date: 2018/07/16
History:
*************************************************/
func OctDec(o int64) (n int64) {
	s := strings.Split(strconv.Itoa(int(o)), "")
	l := len(s)
	i := 0
	d := float64(0)
	for i = 0; i < l; i++ {
		f, err := strconv.ParseFloat(s[i], 10)
		if err != nil {
			log.Println("Octal to decimal error:", err.Error())
			return -1
		}
		d += f * math.Pow(8, float64(l-i-1))
	}
	return int64(d)
}