/**************************************************************
【提交类型】:BUG/新功能: 账单下载
【问题描述】:ftp模块代码
【修改内容】:初次提交
【提交人】:failymao
【评审人】:
***************************************************************/

package ftp

import (
	"os"
	"path"
	"io/ioutil"
	"fmt"
	"time"
	"net"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"platform/config"
	"strings"
	"archive/zip"
	"io"
//	"github.com/howeyc/fsnotify"
)

var sftpClient   *sftp.Client

//初始化连接服务器
func init() {
	var err error
	sftpClient ,err = ConnectSftp()
	if err != nil {
		fmt.Printf("Connection to sftp server failed : %s\n", err)
	}
	//defer sftpClient.Close()
	AboutLogCreat()
}


//sftp连接
func ConnectSftp() (*sftp.Client, error) {

	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)

	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(config.Gconfig.Ftpcfg.Password))
	//auth := []ssh.AuthMethod{ssh.Password(config.Gconfig.Ftpcfg.Password)} [0x553a10]

	clientConfig = &ssh.ClientConfig{
		User:    config.Gconfig.Ftpcfg.Username,
		Auth:    auth,
		Timeout: 30 * time.Second,
		//默认密钥不受信任时，断掉连接，无需操作return nil
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	//连接ssh
	addr = fmt.Sprintf("%s:%d", config.Gconfig.Ftpcfg.Hostname, config.Gconfig.Ftpcfg.Port)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		fmt.Printf("sshClient err :%s\n",err)
		return nil, err
	}

    //连接sftp
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		fmt.Printf("sftpClient err :%s\n",err)
		return nil, err
	}

	return sftpClient, nil
}


//判断目录是否存在
func DirJudge(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}


//创建日志相关目录
func AboutLogCreat() {
	LogDirCreat()
	OptlogDirCreat()
	LogzipDirCreat()
}


//判断日志文件存储目录/log是否存在，不存在则创建
func LogDirCreat() (err error) {
	exist, err := DirJudge(config.Gconfig.Ftpcfg.Logpath)
	if err != nil {
		fmt.Printf("Get dir error : %s\n", err)
		return err
	}
	if exist {
		fmt.Printf("Has dir /log.\n")
	} else {
		fmt.Printf("No dir /log.\n")
		// 创建文件夹
		err := os.Mkdir(config.Gconfig.Ftpcfg.Logpath, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir /log failed :%s\n", err)
			return err
		} else {
			fmt.Printf("mkdir /log success!\n")
		}
	}
	return nil
}

//判断操作日志存储路径/log/opt_log是否存在，不存在则创建
func OptlogDirCreat() (err error) {
	exist, err := DirJudge(config.Gconfig.Ftpcfg.Optlogpathpath)
	if err != nil {
		fmt.Printf("Get dir error : %s\n", err)
		return err
	}
	if exist {
		fmt.Printf("Has dir /log/opt_log.\n")
	} else {
		fmt.Printf("No dir /log/opt_log.\n")
		// 创建文件夹
		err := os.Mkdir(config.Gconfig.Ftpcfg.Optlogpathpath, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir /log/opt_log failed :%s\n", err)
			return err
		} else {
			fmt.Printf("mkdir /log/opt_log success!\n")
		}
	}
	return nil
}


//判断日志压缩包备份目录/log/log_zip_bak是否存在，不存在则创建
func LogzipDirCreat() (err error) {
	exist, err := DirJudge(config.Gconfig.Ftpcfg.Logzippath)
	if err != nil {
		fmt.Printf("Get dir error : %s\n", err)
		return err
	}
	if exist {
		fmt.Printf("Has dir /log/log_zip_bak.\n")
	} else {
		fmt.Printf("No dir /log/log_zip_bak.\n")
		// 创建文件夹
		err := os.Mkdir(config.Gconfig.Ftpcfg.Logzippath, os.ModePerm)
		//err := os.Mkdir("/ftp", os.ModePerm)
		//err := os.Mkdir("/ftp/zip_bak", os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir /log/log_zip_bak failed :%s\n", err)
			return err
		} else {
			fmt.Printf("mkdir /log/log_zip_bak success!\n")
		}
	}
	return nil
}



