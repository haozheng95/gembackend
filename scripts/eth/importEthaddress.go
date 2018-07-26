package eth

import (
	"github.com/gembackend/rpc"
	"strings"
	"strconv"
	"github.com/regcostajr/go-web3/dto"
	"github.com/gembackend/controllers"
	"github.com/regcostajr/go-web3/complex/types"
	"github.com/gembackend/models/eth_query"
)



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
	eth_query_address_token.WalletId = walletId
	eth_query_address_token.Addr = ethAddr
	eth_query_address_token.Added = 1
	eth_query_address_token.UnconfirmAmount = "0"
	eth_query_address_token.Amount = "0"

	transactionParameters.From = ethAddr
	transactionParameters.Data = _tokenSymbol

	for k, v := range tokenmap {

		transactionParameters.To = k
		tokenSymbol, err := connection_web3.Eth.Call(transactionParameters)
		if err != nil {
			log.Error(err)
			continue
		}
		eth_query_address_token.ContractAddr = k
		t, _ := strconv.Atoi(v)
		eth_query_address_token.Decimal = int64(t)
		symbol, _ := rpc.DecodeHexResponseToString(interface2string(tokenSymbol.Result))
		eth_query_address_token.TokenName = symbol

		eth_query_address_token.InsertOneRaw(eth_query_address_token)
	}
}

// 更新token余额
func updateTokenAddress(ethAddr string, transactionParameters *dto.TransactionParameters) {
	for k, v := range tokenmap {
		transactionParameters.To = k
		transactionParameters.Data = types.ComplexString(_tokenBalance + ethAddr[2:])
		tokenBalanceRes, _ := connection_web3.Eth.Call(transactionParameters)
		tokenBalance := formatAmountString(tokenBalanceRes.Result.(string), v)

		eth_query_address_token.Amount = tokenBalance
		eth_query_address_token.Addr = ethAddr
		eth_query_address_token.UnconfirmAmount = "0"
		eth_query_address_token.ContractAddr = k

		eth_query_address_token.UpdateAmount(ethAddr)
	}
}

// 添加eth账户
func addEthAccount(walletId, addr string) {
	eth_query_address = &eth_query.Address{
		WalletId:        walletId,
		Addr:            addr,
		Nonce:           "0",
		Amount:          "0",
		UnconfirmAmount: "0",
		TypeId:          16,
		Decimal:         18,
	}
	eth_query_address.InsertOneRaw(eth_query_address)
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
	eth_query_address.Nonce = nonce
	eth_query_address.Amount = balance
	eth_query_address.UnconfirmAmount = "0"

	// db 操作
	eth_query_address.Update(addr)
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
		eth_query_token_tx.From = from
		eth_query_token_tx.To = to
		eth_query_token_tx.TxHash = tx_hash
		eth_query_token_tx.LogIndex = logindex
		eth_query_token_tx.BlockState = 1
		eth_query_token_tx.IsToken = 1
		eth_query_token_tx.TxState = 1
		eth_query_token_tx.Decimal = tokenDecimal
		eth_query_token_tx.ContractAddr = contractAddr
		intDecimal, _ := strconv.Atoi(tokenDecimal)
		eth_query_token_tx.Amount = formatAmount(amount, intDecimal)

		eth_query_token_tx.GasUsed = eth_query_tx.GasUsed
		eth_query_token_tx.GasPrice = eth_query_tx.GasPrice
		eth_query_token_tx.Fee = eth_query_tx.Fee
		eth_query_token_tx.ConfirmTime = eth_query_tx.ConfirmTime
		eth_query_token_tx.InputData = eth_query_tx.InputData
		eth_query_token_tx.BlockHeight = eth_query_tx.BlockHeight
		eth_query_token_tx.BlockHash = eth_query_tx.BlockHash
		eth_query_token_tx.Nonce = eth_query_tx.Nonce
		eth_query_token_tx.GasLimit = eth_query_tx.GasLimit

		eth_query_token_tx.InsertOneRaw(eth_query_token_tx)

		// token map add
		tokenmap[contractAddr] = tokenDecimal
	}
}

func insertTx(v map[string]interface{}) {
	eth_query_tx.BlockHeight = interface2string(v["blockNumber"])
	eth_query_tx.Amount = format10Decimals(interface2string(v["value"]), 18)
	eth_query_tx.BlockHash = interface2string(v["blockHash"])
	eth_query_tx.BlockState = 1
	eth_query_tx.ConfirmTime = interface2string(v["timeStamp"])
	eth_query_tx.TxState, _ = strconv.Atoi(interface2string(v["txreceipt_status"]))
	eth_query_tx.From = interface2string(v["from"])
	eth_query_tx.To = interface2string(v["to"])
	eth_query_tx.GasLimit = interface2string(v["gas"])
	eth_query_tx.GasPrice = interface2string(v["gasPrice"])
	eth_query_tx.GasUsed = interface2string(v["gasUsed"])
	eth_query_tx.Nonce = interface2string(v["nonce"])
	eth_query_tx.InputData = interface2string(v["input"])
	eth_query_tx.TxHash = interface2string(v["hash"])
	eth_query_tx.Fee = makeFee(eth_query_tx.GasPrice, eth_query_tx.GasUsed)
	// 处理token交易
	if strings.HasPrefix(eth_query_tx.InputData, _TRANSFER) {
		eth_query_tx.IsToken = 1
		insertTokenTx(eth_query_tx.TxHash)
	} else {
		eth_query_tx.IsToken = 0
	}
	eth_query_tx.InsertOneRaw(eth_query_tx)
}

func interface2string(v interface{}) string {
	return v.(string)
}
