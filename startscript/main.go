package main

import (
	_ "github.com/astaxie/beego/config/xml"
	"github.com/gembackend/messagequeue"
)

func main() {
	// 5个线程更新
	//scripts.StartEthupdaterMul(5000000)
	// 外部钱包导入
	//scripts.Main("555", "0xddbd2b932c763ba5b1b7ae3b362eac3e8d40121a")
	p := messagequeue.MakeProducer()
	messagequeue.MakeMessage("222nihao111", p)
	c := messagequeue.MakeConsumer()
	messagequeue.ReadMessage(c)
}
