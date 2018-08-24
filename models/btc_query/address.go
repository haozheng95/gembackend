//software: GoLand
//file: address.go
//time: 2018/8/20 下午4:07
package btc_query

import "time"

type AddressBtc struct {
	Id              int64
	WalletId        string `orm:"index"`
	Addr            string `orm:"unique"`
	TypeId          int
	Nonce           string
	Created         time.Time `orm:"auto_now_add;type(datetime)"`
	Amount          string
	UnconfirmAmount string
	Updated         time.Time `orm:"auto_now;type(datetime)"`
	Decimal         int64     `orm:"default(8)"`
}
