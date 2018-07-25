package main

import (
	_ "github.com/astaxie/beego/config/xml"
	"github.com/gembackend/messagequeue"
	"github.com/gembackend/conf"
	"github.com/gembackend/scripts"
	"flag"
	"github.com/gembackend/gembackendlog"
)

func main() {

	action := flag.String("action", "", "change a action")
	height := flag.Uint64("height", 5000000, "change start height")
	flag.Parse()
	log := gembackendlog.Log

	switch *action {
	case "eth-updater-web3":
		log.Info("eth-updater-web3 start")
		ethUpdaterWeb3(*height)
	case "eth-updater-web3-mul":
		log.Info("eth-updater-web3-mul start")
		scripts.StartEthupdaterMul(*height)
	case "eth-updater-ethscan":
		log.Info("eth-updater-ethscan start")
		ethUpdaterEthscan(*height)
	case "eth-kafka-script":
		log.Info("eth-kafka-script start")
		ethkafkascript()
	default:
		log.Info("no operation was selected")
		log.Info("you can select action")
		log.Info("eth-updater-web3")
		log.Info("eth-updater-web3-mul")
		log.Info("eth-updater-ethscan")
		log.Info("eth-kafka-script")
	}
}

// 处理eth外部钱包导入
func ethkafkascript() {
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

// 单线程更新web3更新程序
func ethUpdaterWeb3(height uint64) {
	updater := scripts.NewEthUpdaterWeb3(height)
	updater.Forever()
}

// 单线程ethscan接口更新程序
func ethUpdaterEthscan(height uint64) {
	updater := scripts.NewEthUpdaterApi(height)
	updater.Forever()
}
