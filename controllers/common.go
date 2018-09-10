package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/astaxie/beego"
	"github.com/gembackend/gembackendlog"
	"github.com/gembackend/messagequeue"
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
)

var (
	defaultTokenGasLimit = "150000"
	defaultGasLimit      = "21000"
	defaultGasPrice      = "90000000000"

	resultResponseErrors = map[int]string{
		2000: "missing parameter",
		2001: "Parameter resolution error",
		2002: "Have no legal power",
		2003: "Query error",
		2004: "User does not exist",
		2005: "param error",
		2006: "User does not exist",
		2007: "checkSign false",
		2008: "map to json error",
		2009: "send raw tx error",
		2010: "This currency is not supported",
		2011: "You can't transfer it to yourself",
		2012: "The returned data is less than one",
	}
	log = gembackendlog.Log
	//默认token
	DefaultToken = [][]string{
		//EOS
		{"0x86fa049857e0209aa7d9e616f7eb3b3b78ecfdb0", "18", "EOS"},
		//TRX
		{"0xf230b790e05390fc8295f4d3f60332c93bed42e2", "6", "TRX"},
		//OMG
		{"0xd26114cd6ee289accf82350c8d8487fedb8a0c07", "18", "OMG"},
		//SNT
		{"0x744d70fdbe2ba4cf95131626614a1763df805b9e", "18", "SNT"},
		//MKR
		{"0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2", "18", "MKR"},
		//ZRX
		{"0xe41d2489571d322189246dafa5ebde1f4699f498", "18", "ZRX"},
		//REP
		{"0xe94327d07fc17907b4db788e5adf2ed424addff6", "18", "REP"},
		//BAT
		{"0x0d8775f648430679a709e98d2b0cb6250d2887ef", "18", "BAT"},
		//SALT
		{"0x4156d3342d5c385a87d264f90653733592000581", "8", "SALT"},
		//FUN
		{"0x419d0d8bdd9af5e606ae2232ed285aff190e711b", "8", "FUN"},
	}
)

const (
	HMAC_KEY = "0642de2eb660d56402fa690b1c5474a4"
)

func checkSign(walletId, sign string) bool {
	h := sha256.New()
	h.Write([]byte(walletId + HMAC_KEY))
	bs := h.Sum(nil)
	hashValue := hex.EncodeToString(bs)
	return strings.Compare(sign, hashValue) == 0
}

func resultResponseMake(result interface{}) interface{} {
	resultResponse := map[string]interface{}{
		"status":  0,
		"result":  nil,
		"version": 1.0,
	}
	resultResponse["result"] = result
	return resultResponse
}

func resultResponseErrorMake(errorCode int, err interface{}) interface{} {
	resultResponse := map[string]interface{}{
		"status":       0,
		"result":       nil,
		"version":      1.0,
		"error":        "Undefined error",
		"error_detail": nil,
	}
	e, check := resultResponseErrors[errorCode]
	if check {
		resultResponse["error"] = e
		resultResponse["error_detail"] = err
	}
	resultResponse["status"] = errorCode
	return resultResponse
}

type ErrorsController struct {
	beego.Controller
}

func (e *ErrorsController) Get() {
	p := e.Ctx.Input.Param(":error_id")
	if strings.Compare(strings.TrimSpace(p), "") == 0 {
		e.Data["json"] = resultResponseMake(nil)
	} else {
		eId, err := strconv.Atoi(p)
		if err != nil {
			e.Data["json"] = resultResponseErrorMake(eId, err.Error())
		}
		e.Data["json"] = resultResponseErrorMake(eId, nil)
	}
	e.ServeJSON()
}

func SubString(d1, d2 string) (r string) {
	m1, _ := decimal.NewFromString(d1)
	m2, _ := decimal.NewFromString(d2)
	m3 := m1.Sub(m2)
	r = m3.String()
	return
}

func SubStringDecimal(d1, d2 string) (r decimal.Decimal) {
	m1, _ := decimal.NewFromString(d1)
	m2, _ := decimal.NewFromString(d2)
	m3 := m1.Sub(m2)
	r = m3
	return
}

func AddString(d1, d2 string) (r string) {
	m1, _ := decimal.NewFromString(d1)
	m2, _ := decimal.NewFromString(d2)
	m3 := m1.Add(m2)
	r = m3.String()
	return
}

func MulString(d1, d2 string) (r string) {
	m1, _ := decimal.NewFromString(d1)
	m2, _ := decimal.NewFromString(d2)
	m3 := m1.Mul(m2)
	r = m3.String()
	return
}

func MulString2(d1, d2 string) (r decimal.Decimal) {
	m1, _ := decimal.NewFromString(d1)
	m2, _ := decimal.NewFromString(d2)

	r = m1.Mul(m2)
	return
}

func DivString(d1, d2 string) (r string) {
	m1, _ := decimal.NewFromString(d1)
	m2, _ := decimal.NewFromString(d2)
	m3 := m1.Div(m2)
	r = m3.String()
	return
}

func SaveForKafka(topicname, kafkaParam string) {
	p := messagequeue.MakeProducer()
	defer p.Close()
	messagequeue.MakeMessage(topicname, kafkaParam, p)
}
