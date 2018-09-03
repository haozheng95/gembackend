package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gembackend/models/btc_query"
	"github.com/gembackend/models/eth_query"
	"strings"
)

//software: GoLand
//file: txinfo.go
//time: 2018/7/31 下午2:36
type TxinfoController struct {
	beego.Controller
}

func (t *TxinfoController) Get() {
	cointype := t.Ctx.Input.Param(":coin_type")
	txhash := t.Input().Get("tx_hash")
	switch cointype {
	case "eth":
		// get now height
		blockHeight := eth_query.GetBlockHeight()
		//log.Debug(blockHeight)
		contract := t.Input().Get("contract_addr")
		if len(contract) > 1 {
			// contract info
			r := eth_query.GetTokenTxinfo(txhash)
			if len(r) > 0 {
				tmp := r[0]
				if tmp.TxState != -1 && strings.Compare(tmp.BlockHeight, "") != 0 && len(tmp.BlockHeight) > 0 {
					tmp.ConfirmNum = SubString(blockHeight, tmp.BlockHeight)
				}
				t.Data["json"] = resultResponseMake(tmp)
			} else {
				t.Data["json"] = resultResponseErrorMake(2012, nil)
			}
		} else {
			// eth info
			r := eth_query.GetTxInfo(txhash)
			if len(r) > 0 {
				tmp := r[0]
				if tmp.TxState != -1 && strings.Compare(tmp.BlockHeight, "") != 0 && len(tmp.BlockHeight) > 0 {
					tmp.ConfirmNum = SubString(blockHeight, tmp.BlockHeight)
				}
				t.Data["json"] = resultResponseMake(tmp)
			} else {
				t.Data["json"] = resultResponseErrorMake(2012, nil)
			}
		}
	case "btc":
		blockHeight := btc_query.CurrBlockNum()
		txInfos := btc_query.GetTxInfo(txhash)
		if len(txInfos) > 0 {
			txInfo := txInfos[0]
			confirmNum := blockHeight - txInfo.BlockNum
			t.Data["json"] = resultResponseMake(map[string]interface{}{
				"confirmnum": confirmNum,
				"currentnum": blockHeight,
				"txinfo":     txInfo,
			})
		} else {
			t.Data["json"] = resultResponseErrorMake(2012, nil)
		}

	default:
		t.Data["json"] = resultResponseErrorMake(2010, nil)
	}
	t.ServeJSON()
}

//type txinfores struct {
//	From         string `json:"from"`
//	To           string `json:"to"`
//	Amount       string `json:"amount"`
//	InputData    string `json:"input_data"`
//	GasLimit     string `json:"gas_limit"`
//	GasPrice     string `json:"gas_price"`
//	GasUsed      string `json:"gas_used"`
//	Fee          string `json:"fee"`
//	TxHash       string `json:"tx_hash"`
//	BlockHash    string `json:"block_hash"`
//	BlockHeight  string `json:"block_height"`
//	TxState      int    `json:"tx_state"`
//	IsToken      int    `json:"is_token"`
//	Logindex     string `json:"logindex"`
//	ContractAddr string `json:"contract_addr"`
//	Decimal      string `json:"decimal"`
//	ConfirmTime  string `json:"confirm_time"`
//	ConfirmNum   string `json:"confirm_num"`
//}
