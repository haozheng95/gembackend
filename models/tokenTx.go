package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"github.com/gembackend/gembackendlog"
)

var log = gembackendlog.Log

type TokenTx struct {
	Id           int64
	From         string    `orm:"index"`
	To           string    `orm:"index"`
	Amount       string
	InputData    string    `orm:"type(text)"`
	Nonce        string
	GasLimit     string
	GasPrice     string
	GasUsed      string
	Fee          string
	TxHash       string    `orm:"index"`
	BlockHash    string    `orm:"index"`
	BlockHeight  string    `orm:"index"`
	ConfirmTime  string    `orm:"index"`
	Created      time.Time `orm:"auto_now_add;type(datetime);index"`
	BlockState   int
	TxState      int
	IsToken      int       `orm:"default(1)"`
	LogIndex     string
	ContractAddr string    `orm:"index"`
	Decimal      string
}

//func (u *TokenTx) TableUnique() [][]string {
//	return [][]string{
//		{"TxHash", "LogIndex"},
//	}
//}

func (Self *TokenTx) DeleteOneRaw(blockHash string) *TokenTx {
	o := orm.NewOrm()
	qs := o.QueryTable(Self)
	num, err := qs.Filter("block_hash", blockHash).Delete()
	if err != nil {
		log.Errorf("token tx delete error %s", err)
	}
	log.Debugf("token tx delete num = %startscript", num)
	return Self
}
func (Self *TokenTx) InsertOneRaw(data *TokenTx) *TokenTx {
	o := orm.NewOrm()
	data.Id = 0
	data.Created = time.Now()
	_, err := o.Insert(data)
	if err != nil {
		log.Errorf("Tx insert error %s", err)
	}
	//log.Debugf("Tx insert id %startscript", id)
	return Self
}

//func (u *TokenTx) TableEngine() string {
//	return "MYISAM"
//}
