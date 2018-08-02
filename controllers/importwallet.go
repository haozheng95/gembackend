package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/gembackend/conf"
	"github.com/gembackend/models/eth_query"
)

type ImportWalletController struct {
	beego.Controller
}

// 钱包导入
func (i *ImportWalletController) Post() {
	m := make(map[string]string)
	err := json.Unmarshal(i.Ctx.Input.RequestBody, &m)
	if err != nil {
		log.Error(err)
		i.Data["json"] = resultResponseErrorMake(2001, err.Error())
		i.ServeJSON()
		return
	}

	walletId, err1 := m["wallet_id"]
	sign, err2 := m["sign"]
	ethAddr, err3 := m["eth_addr"]

	if !err1 || !err2 || !err3 {
		log.Error(err1, err2, err3)
		i.Data["json"] = resultResponseErrorMake(2000, nil)
		i.ServeJSON()
		return
	}

	if !checkSign(walletId, sign) {
		log.Error("checkSign false")
		i.Data["json"] = resultResponseErrorMake(2007, nil)
		i.ServeJSON()
		return
	}

	if eth_query.GetEthAddrExist(ethAddr) {
		i.Data["json"] = resultResponseMake("import success")
		i.ServeJSON()
		return
	}

	// todo 将钱包信息加入kafaka队列
	ethkafka := map[string]interface{}{
		"walletId": walletId,
		"addr":     ethAddr,
	}
	ethkafkaparam, err := json.Marshal(ethkafka)
	if err != nil {
		i.Data["json"] = resultResponseErrorMake(2008, err.Error())
		i.ServeJSON()
		return
	}
	// save data for kafka
	ethtopicname := conf.KafkaimportEthTopicName
	SaveForKafka(ethtopicname, string(ethkafkaparam))
	// add eth address for kafka
	SaveForKafka(conf.KafkagetbalanceParityTopic, ethAddr)

	i.Data["json"] = resultResponseMake("import success! pleases! wait some time")
	i.ServeJSON()
}
