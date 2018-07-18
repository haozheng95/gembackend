package controllers

import (
	"github.com/astaxie/beego"
	"strings"
	"github.com/gembackend/models"
)

type TxrecordController struct {
	beego.Controller
}

func (t *TxrecordController) Get() {
	cointype := t.Ctx.Input.Param(":coin_type")
	walletid := t.Input().Get("wallet_id")
	contract := t.Input().Get("contract_addr")
	page, err := t.GetUint64("page", 0)
	if err != nil {
		t.Data["json"] = resultResponseErrorMake(2005, err.Error())
		t.ServeJSON()
		return
	}
	size, err := t.GetUint64("size", 10)
	if err != nil {
		t.Data["json"] = resultResponseErrorMake(2005, err.Error())
		t.ServeJSON()
		return
	}

	cointype = strings.ToLower(cointype)

	var txs interface{}
	var r int64
	switch cointype {
	case "eth":
		addst := new(models.Address)
		addst.WalletId = walletid
		_, err = addst.SelectAddr()
		if err != nil {
			//用户不存在
			log.Info("no such this user")
			t.Data["json"] = resultResponseErrorMake(2006, err.Error())
			t.ServeJSON()
			return
		}

		if len(contract) != 0 {
			//合约记录
			txs, r = models.GetEthTokenTxrecord(addst.Addr, contract, page*size, size)
		} else {
			//eth记录
			txs, r = models.GetEthTxrecord(addst.Addr, page*size, size)
		}
	}

	t.Data["json"] = resultResponseMake(map[string]interface{}{
		"record": txs,
		"size":r,
		"last":r < int64(size),
	})
	t.ServeJSON()
}
