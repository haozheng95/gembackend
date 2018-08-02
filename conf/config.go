package conf

import (
	"github.com/astaxie/beego/config"
	"path/filepath"
	"runtime"
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
	MysqlPort   string
	MysqlUser   string
	//eth rpc config
	EthRpcHost       string
	EthRpcPort       string
	EthWebsocketPort string
	EthRpcTimeOut    int
	EthRpcSecure     bool
	EthWebsocketUrl  string
	//kafka config
	KafkaHost                   string
	KafkaPort                   string
	KafkaimportEthTopicName     string
	KafkatransactionParityTopic string
	KafkagetbalanceParityTopic  string
)

func init() {
	currentfilepath = GetCurrentFilepath()
	configfiledir = filepath.Dir(currentfilepath)
	configfilename = filepath.Join(configfiledir, configPath)
	iniconf, _ = config.NewConfig(configType, configfilename)
	// init mysql config
	MysqlHost = iniconf.String("mysql::host")
	MysqlPasswd = iniconf.String("mysql::passwd")
	MysqlPort = iniconf.String("mysql::port")
	MysqlUser = iniconf.String("mysql::user")
	// init eth rpc config
	EthRpcHost = iniconf.String("eth_rpc::host")
	EthRpcPort = iniconf.String("eth_rpc::port")
	EthWebsocketPort = iniconf.String("eth_rpc::websocket")
	EthWebsocketUrl = EthRpcHost + ":" + EthWebsocketPort
	EthRpcTimeOut, _ = iniconf.Int("eth_rpc::timeOut")
	EthRpcSecure, _ = iniconf.Bool("eth_rpc::secure")
	// init kafka config
	KafkaHost = iniconf.String("kafka::host")
	KafkaPort = iniconf.String("kafka::port")
	KafkaimportEthTopicName = iniconf.String("kafka::importEthTopic")
	KafkatransactionParityTopic = iniconf.String("kafka::transactionParityTopic")
	KafkagetbalanceParityTopic = iniconf.String("kafka::getbalanceParityTopic")
}

func GetCurrentFilepath() (filename string) {
	_, filename, _, _ = runtime.Caller(0)
	return
}
