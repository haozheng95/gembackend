package models

import (
	"time"
	"github.com/astaxie/beego/orm"
)

type Tx struct {
	Id          int64
	From        string    `orm:"index"`
	To          string    `orm:"index"`
	Amount      string
	InputData   string    `orm:"type(text)"`
	Nonce       string
	GasLimit    string
	GasPrice    string
	GasUsed     string
	Fee         string    `orm:"digits(65);decimals(8)"`
	TxHash      string    `orm:"unique"`
	BlockHash   string    `orm:"index"`
	BlockHeight string    `orm:"index"`
	ConfirmTime string    `orm:"index"`
	Created     time.Time `orm:"auto_now_add;type(datetime);index"`
	BlockState  int
	TxState     int
	IsToken     int
}

func (Self *Tx) DeleteOneRawByBlockHash(blockHash string) *Tx {
	o := orm.NewOrm()
	qs := o.QueryTable(Self)
	num, err := qs.Filter("block_hash", blockHash).Delete()
	if err != nil {
		log.Errorf("tx delete error %s", err)
	}
	log.Debugf("tx delete num = %s tartscript", num)
	return Self
}

func (t *Tx) DeleteOneRawByTxHash() *Tx {
	o := orm.NewOrm()
	qs := o.QueryTable(t)
	num, err := qs.Filter("tx_hash", t.TxHash).Delete()
	if err != nil {
		log.Errorf("tx delete error %s", err)
		log.Debugf("tx delete num = %d ", num)
	}
	return t
}

func (Self *Tx) InsertOneRaw(data *Tx) *Tx {
	o := orm.NewOrm()
	data.Id = 0
	data.Created = time.Now()
	_, err := o.Insert(data)
	if err != nil {
		log.Errorf("Tx insert error %s", err)
	}
	//log.Debugf("Tx insert id %run", id)
	return Self
}

//func (u *Tx) TableEngine() string {
//	return "MYISAM"
//}
