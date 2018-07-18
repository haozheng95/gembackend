// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"github.com/gembackend/controllers"
	"github.com/astaxie/beego/context"
	"github.com/gembackend/hjwt"
	"strings"
	"github.com/astaxie/beego/plugins/cors"
)

func init() {
	// 设置jwt验证
	beego.InsertFilter("/api/*", beego.BeforeStatic, func(ctx *context.Context) {
		if strings.Compare(ctx.Request.URL.Path, "/api/auth") == 0 {
			return
		}

		token, check := ctx.Request.Header["Auth-Token"]
		if check && len(token) > 0 && hjwt.CheckToken(token[0]) {

		} else {
			// 跳转错误页面
			ctx.Redirect(302, "/error/2002")
		}
	})

	// 设置擦车
	beego.InsertFilter("*", beego.BeforeRouter, func(ctx *context.Context) {
		ctx.Output.Header("Cache-Control", "no-cache,no-store")
	})

	// 设置跨域
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))

	beego.Router("/api/asset", &controllers.AssetController{})
	beego.Router("/api/balance", &controllers.BalanceController{})
	beego.Router("/api/address", &controllers.RegisterController{})
	beego.Router("/api/auth", &controllers.AuthController{})
	beego.Router("/api/txs/?:coin_type", &controllers.TxrecordController{})
	// 错误信息返回
	beego.Router("/error/?:error_id", &controllers.ErrorsController{})
}
