package main

import (
	"github.com/astaxie/beego"
	_ "github.com/gembackend/models"
	_ "github.com/gembackend/routers"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run(":9000")
}
