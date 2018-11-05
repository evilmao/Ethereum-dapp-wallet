package eth

/*************************************************
Function: Withdraw
Description: 提现上链
Author:failymao
Date: 2018/7/11
History:
*************************************************/

//调用转账函数，
func Withdraw(serial_num, keypath, password, ExternalAccount *string, transferAmount *int64) (string, error) {
	wd_Hansh, err := TransferAccounts(serial_num, keypath, password, ExternalAccount, transferAmount)
	return wd_Hansh, err
}
