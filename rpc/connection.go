package rpc

import (
	"github.com/gembackend/conf"
	"github.com/regcostajr/go-web3/providers"
	"github.com/regcostajr/go-web3"
)

// Conn pool

var ConnectMap = make(map[string]interface{})
// main ------
func MakeConn() {

	if len(ConnectMap) == 0 {
		log.Debug("add connect -----")
		addConnect()
	}
}

// ----------------------------- public
func addConnect() {
	ConnectMap["eth-web3"] = makeEthConn()
	ConnectMap["eth-web3-original"] = makeEthConnOriginal()
}

func ReMakeAllConn() {
	addConnect()
}

// --------------------------- retry
func ReMakeWeb3Conn() (conn *Web3) {
	conn = makeEthConn()
	ConnectMap["eth-web3"] = conn
	return
}
func ReMakeWeb3ConnOriginal() (conn *web3.Web3) {
	conn = makeEthConnOriginal()
	ConnectMap["eth-web3-original"] = conn
	return
}

// --------------------------  connect
func makeEthConn() (connection *Web3) {
	timeOut := conf.EthRpcTimeOut
	source := conf.EthRpcSecure
	url := conf.EthRpcHost + ":" + conf.EthRpcPort
	connection = NewWeb3(providers.NewHTTPProvider(url, int32(timeOut), source))
	return
}

func makeEthConnOriginal() (connection *web3.Web3) {
	timeOut := conf.EthRpcTimeOut
	source := conf.EthRpcSecure
	url := conf.EthRpcHost + ":" + conf.EthRpcPort
	connection = web3.NewWeb3(providers.NewHTTPProvider(url, int32(timeOut), source))
	return
}