//判断目录是否存在，不存在则创建
func DirCreat(dirpath string) (err error) {
	exist, err := DirJudge(dirpath)
	if err != nil {
		fmt.Printf("Get dir error : %s\n", err)
		return err
	}
	if exist {
		fmt.Printf("Has dir %s.\n", dirpath)
	} else {
		fmt.Printf("No dir %s.\n", dirpath)
		// 创建文件夹
		err := os.Mkdir(dirpath, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir %s failed :%s\n", dirpath, err)
			return err
		} else {
			fmt.Printf("mkdir /log/log_zip_bak success!\n")
		}
	}
	return nil
}


//判断是否是文件夹
func isDir(dirname string) bool  {
	fhandler, err := os.Stat(dirname)
	if ! (err == nil || os.IsExist(err))  {
		return false
	}else {
		return fhandler.IsDir()
	}
}

//判断是否是文件
func isFile(filename string) bool  {
	fhandler, err := os.Stat(filename)
	if ! (err == nil || os.IsExist(err))  {
		return false
	}else if fhandler.IsDir() {
		return false
	}
	return true
}


//获取文件大小
func getFileByteSize(filename string) (bool,int64) {
	if  ! isFile(filename) {
		return false,0
	}
	fhandler, _ := os.Stat(filename)
	return true, fhandler.Size()
}


//获取字符串形式的连续时间戳
func GetTimestamp()(string) {
	//获取当前时间
	current := time.Now().Format("2006-01-02 15:04:05")
	tm := strings.Replace(current, " ", "", -1)
	tm = strings.Replace(tm, "-", "", -1)
	tm = strings.Replace(tm, ":", "", -1)
	return tm
}


//监控操作日志文件大小
//func GetFileSize(filepath string, filemaxsize int64, destDir string) {
//
//	var filename = path.Base(filepath)
//
//	watcher, err := fsnotify.NewWatcher()
//	if err != nil {
//		log.Fatal(err)
//	}
//	done := make(chan bool)
//
//	go func() {
//		for {
//			select {
//			case ev := <-watcher.Event:
//				if ev.IsModify() {
//					_,size := getFileByteSize(ev.Name)
//					log.Println("event:",ev,",byte:",size)
//					//日志大小超过最大值时压缩当前文件
//					if size >= filemaxsize {
//						f, err := os.Open(filepath)
//						defer f.Close()
//						if err != nil {
//							fmt.Printf("Open %s failed, error: %s \n", filepath, err)
//						}
//						files := []*os.File{f}
//						tm := GetTimestamp()
//						dest := strings.Join([]string{filename + "_" + tm + "_bak.zip"}, "")
//						fmt.Printf("dest :%s\n", dest)
//
//						//以时间戳命名压缩文件并存储到目的压缩路径，如：日志压缩包备份目录/log/log_zip_bak
//						err = Compress(files, destDir, dest)
//						if err != nil {
//							fmt.Printf("Compress failed, error :%s\n", err)
//						}
//						//可选：删除超限操作日志文件 || 清空超限操作日志文件（清空后可复用，文件依旧写不用再创建）
//						//删除超限操作日志文件
//						RemoveFile(ev.Name)
//						fmt.Printf("RemoveFile :%s\n", ev.Name)
//
//						/*
//                        //清空超限操作日志文件
//						if  ! EmptiedFile(ev.Name)  {
//							fmt.Println("EmptiedFile failed, error\n", err)
//						}
//						*/
//					}
//				}
//			case err := <-watcher.Error:
//				log.Printf("error: %s\n",err)
//				fmt.Printf("watcher error :%s\n", err)
//			}
//		}
//	}()
//	//err = watcher.Watch(config.Gconfig.Ftpcfg.Optlogpathpath)
//	err = watcher.Watch(filepath)
//	if err != nil {
//		log.Fatal(err)
//	}
//	<-done
//
//	watcher.Close()
//}


