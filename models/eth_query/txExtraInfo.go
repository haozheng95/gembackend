package eth_query

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type TxExtraInfo struct {
	Id          int64
	From        string `orm:"index"`
	To          string `orm:"index"`
	TxHash      string `orm:"unique"`
	Nonce       string
	Amount      string
	TokenAmount string
	Comment     string    `orm:"size(1000)"`
	Created     time.Time `orm:"auto_now_add;type(datetime)"`
}

//func (u *TxExtraInfo) TableEngine() string {
//	return "MYISAM"
//}

func (t *TxExtraInfo) InsertOneRaw() (r *TxExtraInfo, err error) {
	o := orm.NewOrm()
	o.Using(databases)
	t.Id = 0
	t.Created = time.Now()
	_, err = o.Insert(t)
	if err != nil {
		log.Errorf("tx extra info insert error %s", err)
	}
	return
}
