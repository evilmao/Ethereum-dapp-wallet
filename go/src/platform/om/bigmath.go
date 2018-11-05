/*************************************************
Function: str2Big
Description:转换string到bigint之间的转换，出错返回nil
Author:
Date: 2018/07/13
History:
*************************************************/

package om

import "math/big"


func Str2Big(str string) *big.Int {
	symbo := str[:2]
	i := new(big.Int)
	if symbo == `0x` {
		i, _ = i.SetString(str[2:], 16)
	} else {
		i, _ = i.SetString(str, 10)
	}
	return i
}

func Big2Str(i *big.Int, base int) string {
	if base == 10 {
		return i.Text(10)
	}
	if base == 16 {
		return `0x` + i.Text(16)
	}
	return ``
}

// GWei -> big int的单位是GWei
func Gwei2Big(gwei int) *big.Int {
	i := big.NewInt(int64(gwei))
	return  i.Mul(i, big.NewInt(1000000000))
}

// float -> big, float最大精度为9位
func Float2Big(f float64) *big.Int {
	i := big.NewInt(int64(f * 1000000000))
	return i.Mul(i, big.NewInt(	1000000000))
}