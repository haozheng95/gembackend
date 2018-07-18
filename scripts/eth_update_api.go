package scripts

import (
	"github.com/gembackend/rpc"
	"github.com/gembackend/gembackendlog"
	"github.com/gembackend/models"
	"strconv"
	"strings"
	"math"
	"github.com/shopspring/decimal"
	"fmt"
)

type Updater interface {
	Forever()
}

var log = gembackendlog.Log

const (
	_TRANSFER          = "0xa9059cbb"
	_TRANSFER_FROM     = "0x23b872dd"
	_TRANSACTION_TOPIC = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
)

type EthUpdater struct {
	StartHeight       uint64
	TxReceipt         *rpc.Response
	TableBlock        *models.Block
	TableTx           *models.Tx
	TableTokenTx      *models.TokenTx
	TableAddress      *models.Address
	TableTokenAddress *models.TokenAddress
}

func (updater *EthUpdater) Forever() {
	// todo 准备开始的块高度
	//height := MaxInt(updater.StartHeight, updater.TableBlock.BlockHeight)
	height := updater.StartHeight
	for {
		// todo 获取块信息
		blockInfo := updater.getBlockInfo(height)
		dbBlockInfo := updater.TableBlock.SelectRawByHeight(height - 1)
		//// todo 验证块高度
		if dbBlockInfo.Id != 0 && blockInfo.Result["parentHash"] != dbBlockInfo.BlockHash {
			// 开始回滚高度
			log.Warningf("block exception!! will rollback !! except height = %startscript", height-1)
			height, blockInfo = updater.RollBackBlock(height)
		}
		// todo 更新块
		log.Warningf("\n The height %startscript", height)
		if blockInfo.Result != nil {
			updater.BeginUpdateBlockInfo(blockInfo)
			height++
		} else {
			log.Warning("result is nil")
		}
	}
}

func (updater *EthUpdater) getBlockInfo(height uint64) (*rpc.Response) {
	info, err := rpc.Eth_getBlockByNumber(height)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	blockInfo, err := rpc.FormatResponse(&info)
	if err != nil {
		panic(err)
	}
	return blockInfo
}

func (updater *EthUpdater) RollBackBlock(height uint64) (uint64, *rpc.Response) {
	for i := height; i >= updater.StartHeight; i-- {
		blockInfo := updater.getBlockInfo(i)
		dbBlockInfo := updater.TableBlock.SelectRawByHeight(i)
		if dbBlockInfo.BlockHash != blockInfo.Result["hash"] && dbBlockInfo.Id != 0 {
			log.Debugf("delete height %startscript", i)
			updater.TableBlock.DeleteOneRaw(dbBlockInfo.BlockHash)
		} else if dbBlockInfo.BlockHash == blockInfo.Result["hash"] {
			log.Debugf("delete block over!! will recover  i = %startscript", i)
			return i, blockInfo
		} else {
			panic("RollBackBlock Error")
		}
	}
	return height, updater.getBlockInfo(height)
}

func (updater *EthUpdater) formatBlockInfo(info map[string]interface{}) {
	//var err error = nil
	if info["nonce"] != nil {
		updater.TableBlock.Nonce = HexDec(info["nonce"].(string))
	}
	updater.TableBlock.ExtraData = info["extraData"].(string)
	updater.TableBlock.MixHash = info["mixHash"].(string)
	updater.TableBlock.Miner = info["miner"].(string)
	updater.TableBlock.ParentHash = info["parentHash"].(string)
	//updater.TableBlock.TimeStamp, err = strconv.ParseUint(info["timestamp"].(string)[2:], 16, 64)
	updater.TableBlock.TimeStamp = HexDec(info["timestamp"].(string))
	//updater.TableBlock.BlockHeight, err = strconv.ParseUint(info["number"].(string)[2:], 16, 64)
	updater.TableBlock.BlockHeight = HexDec(info["number"].(string))
	updater.TableBlock.BlockHash = info["hash"].(string)
	updater.TableBlock.GasLimit = HexDec(info["gasLimit"].(string)[2:])
	updater.TableBlock.GasUsed = HexDec(info["gasUsed"].(string)[2:])
	updater.TableBlock.Size = HexDec(info["size"].(string)[2:])
}

