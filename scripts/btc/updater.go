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
	"github.com/gembackend/rpc"
	"github.com/op/go-logging"
)

var (
	btcRpc *rpcclient.Client
	log    *logging.Logger
)

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
	//hashs = make([]chainhash.Hash, 0, 100)
	//log.Debug(len(separates))
	filtermap := make(map[interface{}]int)
	for _, v := range separates {
		if v == nil || v.vin == nil {
			//log.Fatal("fatal error for range separates")
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
	log.Debug("hashgrouop len=", len(hashgroup), "cap=", cap(hashgroup))
	return
}

func Main() {
	log.Debug("start")
	num, err := currnum()
	log.Debug("1")
	info, err := getBlockInfo(num)
	txs, err := info.TxHashes()
	log.Debug("2", "len txs === ", len(txs))
	// get curr block all transactions
	alltrans, err := getTxsDetail(txs[:10])
	log.Debug("3")
	//log.Debug(num, err, alltrans)
	alltransCut := separateTxsDet(alltrans)
	log.Debug("4")
	log.Debug(err)
	vinhashgroup := getvinhashs(alltransCut)
	//getvinhashs(alltransCut)
	log.Debug("5")
	// get spend tx
	spendtrans, err := getTxsDetail(vinhashgroup[:3])
	spendtransCut := separateTxsDet(spendtrans)
	log.Debug(len(spendtransCut))
	spendData := separateTxsData(spendtransCut)
	// get unspend tx
	unspendData := separateTxsData(alltransCut)

	// reduce
	datareduce(spendData, unspendData)
}

type reduceData struct {
	from, to interface{}
}

func datareduce(spendData, unspendData map[string]*cutdata) (data []*reduceData) {
	data = make([]*reduceData, len(unspendData))
	for k1, v1 := range unspendData {
		for k2, v2 := range v1.vout {

		}
	}
	return
}

type cutdata struct {
	vin  map[string]interface{}
	vout map[float64]map[string]interface{}
}

//vin ----  map[a6ebc95817c76c486dfa84e1b3ecd407e5feb313f730d3766a845d4699cbfe8a:0]
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
		//log.Debug("len vin debug === ", len(vin))
		if resdata[txid] == nil {
			resdata[txid] = new(cutdata)
		}
		resdata[txid].vin = make(map[string]interface{})
		resdata[txid].vout = make(map[float64]map[string]interface{})
		for _, v1 := range vin {
			v2 := v1.(map[string]interface{})
			if v2["txid"] == nil {
				log.Debug("coinbase ---")
				continue
			}
			delete(v2, "scriptSig")
			delete(v2, "txinwitness")
			delete(v2, "sequence")
			//resdata[txid].vin = append(resdata[txid].vin, v2)

			resdata[txid].vin[v2["txid"].(string)] = v2["vout"]
		}
		log.Debug("vin ---- ", resdata[txid].vin)
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
