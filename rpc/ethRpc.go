package rpc

import (
	"net/http"
	"io/ioutil"
	"github.com/gembackend/gembackendlog"
	"fmt"
	"encoding/json"
	"strings"
	"bytes"
	"github.com/gembackend/conf"
)

var log = gembackendlog.Log

// Ethereum Developer APIs
const (
	_API_KEY                                = "E7BJ9TNPC31ZAIT61K8ZDJM9HRZXV3TWMM"
	eth_blockNumber                         = "https://api.etherscan.io/api?module=proxy&action=eth_blockNumber&apikey=" + _API_KEY
	eth_getBlockByNumber                    = "https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&tag=%s&boolean=true&apikey=" + _API_KEY
	eth_getUncleByBlockNumberAndIndex       = "https://api.etherscan.io/api?module=proxy&action=eth_getUncleByBlockNumberAndIndex&tag=%s&index=%s&apikey=" + _API_KEY
	eth_getBlockTransactionCountByNumber    = "https://api.etherscan.io/api?module=proxy&action=eth_getBlockTransactionCountByNumber&tag=%s&apikey=" + _API_KEY
	eth_getTransactionByHash                = "https://api.etherscan.io/api?module=proxy&action=eth_getTransactionByHash&txhash=%s&apikey=" + _API_KEY
	eth_getTransactionByBlockNumberAndIndex = "https://api.etherscan.io/api?module=proxy&action=eth_getTransactionByBlockNumberAndIndex&tag=%s&index=%s&apikey=" + _API_KEY
	eth_getTransactionCount                 = "https://api.etherscan.io/api?module=proxy&action=eth_getTransactionCount&address=%s&tag=latest&apikey=" + _API_KEY
	eth_sendRawTransaction                  = "https://api.etherscan.io/api?module=proxy&action=eth_sendRawTransaction&hex=%s&apikey=" + _API_KEY
	eth_getTransactionReceipt               = "https://api.etherscan.io/api?module=proxy&action=eth_getTransactionReceipt&txhash=%s&apikey=" + _API_KEY
	eth_call                                = "https://api.etherscan.io/api?module=proxy&action=eth_call&to=%s&data=%s&tag=latest&apikey=" + _API_KEY
	eth_getCode                             = "https://api.etherscan.io/api?module=proxy&action=eth_getCode&address=%s&tag=latest&apikey=" + _API_KEY
	eth_getStorageAt                        = "https://api.etherscan.io/api?module=proxy&action=eth_getStorageAt&address=%s&position=%s&tag=latest&apikey=" + _API_KEY
	eth_gasPrice                            = "https://api.etherscan.io/api?module=proxy&action=eth_gasPrice&apikey=" + _API_KEY
	eth_estimateGas                         = "https://api.etherscan.io/api?module=proxy&action=eth_estimateGas&to=%s&value=0xff22&gasPrice=%s&gas=%s&apikey=" + _API_KEY
	eth_getMulBalance                       = "https://api.etherscan.io/api?module=account&action=balancemulti&address=%s&tag=latest&apikey=" + _API_KEY
	eth_getBalance                          = "https://api.etherscan.io/api?module=account&action=balance&address=%s&tag=latest&apikey=" + _API_KEY
)

// 获取token信息
const (
	_tokenName     = "0x06fdde03"
	_tokenSymbol   = "0x95d89b41"
	_tokenDecimals = "0x313ce567"
	_tokenBalance  = "0x70a08231000000000000000000000000"
)

type Response struct {
	Jsonrpc string
	Id      int
	Result  map[string]interface{}
	Error   map[string]interface{}
}

type TokenResonse struct {
	Jsonrpc string
	Id      int
	Result  string
	Error   map[string]interface{}
}

func test() {
	//s, _ := Eth_getTokenBalance("0xf230b790e05390fc8295f4d3f60332c93bed42e2", "5e5978030b2c74e3fa5aa0ef8a40da3912e00e93")
	//s, _ := Eth_getTokenName("0xf230b790e05390fc8295f4d3f60332c93bed42e2")
	s, _ := Eth_getTokenSymbol("0xf230b790e05390fc8295f4d3f60332c93bed42e2")
	e, _ := FormatTokenResponse(s)
	fmt.Println(DecodeHexResponseToString(e.Result))

}

