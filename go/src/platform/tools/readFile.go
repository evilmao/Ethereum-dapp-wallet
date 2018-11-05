/*****************************************************************************
File name: readFile.go
Description:读取文件
	1.据文件路径名称，读取文件内容，并将文件内容转换成string型返回
Author: failymao
Version: V1.0
Date: 2018/06/27
History:
*****************************************************************************/
package tools

import (
	"io/ioutil"
	"os"
	"platform/om/log"
)

func ReadFile(filename string) (string, error) {
	// 如果文件不存在，则返回错误
	_, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err){
		log.Error("File does not exist.", err)
		return "", err
	}
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error("fail to read file:", err)
		return "", err
	}
	return string(buf), err
}