func (updater *EthUpdater) BeginUpdateBlockInfo(response *rpc.Response) {
	//	插入块信息
	info := response.Result
	updater.formatBlockInfo(info)
	updater.disposeTransaction(info["transactions"])
	updater.TableBlock.InsertOneRaw(updater.TableBlock)
}
func (updater *EthUpdater) disposeTransaction(v interface{}) {
	ws := v.([]interface{})
	for _, k := range ws {

		updater.formatTransaction(k.(map[string]interface{}))

		updater.formatTransactionOther()

		if strings.HasPrefix(updater.TableTx.InputData, _TRANSFER) {
			updater.TableTx.IsToken = 1
			updater.formatTokenTransaction()
		} else {
			updater.TableTx.IsToken = 0
		}
		updater.TableTx.InsertOneRaw(updater.TableTx)

		// 更新用户以太坊信息
		updater.disposeUpdateEthInfo(updater.TableTx.From)
		updater.disposeUpdateEthInfo(updater.TableTx.To)
	}

}
func (updater *EthUpdater) formatTransaction(s map[string]interface{}) {

	updater.TableTx.Nonce = HexDec(s["nonce"].(string)[2:])
	updater.TableTx.GasLimit = HexDec(s["gas"].(string)[2:])
	updater.TableTx.GasPrice = HexDec(s["gasPrice"].(string)[2:])

	updater.TableTx.From = s["from"].(string)
	if s["to"] != nil {
		updater.TableTx.To = s["to"].(string)
	}
	updater.TableTx.BlockHash = updater.TableBlock.BlockHash
	updater.TableTx.BlockHeight = updater.TableBlock.BlockHeight
	updater.TableTx.ConfirmTime = updater.TableBlock.TimeStamp
	updater.TableTx.BlockState = 1
	updater.TableTx.InputData = s["input"].(string)
	updater.TableTx.Amount = updater.FormatAmount(s["value"].(string), 18)
	updater.TableTx.TxHash = s["hash"].(string)

}
func (updater *EthUpdater) FormatAmount(s string, i int) string {
	var str string
	if len(s) > 2 && strings.HasPrefix(strings.ToLower(s[:2]), "0x") {
		str = s[2:]
	} else {
		str = s
	}
	amount := HexDec(str)
	return updater.format10Decimals(amount, i)
}

func (updater *EthUpdater) format10Decimals(amount string, i int) string {
	tempFloat, _ := decimal.NewFromString(amount)
	t := decimal.NewFromFloat(1.0 * math.Pow(10, float64(i)))
	tempFloat = tempFloat.Div(t)
	return tempFloat.String()
}

func (updater *EthUpdater) formatTransactionOther() {
	res, err := rpc.Eth_getTransactionReceipt(updater.TableTx.TxHash)
	if err != nil {
		log.Errorf("Eth_getTransactionReceipt Error %s", err)
		return
	}
	updater.TxReceipt, err = rpc.FormatResponse(&res)
	if err != nil {
		log.Errorf("FormatResponse Error %s", err)
		return
	}
	if updater.TxReceipt.Error != nil {
		log.Errorf("Response Error %s", updater.TxReceipt.Error)
		return
	}

	updater.FormatReceipt()
}
func (updater *EthUpdater) FormatReceipt() {
	info := updater.TxReceipt.Result
	updater.TableTx.GasUsed = HexDec(info["gasUsed"].(string)[2:])
	updater.TableTx.Fee = updater.MakeFee()
}
func (updater *EthUpdater) MakeFee() string {
	a, _ := decimal.NewFromString(updater.TableTx.GasPrice)
	b, _ := decimal.NewFromString(updater.TableTx.GasUsed)
	a = a.Div(decimal.NewFromFloat(1.0 * math.Pow(10, 18)))
	c := a.Mul(b)
	//t := decimal.NewFromFloat(1.0 * math.Pow(10, 18))
	//startscript := c.Div(t)
	//return startscript.String()
	return c.String()
}
func (updater *EthUpdater) formatTokenTransaction() {
	updater.TableTokenTx.From = updater.TableTx.From
	updater.TableTokenTx.ContractAddr = updater.TableTx.To
	updater.TableTokenTx.InputData = updater.TableTx.InputData
	updater.TableTokenTx.Nonce = updater.TableTx.Nonce
	updater.TableTokenTx.GasUsed = updater.TableTx.GasUsed
	updater.TableTokenTx.GasPrice = updater.TableTx.GasPrice
	updater.TableTokenTx.GasLimit = updater.TableTx.GasLimit
	updater.TableTokenTx.Fee = updater.TableTx.Fee
	updater.TableTokenTx.TxHash = updater.TableTx.TxHash
	updater.TableTokenTx.BlockHeight = updater.TableBlock.BlockHeight
	updater.TableTokenTx.BlockHash = updater.TableBlock.BlockHash
	updater.TableTokenTx.ConfirmTime = updater.TableBlock.TimeStamp
	updater.AnalysisTokenLog()
}
func (updater *EthUpdater) AnalysisTokenLog() {
	info := updater.TxReceipt.Result["logs"].([]interface{})
	for _, v := range info {
		t := v.(map[string]interface{})
		t1 := t["topics"].([]interface{})
		if strings.Compare(t1[0].(string), _TRANSACTION_TOPIC) == 0 {
			updater.TableTokenTx.To = t1[2].(string)
			dec, _ := rpc.Eth_getTokenDecimals(updater.TableTokenTx.ContractAddr)
			f, _ := rpc.FormatTokenResponse(dec)
			updater.TableTokenTx.Decimal = HexDec(f.Result[2:])
			updater.TableTokenTx.LogIndex = HexDec(t["logIndex"].(string))
			i, _ := strconv.Atoi(updater.TableTokenTx.Decimal)
			updater.TableTokenTx.Amount = updater.FormatAmount(t["data"].(string), i)
			updater.TableTokenTx.TxState = 1
			updater.TableTokenTx.BlockState = 1
			updater.TableTokenTx.IsToken = 1
			// 添加表数据
			updater.TableTokenTx.InsertOneRaw(updater.TableTokenTx)

			// 更新token用户信息
			updater.disposeUpdateEthTokenInfo(updater.TableTokenTx.From, updater.TableTokenTx.ContractAddr)
			updater.disposeUpdateEthTokenInfo(updater.TableTokenTx.To, updater.TableTokenTx.ContractAddr)
		}
	}
}

