//software: GoLand
//file: tradingparticulars.go
//time: 2018/8/23 下午2:51
package btc_query

import (
	"time"
)

type TradingParticulars struct {
	Id        int64
	From      string
	To        string
	Txid      string    `orm:"unique"`
	Updated   time.Time `orm:"auto_now;type(datetime)"`
	BlockHash string    `orm:"index"`
	BlockNum  string
	Confirm   string `orm:"index"`
}
