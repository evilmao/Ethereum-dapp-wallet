/*****************************************************************************
File name: returncode.go
Description: 错误信息
Author: failymao
Version: V1.0
Date: 2018/06/14
History:
*****************************************************************************/
package web

const(
	SUBCODE_INSUFFICIENT = "INSUFFICIENT-ACCOUNT-BALANCE"
)

const(
	SUBMSG_SUCCESS = "转账成功"
	SUBMSG_INSUFFICIENT = "账户余额不足"
	WITHDRAW_SUCCESS    = "提现成功"
)

const(
	CODE_SUCCESS = "10000"
	CODE_UNKNOW = "20000"
	CODE_UNSUPPORT = "20001"
	CODE_PARATYPEERR = "20002"
	CODE_PARAERR = "20003"
	CODE_NOTEXISTACCOUNT = "20004"
	CODE_EXISTACCOUNT    = "20005"
	CODE_NOUSERPREPARED  = "20006"
	CODE_INTERLACCOUNT   = "20007" //内部账户错误代码
	CODE_NEG_AMOUNT		 = "20009"
)

const(
	MSG_SUCCESS = "success"
	MSG_UNKNOW = "unknow-error"
	MSG_UNSUPPORT = "unsupport-method"
	MSG_PARATYPEERR = "Parameter-type-error"
	MSG_INVALID_NAME = "invalid-username"
	MSG_INVALID_TYPE = "invalid-usertype"
	MSG_INVALID_ACCOUNT = "invalid-account"
	MSG_NOTEXIST_TOACCOUNT = "notexist-to-account"
	MSG_NOTEXIST_FROMACCOUNT = "notexist-from-account"
	MSG_EXIST_ACCOUNT = "exist-account"
	MSG_NOPREPARED_USER = "no-user-avalable"
	MSG_UNSUPPORTED_CURRENCY = "unsupported-currency"
	MSG_NEGATIVE_AMOUNT = "negative-transfer-amount"
	MSG_INTERNAL_ACCOUNT     = "internal_account" //内部账户提示信息
)

func GetReturnCode(code string) string {
	switch code{
	case MSG_SUCCESS:
		return CODE_SUCCESS
	case MSG_UNKNOW:
		return CODE_UNKNOW
	case MSG_UNSUPPORT:
		return CODE_UNSUPPORT
	case MSG_PARATYPEERR:
		return CODE_PARATYPEERR
	case MSG_INVALID_NAME:
		return CODE_PARAERR
	case MSG_INVALID_TYPE:
		return CODE_PARAERR
	case MSG_INVALID_ACCOUNT:
		return CODE_PARAERR
	case MSG_NOTEXIST_TOACCOUNT:
		return CODE_NOTEXISTACCOUNT
	case MSG_NOTEXIST_FROMACCOUNT:
		return CODE_NOTEXISTACCOUNT
	case MSG_EXIST_ACCOUNT:
		return CODE_EXISTACCOUNT
	case MSG_NOPREPARED_USER:
		return CODE_NOUSERPREPARED
	case MSG_NEGATIVE_AMOUNT:
		return CODE_NEG_AMOUNT
	case MSG_INTERNAL_ACCOUNT:               //提现是内部账户提示信息
		return CODE_INTERLACCOUNT
	default:
		return CODE_UNKNOW
	}
	return CODE_UNKNOW
}

