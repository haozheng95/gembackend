package controllers

import (
	"github.com/astaxie/beego"
	"strings"
	"strconv"
)

var (
	resultResponseErrors = map[int]string{
		2000: "missing parameter",
		2001: "Parameter resolution error",
		2002: "Have no legal power",
		2003: "Query error",
		2004: "User does not exist",
	}
)

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
