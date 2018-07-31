package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/gembackend/models/eth_query"
	"github.com/gembackend/rpc"
	"strings"
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
		s.Data["json"] = resultResponseErrorMake(2001, err.Error())
		s.ServeJSON()
		return
	}
	amount, ck := m["amount"]
	if !ck {
		s.Data["json"] = resultResponseErrorMake(2000, "nonce")
		s.ServeJSON()
		return
	}
	raw, ck := m["raw"]
	if !ck {
		s.Data["json"] = resultResponseErrorMake(2000, "raw")
		s.ServeJSON()
		return
	}
	fee, ck := m["fee"]
	if !ck {
		s.Data["json"] = resultResponseErrorMake(2000, "fee")
		s.ServeJSON()
		return
	}
	from, ck := m["from"]
	if !ck {
		s.Data["json"] = resultResponseErrorMake(2000, "from")
		s.ServeJSON()
		return
	}
	to, ck := m["to"]
	if !ck {
		s.Data["json"] = resultResponseErrorMake(2000, "to")
		s.ServeJSON()
		return
	}
	note, ck := m["note"]
	if !ck {
		note = ""
	}

	rpc.MakeConn()
	switch cointype {
	case "eth":
		// error assert
		if strings.Compare(from, to) == 0 {
			s.Data["json"] = resultResponseErrorMake(2011, nil)
			s.ServeJSON()
			return
		}

		conn, ck := rpc.ConnectMap["eth-web3"]
		var web3conn *rpc.Web3
		if ck {
			//log.Debug("--------", )
			web3conn = conn.(*rpc.Web3)
		} else {
			// get conn error
			web3conn = rpc.ReMakeWeb3Conn()
		}
		nonce, ck := m["nonce"]
		if !ck {
			s.Data["json"] = resultResponseErrorMake(2000, "nonce")
			s.ServeJSON()
			return
		}
		gaslimit, ck := m["gaslimit"]
		if !ck {
			s.Data["json"] = resultResponseErrorMake(2000, "gaslimit")
			s.ServeJSON()
			return
		}
		gasprice, ck := m["gasprice"]
		if !ck {
			s.Data["json"] = resultResponseErrorMake(2000, "gasprice")
			s.ServeJSON()
			return
		}
		dec, ck := m["dec"]
		if !ck {
			s.Data["json"] = resultResponseErrorMake(2000, "dec")
			s.ServeJSON()
			return
		}
		eth_amount := amount
		token_amount := "0"
		is_token := 0
		contractaddr, ok := m["contract_addr"]
		if ok && len(contractaddr) > 1 {
			// token
			eth_amount = "0"
			token_amount = amount
			is_token = 1
		}
		txhash, err := web3conn.Eth.SendRawTransaction([]string{raw})
		// dispose error
		// db operation
		log.Error(err)
		// todo extra info table operation
		if err == nil {
			// Combining data
			st1 := eth_query.TxExtraInfo{
				From:        from,
				To:          to,
				Amount:      eth_amount,
				TokenAmount: token_amount,
				Comment:     note,
				Nonce:       nonce,
				TxHash:      txhash,
			}
			_, err = st1.InsertOneRaw()
			if err != nil {
				// db error
				log.Fatal(fee, from, to, note, nonce, gaslimit, gasprice, eth_amount, token_amount, is_token, txhash)
				log.Fatalf("db error: %s", err)
			}
			st2 := eth_query.Tx{
				From:      from,
				To:        to,
				Amount:    eth_amount,
				TxHash:    txhash,
				InputData: raw,
				Nonce:     nonce,
				GasLimit:  gaslimit,
				GasPrice:  gasprice,
				TxState:   -1,
				IsToken:   is_token,
				Fee:       fee,
			}
			st2.InsertOneRaw(&st2)
			if is_token == 1 {
				st3 := eth_query.TokenTx{
					From:         from,
					To:           to,
					ContractAddr: contractaddr,
					Amount:       token_amount,
					InputData:    raw,
					Nonce:        nonce,
					GasLimit:     gaslimit,
					GasPrice:     gasprice,
					Fee:          fee,
					IsToken:      is_token,
					TxHash:       txhash,
					TxState:      -1,
					Decimal:      dec,
				}

				st3.InsertOneRaw(&st3)
				// update token address
				eth_query.UpdateTokenAddress(token_amount, from, contractaddr)
			}

			// update address
			unconfirmAmount := AddString(eth_amount, fee)
			eth_query.UpdateAddress(unconfirmAmount, from)

			s.Data["json"] = resultResponseMake(txhash)
		} else {
			s.Data["json"] = resultResponseErrorMake(2009, err.Error())
		}
	default:
		//error
		s.Data["json"] = resultResponseErrorMake(2010, nil)
	}
	s.ServeJSON()
}
