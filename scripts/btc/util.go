//software: GoLand
//file: util.go
//time: 2018/8/22 上午10:43
package btc

import (
	"encoding/hex"
	"encoding/json"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/gembackend/models/btc_query"
	"github.com/shopspring/decimal"
	"math"
	"time"
)

func AccountUpdateMul(addrs []string) {
	for _, v := range addrs {
		AccountUpdate(v)
	}
}

func MapToSlice(input []map[string]float64) (output []string) {
	output = make([]string, 0, len(input))
	for _, v1 := range input {
		for k1, _ := range v1 {
			output = append(output, k1)
		}
	}
	return
}

func AccountUpdate(addr string) {
	amount, _ := decimal.NewFromString("0")
	amounts := btc_query.GetAllUnspent(addr)
	for _, value := range amounts {
		v, _ := decimal.NewFromString(value.Value)
		amount.Add(v)
	}
	log.Debug(amount.String())
	err := btc_query.UpdateAddr(addr, amount.String())
	if err != nil {
		log.Warning(err)
	}
}

func makeChHash(str string) (hash chainhash.Hash) {
	h, _ := hex.DecodeString(str)
	for j, v3 := range h {
		hash[chainhash.HashSize-1-j] = v3
	}
	return
}

func NewTradeCollections(output, input, blockhash, txid, addr, fee string,
	height, confirmtime int64, pay int, value float64) (st *btc_query.TradeCollection) {
	st = new(btc_query.TradeCollection)
	st.TotalInput = input
	st.TotalOutput = output
	st.Height = height
	st.BlockHash = blockhash
	st.Txid = txid
	st.Addr = addr
	st.ConfirmTime = confirmtime
	st.Pay = pay
	st.Fee = fee
	st.Updated = time.Now()
	st.Value = value
	return
}

func NewUnSpentVoutSt(txid, addr, value, blockhash string, height, index int64) (st *btc_query.UnspentVout) {
	st = new(btc_query.UnspentVout)
	st.Txid = txid
	st.Address = addr
	st.Index = index
	st.Value = value
	st.Updated = time.Now()
	st.BlockHash = blockhash
	st.Height = height
	return
}

func MakeTradingParticulars(totalinput, totaloutput, fee string, vin, vout interface{}, from1 []string,
	to1 []string, txid, blockhash string, Confirm, blocknum int64) *btc_query.TradingParticulars {
	from2, err := json.Marshal(from1)
	if err != nil {
		log.Warning("from2", err)
		return nil
	}
	to2, err := json.Marshal(to1)
	if err != nil {
		log.Warning("to2", err)
		return nil
	}
	vin1, err := json.Marshal(vin)
	if err != nil {
		log.Warning("vin", err)
		return nil
	}
	vout1, err := json.Marshal(vout)
	if err != nil {
		log.Warning("vout", err)
		return nil
	}
	log.Debug("vin======", string(vin1))
	log.Debug("vout=====", string(vout1))
	from := string(from2)
	to := string(to2)

	st := &btc_query.TradingParticulars{
		From: from, To: to, Txid: txid, BlockHash: blockhash, Confirm: Confirm, BlockNum: blocknum, Updated: time.Now(),
		Vin: string(vin1), Vout: string(vout1), TotalInput: totalinput, TotalOut: totaloutput, Fee: fee,
	}
	return st
}

func eachaddress(tovalue, fromvalue, blockhash, txid, fee string, height,
	confirmtime int64) func(addresses []map[string]float64, pay int) (
	tradeCollections []*btc_query.TradeCollection) {
	return func(addresses []map[string]float64, pay int) (tradeCollections []*btc_query.TradeCollection) {
		tradeCollections = make([]*btc_query.TradeCollection, 0, len(addresses))
		for _, v1 := range addresses {
			for k2, v2 := range v1 {
				tradeCollections = append(tradeCollections, NewTradeCollections(tovalue,
					fromvalue, blockhash, txid, k2, fee, height, confirmtime, pay, v2))

			}
		}
		return
	}
}

func Subfloat(d1, d2 float64) (r string) {
	m1 := decimal.NewFromFloat(d1)
	m2 := decimal.NewFromFloat(d2)
	m3 := m1.Sub(m2)
	r = m3.String()
	return
}

func floatToString(value float64) string {
	return decimal.NewFromFloat(value).String()
}

func Decimal(value float64) float64 {
	return math.Trunc(value*1e8) * 1e-8
}

var CoolStr string = `        
         __           __              _     
   _____/ /_______   / /_____  ____  (_)____
  / ___/ //_/ ___/  / __/ __ \/ __ \/ / ___/
 (__  ) ,< / /     / /_/ /_/ / /_/ / / /__  
/____/_/|_/_/      \__/\____/ .___/_/\___/  
                           /_/
`
