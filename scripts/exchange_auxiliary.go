package scripts

import (
	"github.com/gembackend/models/exchange"
	"github.com/gembackend/rpc"
	"github.com/regcostajr/go-web3"
	"github.com/regcostajr/go-web3/complex/types"
	"github.com/regcostajr/go-web3/dto"
	"math/big"
	"time"
)

//software: GoLand
//file: exchange_auxiliary.go
//time: 2018/8/1 下午4:22

func AuxiliaryMain() {
	rpc.MakeConn()
	auxiliaryEth()
}

// get eth gasprice and gaslimit
func auxiliaryEth() {
	to := "0x962c2faad4bc2321b896b79ffc9d362295328ca5"
	data := "0xa9059cbb"

	web3Conn, ok := rpc.ConnectMap["eth-web3-original"]
	if !ok {
		web3Conn = rpc.ReMakeWeb3ConnOriginal()
	}
	conn := web3Conn.(*web3.Web3)
	// interval == 10s
	ticker := time.NewTicker(time.Second * interval)
	//ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	transaction := new(dto.TransactionParameters)
	transaction.To = to
	transaction.From = to
	transaction.Data = types.ComplexString(data)

	for range ticker.C {
		rpc.MakeConn()
		gasprice, err := conn.Eth.GetGasPrice()
		//log.Debug(gasprice)
		if err != nil {
			log.Errorf("gasprice error %s", err)
			continue
		}
		gaslimit, err := conn.Eth.EstimateGas(transaction)
		//log.Debug(gaslimit, err)
		if err != nil {
			log.Errorf("gaslimit error %s", err)
			continue
		}
		auxiliaryEthDb(gaslimit, gasprice)
	}

}

// update eth: gasprice and gaslimit with eth for main chain
func auxiliaryEthDb(gaslimit, gasprice *big.Int) {
	coin := "eth"
	dec := "10"
	for i := 0; i < 18; i++ {
		dec += "0"
	}
	gaslimitStr := gaslimit.String()
	gaspriceStr := gasprice.String()
	fee := DivString(MulString(gaslimitStr, gaspriceStr), dec).String()
	log.Infof("fee: %s", fee)
	exchange.UpdateMainChainGasPriceAndGasLimit(coin, gaslimitStr, gaspriceStr, fee)
}
