//software: GoLand
//file: unspentvout.go
//time: 2018/9/4 上午11:35
package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gembackend/models/btc_query"
)

type UnspentVout struct {
	beego.Controller
}

func (u *UnspentVout) Get() {
	coin := u.Ctx.Input.Param(":coin_type")
	walletId := u.Input().Get("wallet_id")
	var res interface{}
	switch coin {
	case "btc":
		res = btc_query.GetUnspent(walletId)
	default:
		res = resultResponseErrorMake(2010, nil)
	}
	u.Data["json"] = resultResponseMake(res)
	u.ServeJSON()
}
