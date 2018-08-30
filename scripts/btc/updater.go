//software: GoLand
//file: updater.go
//time: 2018/8/21 下午2:18
package btc

import (
	"encoding/json"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/gembackend/gembackendlog"
	"github.com/gembackend/models/btc_query"
	"github.com/gembackend/rpc"
	"github.com/op/go-logging"
	"time"
)

var (
	btcRpc *rpcclient.Client
	log    *logging.Logger
)

const interval = 50
const beginHeight = 538690
const sleeptime = 30

func init() {
	btcrpc, ok := rpc.ConnectMap["btc-conn"]
	if ok && btcrpc != nil {
		btcRpc = btcrpc.(*rpcclient.Client)
	} else {
		btcRpc = rpc.ReMakeBtcConn()
	}
	log = gembackendlog.Log
}

func currnum() (num int64, err error) {
	num, err = btcRpc.GetBlockCount()
	if err != nil {
		log.Fatalf("get curr num error: %s", err)
	}
	return
}

func getBlockInfo(num int64) (block *wire.MsgBlock, err error) {
	hash, err := btcRpc.GetBlockHash(num)
	if err != nil {
		log.Fatalf("get block hash error %s", err)
		return
	}
	block, err = btcRpc.GetBlock(hash)
	if err != nil {
		log.Fatalf("get block error %s", err)
		//block.TxHashes()
	}
	return
}

func getTxsDetail(hashs []chainhash.Hash) (response []map[string]interface{}, err error) {
	params := make([][]interface{}, len(hashs))
	method := "getrawtransaction"
	detail := 1
	for i, v := range hashs {
		params[i] = []interface{}{method, v.String(), detail}
	}
	//log.Debug(params)
	body := rpc.Batch(params)
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("get tx detail error %s", err)
		log.Fatal("==========", string(body))
	}
	return
}

type separate struct{ blockhash, blocktime, txid, hash, vout, vin interface{} }

func separateTxsDet(response []map[string]interface{}) (result1 []*separate) {
	//log.Debug("=====", len(response))
	result1 = make([]*separate, len(response))
	i := 0
	for _, v := range response {
		if v == nil || v["result"] == nil {
			continue
		}
		result := v["result"].(map[string]interface{})
		blockhash := result["blockhash"]
		blocktime := result["blocktime"]
		txid := result["txid"]
		hash := result["hash"]
		vout := result["vout"]
		vin := result["vin"]
		result1[i] = &separate{blockhash, blocktime, txid, hash, vout, vin}
		i++
		//log.Debug("blockhash=", blockhash)
		//log.Debug("blocktime=", blocktime)
		//log.Debugf("txid=%s", txid)
		//log.Debugf("hash=%s", hash)
		//log.Debugf("vout=%s", vout)
		//log.Debugf("vin=%s", vin)
	}
	return
}

func getvinhashs(separates []*separate) (hashgroup []chainhash.Hash) {
	//hashs := make([]chainhash.Hash, 0, 100)
	//log.Debug(len(separates))
	filtermap := make(map[interface{}]int)
	for _, v := range separates {
		if v == nil || v.vin == nil {
			log.Fatal("fatal error for range separates")
			continue
		}
		vin := v.vin.([]interface{})
		//log.Debug("hashgroup size=", len(hashgroup), "len(vin)=", len(vin))
		for _, v1 := range vin {
			//log.Debug("v1 ====", v1)
			v2 := v1.(map[string]interface{})
			if v2["txid"] != nil {
				filtermap[v2["txid"]] = filtermap[v2["txid"]] + 1
				//hashs = append(hashs, makeChHash(v2["txid"].(string)))
				//log.Debug("txid = ", hashs[i].String())
				//log.Debug("txid = <-->", v2["txid"].(string))
			}
		}
	}
	log.Debug("hashgroup---------------------------")
	hashgroup = make([]chainhash.Hash, 0, len(filtermap))
	for k, v := range filtermap {
		log.Debug(k, "====", v)
		hashgroup = append(hashgroup, makeChHash(k.(string)))
	}
	//hashgroup = hashs
	log.Debug("hashgrouop len=", len(hashgroup), "cap=", cap(hashgroup))
	return
}

