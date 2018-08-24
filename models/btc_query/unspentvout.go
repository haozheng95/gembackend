//software: GoLand
//file: unspentvout.go
//time: 2018/8/23 下午2:59
package btc_query

import "time"

type UnspentVout struct {
	Id        int64
	Txid      string `orm:"index"`
	Spent     int    `orm:"default(0)"`
	SpentTxid string `orm:"index"`
	Index     int64
	Value     string
	Address   string    `orm:"index"`
	Updated   time.Time `orm:"auto_now;type(datetime)"`
}