//清空文件
func EmptiedFile(filename string) bool  {
	fmt.Printf("EmptiedFile filename :%s\n", filename)
	fileName,err := os.Create(filename)
	defer fileName.Close()
	if err != nil {
		return false
	}
	fmt.Fprint(fileName,"")
	return true
}


//将超过规定大小的操作日志文件压缩到/log/log_zip_bak
func Compress(files []*os.File, destDir string, dest string) error {
	d, _ := os.Create(path.Join(destDir, dest))
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {
		err := compress(file, "", w)
		if err != nil {
			return err
		}
	}
	return nil
}

//压缩文件
func compress(file *os.File, prefix string, zipw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, zipw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zipw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

/*
//解压缩
func DeCompress(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		filename := dest + file.Name
		err = os.MkdirAll(getDir(filename), 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		w.Close()
		rc.Close()
	}
	return nil
}


//获取文件夹名
func getDir(path string) string {
	return subString(path, 0, strings.LastIndex(path, "/"))
}


//提取路径中文件夹名
func subString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("Start is wrong")
	}

	if end < start || end > length {
		panic("End is wrong")
	}

	return string(rs[start:end])
}
*/

//文件上传
func UploadFile(localFilePath string, remoteDir string) {
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		fmt.Printf("UploadFile failed : %s\n", err)
		return
	}

	defer srcFile.Close()

	//dstFile, err := sftpClient.Create(path.Join(remoteDir, remoteFileName))
	dstFile, err := sftpClient.Create(remoteDir)
	if err != nil {
		fmt.Printf("Dst error: %s , path: %s \n", err, remoteDir)
	}

	defer dstFile.Close()

	buf, err := ioutil.ReadAll(srcFile)

	if err != nil {
		panic(err)
	}

	dstFile.Write(buf)
	fmt.Printf("Uploaded file: %s\n", localFilePath)
}


//文件夹上传
func UploadDirectory(localPath string, remotePath string) {
	localFiles, err := ioutil.ReadDir(localPath)

	if err != nil {
		panic(err)
	}

	var remoteDirName = path.Base(localPath)

	sftpClient.Mkdir(path.Join(remotePath, remoteDirName))
	fmt.Printf("Make remote dir: %s\n", path.Join(remotePath, remoteDirName))

	for _, backupDir := range localFiles {

		localFilePath := path.Join(localPath, backupDir.Name())
		remoteDirPath := path.Join(remotePath, remoteDirName)
		remoteFilePath := path.Join(remotePath, remoteDirName, backupDir.Name())

        //文件夹上传：遍历文件夹，上传文件夹--遍历至文件夹最后一层上传文件
		if backupDir.IsDir() {

			exist,err := DirJudge(remoteDirPath)
			if err != nil {
				fmt.Printf("Get dir error : %s\n", err)
				return
			}
			if ! exist {
				sftpClient.Mkdir(remoteDirPath)
				if err != nil {
					fmt.Printf("mkdir %s failed :%s\n", remoteDirPath, err)
					return
				}
			}

			exist,err = DirJudge(remoteFilePath)
			if err != nil {
				fmt.Printf("Get dir error : %s\n", err)
				return
			}
			if ! exist {
				sftpClient.Mkdir(remoteFilePath)
				if err != nil {
					fmt.Printf("mkdir %s failed :%s\n", remoteFilePath, err)
					return
				}
			}

			UploadDirectory(localFilePath, path.Join(remotePath, remoteDirName))
		} else {

			UploadFile(path.Join(localPath, backupDir.Name()), remoteFilePath)
		}
	}
}

