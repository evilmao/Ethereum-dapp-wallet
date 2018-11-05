package ftp

import (
	"fmt"
	"testing"
	"os"
	"strconv"
	"io"
	"strings"
	"time"
)

func TestAboutLogCreat(*testing.T) {
	var err error
	sftpClient, err = ConnectSftp()
	if err != nil {
		fmt.Printf("Connection to sftp server failed : %s\n", err)
	}
}

func TestDirCreat(*testing.T) {
	DirCreat("/ftptest")
	DirCreat("/ftptest/dir")
	os.Create("/ftptest/dir/1.txt")
	dstFile, err := os.Create("/ftptest/testfile.txt")
	if err != nil {
		fmt.Printf("Dst error: %s , path: /ftptest/testfile.txt \n", err)
	}
	defer dstFile.Close()
}

func TestUploadDirectory(*testing.T) {

	localFilePath := "/ftptest"
	remoteFilePath := "/log/opt_log"

	UploadDirectory(localFilePath, remoteFilePath)
	fmt.Println("TestUploadDirectory complete.")
}

func TestDownloadDirectory(*testing.T) {
	localFilePath := "/ftptest"
	remoteFilePath := "/log/opt_log"

	DownloadDirectory(localFilePath, remoteFilePath)
	fmt.Println("TestDownloadDirectory complete.")
}

func TestRemoveWhatever(*testing.T) {
	remoteFilePath := "/log/opt_log"
	fileorDirName := "ftptest"

	RemoveWhatever(fileorDirName, remoteFilePath)
	fmt.Println("TestRemoveWhatever complete.")
}

//func TestGetTmpDirSize(*testing.T) {
//	DirCreat("/ftptest")
//	dstFile, err := os.Create("/ftptest/testfile.txt")
//	if err != nil {
//		fmt.Printf("Dst error: %s , path: /ftptest/testfile.txt \n", err)
//	}
//	//文件大小实时监控和压缩
//	GetFileSize("/ftptest/testfile.txt", int64(1*1024), "/log/log_zip_bak") //超过大小时压缩到指定目录 OK
//	fmt.Printf("TestGetFileSize complete.\n")
//	defer dstFile.Close()
//}

func TestWriteToFile(*testing.T) {
	DirCreat("/ftptest")
	DirCreat("/ftptest/compressdir")

	//文件夹压缩
	dstFile, err := os.Create("/ftptest/compressdir/compressTest.txt")
	files := []*os.File{dstFile}
	//tm = GetTimestamp()
	dest := strings.Join([]string{"compressTest.txt" + "_" + GetTimestamp() + "_bak.zip"}, "")
	fmt.Printf("dest :%s\n", dest)
	//以时间戳命名压缩文件夹并存储到目的压缩路径，如：日志压缩包备份目录/log/log_zip_bak
	err = Compress(files, "/log/log_zip_bak", dest)
	if err != nil {
		fmt.Printf("Compress failed, error :%s\n", err)
	}

	//err = DeCompress(dest, "/log/log_zip_bak")
	//if err != nil {
	//	fmt.Printf("DeCompress failed, error :%s\n", err)
	//}

	dstFile, err = os.Create("/ftptest/testfile.txt")
	if err != nil {
		fmt.Printf("Dst error: %s , path: /ftptest/testfile.txt \n", err)
	}
	for i:=0; i<100; i++ {
		wireteString := "文件写入1234567890~!@#$%^&*()_+{}|[]qwertyuiopasdfghjklzxcvbnm\n----------" + strconv.Itoa(i) +"----------\n"
		_, err := io.WriteString(dstFile, wireteString) //写入文件(字符串)
		if err != nil {
			fmt.Printf("TestWriteToFile failed :%s\n", err)
		}
	}

	_,size := getFileByteSize("/ftptest/testfile.txt")
	fmt.Printf("/ftptest/testfile.txt size :%d\n", size)
	//日志大小超过最大值时压缩当前文件
	if size >=  int64(1*1024) {
		f, err := os.Open("/ftptest/testfile.txt")
		defer f.Close()
		if err != nil {
			fmt.Printf("Open /ftptest/testfile.txt failed, error: %s\n", err)
		}
		files := []*os.File{f}
		tm := GetTimestamp()
		dest := strings.Join([]string{"testfile" + "_" + tm + "_bak.zip"}, "")
		fmt.Printf("dest :%s\n", dest)

		////以时间戳命名压缩文件夹并存储到目的压缩路径，如：日志压缩包备份目录/log/log_zip_bak
		//err = Compress(files, "/log/log_zip_bak", dest)
		//if err != nil {
		//	fmt.Printf("Compress failed, error :%s\n", err)
		//}
		//
		//files = []*os.File{f}
		////tm = GetTimestamp()
		//dest = strings.Join([]string{"testfile" + "_" + tm + "_bak.zip"}, "")
		//fmt.Printf("dest :%s\n\n", dest)

		//以时间戳命名压缩文件并存储到目的压缩路径，如：日志压缩包备份目录/log/log_zip_bak
		err = Compress(files, "/log/log_zip_bak", dest)
		if err != nil {
			fmt.Printf("Compress failed, error :%s\n", err)
		}
		//可选：删除超限操作日志文件 || 清空超限操作日志文件（清空后可复用，文件依旧写不用再创建）
		//删除超限操作日志文件
		res := EmptiedFile("/ftptest/testfile.txt")
		if !res {
			fmt.Printf("EmptiedFile failed\n")
		}
		//time.Sleep(time.Duration(1)*time.Second)
		//RemoveFile("/ftptest/testfile.txt")
	}
}


func TestRemoveTestAbout(*testing.T) {
	time.Sleep(time.Duration(2)*time.Second)

	err := os.RemoveAll("/ftptest")
	if nil != err {
		fmt.Printf("RemoveAll /ftptest err :%s\n", err)
	}

	err = os.RemoveAll("/log")
	if nil != err {
		fmt.Printf("RemoveAll /log err :%s\n", err)
	}
}


