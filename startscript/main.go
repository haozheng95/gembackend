package main

import (
	_ "github.com/astaxie/beego/config/xml"
	"github.com/gembackend/scripts"
)

func main() {
	scripts.StartEthupdaterMul(5000000)
	//fmt.Println(rpc.Eth_getTxList("0xaaa5517cc033189da19d88f20b2d68085e49c259"))
	//scripts.GetAllTxList("0xaaa5517cc033189da19d88f20b2d68085e49c259")
	//scripts.Main("555", "0xddbd2b932c763ba5b1b7ae3b362eac3e8d40121a")
}
