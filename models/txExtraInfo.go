package models

import "time"

type TxExtraInfo struct {
	Id          int64
	From        string    `orm:"index"`
	To          string    `orm:"index"`
	TxHash      string    `orm:"unique"`
	Nonce       int64
	Amount      string
	TokenAmount string
	Comment     string    `orm:"size(1000)"`
	Created     time.Time `orm:"auto_now_add;type(datetime)"`
}

//func (u *TxExtraInfo) TableEngine() string {
//	return "MYISAM"
//}
