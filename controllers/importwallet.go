package controllers

import (
	"github.com/astaxie/beego"
	"encoding/json"
	"github.com/gembackend/models"
	"github.com/gembackend/messagequeue"
	"github.com/gembackend/conf"
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

	if models.GetEthAddrExist(ethAddr) {
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

	ethtopicname := conf.KafkaimportEthTopicName
	p := messagequeue.MakeProducer()
	defer p.Close()
	messagequeue.MakeMessage(ethtopicname, string(ethkafkaparam), p)
	i.Data["json"] = resultResponseMake("import success! pleases! wait some time")
	i.ServeJSON()
}
