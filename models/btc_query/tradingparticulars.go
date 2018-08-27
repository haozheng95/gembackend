//software: GoLand
//file: tradingparticulars.go
//time: 2018/8/23 下午2:51
package btc_query

import (
	"time"
)

type TradingParticulars struct {
	Id         int64
	From       string    `orm:"type(text)"`
	To         string    `orm:"type(text)"`
	Txid       string    `orm:"unique"`
	Updated    time.Time `orm:"auto_now;type(datetime)"`
	BlockHash  string    `orm:"index"`
	BlockNum   int64
	Confirm    int64 `orm:"index"`
	TotalInput string
	TotalOut   string
	Fee        string
	Vin        string `orm:"type(text)"`
	Vout       string `orm:"type(text)"`
}

func MakeTradingParticulars(from, to, txid, blockhash string, Confirm, blocknum int64) *TradingParticulars {
	st := &TradingParticulars{
		From: from, To: to, Txid: txid, BlockHash: blockhash, Confirm: Confirm, BlockNum: blocknum, Updated: time.Now(),
	}
	return st
}
