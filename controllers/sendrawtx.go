package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/gembackend/conf"
	"github.com/gembackend/models/btc_query"
	"github.com/gembackend/models/eth_query"
	"github.com/gembackend/rpc"
	"github.com/gembackend/scripts/btc"
	"github.com/shopspring/decimal"
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

	//rpc.MakeConn()
	switch cointype {
	case "eth":
		// error assert
		if strings.Compare(from, to) == 0 {
			s.Data["json"] = resultResponseErrorMake(2011, nil)
			s.ServeJSON()
			return
		}

		var conn interface{}
		var ck bool
		var web3conn *rpc.Web3
		// to lower
		from = strings.ToLower(from)
		to = strings.ToLower(to)

		if conf.RunMode == "node" {
			conn, ck = rpc.ConnectMap["eth-web3"]
			if ck {
				web3conn = conn.(*rpc.Web3)
			} else {
				web3conn = rpc.ReMakeWeb3Conn()
			}
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
		var txhash string
		var err error
		if conf.RunMode == "node" {
			txhash, err = web3conn.Eth.SendRawTransaction([]string{raw})
		} else {
			txhash = rpc.Eth_sendRawTransaction(raw)
			if txhash == "" || len(txhash) == 0 {
				err = errors.New("send error")
			}
		}
		// dispose error
		// db operation
		log.Error(err)
		// todo extra info table operation
		if err == nil {
			// Save for kafka
			// wait:parity web socket don't support getTraction func
			if conf.RunMode != "node" {
				ethtopicname := conf.KafkaSendRawTopic
				b, _ := json.Marshal(map[string]interface{}{
					"hash":     txhash,
					"is_token": is_token,
				})
				fmt.Println(string(b))
				SaveForKafka(ethtopicname, string(b))
			}
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
	case "btc":
		//vinstr, _ := m["vin"]
		//voutstr, _ := m["vout"]
		//change, _ := m["change"]
		//var btcRpc *rpcclient.Client
		//btcrpc, ok := rpc.ConnectMap["btc-conn"]
		//if ok && btcrpc != nil {
		//	btcRpc = btcrpc.(*rpcclient.Client)
		//} else {
		//	btcRpc = rpc.ReMakeBtcConn()
		//}
		//btcRpc.CreateRawTransaction()
	default:
		//error
		s.Data["json"] = resultResponseErrorMake(2010, nil)
	}
	s.ServeJSON()
}
func decodeVinStr(vinstr string) (result []*btc_query.UnspentVout, err error) {
	result = make([]*btc_query.UnspentVout, 0, 10)
	if err = json.Unmarshal([]byte(vinstr), &result); err != nil {
		log.Fatal("decode vin str error ====", err)
	}
	return
}

//NewTradeCollections(output, input, blockhash, txid, addr, fee string, height, confirmtime int64, pay int, value float64) (st *btc_query.TradeCollection) {
func conversionunspent(vin []*btc_query.UnspentVout, txhash, in, out, fee string, pay int) []*btc_query.TradeCollection {
	result := make([]*btc_query.TradeCollection, 0, len(vin))
	for _, unspent := range vin {
		value, _ := decimal.NewFromString(unspent.Value)
		fvalue, _ := value.Float64()
		temp := btc.NewTradeCollections(out, in, "", txhash, unspent.Address, fee, 0, 0, pay, fvalue)
		result = append(result, temp)
	}
	return result
}

func TestDecodeVinStr() {
	testStr := `[ { "Id": 20674, "Txid": "b67c084194190c7c560ef6ba43f3877b10cce6098bcc353e2f1ef5868b10e8ed", "Spent": 0, "SpentTxid": "", "Index": 0, "Value": "0.0033999999999999998105681964233326652902178466320037841796875", "Address": "1Bd1vnozJKtBkVM1CFLWbXvb2AudTmPY3U", "Updated": "2018-08-28T03:16:35+08:00", "BlockHash": "0000000000000000005c9959b3216f8640f94ec96edea69fe12ad7dee8b74e92", "Height": 500001 }, { "Id": 20679, "Txid": "9af162a777bbbaf95c6afed6f05a1fc78cdae3b2868e516b7c9bbf8751b1b402", "Spent": 0, "SpentTxid": "", "Index": 1, "Value": "0.321171649999999975211295577537384815514087677001953125", "Address": "17WYBJEpR3KiwRRVTAtvCCTiSsRHFgCc9c", "Updated": "2018-08-28T03:16:35+08:00", "BlockHash": "0000000000000000005c9959b3216f8640f94ec96edea69fe12ad7dee8b74e92", "Height": 500001 } ]`
	decodeVinStr(testStr)
}
