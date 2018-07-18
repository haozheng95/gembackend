package controllers

import (
	"github.com/astaxie/beego"
	"crypto/sha256"
	"strings"
	"encoding/json"
	"github.com/gembackend/gembackendlog"
	"github.com/gembackend/models"
	"strconv"
	"encoding/hex"
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
		r.Data["json"] = err.Error()
		r.ServeJSON()
		return
	}

	walletId, err1 := m["wallet_id"]
	sign, err2 := m["sign"]
	ethAddr, err3 := m["eth_addr"]

	if !err1 || !err2 || !err3 {
		log.Error(err1, err2, err3)
		r.Data["json"] = "missing params"
		r.ServeJSON()
		return
	}

	if !checkSign(walletId, sign) {
		log.Error("checkSign false")
		r.Data["json"] = "checkSign false"
		r.ServeJSON()
		return
	}
	// 添加eth地址
	addressTable := &models.Address{
		WalletId:        walletId,
		Addr:            ethAddr,
		Nonce:           "0",
		Amount:          "0",
		UnconfirmAmount: "0",
		TypeId:          4,
	}
	addressTable.InsertOneRaw(addressTable)
	// 添加token地址
	addressTokenTable := &models.TokenAddress{
		WalletId:        walletId,
		Addr:            ethAddr,
		Amount:          "0",
		UnconfirmAmount: "0",
		Added:           1,
	}
	for _, v := range defaultToken {
		addressTokenTable.ContractAddr = v[0]
		dec, _ := strconv.Atoi(v[1])
		addressTokenTable.Decimal = int64(dec)
		addressTokenTable.InsertOneRaw(addressTokenTable)
	}

	r.Data["json"] = resultResponseMake("success")
	r.ServeJSON()
}

func checkSign(walletId, sign string) bool {
	h := sha256.New()
	h.Write([]byte(walletId + HMAC_KEY))
	bs := h.Sum(nil)
	hashValue := hex.EncodeToString(bs)
	return strings.Compare(sign, hashValue) == 0
}

const HMAC_KEY = "0642de2eb660d56402fa690b1c5474a4"

var log = gembackendlog.Log
//默认token
var defaultToken = [][]string{
	//EOS
	{"0x86fa049857e0209aa7d9e616f7eb3b3b78ecfdb0", "18"},
	//TRX
	{"0xf230b790e05390fc8295f4d3f60332c93bed42e2", "6"},
	//OMG
	{"0xd26114cd6ee289accf82350c8d8487fedb8a0c07", "18"},
	//SNT
	{"0x744d70fdbe2ba4cf95131626614a1763df805b9e", "18"},
	//MKR
	{"0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2", "18"},
	//ZRX
	{"0xe41d2489571d322189246dafa5ebde1f4699f498", "18"},
	//REP
	{"0xe94327d07fc17907b4db788e5adf2ed424addff6", "18"},
	//BAT
	{"0x0d8775f648430679a709e98d2b0cb6250d2887ef", "18"},
	//SALT
	{"0x4156d3342d5c385a87d264f90653733592000581", "8"},
	//FUN
	{"0x419d0d8bdd9af5e606ae2232ed285aff190e711b", "8"},
}
