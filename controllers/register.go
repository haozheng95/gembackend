package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/gembackend/models/btc_query"
	"github.com/gembackend/models/eth_query"
	"strconv"
)

type RegisterController struct {
	beego.Controller
}

//创建账户 添加十个默认token
func (r *RegisterController) Post() {
	m := make(map[string]string)
	err := json.Unmarshal(r.Ctx.Input.RequestBody, &m)
	if err != nil {
		log.Error(err)
		r.Data["json"] = resultResponseErrorMake(2001, err.Error())
		r.ServeJSON()
		return
	}

	walletId, err1 := m["wallet_id"]
	sign, err2 := m["sign"]
	ethAddr, err3 := m["eth_addr"]
	btcAddr, err4 := m["btc_addr"]
	if !err1 || !err2 || !err3 || !err4 {
		log.Error(err1, err2, err3)
		r.Data["json"] = resultResponseErrorMake(2000, nil)
		r.ServeJSON()
		return
	}

	if !checkSign(walletId, sign) {
		log.Error("checkSign false")
		r.Data["json"] = resultResponseErrorMake(2007, nil)
		r.ServeJSON()
		return
	}
	// 添加eth地址
	addressTable := &eth_query.Address{
		WalletId:        walletId,
		Addr:            ethAddr,
		Nonce:           "0",
		Amount:          "0",
		UnconfirmAmount: "0",
		TypeId:          4,
		Decimal:         18,
	}
	addressTable.InsertOneRaw(addressTable)
	// 添加token地址
	addressTokenTable := &eth_query.TokenAddress{
		WalletId:        walletId,
		Addr:            ethAddr,
		Amount:          "0",
		UnconfirmAmount: "0",
		Added:           1,
	}
	for _, v := range DefaultToken {
		addressTokenTable.ContractAddr = v[0]
		dec, _ := strconv.Atoi(v[1])
		addressTokenTable.Decimal = int64(dec)
		addressTokenTable.TokenName = v[2]
		addressTokenTable.InsertOneRaw(addressTokenTable)
	}

	// add eth address for kafka
	//SaveForKafka(conf.KafkagetbalanceParityTopic, ethAddr)
	// add btc addr
	if res := MakeAddressBtc(walletId, btcAddr); len(res) > 0 {
		btc_query.InsertAddress(res)
	}

	r.Data["json"] = resultResponseMake("success")
	r.ServeJSON()
}

func MakeAddressBtc(walletId, btcAddr string) (res []*btc_query.AddressBtc) {
	addrs := make([]string, 0, 10)
	if err := json.Unmarshal([]byte(btcAddr), &addrs); err != nil {
		log.Warning(err)
	}
	if len(addrs) > 0 {
		res = make([]*btc_query.AddressBtc, 0, len(addrs))
		for _, addr := range addrs {
			//walletId, addr, amount, unconfirmamount string, typeId int
			tmp := btc_query.NewAddress(walletId, addr, "0", "0", 4)
			res = append(res, tmp)
		}
	}
	return
}
