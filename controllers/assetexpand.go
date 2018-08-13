//software: GoLand
//file: assetexpand.go
//time: 2018/8/13 下午2:53
package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/gembackend/models/eth_query"
)

type AssetExpandController struct {
	beego.Controller
}

/**
param = {
"eth":["","contractaddr", "contractaddr"],
}
*/
func (a *AssetExpandController) Get() {
	param := a.Input().Get("param")
	ethAddr := a.Input().Get("eth_addr")
	result := make(map[string]interface{})

	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(param), &m)
	if err != nil {
		// dispose err
		a.Data["json"] = resultResponseErrorMake(2001, err.Error())
		a.ServeJSON()
		return
	}
	for k, v := range m {
		switch k {
		case "eth":
			v1 := v.([]interface{})
			result1 := make([]interface{}, len(v1))
			for i, v2 := range v1 {
				v3 := v2.(string)
				if len(v3) > 0 {
					//dispose token
					tx, r := eth_query.GetEthTokenTxrecord(ethAddr, v3, 0, 1)
					if r > 0 {
						//result1 = append(result1, map[string]interface{}{v3: tx[0]})
						result1[i] = map[string]interface{}{v3: tx[0]}
					}
				} else {
					//dispose eth
					tx, r := eth_query.GetEthTxrecord(ethAddr, 0, 1)
					if r > 0 {
						//result1 = append(result1, map[string]interface{}{"eth": tx[0]})
						result1[i] = map[string]interface{}{"eth": tx[0]}
					}
				}
			}
			result["eth"] = result1
		}
	}

	a.Data["json"] = resultResponseMake(result)
	a.ServeJSON()
}
