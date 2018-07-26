package scripts

import (
	"github.com/gorilla/websocket"
	"time"
	"net/url"
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"fmt"
	"encoding/json"
)

var (
	huobisubstr = `{"sub": "market.%s.kline.1min","id": "id10"}`
)

func Huobiwebsocker() {
	u := url.URL{Scheme: "wss", Host: "api.huobi.pro", Path: "/ws"}
	log.Infof("connecting to %s", u.String())
	// 在这里可以设置代理
	c, r, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("dial: %s", err)
	}
	wg.Add(2)
	go func() {
		for true {
			msgtype, msg, err := c.ReadMessage()
			if err != nil {
				log.Errorf("read error: %s", err)
				wg.Done()
				continue
			}
			log.Infof("msg-type: %d", msgtype)
			dat := huobiunzipmsg(msg)
			huobidisposebyte(c, dat)
		}
	}()

	log.Infof("response: %s", r)
	defer c.Close()
	defer c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

	// keep heart
	ticker := time.NewTicker(time.Second * 15)
	defer ticker.Stop()
	go func() {
		echo := make(map[string]interface{})
		for t := range ticker.C {
			echo["ping"] = t.Unix()
			jsonBytes, _ := json.Marshal(echo)
			//fmt.Println(t, string(jsonBytes))
			c.WriteMessage(websocket.TextMessage, jsonBytes)
		}
	}()

	// sub
	huobisub(c)
	wg.Wait()
	log.Warning("huobiwebsocket exit!!!!")
}

func huobisub(c *websocket.Conn) {
	substr := fmt.Sprintf(huobisubstr, "btcusdt")
	fmt.Println(substr)
	c.WriteMessage(websocket.TextMessage, []byte(substr))
}

func huobiunzipmsg(msg []byte) (s []byte) {
	rdata := bytes.NewReader(msg)
	r, err := gzip.NewReader(rdata)
	if err != nil {
		log.Errorf("huobiunzipmsg setp-1 error: %s", err)
		return
	}
	s, err = ioutil.ReadAll(r)
	if err != nil {
		log.Errorf("huobiunzipmsg setp-2 error: %s", err)
		return
	}
	return
}

func huobidisposebyte(c *websocket.Conn, s []byte) {
	var dat map[string]interface{}
	var echo map[string]interface{}
	echo = make(map[string]interface{})
	err := json.Unmarshal([]byte(s), &dat)
	if err != nil {
		return
	}
	ping, ok := dat["ping"]
	if ok {
		echo["pong"] = ping
		jsonBytes, _ := json.Marshal(echo)
		//fmt.Println(string(jsonBytes))
		c.WriteMessage(websocket.TextMessage, jsonBytes)
		return
	}
	pong, ok := dat["pong"]
	if ok {
		echo["ping"] = pong
		jsonBytes, _ := json.Marshal(echo)
		//fmt.Println(string(jsonBytes))
		c.WriteMessage(websocket.TextMessage, jsonBytes)
		return
	}

}
