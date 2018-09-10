//software: GoLand
//file: btcutil.go
//time: 2018/8/21 下午12:13
package rpc

import (
	"bytes"
	"encoding/json"
	"github.com/gembackend/conf"
	"io/ioutil"
	"net/http"
)

func SentBtcRawTraction(raw string) (res []map[string]interface{}) {
	method := "sendrawtransaction"
	params := [][]interface{}{{method, raw}}
	body := Batch(params)
	//log.Debug(string(body))
	if err := json.Unmarshal(body, &res); err == nil {
		log.Debug(res)
		return
	} else {
		log.Warning(err)
	}
	return
}

func Batch(rpcCalls [][]interface{}) (body []byte) {
	batchData := make([]map[string]interface{}, len(rpcCalls))
	for i, v := range rpcCalls {
		batchData[i] = map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  v[0],
			"params":  v[1:],
			"id":      i + 1,
		}
	}
	//log.Debug(batchData)
	postData, err := json.Marshal(batchData)
	//log.Info(string(postData))

	if err != nil {
		log.Fatal(err)
		log.Fatalf("error data = %s", string(postData))
	}
	reader := bytes.NewReader(postData)
	url := "http://" + conf.BtcUser + ":" + conf.BtcPass + "@"
	url += conf.BtcHost + ":" + conf.BtcPort

	resp, err := http.Post(url, "application/json", reader)
	if err != nil {
		log.Fatalf("request errpr = %s", err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("get body error = %s", err)
		log.Fatal(string(body))
	}
	//log.Debug(string(body))
	return
}