func checkheight(num int64) int64 {
again:
	block, _ := currnum()
	db := btc_query.CurrBlockNum()
	log.Debug("db height     === ", db)
	log.Debug("block height  === ", block)
	log.Debug("update height === ", num)

	if db < num {
		db = num
	}

	if db < block {
		num = db + 1
	} else if num > block {
		log.Debug("pending start ========")
		time.Sleep(sleeptime)
		log.Debug("pending over =========")
		goto again
	}
	// debug
	//num = 538694
	return num
}

func checkblockhash(num int64) *wire.MsgBlock {
	for {
		info, _ := getBlockInfo(num)
		previousHash := info.Header.PrevBlock.String()
		num--
		dbhash := btc_query.Getblockhash(num)
		if dbhash == "" {
			return info
		}
		if previousHash == dbhash {
			return info
		} else {
			log.Debug("roll back height ==== ", num)
			btc_query.Deleteblockhash(dbhash)
		}
	}

}

func Main() {
	for {
		//log.Debug(Decimal(152.5419366900003))
		//return
		start(beginHeight)
		log.Debug("=========> over <===========")
	}
}

func start(begin int64) {
	log.Debug("start")

	var (
		blockHash    string
		blockNumber  int64
		previousHash string
		confirmTime  int64
		blockNonce   uint32
	)

	num := checkheight(begin)
	log.Debug("1")
	log.Debug("------", num)

	info := checkblockhash(num)
	// block info -----
	blockHash = info.Header.BlockHash().String()
	blockNumber = num
	previousHash = info.Header.PrevBlock.String()
	confirmTime = info.Header.Timestamp.Unix()
	blockNonce = info.Header.Nonce
	log.Debug("block hash          ====", blockHash)
	log.Debug("block number        ====", blockNumber)
	log.Debug("block previous hash ====", previousHash)
	log.Debug("confirm time        ====", confirmTime)
	log.Debug("block nonce         ====", blockNonce)
	// ---------------
	//btcBlock := NewBlockBtc(height int64, blockhash string, prehash string, confirm int64, nonce uint32)
	btcBlock := btc_query.NewBlockBtc(blockNumber, blockHash, previousHash, confirmTime, blockNonce)
	btcBlock.Insert()
	// ---------------
	txs, err := info.TxHashes()

	log.Debug("2", "len txs === ", len(txs))
	// get curr block all transactions
	ltxs := len(txs)
	alltrans := make([]map[string]interface{}, 0, ltxs)
	if ltxs > interval {
		i := 0
		end := interval
		for i < ltxs {
			log.Debug("block alltrans start ===", i)
			log.Debug("block alltrans end   ===", end)
			tslice, _ := getTxsDetail(txs[i:end])
			alltrans = append(alltrans, tslice...)
			i += interval
			end += interval
			if end < ltxs+1 {
				// nothing
			} else {
				end = ltxs
			}
		}
	} else {
		alltrans, err = getTxsDetail(txs)
	}
	log.Debug("3")
	//log.Debug(num, err, alltrans)
	alltransCut := separateTxsDet(alltrans)
	log.Debug("4")
	log.Debug(err)
	vinhashgroup := getvinhashs(alltransCut)
	//getvinhashs(alltransCut)
	// get spend tx
	lvintxs := len(vinhashgroup)
	log.Debug("5", "lvintxs ======", lvintxs)
	spendtrans := make([]map[string]interface{}, 0, lvintxs)
	if lvintxs > interval {
		vinstart := 0
		vinend := interval
		for vinstart < lvintxs {
			log.Debug("vinhashgroup start ====", vinstart)
			log.Debug("vinhashgroup end   ====", vinend)
			tslice, _ := getTxsDetail(vinhashgroup[vinstart:vinend])
			spendtrans = append(spendtrans, tslice...)
			vinstart += interval
			vinend += interval
			if vinend > lvintxs {
				vinend = lvintxs
			}
		}
	} else {
		spendtrans, err = getTxsDetail(vinhashgroup)
	}

	spendtransCut := separateTxsDet(spendtrans)
	spendData := separateTxsData(spendtransCut)
	// get unspend tx
	unspendData := separateTxsData(alltransCut)

	// reduce
	result := datareduce(spendData, unspendData)
	//for i, v := range result {
	//	log.Debug(i)
	//	log.Debug("txid ==== ", v.txid)
	//	log.Debug("from ==== ", v.from)
	//	log.Debug("to   ==== ", v.to)
	//}
	// last : format and db operation
	formatreduce(result, blockHash, blockNumber, confirmTime)
}