//删除文件或文件夹
func RemoveWhatever(fileorDirName string, remoteDir string) {
	whateverName := path.Base(fileorDirName)
	remoteFileName := path.Join(remoteDir, whateverName)
	fileInfo, err := sftpClient.Stat(remoteFileName)

	a := fileInfo.Name()
	fmt.Printf(a)

	if err != nil {
		fmt.Printf("RemoveWhatever failed : %s\n", err)
		return
	}

	if fileInfo.IsDir() {
		RemoveDirectory(remoteFileName, remoteDir)      //删除文件夹
	} else {
		RemoveFile(remoteFileName)                      //删除文件
	}
}

//删除文件
func RemoveFile(filename string) {
	err := sftpClient.Remove(filename)
	if err != nil {
		fmt.Printf("Can not remove file %s ：%s\n", filename, err)
	} else {
		fmt.Printf("RemoveFile success! \n")
	}
}

//删除文件夹
func RemoveDirectory(dir string, remoteDir string) {
	fmt.Printf("--------RemoveDirectory-----------\n")
	fileInfos, err := sftpClient.ReadDir(dir)

	if err != nil {
		fmt.Printf("ReadDir err: %s", err)
	}

	for _, fileInfo := range fileInfos {

		if fileInfo.IsDir() {
			fmt.Println("Remove directory: ", path.Join(dir, fileInfo.Name()))
			//RemoveDirectory(path.Join(dir, fileInfo.Name()), remoteDir)
			RemoveDirectory( path.Join(dir, fileInfo.Name()), remoteDir)
		} else {
			RemoveFile(path.Join(dir, fileInfo.Name()))
		}
	}

	sftpClient.RemoveDirectory(dir)
}


//下载文件
func DownloadFile(localDir string, remoteFilePath string) {
	srcFile, err := sftpClient.Open(remoteFilePath)
	if err != nil {
		fmt.Printf("DownloadFile failed : %s\n", err)
	}
	defer srcFile.Close()

	var localFileName = path.Base(remoteFilePath)

	dstFile, err := os.Create(localDir)
	if err != nil {
		fmt.Printf("Dst error: %s , path: %s \n", err, path.Join(localDir, localFileName))
	}
	defer dstFile.Close()

	if _, err = srcFile.WriteTo(dstFile); err != nil {
		fmt.Printf("DownloadFile error: %s\n", err)
	}

	fmt.Printf("Copy file from %s finished! \n", remoteFilePath)
}

//下载文件夹
func DownloadDirectory(localDir string, remoteFilePath string) {

	fileInfos, err := sftpClient.ReadDir(remoteFilePath)

	if err != nil {
		fmt.Printf("DownloadDirectory err :%s\n", err)
		panic(err)
	}

	var localDirName = path.Base(remoteFilePath)

	for _, fileInfo := range fileInfos {

		localDirPath := path.Join(localDir, localDirName)
		localFilePath := path.Join(localDir, localDirName, fileInfo.Name())

		if fileInfo.IsDir() {

			exist,err := DirJudge(localDirPath)
			if err != nil {
				fmt.Printf("Get dir error : %s\n", err)
				return
			}
			if ! exist {
				err := os.Mkdir(localDirPath, os.ModePerm)
				if err != nil {
					fmt.Printf("mkdir %s failed :%s\n", localDirPath, err)
					return
				}
			}

			exist,err = DirJudge(localFilePath)
			if err != nil {
				fmt.Printf("Get dir error : %s\n", err)
				return
			}
			if ! exist {
				err = os.Mkdir(localFilePath, os.ModePerm)
				if err != nil {
					fmt.Printf("mkdir %s failed :%s\n", localFilePath, err)
					return
				}
			}

			DownloadDirectory(localDirPath, path.Join(remoteFilePath, fileInfo.Name()))
		} else {
			DownloadFile(path.Join(localDir, localDirName, fileInfo.Name()), path.Join(remoteFilePath, fileInfo.Name()))
		}
	}
}

//断点续传