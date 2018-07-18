package conf

import (
	"github.com/astaxie/beego/config"
	"runtime"
	"path/filepath"
)

const (
	configType = "ini"
	configPath = "config.ini"
)

var (
	iniconf         config.Configer
	currentfilepath string
	configfiledir   string
	configfilename  string
	//mysql config
	MysqlHost   string
	MysqlPasswd string
	MysqlPost   string
	MysqlUser   string
	//eth rpc config
	EthRpcHost    string
	EthRpcPort    string
	EthRpcTimeOut int
	EthRpcSecure  bool
)

func init() {
	currentfilepath = GetCurrentFilepath()
	configfiledir = filepath.Dir(currentfilepath)
	configfilename = filepath.Join(configfiledir, configPath)
	iniconf, _ = config.NewConfig(configType, configfilename)
	// init mysql config
	MysqlHost = iniconf.String("mysql::host")
	MysqlPasswd = iniconf.String("mysql::passwd")
	MysqlPost = iniconf.String("mysql::post")
	MysqlUser = iniconf.String("mysql::user")
	// init eth rpc config
	EthRpcHost = iniconf.String("eth_rpc::host")
	EthRpcPort = iniconf.String("eth_rpc::port")
	EthRpcTimeOut, _ = iniconf.Int("eth_rpc::timeOut")
	EthRpcSecure, _ = iniconf.Bool("eth_rpc::secure")
}


func GetCurrentFilepath() (filename string) {
	_, filename, _, _ = runtime.Caller(0)
	return
}