func formatreduce(data []*reduceData, blockhash string, blocknum, confirmTime int64) {
	unspentvouts := make([]*btc_query.UnspentVout, 0, len(data))
	tradeCollections := make([]*btc_query.TradeCollection, 0, len(data))
	tradingParticulars := make([]*btc_query.TradingParticulars, 0, len(data))
	allAddress := make([]string, 0, len(data))

	for _, v := range data {
		sw := false
		vout := v.to.(map[float64]map[string]interface{})
		txid := v.txid
		// db delete
		btc_query.Deletetxrecord(txid)
		// --------------------
		var amount float64 = 0
		var faddresses []string
		var vin []interface{}
		if v.from == nil {
			goto vinend
		}
		vin = v.from.([]interface{})
		faddresses = make([]string, 0, len(vin))
		for _, v := range vin {
			v1 := v.(map[string]interface{})
			amount += v1["value"].(float64)
			vintxid := v1["txid"].(string)
			index := v1["n"].(float64)

			// db operation
			btc_query.UpdateUnspentVout(vintxid, floatToString(index), txid)
			// make from addresses
			addresses := v1["addresses"].([]interface{})
			for _, v2 := range addresses {
				addr := v2.(string)
				if btc_query.GetBtcAddrExist(addr) {
					sw = true
					faddresses = append(faddresses, addr)
				}
			}
		}
	vinend:
		formatVout := make([]interface{}, 0, len(vout))
		taddresses := make([]string, 0, len(vout))
		var tovalue float64 = 0
		for _, v1 := range vout {
			//"txid", "index", "value", "address"
			formatVout = append(formatVout, v1)
			v2 := v1
			value := v2["value"].(float64)
			tovalue += value
			n := v2["n"].(float64)
			var addresses []interface{}
			if v2["addresses"] != nil {
				addresses = v2["addresses"].([]interface{})
			}

			for _, v3 := range addresses {
				addr := v3.(string)
				if btc_query.GetBtcAddrExist(addr) {
					taddresses = append(taddresses, addr)
					unspentvouts = append(unspentvouts, NewUnSpentVoutSt(txid, addr, floatToString(value), blockhash, blocknum, int64(n)))
					sw = true
				}
				//btc_query.InsertUnspentVout(txid, floatToString(n), floatToString(value), addr)
			}
		}

		if !sw {
			continue
		}
		amount = Decimal(amount)
		tovalue = Decimal(tovalue)
		fee := Subfloat(amount, tovalue)
		if len(vin) == 0 {
			log.Debug("coinbase vin ====", vin)
			fee = "0"
		}
		log.Debug("from value ====", amount)
		log.Debug("to   value ====", tovalue)
		log.Debug("fee  value ====", fee)
		log.Debug("txid value ====", txid)
		if tovalue > amount && len(vin) != 0 {
			log.Info("error ========")
			log.Debug("vin  ========", vin)
			log.Debug("vout ========", vout)
			panic("error ======")
		}
		//NewTradeCollections(output, input, blockhash, txid, addr, fee string, height, confirmtime int64, pay int)
		closure := eachaddress(floatToString(tovalue), floatToString(amount), blockhash, txid, fee, blocknum, confirmTime)
		tradeCollections = append(tradeCollections, closure(faddresses, 1)...)
		tradeCollections = append(tradeCollections, closure(taddresses, 0)...)
		// mk tx
		//from1 []string, to1 []string, txid, blockhash, Confirm  , blocknum
		//totalinput, totaloutput, fee string, vin, vout interface{},
		tradingParticulars = append(tradingParticulars,
			MakeTradingParticulars(floatToString(amount), floatToString(tovalue),
				fee, vin, formatVout, faddresses, taddresses, txid, blockhash, confirmTime, blocknum))

		// all address
		allAddress = append(allAddress, faddresses...)
		allAddress = append(allAddress, taddresses...)
	}
	// db
	err := btc_query.InsertMulUnspentVout(unspentvouts)
	if err == nil {
		err = btc_query.InsertMulTradeCollection(tradeCollections)
	} else {
		log.Fatalf("db step1 error %s", err)
	}
	if err == nil {
		err = btc_query.InsertMulTradingParticulars(tradingParticulars)
	} else {
		log.Fatalf("db step2 error == %s", err)
	}

	AccountUpdateMul(allAddress)
	log.Debug("related addr len ===", len(allAddress))
}

