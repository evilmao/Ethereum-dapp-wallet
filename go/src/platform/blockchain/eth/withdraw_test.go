/*****************************************************************************
File name: withdraw_test.go
Description: 链上提现测试
Author: failymao
Version: V1.0
Date: 2018/07/23
History:
*****************************************************************************/

package eth

// 调用TransferAccountsV1函数
func WithdrawV1( keypath, password, ExternalAccount *string, transferAmount *int64) int64{
	//获取提现后的代币数量
	eth_num := TransferAccountsV1(keypath, password, ExternalAccount , transferAmount)
	return  eth_num

}