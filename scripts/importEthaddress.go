package scripts

import (
	"github.com/gembackend/rpc"
	"strings"
	"github.com/gembackend/models"
	"strconv"
	"github.com/regcostajr/go-web3"
	"github.com/gembackend/conf"
	"github.com/regcostajr/go-web3/providers"
	"github.com/regcostajr/go-web3/dto"
	"github.com/gembackend/controllers"
	"github.com/regcostajr/go-web3/complex/types"
)

var (
	models_tx            *models.Tx
	models_token_tx      *models.TokenTx
	models_address_token *models.TokenAddress
	models_address       *models.Address
	connection_web3      *web3.Web3
	tokenmap             map[string]string
)

func init() {
	models_tx = new(models.Tx)
	models_token_tx = new(models.TokenTx)
	models_address_token = new(models.TokenAddress)
	models_address = new(models.Address)

	url := conf.EthRpcHost + ":" + conf.EthRpcPort
	timeOut := conf.EthRpcTimeOut
	source := conf.EthRpcSecure
	connection_web3 = web3.NewWeb3(providers.NewHTTPProvider(url, int32(timeOut), source))
	initTokenmap()
}

func Main(walletid, address string) {
	// 添加eth账户
	addEthAccount(walletid, address)
	updateEthBalance(address)

	GetAllTxList(address)
	transactionParameters := new(dto.TransactionParameters)
	addTokenAddress(walletid, address, transactionParameters)
	updateTokenAddress(address, transactionParameters)
	// 重制tokenmap
	initTokenmap()
}

func GetAllTxList(address string) {
	z := 0
	for i := 1; i > 0; i++ {
		r, err := rpc.Eth_getTxList(address, i)
		if err != nil {
			return
		}
		resmap, err := rpc.FormatResponseMap(&r)
		if err != nil {
			log.Errorf("getAllTxList FormatResponseMap error %s", err)
			return
		}
		if strings.Compare(resmap["message"].(string), "OK") != 0 {
			log.Infof("address synchronization is complete,%d ---- %s", z, address)
			return
		}
		relist := resmap["result"].([]interface{})
		for _, v := range relist {
			insertTx(v.(map[string]interface{}))
			z++
		}
	}
}

// 初始化token map
func initTokenmap() {
	tokenmap = make(map[string]string, 10)
	for _, v := range controllers.DefaultToken {
		tokenmap[v[0]] = v[1]
	}
}

// 添加地址
func addTokenAddress(walletId, ethAddr string, transactionParameters *dto.TransactionParameters) {
	// 添加token地址
	models_address_token.WalletId = walletId
	models_address_token.Addr = ethAddr
	models_address_token.Added = 1
	models_address_token.UnconfirmAmount = "0"
	models_address_token.Amount = "0"

	transactionParameters.From = ethAddr
	transactionParameters.Data = _tokenSymbol

	for k, v := range tokenmap {

		transactionParameters.To = k
		tokenSymbol, err := connection_web3.Eth.Call(transactionParameters)
		if err != nil {
			log.Error(err)
			continue
		}
		models_address_token.ContractAddr = k
		t, _ := strconv.Atoi(v)
		models_address_token.Decimal = int64(t)
		symbol, _ := rpc.DecodeHexResponseToString(interface2string(tokenSymbol.Result))
		models_address_token.TokenName = symbol

		models_address_token.InsertOneRaw(models_address_token)
	}
}

// 更新token余额
func updateTokenAddress(ethAddr string, transactionParameters *dto.TransactionParameters) {
	for k, v := range tokenmap {
		transactionParameters.To = k
		transactionParameters.Data = types.ComplexString(_tokenBalance + ethAddr[2:])
		tokenBalanceRes, _ := connection_web3.Eth.Call(transactionParameters)
		tokenBalance := formatAmountString(tokenBalanceRes.Result.(string), v)

		models_address_token.Amount = tokenBalance
		models_address_token.Addr = ethAddr
		models_address_token.UnconfirmAmount = "0"
		models_address_token.ContractAddr = k

		models_address_token.UpdateAmount(ethAddr)
	}
}

