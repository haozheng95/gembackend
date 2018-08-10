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
	KafkaSendRawTopic           string //ethscan
	KafkaTxRecordTopic          string //ethscan
	//jwt config
	JwtGodKey     string
	JwtExpiration int64
	JwtKey        string
	//init server config
	RunMode string
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
	KafkaSendRawTopic = iniconf.String("kafka::KafkaSendRawTopic")
	KafkaTxRecordTopic = iniconf.String("kafka::KafkaTxRecordTopic") // eth
	// init jwt
	JwtGodKey = iniconf.String("jwt::JwtGodKey")
	JwtKey = iniconf.String("jwt::JwtKey")
	JwtExpiration, _ = iniconf.Int64("jwt::JwtExpiration")
	// init app config
	RunMode = iniconf.String("app::mode")
}

func GetCurrentFilepath() (filename string) {
	_, filename, _, _ = runtime.Caller(0)
	return
}