func (updater *EthUpdater) disposeUpdateEthInfo(addr string) {
	//更新用户信息
	//updater.TableAddress.Update(updater.TableTx.From)
	//updater.TableAddress.Update(updater.TableTx.To)

	r1, _ := rpc.Eth_getTransactionCount(addr)
	f1, err := rpc.FormatTokenResponse(r1)
	if err != nil {
		log.Errorf("get nonce error %s", err)
	}
	a1 := rpc.Eth_getBalance(updater.TableTx.From)
	nonce := HexDec(f1.Result)
	amountTmp, _ := rpc.FormatTokenResponse(a1)
	amount := updater.format10Decimals(amountTmp.Result, 18)
	updater.TableAddress.Nonce = nonce
	updater.TableAddress.Amount = amount
	updater.TableAddress.UnconfirmAmount = "0"
	updater.TableAddress.Update(addr)
}
func (updater *EthUpdater) disposeUpdateEthTokenInfo(addr string, contractAddr string) {
	// 获取token精度
	decimalTemp, _ := rpc.Eth_getTokenDecimals(contractAddr)
	d, _ := rpc.FormatTokenResponse(decimalTemp)
	dec := HexDec(d.Result)
	decInt, _ := strconv.Atoi(dec)
	// 获取用户余额
	amountResponse, err := rpc.Eth_getTokenBalance(contractAddr, addr[2:])
	if err != nil {
		log.Warningf("Eth_getTokenBalance error %s", err)
	}
	amountTemp, _ := rpc.FormatTokenResponse(amountResponse)

	amount := updater.FormatAmount(amountTemp.Result, decInt)
	// db操作

	updater.TableTokenAddress.Decimal = int64(decInt)
	updater.TableTokenAddress.Amount = amount
	updater.TableTokenAddress.UnconfirmAmount = "0"
	updater.TableTokenAddress.ContractAddr = contractAddr
	updater.TableTokenAddress.Addr = addr
	// update
	updater.TableTokenAddress.Update(updater.TableTokenAddress.Addr)
}

func MaxInt(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func MaxIntByString(a uint64, b string) uint64 {
	c , _ := strconv.ParseUint(b, 10, 64)
	if a > c {
		return a
	}
	return c
}

func NewEthUpdater(startHeight uint64) *EthUpdater {
	u := new(EthUpdater)
	u.StartHeight = startHeight
	u.TableBlock = new(models.Block)
	u.TableTx = new(models.Tx)
	u.TableTokenTx = new(models.TokenTx)
	u.TableAddress = new(models.Address)
	u.TableTokenAddress = new(models.TokenAddress)
	return u
}

func HexDec(h string) (n string) {
	//log.Debugf("------- %s", h)
	if len(h) > 2 && strings.HasPrefix(strings.ToLower(h[:2]), "0x") {
		h = h[2:]
	} else if strings.Compare(h, "0x") == 0 {
		h = "0"
	}

	s := strings.Split(strings.ToUpper(h), "")
	l := len(s)
	i := 0
	d := decimal.NewFromFloat(0)
	hex := map[string]string{"A": "10", "B": "11", "C": "12", "D": "13", "E": "14", "F": "15"}
	for i = 0; i < l; i++ {
		c := s[i]
		if v, ok := hex[c]; ok {
			c = v
		}
		f, err := strconv.ParseFloat(c, 10)
		if err != nil {
			fmt.Println(h)
			log.Error(err)
			return decimal.NewFromFloat(-1).String()
		}
		d = d.Add(decimal.NewFromFloat(f * math.Pow(16, float64(l-i-1))))
	}
	return d.String()
}
