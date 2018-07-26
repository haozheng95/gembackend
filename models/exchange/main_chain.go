package exchange

import (
	"time"
	"github.com/astaxie/beego/orm"
	"os"
)

type MainChain struct {
	Id       int64
	Name     string    `orm:"unique"`
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
	os.Unsetenv(databases)
	qs := o.QueryTable(e)
	p := orm.Params{
		"cny":     e.Cny,
		"usdt":    e.Usdt,
		"updated": time.Now(),
	}
	qs.Filter("name", e.Name).Update(p)
	return e
}