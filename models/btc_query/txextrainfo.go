//software: GoLand
//file: txextrainfo.go
//time: 2018/9/6 下午5:44
package btc_query

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type TxExtraInfoBtc struct {
	Id       int64
	WalletId string `orm:"index"`
	Vin      string `orm:"size(1000)"`
	To       string `orm:"index"`
	Change   string `orm:"index"`
	TxHash   string `orm:"unique"`
	Amount   string
	Comment  string    `orm:"size(1000)"`
	Created  time.Time `orm:"auto_now_add;type(datetime)"`
}

func NewTxExtrainfo(walletId, vin, to, change, txhash, amount, comment string) *TxExtraInfoBtc {
	return &TxExtraInfoBtc{WalletId: walletId, Vin: vin, To: to, Change: change,
		TxHash: txhash, Amount: amount, Comment: comment, Created: time.Now()}
}

func (t *TxExtraInfoBtc) Insert() {
	o := orm.NewOrm()
	o.Using(databases)
	if _, err := o.Insert(t); err != nil {
		log.Warning(err)
	}
}
