//software: GoLand
//file: parity_websocket.go
//time: 2018/8/2 下午2:46

//wait:parity web socket don't support getTraction func
package eth

import (
	"encoding/json"
	"fmt"
	"github.com/gembackend/conf"
	"github.com/gorilla/websocket"
	"net/url"
)

var (
	subEthGetBalance  = `{"method":"parity_subscribe","params":["eth_getBalance",["%s","latest"]],"id":1,"jsonrpc":"2.0"}`
	subEthBlockNumber = `{"method":"parity_subscribe","params":["eth_blockNumber",[]],"id":1,"jsonrpc":"2.0"}`
)

var (
	recordSub        = make(map[string]string)
	addressRecordSub = make(map[string]string)
	unSub            = `{"method":"parity_unsubscribe","params":["%s"],"id":1,"jsonrpc":"2.0"}`
	parityUrl        = url.URL{Scheme: "ws", Host: conf.EthWebsocketUrl, Path: ""}
	contrast         = make(chan string)

	formatsUnSubString = func(str string) string {
		return fmt.Sprintf(unSub, str)
	}
)

func ParityWebSocketStart() {
	log.Info(parityUrl.String())
	c, r, err := websocket.DefaultDialer.Dial(parityUrl.String(), nil)
	if err != nil {
		log.Fatalf("dial: %s", err)
		return
	}
	log.Infof("response: %s", r)
	defer c.Close()
	defer c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

	wg.Add(1)
	go func(conn *websocket.Conn) {
		for true {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Error(err)
				continue
			}
			log.Debug(string(msg))
			// dispose msg
			disposed := disposemsg(msg)
			// save result id
			resultId, ok := disposed["result"]
			if ok {
				// save and go on sub
				recordSub[resultId.(string)] = <-contrast
				log.Debug(recordSub)
			} else if params, ok := disposed["params"]; ok {
				// dispose result
				disposeResult(params)
			}
		}
	}(c)
	subscribe(c)

	// may be naver run this
	defer unsubscribe(c)
	wg.Wait()
}
func disposeResult(params interface{}) {
	switch params.(type) {
	case map[string]interface{}:
		data := params.(map[string]interface{})
		method, ok := recordSub[data["subscription"].(string)]
		if ok {
			switch method {
			case "eth_blockNumber":
				log.Debug(data, method)
			case "eth_getBalance":
				log.Debug(data, method)
			case "eth_getTransaction":
			}
		}
	default:
		log.Warning("no support type")
	}
}

func disposemsg(bytes []byte) (result map[string]interface{}) {
	err := json.Unmarshal(bytes, &result)
	if err != nil {
		log.Error(err)
	}
	return
}
func unsubscribe(conn *websocket.Conn) {
	for k, _ := range recordSub {
		conn.WriteMessage(websocket.TextMessage, []byte(formatsUnSubString(k)))
	}
}
func subscribe(conn *websocket.Conn) {
	// block number
	conn.WriteMessage(websocket.TextMessage, []byte(subEthBlockNumber))
	contrast <- "eth_blockNumber"

	// setp 2---- chain by this format
	contrast <- "eth_getBalance"

	log.Info("subscribe finish")
}
