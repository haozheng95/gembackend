package main

import (
	_ "github.com/astaxie/beego/config/xml"
	"github.com/gembackend/scripts"
)

func main() {
	scripts.StartEthupdaterMul(5000000)
}