// 添加eth账户
func addEthAccount(walletId, addr string) {
	models_address = &models.Address{
		WalletId:        walletId,
		Addr:            addr,
		Nonce:           "0",
		Amount:          "0",
		UnconfirmAmount: "0",
		TypeId:          16,
		Decimal:         18,
	}
	models_address.InsertOneRaw(models_address)
}

// 更新eth余额
func updateEthBalance(addr string) {
	userbalance, err := connection_web3.Eth.GetBalance(addr, _tag)
	if err != nil {
		log.Errorf("address balance format error %s", err)
		return
	}

	balance := format10Decimals(userbalance.String(), 18)

	usernonce, err := connection_web3.Eth.GetTransactionCount(addr, _tag)

	if err != nil {
		log.Errorf("address nonce format error %s", err)
		return
	}

	nonce := usernonce.String()
	models_address.Nonce = nonce
	models_address.Amount = balance
	models_address.UnconfirmAmount = "0"

	// db 操作
	models_address.Update(addr)
}

func insertTokenTx(tx_hash string) {
	transactionReceiptInfo, err := connection_web3.Eth.GetTransactionReceipt(tx_hash)
	if err != nil {
		log.Error(err)
		return
	}
	transactionParameters := new(dto.TransactionParameters)
	for _, v := range transactionReceiptInfo.Logs {
		from, to, amount, logindex := AnalysisTokenLog(v)
		if logindex == "" {
			continue
		}
		contractAddr := v.Address
		transactionParameters.From = contractAddr
		transactionParameters.To = contractAddr
		transactionParameters.Data = _tokenDecimals
		tokenDecimalRes, _ := connection_web3.Eth.Call(transactionParameters)

		tokenDecimal := HexDec(interface2string(tokenDecimalRes.Result))
		models_token_tx.From = from
		models_token_tx.To = to
		models_token_tx.TxHash = tx_hash
		models_token_tx.LogIndex = logindex
		models_token_tx.BlockState = 1
		models_token_tx.IsToken = 1
		models_token_tx.TxState = 1
		models_token_tx.Decimal = tokenDecimal
		models_token_tx.ContractAddr = contractAddr
		intDecimal, _ := strconv.Atoi(tokenDecimal)
		models_token_tx.Amount = formatAmount(amount, intDecimal)

		models_token_tx.GasUsed = models_tx.GasUsed
		models_token_tx.GasPrice = models_tx.GasPrice
		models_token_tx.Fee = models_tx.Fee
		models_token_tx.ConfirmTime = models_tx.ConfirmTime
		models_token_tx.InputData = models_tx.InputData
		models_token_tx.BlockHeight = models_tx.BlockHeight
		models_token_tx.BlockHash = models_tx.BlockHash
		models_token_tx.Nonce = models_tx.Nonce
		models_token_tx.GasLimit = models_tx.GasLimit

		models_token_tx.InsertOneRaw(models_token_tx)

		// token map add
		tokenmap[contractAddr] = tokenDecimal
	}
}

func insertTx(v map[string]interface{}) {
	models_tx.BlockHeight = interface2string(v["blockNumber"])
	models_tx.Amount = format10Decimals(interface2string(v["value"]), 18)
	models_tx.BlockHash = interface2string(v["blockHash"])
	models_tx.BlockState = 1
	models_tx.ConfirmTime = interface2string(v["timeStamp"])
	models_tx.TxState, _ = strconv.Atoi(interface2string(v["txreceipt_status"]))
	models_tx.From = interface2string(v["from"])
	models_tx.To = interface2string(v["to"])
	models_tx.GasLimit = interface2string(v["gas"])
	models_tx.GasPrice = interface2string(v["gasPrice"])
	models_tx.GasUsed = interface2string(v["gasUsed"])
	models_tx.Nonce = interface2string(v["nonce"])
	models_tx.InputData = interface2string(v["input"])
	models_tx.TxHash = interface2string(v["hash"])
	models_tx.Fee = makeFee(models_tx.GasPrice, models_tx.GasUsed)
	// 处理token交易
	if strings.HasPrefix(models_tx.InputData, _TRANSFER) {
		models_tx.IsToken = 1
		insertTokenTx(models_tx.TxHash)
	} else {
		models_tx.IsToken = 0
	}
	models_tx.InsertOneRaw(models_tx)
}

func interface2string(v interface{}) string {
	return v.(string)
}
