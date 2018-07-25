package scripts

import (
	"github.com/gembackend/models"
	"github.com/regcostajr/go-web3"
	"time"
	"github.com/regcostajr/go-web3/providers"
	"fmt"
	"github.com/gembackend/rpc"
	"strings"
	"github.com/regcostajr/go-web3/dto"
	"strconv"
	"github.com/gembackend/conf"
	"github.com/regcostajr/go-web3/complex/types"
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
	// 获取数据库高度
	updaterWeb3.TableBlock = updaterWeb3.TableBlock.SelectMaxHeight()
	dbHehght := updaterWeb3.TableBlock.BlockHeight
	for height := MaxIntByString(updaterWeb3.StartHeight, dbHehght); ; {
		rpcHeight, err := updaterWeb3.connection.Eth.GetBlockNumber()
		if err != nil {
			log.Errorf("web3 rpcHeight error %s", err)
			return
		}

		log.Infof("db height %d ==== rpc height %s", height, rpcHeight.String())

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

			log.Infof("block update success %d", height)
			height++
		} else {
			//Have a rest
			log.Info("block pending")
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
	transactionReceiptInfo, err := updaterWeb3.connection.Eth.GetTransactionReceipt(transaction["hash"].(string))
	if err != nil {
		log.Error(err)
		return
	}
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
			// 判断是否是相关eth地址
			booltokenfrom := models.GetEthAddrExist(updaterWeb3.TableTokenTx.From)
			booltokento := models.GetEthAddrExist(updaterWeb3.TableTokenTx.To)
			// debug
			booltokenfrom = true
			booltokento = true
			if booltokenfrom || booltokento {
				// 更新用户token信息
				if booltokenfrom {
					updaterWeb3.disposeusertoken(updaterWeb3.TableTokenTx.From, tokenDecimal, transactionParameters)
				}

				if booltokento {
					updaterWeb3.disposeusertoken(updaterWeb3.TableTokenTx.To, tokenDecimal, transactionParameters)
				}
				// 数据库操作
				updaterWeb3.TableTokenTx.DeleteOneRawByHashAndLogindex(updaterWeb3.TableTokenTx.TxHash) //删除接口插入的记录
				updaterWeb3.TableTokenTx.InsertOneRaw(updaterWeb3.TableTokenTx)
			}

		}
	} else {
		updaterWeb3.TableTx.IsToken = 0
	}
	// 判断是否是相关eth地址
	boolfrom := models.GetEthAddrExist(updaterWeb3.TableTx.From)
	boolto := models.GetEthAddrExist(updaterWeb3.TableTx.To)
	// debug
	boolto = true
	boolfrom = true
	if boolfrom || boolto {
		if boolfrom {
			//log.Infof("eth from %s", updaterWeb3.TableTx.From)
			updaterWeb3.disposeuserbalance(updaterWeb3.TableTx.From)
		}
		if boolto {
			//log.Infof("eth to %s", updaterWeb3.TableTx.To)
			updaterWeb3.disposeuserbalance(updaterWeb3.TableTx.To)
		}
		updaterWeb3.TableTx.DeleteOneRawByTxHash() //删除当前hash
		updaterWeb3.TableTx.InsertOneRaw(updaterWeb3.TableTx)
	}
}

func (updaterWeb3 *EthUpdaterWeb3) disposeusertoken(useraddr string, tokenDecimal string, parameters *dto.TransactionParameters) {
	parameters.Data = types.ComplexString(_tokenBalance + useraddr[2:])
	tokenBalanceRes, _ := updaterWeb3.connection.Eth.Call(parameters)
	tokenBalance := formatAmountString(tokenBalanceRes.Result.(string), tokenDecimal)

	updaterWeb3.TableTokenAddress.Amount = tokenBalance
	updaterWeb3.TableTokenAddress.Addr = useraddr
	updaterWeb3.TableTokenAddress.UnconfirmAmount = "0"
	updaterWeb3.TableTokenAddress.ContractAddr = parameters.To
	// 数据库更新
	updaterWeb3.TableTokenAddress.Update(useraddr)
}
func (updaterWeb3 *EthUpdaterWeb3) disposeuserbalance(addr string) {
	userbalance, err := updaterWeb3.connection.Eth.GetBalance(addr, _tag)
	if err != nil {
		log.Errorf("address balance format error %s-----%s", err, addr)
		return
	}

	balance := format10Decimals(userbalance.String(), 18)

	usernonce, err := updaterWeb3.connection.Eth.GetTransactionCount(addr, _tag)

	if err != nil {
		log.Errorf("address nonce format error %s", err)
		return
	}

	nonce := usernonce.String()
	updaterWeb3.TableAddress.Nonce = nonce
	updaterWeb3.TableAddress.Amount = balance
	updaterWeb3.TableAddress.UnconfirmAmount = "0"

	// db 操作
	updaterWeb3.TableAddress.Update(addr)
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
	url := conf.EthRpcHost + ":" + conf.EthRpcPort
	u.connection = web3.NewWeb3(providers.NewHTTPProvider(url, int32(timeOut), source))
	u.parityParam = map[string]interface{}{
		"id":      1,
		"jsonrpc": "2.0",
	}
	return u
}
