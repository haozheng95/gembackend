package controllers

import (
	"github.com/astaxie/beego"
	"strings"
	"strconv"
	"github.com/gembackend/gembackendlog"
	"crypto/sha256"
	"encoding/hex"
)

var (
	resultResponseErrors = map[int]string{
		2000: "missing parameter",
		2001: "Parameter resolution error",
		2002: "Have no legal power",
		2003: "Query error",
		2004: "User does not exist",
		2005: "param error",
		2006: "User does not exist",
	}
	log = gembackendlog.Log
	//默认token
	defaultToken = [][]string{
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
