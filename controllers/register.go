package controllers

import (
	"github.com/astaxie/beego"
	"encoding/json"
	"strconv"
	"github.com/gembackend/models/eth_query"
)

type RegisterController struct {
	beego.Controller
}

//创建账户 添加十个默认token
func (r *RegisterController) Post() {
	m := make(map[string]string)
	err := json.Unmarshal(r.Ctx.Input.RequestBody, &m)
	if err != nil {
		log.Error(err)
		r.Data["json"] = resultResponseErrorMake(2001, err.Error())
		r.ServeJSON()
		return
	}

	walletId, err1 := m["wallet_id"]
	sign, err2 := m["sign"]
	ethAddr, err3 := m["eth_addr"]

	if !err1 || !err2 || !err3 {
		log.Error(err1, err2, err3)
		r.Data["json"] = resultResponseErrorMake(2000, nil)
		r.ServeJSON()
		return
	}

	if !checkSign(walletId, sign) {
		log.Error("checkSign false")
		r.Data["json"] = resultResponseErrorMake(2007, nil)
		r.ServeJSON()
		return
	}
	// 添加eth地址
	addressTable := &eth_query.Address{
		WalletId:        walletId,
		Addr:            ethAddr,
		Nonce:           "0",
		Amount:          "0",
		UnconfirmAmount: "0",
		TypeId:          4,
		Decimal:         18,
	}
	addressTable.InsertOneRaw(addressTable)
	// 添加token地址
	addressTokenTable := &eth_query.TokenAddress{
		WalletId:        walletId,
		Addr:            ethAddr,
		Amount:          "0",
		UnconfirmAmount: "0",
		Added:           1,
	}
	for _, v := range DefaultToken {
		addressTokenTable.ContractAddr = v[0]
		dec, _ := strconv.Atoi(v[1])
		addressTokenTable.Decimal = int64(dec)
		addressTokenTable.TokenName = v[2]
		addressTokenTable.InsertOneRaw(addressTokenTable)
	}

	r.Data["json"] = resultResponseMake("success")
	r.ServeJSON()
}
