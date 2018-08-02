package scripts

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/gembackend/models/exchange"
	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/url"
	"strings"
	"time"
)

var (
	huobisubstr = `{"sub": "market.%s.kline.15min","id": "id10"}`
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
				log.Infof("msg-type: %d", msgtype)
				wg.Done()
				continue
			}
			//log.Infof("msg-type: %d", msgtype)
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
	log.Warning("wg == 0 !!@huobiwebsocket exit!!!!")
}

func huobisub(c *websocket.Conn) {

	for _, v := range exchange.GetAllTokenName() {
		substr := fmt.Sprintf(huobisubstr, strings.ToLower(v.TokenName)+baseCoin)
		//substr := fmt.Sprintf(huobisubstr, strings.ToLower("bat") + baseCoin)
		fmt.Println(substr, v)
		c.WriteMessage(websocket.TextMessage, []byte(substr))
	}
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
	//log.Debug(string(s))
	return
}

func huobidisposebyte(c *websocket.Conn, s []byte) {
	var dat map[string]interface{}
	var echo map[string]interface{}
	echo = make(map[string]interface{})
	err := json.Unmarshal(s, &dat)
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

	ch, ok := dat["ch"]
	if !ok {
		return
	}
	//log.Debugf("ch: %s", ch)
	tick, ok := dat["tick"]
	if !ok {
		return
	}

	//log.Debugf("tick: %s",tick)
	//log.Debug(dat)
	coin := splitch(ch)
	log.Debugf("coin: %s", coin)
	price := filttick(tick)
	st := new(exchange.EthToken)
	usdtcny := exchange.GetMainChainCny(baseCoinFullName)
	st.TokenName = strings.ToUpper(coin) //must to upper
	st.Cny = price.String()
	st.Usdt = DivString(st.Cny, usdtcny).String()
	st.UpdateCnyAndUsdt()
}

func splitch(ch interface{}) string {
	chs := ch.(string)
	s := strings.Split(chs, ".")
	return strings.Replace(s[1], baseCoin, "", -1)
}

func filttick(tick interface{}) (price decimal.Decimal) {
	tickmap := tick.(map[string]interface{})
	ticclose, ok := tickmap["close"]
	if !ok {
		//log.Debugf("close ---- %f", ticclose)
		log.Error("filttick error")
		log.Error(tick)
		price = decimal.New(0, 0)
		return
	}
	//log.Debug(ticclose.(float64))

	price = decimal.NewFromFloat(ticclose.(float64))
	return
}
