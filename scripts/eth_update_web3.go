package scripts

import (
	"github.com/gembackend/models"
	"github.com/regcostajr/go-web3"
	"time"
	"github.com/regcostajr/go-web3/providers"
	"fmt"
	"github.com/gembackend/rpc"
	"strings"
	"github.com/shopspring/decimal"
	"math"
	"github.com/regcostajr/go-web3/dto"
	"strconv"
	"github.com/gembackend/conf"
)

const (
	_tokenName     = "0x06fdde03"
	_tokenSymbol   = "0x95d89b41"
	_tokenDecimals = "0x313ce567"
	_tokenBalance  = "0x70a08231000000000000000000000000"
)

type EthUpdaterWeb3 struct {
	StartHeight       uint64
	TableBlock        *models.Block
	TableTx           *models.Tx
	TableTokenTx      *models.TokenTx
	TableAddress      *models.Address
	TableTokenAddress *models.TokenAddress
	connection        *web3.Web3
	parityParam       map[string]interface{}
	rpcRes            *rpc.Response
}

func (updaterWeb3 *EthUpdaterWeb3) Forever() {
	// todo 对比最新块的高度和数据库中的高度，小于最新高度就更新，否则就等待
	// 获取数据库高度
	updaterWeb3.TableBlock = updaterWeb3.TableBlock.SelectMaxHeight()
	dbHehght := updaterWeb3.TableBlock.BlockHeight
	for height := MaxIntByString(updaterWeb3.StartHeight, dbHehght); ; {
		rpcHeight, err := updaterWeb3.connection.Eth.GetBlockNumber()
		if err != nil {
			log.Errorf("web3 rpcHeight error %s", err)
			return
		}

		//rpc高度大于数据库高度 就更新
		if rpcHeight.Uint64() > height {

			//todo 获取块交易
			hexHeight := fmt.Sprintf("0x%x", height)
			updaterWeb3.parityParam["method"] = "eth_getBlockByNumber"
			updaterWeb3.parityParam["params"] = []interface{}{hexHeight, true}
			blockInfoString := rpc.HttpPost(updaterWeb3.parityParam)
			updaterWeb3.rpcRes, err = rpc.FormatResponse(&blockInfoString)
			if err != nil {
				log.Errorf("FormatResponse error %s", err)
				panic(err)
			}
			//todo 处理块信息
			updaterWeb3.disposeBlockInfo()
			//todo 处理块内的交易信息
			updaterWeb3.disposeTransactions()
			height++
		} else {
			//Have a rest
			time.Sleep(time.Second * 5)
		}

	}
}
func (updaterWeb3 *EthUpdaterWeb3) disposeBlockInfo() {
	//格式化block数据
	result := updaterWeb3.rpcRes.Result
	updaterWeb3.TableBlock.BlockHeight = HexDec(result["number"].(string))
	updaterWeb3.TableBlock.TimeStamp = HexDec(result["timestamp"].(string))
	updaterWeb3.TableBlock.Nonce = HexDec(result["nonce"].(string))
	updaterWeb3.TableBlock.BlockHash = result["hash"].(string)
	updaterWeb3.TableBlock.GasLimit = HexDec(result["gasLimit"].(string))
	updaterWeb3.TableBlock.GasUsed = HexDec(result["gasUsed"].(string))
	updaterWeb3.TableBlock.Size = HexDec(result["size"].(string))
	updaterWeb3.TableBlock.ParentHash = result["parentHash"].(string)
	updaterWeb3.TableBlock.Miner = result["miner"].(string)
	updaterWeb3.TableBlock.MixHash = result["mixHash"].(string)
	updaterWeb3.TableBlock.ExtraData = result["extraData"].(string)
	//插入block库
	updaterWeb3.TableBlock.InsertOneRaw(updaterWeb3.TableBlock)
}
func (updaterWeb3 *EthUpdaterWeb3) disposeTransactions() {
	result, ok := updaterWeb3.rpcRes.Result["transactions"]
	if !ok {
		log.Error("get result transactions error")
		panic("transactions Error")
	}
	transactions := result.([]interface{})
	for _, v := range transactions {
		transaction := v.(map[string]interface{})
		updaterWeb3.disposeTransaction(transaction)
	}
}
func (updaterWeb3 *EthUpdaterWeb3) disposeTransaction(transaction map[string]interface{}) {
	transactionReceiptInfo, _ := updaterWeb3.connection.Eth.GetTransactionReceipt(transaction["hash"].(string))

	updaterWeb3.TableTx.Nonce = HexDec(transaction["nonce"].(string))
	updaterWeb3.TableTx.GasLimit = HexDec(transaction["gas"].(string))
	updaterWeb3.TableTx.Amount = formatAmount(transaction["value"].(string), 18)
	updaterWeb3.TableTx.GasPrice = HexDec(transaction["gasPrice"].(string))
	updaterWeb3.TableTx.InputData = transaction["input"].(string)
	updaterWeb3.TableTx.TxHash = transaction["hash"].(string)

	updaterWeb3.TableTx.ConfirmTime = updaterWeb3.TableBlock.TimeStamp

	updaterWeb3.TableTx.From = transaction["from"].(string)
	if transaction["to"] != nil {
		updaterWeb3.TableTx.To = transaction["to"].(string)
	} else {
		updaterWeb3.TableTx.To = ""
	}
	updaterWeb3.TableTx.BlockHeight = transactionReceiptInfo.BlockNumber.String()
	updaterWeb3.TableTx.BlockHash = transactionReceiptInfo.BlockHash
	updaterWeb3.TableTx.GasUsed = transactionReceiptInfo.GasUsed.String()
	updaterWeb3.TableTx.Fee = makeFee(updaterWeb3.TableTx.GasPrice, updaterWeb3.TableTx.GasUsed)

	updaterWeb3.TableTx.BlockState = 1
	if transactionReceiptInfo.Status {
		updaterWeb3.TableTx.TxState = 1
	} else {
		updaterWeb3.TableTx.TxState = 0
	}

	if strings.HasPrefix(updaterWeb3.TableTx.InputData, _TRANSFER) {
		updaterWeb3.TableTx.IsToken = 1

		transactionParameters := new(dto.TransactionParameters)
		//处理token
		for _, v := range transactionReceiptInfo.Logs {
			from, to, amount, logindex := AnalysisTokenLog(v)
			if logindex == "" {
				continue
			}
			// todo 获取token精度 from 数据库 or 节点
			contractAddr := v.Address
			transactionParameters.From = contractAddr
			transactionParameters.To = contractAddr
			transactionParameters.Data = _tokenDecimals
			tokenDecimalRes, _ := updaterWeb3.connection.Eth.Call(transactionParameters)
			transactionParameters.Data = _tokenSymbol
			//tokenSymbol,_ := updaterWeb3.connection.Eth.Call(transactionParameters)

			tokenDecimal := HexDec(tokenDecimalRes.Result.(string))
			//fmt.Println(tokenDecimal, tokenSymbol)
			updaterWeb3.TableTokenTx.From = from
			updaterWeb3.TableTokenTx.To = to
			updaterWeb3.TableTokenTx.TxHash = updaterWeb3.TableTx.TxHash
			updaterWeb3.TableTokenTx.LogIndex = logindex
			updaterWeb3.TableTokenTx.BlockState = updaterWeb3.TableTx.BlockState
			updaterWeb3.TableTokenTx.GasUsed = updaterWeb3.TableTx.GasUsed
			updaterWeb3.TableTokenTx.GasPrice = updaterWeb3.TableTx.GasPrice
			updaterWeb3.TableTokenTx.Fee = updaterWeb3.TableTx.Fee
			updaterWeb3.TableTokenTx.ConfirmTime = updaterWeb3.TableTx.ConfirmTime
			updaterWeb3.TableTokenTx.InputData = updaterWeb3.TableTx.InputData
			updaterWeb3.TableTokenTx.IsToken = 1
			updaterWeb3.TableTokenTx.TxState = 1
			updaterWeb3.TableTokenTx.BlockHeight = updaterWeb3.TableTx.BlockHeight
			updaterWeb3.TableTokenTx.BlockHash = updaterWeb3.TableTx.BlockHash
			updaterWeb3.TableTokenTx.Decimal = tokenDecimal
			updaterWeb3.TableTokenTx.ContractAddr = contractAddr
			updaterWeb3.TableTokenTx.Nonce = updaterWeb3.TableTx.Nonce
			updaterWeb3.TableTokenTx.GasLimit = updaterWeb3.TableTx.GasLimit
			intDecimal, _ := strconv.Atoi(tokenDecimal)
			updaterWeb3.TableTokenTx.Amount = formatAmount(amount, intDecimal)
			// 数据库操作
			updaterWeb3.TableTokenTx.InsertOneRaw(updaterWeb3.TableTokenTx)
		}
	} else {
		updaterWeb3.TableTx.IsToken = 0
	}
	updaterWeb3.TableTx.InsertOneRaw(updaterWeb3.TableTx)
}

