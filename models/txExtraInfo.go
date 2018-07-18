package models

import "time"

type TxExtraInfo struct {
	Id          int64
	From        string    `orm:"index"`
	To          string    `orm:"index"`
	TxHash      string    `orm:"unique"`
	Nonce       int64
	Amount      float64   `orm:"digits(65);decimals(8)"`
	TokenAmount float64   `orm:"digits(65);decimals(8)"`
	Comment     string    `orm:"size(1000)"`
	Created     time.Time `orm:"auto_now_add;type(datetime)"`
}
