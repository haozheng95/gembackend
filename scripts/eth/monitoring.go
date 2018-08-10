//software: GoLand
//file: monitoring.go
//time: 2018/8/10 上午10:34
package eth

import (
	"github.com/gembackend/models/eth_query"
	"github.com/gembackend/rpc"
	"math"
	"strconv"
	"time"
)

var (
	ethdec    = strconv.FormatFloat(math.Pow(10, 18), 'f', -1, 64)
	debugMole = true
)

func Monitoring(hash string, isToken bool) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	hash = "0x6d5536df432f8b7e9ffeecfc91a82cca611015236f5df22c8279a63a8eed287f"
	for range ticker.C {
		ethscanResponse, err := rpc.Eth_getTransactionReceipt(hash)
		//log.Debug(ethscanResponse)
		if err != nil {
			log.Error(err)
			continue
		}
		r, err := rpc.FormatResponse(&ethscanResponse)
		if err != nil || r.Result == nil {
			log.Error(err, ethscanResponse)
			continue
		}
		dbtx, err := eth_query.GetTxOneRawByHash(hash)
		if err != nil && !debugMole {
			log.Error(err)
			time.Sleep(time.Second)
			return
		}
		from, to, status, gasUsed, blockHash, blockNumber, logs := fetch(r)
		fee := MulString(gasUsed, dbtx.GasPrice)
		fee = DivString(fee, ethdec)
		log.Debugf("from=%s", from)
		log.Debugf("to=%s", to)
		log.Debugf("status=%s", status)
		log.Debugf("gasUsed=%s", gasUsed)
		log.Debugf("blockHash=%s", blockHash)
		log.Debugf("blockNumber=%s", blockNumber)
		log.Debugf("logs=%v", logs)
		log.Debugf("fee=%s", fee)
		log.Debugf("gasPrice=%s", dbtx.GasPrice)
		log.Debug(ethscanResponse)
		eth_query.UpdateTxOneRawByHash(gasUsed, fee, status, blockHash, blockNumber, hash)
		UpdateEthAccount(from, to)

		if isToken {
			tokenSts := eth_query.GetTokenTxinfo(hash)
			if len(tokenSts) < 1 {
				log.Error("no this token tx, txhash=", hash)
			} else {
				tokenSt := tokenSts[0]
				confirmTime := strconv.FormatInt(time.Now().Unix(), 10)
				created := strconv.FormatInt(tokenSt.Created.Unix(), 10)
				log.Debug(tokenSt.Created)
				tokenDecimal := tokenSt.Decimal
				deleteTokenTx := false
				for _, v1 := range logs {
					v2 := v1.(map[string]interface{})
					topics1, ok := v2["topics"]
					if ok {
						topics2 := topics1.([]interface{})
						if len(topics2) == 3 && topics2[0].(string) == _TRANSACTION_TOPIC {
							logindex := HexDec(v2["logIndex"].(string))
							err = eth_query.InsertTokenTx(tokenSt.From, tokenSt.To, tokenSt.Amount,
								tokenSt.InputData, tokenSt.Nonce, tokenSt.GasLimit, tokenSt.GasPrice,
								gasUsed, fee, hash, blockHash, confirmTime,
								created, "1", "1", logindex, tokenSt.ContractAddr, tokenDecimal)
							if err == nil {
								deleteTokenTx = true
							} else {
								deleteTokenTx = false
							}
						}
					}
				}
				if deleteTokenTx {
					eth_query.DeleteTokenTxbyHashWhere(hash)
				}
				// update token account
				if eth_query.GetEthAddrExist(tokenSt.From) || debugMole {
					UpdateTokenAccount(tokenSt.From[2:], tokenSt.ContractAddr, tokenSt.Decimal)
				}
				if eth_query.GetEthAddrExist(tokenSt.To) || debugMole {
					UpdateTokenAccount(tokenSt.To[2:], tokenSt.ContractAddr, tokenSt.Decimal)
				}
			}
		}

		return
	}
}
func UpdateEthAccount(s string, s2 string) {
keep:
	res := rpc.Eth_getMulBalance(s, s2)
	log.Debug(res)
	r, err := rpc.FormatResponseMap(&res)
	data, ok := r["result"]
	if err != nil || !ok {
		log.Debug(err)
		log.Error(r)
		goto keep
	}
	data2 := data.([]interface{})
	for _, v := range data2 {
		v1 := v.(map[string]interface{})
		addr, ok := v1["account"]
		balance, ok1 := v1["balance"]
		if !ok || !ok1 {
			log.Error(data2)
			continue
		}
		addr1 := addr.(string)
		if eth_query.GetEthAddrExist(addr1) || debugMole {
		keep1:
			res1, err := rpc.Eth_getTransactionCount(addr1)
			if err != nil {
				time.Sleep(time.Second)
				goto keep1
			}
			res2, err := rpc.FormatResponseMap(&res1)
			if err != nil {
				time.Sleep(time.Second)
				goto keep1
			}
			res3, ok := res2["result"]
			if !ok {
				time.Sleep(time.Second)
				goto keep1
			}
			nonce := HexDec(res3.(string))
			amount := DivString(balance.(string), ethdec)
			//log.Debug(nonce, amount)
			err = eth_query.UpdateAddressAmount("0", addr1, nonce, amount)
			if err != nil {
				log.Error(err)
				time.Sleep(time.Second)
				goto keep1
			}
		}
	}
}

func UpdateTokenAccount(s string, s2 string, s3 string) {
	dec, err := strconv.Atoi(s3)
	if err != nil {
		log.Errorf("strconv.Atoi(s3) error := %s", err)
		return
	}
	s3 = strconv.FormatFloat(math.Pow(10, float64(dec)), 'f', -1, 64)
keep:
	response, _ := rpc.Eth_getTokenBalance(s2, s)
	log.Debug("token Balance response:=", response)
	res, err := rpc.FormatTokenResponse(response)
	if err != nil || res.Result == "" {
		time.Sleep(time.Second * 5)
		goto keep
	}
	amount := DivString(HexDec(res.Result), s3)
	err = eth_query.UpdateTokenAddressAmount(amount, "0", s, s2)
	if err != nil {
		log.Warning(err)
		time.Sleep(time.Second * 5)
		goto keep
	}
}

func fetch(r *rpc.Response) (from, to, status, gasUsed, blockHash, blockNumber string, logs []interface{}) {
	for k, v := range r.Result {
		switch k {
		case "from":
			from = v.(string)
		case "to":
			to = v.(string)
		case "status":
			status = HexDec(v.(string))
		case "gasUsed":
			gasUsed = HexDec(v.(string))
		case "blockHash":
			blockHash = v.(string)
		case "blockNumber":
			blockNumber = HexDec(v.(string))
		case "logs":
			logs = v.([]interface{})
		}
	}
	return
}
