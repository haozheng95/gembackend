package exchange

import (
	"time"
	"github.com/astaxie/beego/orm"
)

type MainChain struct {
	Id       int64
	Name     string    `orm:"index"`
	FullName string    `orm:"unique"`
	Decimal  int64
	Created  time.Time `orm:"auto_now_add;type(datetime)"`
	GasLimit string
	GasPrice string
	Fee      string
	Cny      string
	Usdt     string
	Updated  time.Time `orm:"auto_now;type(datetime)"`
}

func (e *MainChain) UpdateCnyAndUsdt() *MainChain {
	o := orm.NewOrm()
	o.Using(databases)
	qs := o.QueryTable(e)
	p := orm.Params{
		"cny":     e.Cny,
		"usdt":    e.Usdt,
		"updated": time.Now(),
	}
	qs.Filter("full_name", e.FullName).Update(p)
	return e
}

func (e *MainChain) SelectCny() *MainChain {
	o := orm.NewOrm()
	o.Using(databases)
	qs := o.QueryTable(e)
	err := qs.Filter("full_name", e.FullName).One(e)
	if err != nil {
		log.Errorf("select cny error :%s", err)
	}
	return e
}
