package controllers

import (
	"github.com/astaxie/beego"
	"encoding/json"
	"github.com/gembackend/rpc"
)

type SendRawTx struct {
	beego.Controller
}

func (s *SendRawTx) Post() {
	cointype := s.Ctx.Input.Param(":coin_type")
	m := make(map[string]string)
	err := json.Unmarshal(s.Ctx.Input.RequestBody, &m)
	if err != nil {
		//error
	}
	amount := m["amount"]
	raw := m["raw"]
	fee := m["fee"]
	from := m["from"]
	to := m["to"]
	note := m["note"]

	rpc.MakeConn()
	switch cointype {
	case "eth":
		conn, ck := rpc.ConnectMap["eth-web3"]
		var web3conn *rpc.Web3
		if ck {
			//log.Debug("--------", )
			web3conn = conn.(*rpc.Web3)
		} else {
			// get conn error
			web3conn = rpc.ReMakeWeb3Conn()
		}
		nonce := m["nonce"]
		gaslimit := m["gaslimit"]
		gasprice := m["gasprice"]
		eth_amount := amount
		token_amount := "0"
		is_token := 0
		contractaddr, ck := m["contract_addr"]
		if ck && len(contractaddr) > 1 {
			// token
			eth_amount = "0"
			token_amount = amount
			is_token = 1
		}
		txhash, err := web3conn.Eth.SendRawTransaction(raw)
		// dispose error
		// db operation
		log.Debug(fee, from, to, note, nonce, gaslimit, gasprice, eth_amount, token_amount, is_token, txhash)
		log.Debug(err)

	default:
		//error

	}
}