type reduceData struct {
	from, to interface{}
	txid     string
}

func datareduce(spendData, unspendData map[string]*cutdata) (data []*reduceData) {
	data = make([]*reduceData, 0, len(unspendData))
	for k1, v1 := range unspendData {
		st := new(reduceData)
		st.txid = k1
		st.to = v1.vout

		tslice := make([]interface{}, 0, len(v1.vin))
		for k2, v2 := range v1.vin {
			tvin, ok := spendData[k2]
			//log.Debug("k2 ==========", k2)
			if ok {
				//tvout, ok := tvin.vout[v2.(float64)]
				for _, v3 := range v2 {
					tvout, ok := tvin.vout[v3.(float64)]
					//log.Debug("vin ==== ok === ", ok)
					//log.Debug("================", v3)
					if ok {
						tvout["txid"] = k2 // add vin txid
						tslice = append(tslice, tvout)
					}
				}

			}
		}
		if len(tslice) > 0 {
			st.from = tslice
		}
		data = append(data, st)
	}

	return
}

type cutdata struct {
	vin  map[string][]interface{}
	vout map[float64]map[string]interface{}
}

//vin ----  map[a6ebc95817c76c486dfa84e1b3ecd407e5feb313f730d3766a845d4699cbfe8a:[0,1]]
//vout==== map[0:map[value:0.48 n:0 addresses:[3KEJwUXK7eXyTSdG5VbffTa1v6KuHW3AXT]]]
func separateTxsData(separates []*separate) (resdata map[string]*cutdata) {
	resdata = make(map[string]*cutdata)
	for _, v := range separates {
		if v == nil {
			continue
		}
		vin := v.vin.([]interface{})
		vout := v.vout.([]interface{})
		txid := v.txid.(string)
		log.Debug("len vout debug === ", len(vout))
		if resdata[txid] == nil {
			resdata[txid] = new(cutdata)
		}
		resdata[txid].vin = make(map[string][]interface{})
		resdata[txid].vout = make(map[float64]map[string]interface{})
		for _, v1 := range vin {
			v2 := v1.(map[string]interface{})
			log.Debug(v2)
			if v2["txid"] == nil {
				log.Debug("coinbase ---")
				continue
			}
			delete(v2, "scriptSig")
			delete(v2, "txinwitness")
			delete(v2, "sequence")
			//resdata[txid].vin = append(resdata[txid].vin, v2)

			resdata[txid].vin[v2["txid"].(string)] = append(resdata[txid].vin[v2["txid"].(string)], v2["vout"])
		}
		//log.Debug("vin ---- ", resdata[txid].vin)
		for _, v3 := range vout {
			v4 := v3.(map[string]interface{})
			v5 := v4["scriptPubKey"].(map[string]interface{})
			addresses, ok := v5["addresses"]
			datamap := make(map[string]interface{})
			datamap["value"] = v4["value"]
			datamap["n"] = v4["n"]
			if ok {
				//log.Debug("addresses ----", addresses)
				datamap["addresses"] = addresses
			}
			//resdata[txid].vout = append(resdata[txid].vout, datamap)
			if resdata[txid].vout[v4["n"].(float64)] == nil {
				resdata[txid].vout[v4["n"].(float64)] = datamap
			}
			//log.Debug("====", v4)
			//log.Debug("vout====", resdata[txid].vout)
		}
	}
	return
}
