/*****************************************************************************
File name: encode.go
Description: web模块消息封装
Author: failymao
Version: V1.0
Date: 2018/06/14
History:
*****************************************************************************/
package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func transferencode(c *gin.Context, code transfercode, retmsg string){
	if SUBMSG_SUCCESS == code.sub_code{
		c.JSON(http.StatusOK, gin.H{
			"msg": retmsg,
			"code": GetReturnCode(retmsg),
		})
	}else {
		c.JSON(http.StatusOK, gin.H{
			"sub_msg" :code.sub_msg,
			"sub_code": code.sub_code,
			"msg": retmsg,
			"code": GetReturnCode(retmsg),
		})
	}

}

func accountencode(c *gin.Context, reqInfo AccountBody, retmsg string, account string, password string){
	if userENT == reqInfo.Usertype{
		c.JSON(http.StatusOK, gin.H{
			"password" :password,
			"account": account,
			"msg": retmsg,
			"code": GetReturnCode(retmsg),
		})
	}else {
		c.JSON(http.StatusOK, gin.H{
			"account": account,
			"msg": retmsg,
			"code": GetReturnCode(retmsg),
		})
	}
}

func queryencode(c *gin.Context, retmsg string, banlance int64) {
	c.JSON(http.StatusOK, gin.H{
		"balance": banlance,
		"msg": retmsg,
		"code": GetReturnCode(retmsg),
	})
}

func cancellationcode(c *gin.Context, retmsg string, username string) {
	c.JSON(http.StatusOK, gin.H{
		"Username": username,
		"msg": retmsg,
		"code": GetReturnCode(retmsg),
	})
}

func errencode(c *gin.Context, retmsg string){
		c.JSON(http.StatusOK, gin.H{
			"msg": retmsg,
			"code": GetReturnCode(retmsg),
		})
}

//提现encode---返回客户端消息
func withdrawencode(c *gin.Context, code withdrawcode, retmsg string) {
	if WITHDRAW_SUCCESS == code.sub_code {
		c.JSON(http.StatusOK, gin.H{
			"msg":  retmsg,
			"code": GetReturnCode(retmsg),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"sub_msg":  code.sub_msg,
			"sub_code": code.sub_code,
			"msg":      retmsg,
			"code":     GetReturnCode(retmsg),
		})
	}

}
