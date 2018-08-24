//software: GoLand
//file: tradecollection.go
//time: 2018/8/23 下午2:46
package btc_query

import "time"

type TradeCollection struct {
	Id      int64
	Addr    string    `orm:"index"`
	Txid    string    `orm:"index"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
	Amount  string
	Fee     string
	Pay     int
}
