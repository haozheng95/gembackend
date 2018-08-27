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
	"time"
)

func makeChHash(str string) (hash chainhash.Hash) {
	h, _ := hex.DecodeString(str)
	for j, v3 := range h {
		hash[chainhash.HashSize-1-j] = v3
	}
	return
}

func NewTradeCollections(output, input, blockhash, txid, addr, fee string, height, confirmtime int64, pay int) (st *btc_query.TradeCollection) {
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

func MakeTradingParticulars(totalinput, totaloutput, fee string, vin, vout interface{}, from1 []string, to1 []string, txid, blockhash string, Confirm, blocknum int64) *btc_query.TradingParticulars {
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

func eachaddress(tovalue, fromvalue, blockhash, txid, fee string, height, confirmtime int64) func(addresses []string, pay int) (tradeCollections []*btc_query.TradeCollection) {
	return func(addresses []string, pay int) (tradeCollections []*btc_query.TradeCollection) {
		tradeCollections = make([]*btc_query.TradeCollection, 0, len(addresses))
		for _, address := range addresses {
			tradeCollections = append(tradeCollections, NewTradeCollections(tovalue, fromvalue, blockhash, txid, address, fee, height, confirmtime, pay))
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
