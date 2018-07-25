package main

import (
	_ "github.com/astaxie/beego/config/xml"
	"github.com/gembackend/messagequeue"
	"github.com/gembackend/conf"
	"github.com/gembackend/scripts"
)

func main() {
	// 5个线程更新
	//scripts.StartEthupdaterMul(5000000)
	// 外部钱包导入
	//scripts.Main("555", "0xddbd2b932c763ba5b1b7ae3b362eac3e8d40121a")
	//p := messagequeue.MakeProducer()
	//defer p.Close()
	//messagequeue.MakeMessage("testGo", "223eee2nihao111", p)

	//readmessage(pcs[0], func(){
	//	fmt.Println("okokok")
	//})
	readmessageForimportwallet()
}

// 处理eth外部钱包导入
func readmessageForimportwallet() {
	c := messagequeue.MakeConsumer()
	r := make(chan interface{})
	defer c.Close()
	defer close(r)
	go func(r <-chan interface{}) {
		for z := range r {
			t := z.(map[string]interface{})
			scripts.Main(t["walletId"].(string), t["addr"].(string))
		}
	}(r)

	pcs := messagequeue.MakePcs(c, conf.KafkaimportEthTopicName)
	// 优化点 处理多个分区 此版本只处理第一个分区
	messagequeue.ReadMessage(pcs[0], messagequeue.DisJsonFunc, r)
}
