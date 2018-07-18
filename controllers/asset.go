package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gembackend/models"
	"strconv"
)

type AssetController struct {
	beego.Controller
}

func (a *AssetController) Get() {

	addr := a.Input().Get("user_addr")
	begin := a.Input().Get("begin_page")
	size := a.Input().Get("size")
	beginInt, err := strconv.Atoi(begin)

	if err != nil {
		a.Data["json"] = resultResponseErrorMake(2001, err.Error())
		a.ServeJSON()
		return
	}

	sizeInt, err := strconv.Atoi(size)

	if err != nil {
		a.Data["json"] = resultResponseErrorMake(2001, err.Error())
		a.ServeJSON()
		return
	}

	o := orm.NewOrm()
	var address models.Address
	qs := o.QueryTable(address)
	qs.Filter("addr", addr).One(&address)
	var addressResult interface{}

	if address.Id == 0 {
		addressResult = nil
	} else {
		addressResult = address
	}

	var t []*models.TokenAddress
	qs = o.QueryTable("token_address")
	_, err = qs.Filter("addr", addr).Limit(sizeInt, (beginInt-1)*sizeInt).All(&t)
	if err != nil {
		a.Data["json"] = resultResponseErrorMake(2001, err.Error())
	} else {
		a.Data["json"] = resultResponseMake(map[string]interface{}{
			"eth":   addressResult,
			"token": t,
		})
	}
	a.ServeJSON()
}
