package config

import (
	"os"
	"github.com/BurntSushi/toml"
	"log"
	"runtime"
)
var MAXCPUNUM int = 0

func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	MAXCPUNUM = nuCPU
}

var Gconfig *Config = new(Config)

func init() {
	gopath := os.Getenv("GOPATH")
	if "" == gopath {

		log.Fatal("env GOPATH not set:[%s]", gopath)
	}

	gopath = gopath + string(os.PathSeparator) + "src" +
		string(os.PathSeparator) + "platform" +
		string(os.PathSeparator) + "config" +
		string(os.PathSeparator) + "config.toml"

	if _, err := toml.DecodeFile(gopath, Gconfig); err != nil {
		log.Fatal("Parse config file:[%s]", err.Error())
	}
}

type Config struct {
	UpdateInt	int		 `toml:"updateinterval"`
	UsrPrepared int 	 `toml:"preparecount"`
	Mysqlcfg 	MysqlCfg `toml:"mysql"`
	Ethcfg		EthCfg	 `toml:"eth"`
	Servercfg	ServerCfg`toml:"server"`
	Logcfg		LogCfg   `toml:"log"`
	Billcfg		BillCfg  `toml:"bill"`
	Httpscfg    Httpscfg `toml:"https"` //https证书映射
	Ftpcfg      FtpCfg   `toml:"ftp"`   //ftp连接映射
}

type MysqlCfg struct {
	Username	string `toml:"username"`
	Password	string `toml:"password"`
	Charset		string `toml:"charset"`
	Database	string `toml:"database"`
	Port		int	   `toml:"port"`
	Hostname	string `toml:"hostname"`
}

type EthCfg struct {
	RPCUrl 				string `toml:"rpcurl"`
	ContractAccount		string `toml:"contractAccount"`
	Keystoredir 		string `toml:"keystoredir"`
	EnterpriseAccount 	string `toml:"enterpriseAccount"`
	EnterPrisePasswd	string `toml:"enterPrisePasswd"`
	EnterPriseKeypath	string `toml:"enterPriseKey"`
	Transflag	        bool   `toml:"transflag"`
	GasPrice			int64  `toml:"gasprice"`
	MultiplierPrice     int64   `toml:"multiplierprice"` //gasPrice增加倍数
	TimeOut		     	int64  `toml:"timeout"`//交易超时时间分钟
	TimeoutTimes		int64  `toml:"timeouttimes"`//交易重复发送次数
	ConfirmeTimes		int64  `toml:"confirmetimes"` //#交易确认次数

}

type ServerCfg struct {
	ReadTimeout		    int `toml:"readtimeout"`
	WriteTimeout		int `toml:"writetimeout"`
	Port				int `toml:"port"`
	MaxConn				int `toml:"maxconn"`
}

type LogCfg struct {
	LogLev 		string	 `toml:"loglev"`
	LogMax		int		 `toml:"logmax"`
	LogPath		string	 `toml:"logpath"`
}

type BillCfg struct {
	BillPath 	string `toml:"path"`
	BillMax		int	   `toml:"maxsize"`
}

//https证书路径配置
type Httpscfg struct {
	Keypath  string `toml:"sslkeypath"`
	Certpath string `toml:"sslcertpath"`
}

//ftp服务器配置
type FtpCfg struct {
    Username        string   `toml:"username"`
    Password        string 	 `toml:"password"`
    Hostname        string   `toml:"hostname"`
    Port            int	     `toml:"port"`
    Logpath         string	 `toml:"logpath"`
    Optlogpathpath  string	 `toml:"optlogpath"`
    Logzippath      string	 `toml:"logzippath"`
}
