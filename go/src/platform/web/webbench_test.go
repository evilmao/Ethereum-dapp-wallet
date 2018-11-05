/*****************************************************************************
File name: server.go
Description: web server
Author:failymao
Version: V1.0
Date: 2018/06/14
History:
*****************************************************************************/
package web

import (
	"net/http"
	"strings"
	"testing"
)


func Benchmark_Post_1(b *testing.B) {
	for i := 0; i < b.N; i++ { //use b.N for looping
		data :=`{"SourceUser":"Haxima006","DestUser":"test001","TransferAmount":"10"}`
		request, _ := http.NewRequest("POST", "http://192.168.0.191:10086/wallet?method=TransferAccounts", strings.NewReader(data))
		http.DefaultClient.Do(request)

	}

}