func Eth_getTokenSymbol(contract string) (string, error) {
	url := fmt.Sprintf(eth_call, contract, _tokenSymbol)
	log.Debug(url)
	s, e := httpGet(url)
	if e != nil {
		log.Error(e)
	}
	log.Info(s)
	return s, e
}

func Eth_getTokenDecimals(contract string) (string, error) {
	url := fmt.Sprintf(eth_call, contract, _tokenDecimals)
	log.Debug(url)
	s, e := httpGet(url)
	if e != nil {
		log.Error(e)
	}
	log.Info(s)
	return s, e

}

func Eth_getBalance(s string) string {
	url := fmt.Sprintf(eth_getBalance, s)
	log.Debug(url)
	str, e := httpGet(url)
	if e != nil {
		log.Error(e)
	}
	log.Info(str)
	return str
}

// 获取多个eth地址的余额
func Eth_getMulBalance(s ... string) string {
	s1 := strings.Join(s, ",")
	url := fmt.Sprintf(eth_getMulBalance, s1)
	log.Debug(url)
	str, e := httpGet(url)
	if e != nil {
		log.Error(e)
	}
	log.Info(str)
	return str
}

func Eth_getTokenName(contract string) (string, error) {
	url := fmt.Sprintf(eth_call, contract, _tokenName)
	log.Debug(url)
	s, e := httpGet(url)
	if e != nil {
		log.Error(e)
	}
	log.Info(s)
	return s, e
}

// 获取用户的eoken余额
func Eth_getTokenBalance(contract, addr string) (string, error) {
	url := fmt.Sprintf(eth_call, contract, _tokenBalance+addr)
	log.Debug(url)
	s, e := httpGet(url)
	if e != nil {
		log.Error(e)
	}
	log.Info(s)
	return s, e
}

func FormatTokenResponse(s string) (*TokenResonse, error) {
	x := []byte(s)
	r := &TokenResonse{}
	err := json.Unmarshal(x, &r)

	if err != nil {
		log.Error(err)
	}

	return r, err
}

func FormatResponse(s *string) (*Response, error) {
	x := []byte(*s)
	r := &Response{}
	err := json.Unmarshal(x, &r)

	if err != nil {
		log.Error(err)
	}

	return r, err
}

func FormatResponseMap(s *string) (map[string]interface{}, error) {
	x := []byte(*s)
	r := make(map[string]interface{})
	err := json.Unmarshal(x, &r)

	if err != nil {
		log.Error(err)
	}

	return r, err
}

func Eth_getTransactionReceipt(hash string) (string, error) {
	url := fmt.Sprintf(eth_getTransactionReceipt, hash)
	log.Debug(url)
	s, e := httpGet(url)
	if e != nil {
		log.Error(e)
	}
	log.Info(s)
	return s, e

}

func Eth_getTransactionCount(address string) (string, error) {
	url := fmt.Sprintf(eth_getTransactionCount, address)
	log.Debug(url)
	s, e := httpGet(url)
	if e != nil {
		log.Error(e)
	}
	log.Info(s)
	return s, e
}

func Eth_getTransactionByHash(hash string) (string, error) {
	url := fmt.Sprintf(eth_getTransactionByHash, hash)
	log.Debug(url)
	s, e := httpGet(url)
	if e != nil {
		log.Error(e)
	}
	log.Info(s)
	return s, e
}

func Eth_getBlockByNumber(num uint64) (string, error) {
	h := fmt.Sprintf("%x", num)
	url := fmt.Sprintf(eth_getBlockByNumber, h)
	log.Debug(url)
	s, e := httpGet(url)
	if e != nil {
		log.Error(e)
	}
	log.Info(s)
	return s, e
}

func Eth_blockNumber() (string, error) {
	s, e := httpGet(eth_blockNumber)
	if e != nil {
		log.Error(e)
	}
	log.Info(s)
	return s, e
}

func httpGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Error(err)
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return "", err
	}

	return string(body), nil
}

func HttpPost(song map[string]interface{}) string {

	bytesData, err := json.Marshal(song)
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(bytesData)
	url := "http://" + conf.EthRpcHost + conf.EthRpcPort

	resp, err := http.Post(url, "application/json", reader)

	if err != nil {
		log.Error(err)
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Error(err)
		panic(err)
	}

	return string(body)
}