func FormatAddress(addr string) string {
	k := false
	s1 := "0x"
	for _, v := range addr[2:] {
		if v != 48 {
			k = true
		}
		if k {
			s1 += string(v)
		}
	}

	l := len(s1)
	if l < 42 {
		t := ""
		for i := 0; i < 42-l; i++ {
			t += "0"
		}
		s1 = s1[:2] + t + s1[2:]
	}

	if len(s1) > 42 {
		log.Errorf("format addr error %s", addr)
		panic("format addr error %s")
	}
	return s1
}

func AnalysisTokenLog(logs dto.TransactionLogs) (from, to, amount, logindex string) {
	if len(logs.Topics) > 2 && strings.Compare(logs.Topics[0], _TRANSACTION_TOPIC) == 0 {
		to = FormatAddress(logs.Topics[2])
		from = FormatAddress(logs.Topics[1])

		amount = logs.Data
		logindex = logs.LogIndex.String()
	}
	return
}

func formatAmount(s string, d int) string {
	var str string
	if len(s) > 2 && strings.HasPrefix(strings.ToLower(s[:2]), "0x") {
		str = s[2:]
	} else {
		str = s
	}
	amount := HexDec(str)
	return format10Decimals(amount, d)
}

func format10Decimals(amount string, i int) string {
	tempFloat, _ := decimal.NewFromString(amount)
	t := decimal.NewFromFloat(1.0 * math.Pow(10, float64(i)))
	tempFloat = tempFloat.Div(t)
	return tempFloat.String()
}

func makeFee(gasprice, gasused string) string {
	a, _ := decimal.NewFromString(gasprice)
	b, _ := decimal.NewFromString(gasused)
	a = a.Div(decimal.NewFromFloat(1.0 * math.Pow(10, 18)))
	c := a.Mul(b)
	return c.String()
}

func NewEthUpdaterWeb3(startHeight uint64) *EthUpdaterWeb3 {
	u := new(EthUpdaterWeb3)
	u.StartHeight = startHeight
	u.TableBlock = new(models.Block)
	u.TableTx = new(models.Tx)
	u.TableTokenTx = new(models.TokenTx)
	u.TableAddress = new(models.Address)
	u.TableTokenAddress = new(models.TokenAddress)
	timeOut := conf.EthRpcTimeOut
	source := conf.EthRpcSecure
	url := conf.EthRpcHost + conf.EthRpcPort
	u.connection = web3.NewWeb3(providers.NewHTTPProvider(url, int32(timeOut), source))
	u.parityParam = map[string]interface{}{
		"id":      1,
		"jsonrpc": "2.0",
	}
	return u
}
