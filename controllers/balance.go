package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/shopspring/decimal"
	"github.com/gembackend/models/eth_query"
)

type BalanceController struct {
	beego.Controller
}

// 获取eth余额 和合约地址余额
func (b *BalanceController) Get() {

	addr := b.Input().Get("user_addr")
	contract := b.Input().Get("contract_addr")
	var r eth_query.Address
	o := orm.NewOrm()
	qs := o.QueryTable("address")
	err := qs.Filter("addr", addr).One(&r)
	if err != nil {
		b.Data["json"] = resultResponseErrorMake(2003, err.Error())
		b.ServeJSON()
		return
	}
	amount1, err := decimal.NewFromString(r.Amount)
	amount2, err := decimal.NewFromString(r.UnconfirmAmount)

	amount := amount1.Sub(amount2)

	var c eth_query.TokenAddress
	var contractResult interface{} = nil
	if contract != "" {
		qs = o.QueryTable(c)
		err = qs.Filter("addr", addr).Filter("contract_addr", contract).One(&c)
		if err != nil {
			contractResult = err.Error()
		} else {
			tokenAmount1, _ := decimal.NewFromString(c.Amount)
			tokenAmount2, _ := decimal.NewFromString(c.UnconfirmAmount)
			tokenAmount := tokenAmount1.Sub(tokenAmount2)
			contractResult = map[string]interface{}{
				"contract_addr": contract,
				"decimal":       c.Decimal,
				"amount":        tokenAmount,
			}
		}
	}

	b.Data["json"] = resultResponseMake(map[string]interface{}{
		"eth_data": map[string]interface{}{
			"amount": amount,
			"nonce":  r.Nonce,
		},
		"token_data": contractResult,
	})
	b.ServeJSON()
}
