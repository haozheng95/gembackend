package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gembackend/hjwt"
	"github.com/astaxie/beego/orm"
)

type AuthController struct {
	beego.Controller
}

func (a *AuthController) Get() {

	walletId := a.Input().Get("wallet_id")
	// from default eth_query
	o := orm.NewOrm()
	qs := o.QueryTable("address")
	n := qs.Filter("wallet_id", walletId).Exist()
	if n {
		token := hjwt.GenToken()
		a.Data["json"] = resultResponseMake(token)
	}else {
		a.Data["json"] = resultResponseErrorMake(2004, nil)
	}
	a.ServeJSON()
}
