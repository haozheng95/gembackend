package models

import (
	"time"
	"github.com/astaxie/beego/orm"
)

type Address struct {
	Id              int64
	WalletId        string    `orm:"unique"`
	Addr            string    `orm:"unique"`
	TypeId          int
	Nonce           string
	Created         time.Time `orm:"auto_now_add;type(datetime)"`
	Amount          string
	UnconfirmAmount string
	Updated         time.Time `orm:"auto_now;type(datetime)"`
	Decimal         int64     `orm:"default(18)"`
}

func (address *Address) Update(s string) *Address {
	qs := o.QueryTable(address)
	p := orm.Params{"nonce": address.Nonce,
		"amount": address.Amount,
		"unconfirm_amount": address.UnconfirmAmount,
		"updated": time.Now()}
	qs.Filter("addr", s).Update(p)
	return address
}

func (Self *Address) InsertOneRaw(data *Address) *Address {
	data.Id = 0
	data.Created = time.Now()
	id, err := o.Insert(data)
	if err != nil {
		log.Errorf("Address insert error %s", err)
	}
	log.Debugf("Address insert id %startscript", id)
	return Self
}

func (address *Address) SelectAddr() (*Address, error) {
	qs := o.QueryTable(address)
	err := qs.Filter("wallet_id", address.WalletId).One(address)
	if err != nil{
		log.Errorf("select addr error %s", err)
	}
	return address, err
}

func (u *Address) TableEngine() string {
	return "MYISAM"
}