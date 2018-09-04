package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gembackend/models/btc_query"
	"github.com/gembackend/models/eth_query"
	"strings"
)

type TxrecordController struct {
	beego.Controller
}

func (t *TxrecordController) Get() {
	cointype := t.Ctx.Input.Param(":coin_type")
	walletid := t.Input().Get("wallet_id")
	contract := t.Input().Get("contract_addr")
	page, err := t.GetUint64("begin_page", 0)
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
		addst := new(eth_query.Address)
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
			txs, r = eth_query.GetEthTokenTxrecord(addst.Addr, contract, page*size, size)
		} else {
			//eth记录
			txs, r = eth_query.GetEthTxrecord(addst.Addr, page*size, size)
		}
		// save kafka
		//SaveForKafka(conf.KafkaTxRecordTopic, addst.Addr)
	case "btc":
		txs, r = btc_query.GetTxs(walletid, int(size), int(page*size))
	}

	t.Data["json"] = resultResponseMake(map[string]interface{}{
		"record": txs,
		"size":   r,
		"last":   r < int64(size),
	})
	t.ServeJSON()
}
