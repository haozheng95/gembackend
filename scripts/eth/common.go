package eth

import (
	"fmt"
	"github.com/gembackend/conf"
	"github.com/gembackend/gembackendlog"
	"github.com/gembackend/models/eth_query"
	"github.com/regcostajr/go-web3"
	"github.com/regcostajr/go-web3/providers"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type Updater interface {
	Forever()
}

const (
	_TRANSFER          = "0xa9059cbb"
	_TRANSFER_FROM     = "0x23b872dd"
	_TRANSACTION_TOPIC = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	_tokenName         = "0x06fdde03"
	_tokenSymbol       = "0x95d89b41"
	_tokenDecimals     = "0x313ce567"
	_tokenBalance      = "0x70a08231000000000000000000000000"
	_tag               = "latest"
)

var (
	log = gembackendlog.Log
	wg  sync.WaitGroup
)

var (
	eth_query_tx            *eth_query.Tx
	eth_query_token_tx      *eth_query.TokenTx
	eth_query_address_token *eth_query.TokenAddress
	eth_query_address       *eth_query.Address
	connection_web3         *web3.Web3
	tokenmap                map[string]string
)

func init() {
	eth_query_tx = new(eth_query.Tx)
	eth_query_token_tx = new(eth_query.TokenTx)
	eth_query_address_token = new(eth_query.TokenAddress)
	eth_query_address = new(eth_query.Address)

	url := conf.EthRpcHost + ":" + conf.EthRpcPort
	timeOut := conf.EthRpcTimeOut
	source := conf.EthRpcSecure
	connection_web3 = web3.NewWeb3(providers.NewHTTPProvider(url, int32(timeOut), source))
	initTokenmap()
}

func Goid() int {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic recover:panic info:%v", err)
		}
	}()

	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
