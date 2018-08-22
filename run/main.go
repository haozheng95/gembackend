package main

import (
	"flag"
	"fmt"
	_ "github.com/astaxie/beego/config/xml"
	"github.com/gembackend/conf"
	"github.com/gembackend/gembackendlog"
	"github.com/gembackend/messagequeue"
	"github.com/gembackend/models"
	_ "github.com/gembackend/models"
	"github.com/gembackend/scripts"
	"github.com/gembackend/scripts/btc"
	"github.com/gembackend/scripts/eth"
	"os"
	"os/signal"
)

func main() {
	log := gembackendlog.Log
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func(interrupt chan os.Signal) {
		defer close(interrupt)
		msg := <-interrupt
		log.Warningf("exit message %s", msg)
		log.Warning("program exit....... ")
		os.Exit(0)
	}(interrupt)
	//createtestdata()
	//eth.StartEthApiMul(5000000)
	//createtestdata()
	//eth.StartEthupdaterMul(6117492)
	btc.Main()
	os.Exit(0)

	action := flag.String("action", "", "change a action")
	height := flag.Uint64("height", 5000000, "change start height")
	flag.Parse()

	switch *action {
	case "eth-updater-web3":
		log.Info("eth-updater-web3 start")
		ethUpdaterWeb3(*height)
	case "eth-updater-web3-mul":
		log.Info("eth-updater-web3-mul start")
		eth.StartEthupdaterMul(*height)
	case "eth-updater-ethscan":
		log.Info("eth-updater-ethscan start")
		ethUpdaterEthscan(*height)
	case "eth-kafka-script":
		log.Info("eth-kafka-script start")
		ethkafkascript()
	case "feixiaohao":
		log.Info("feixiaohao start")
		feixiaohaoapi()
	case "huobiwebsocket":
		log.Info("huobiwebsocket start")
		huobiwebsocket()
	case "createtestdata":
		log.Info("createtestdata start")
		createtestdata()
	case "auxiliary":
		log.Info("auxiliary start")
		auxiliary()
	case "txMonitoring":
		log.Info("txMonitoring start")
		txMonitoring()
	default:
		log.Info("no operation was selected")
		log.Info("you can select action")
		log.Info("eth-updater-web3")
		log.Info("eth-updater-web3-mul")
		log.Info("eth-updater-ethscan")
		log.Info("eth-kafka-script")
		log.Info("feixiaohao")
		log.Info("createtestdata")
		log.Info("huobiwebsocket")
		log.Info("auxiliary")
	}
}

// 区块浏览器模式导入交易记录
func txMonitoring() {
	c := messagequeue.MakeConsumer()
	r := make(chan interface{})
	defer c.Close()
	defer close(r)

	go func(r <-chan interface{}) {
		for z := range r {
			t := z.(map[string]interface{})
			is_token := false
			if t["is_token"].(float64) > 0 {
				is_token = true
			}
			fmt.Println("txMonitoring, ----", z)
			go eth.Monitoring(t["hash"].(string), is_token)
		}
	}(r)
	pcs := messagequeue.MakePcs(c, conf.KafkaSendRawTopic)
	messagequeue.ReadMessage(pcs[0], messagequeue.DisJsonFunc, r)
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
			eth.Main(t["walletId"].(string), t["addr"].(string))
		}
	}(r)

	pcs := messagequeue.MakePcs(c, conf.KafkaimportEthTopicName)
	// 优化点 处理多个分区 此版本只处理第一个分区
	messagequeue.ReadMessage(pcs[0], messagequeue.DisJsonFunc, r)
}

// 单线程更新web3更新程序
func ethUpdaterWeb3(height uint64) {
	updater := eth.NewEthUpdaterWeb3(height)
	updater.Forever()
}

// 单线程ethscan接口更新程序
func ethUpdaterEthscan(height uint64) {
	updater := eth.NewEthUpdaterApi(height)
	updater.Forever()
}

// 火币websocket启动
func huobiwebsocket() {
	scripts.Huobiwebsocker()
}

// 非小号获取价格启动
func feixiaohaoapi() {
	scripts.FeixiaohaoStart()
}

//创建测试数据
func createtestdata() {
	//models.CreateTable()
	models.AutoInsertData("exchange", "eth_token")
	models.AutoInsertData("exchange", "main_chain")
	models.AutoInsertData("eth_query", "token_tx")
	models.AutoInsertData("eth_query", "tx")
	models.AutoInsertData("eth_query", "tx_extra_info")
	models.AutoInsertData("eth_query", "address")
}

//gaslimit和gasprice更新
func auxiliary() {
	scripts.AuxiliaryMain()
}
