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
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/gembackend/controllers"
	"github.com/gembackend/hjwt"
	"strings"
)

func init() {
	// 设置擦车
	beego.InsertFilter("*", beego.BeforeRouter, func(ctx *context.Context) {
		ctx.Output.Header("Access-Control-Allow-Origin", "*")
		ctx.Output.Header("Cache-Control", "no-cache,no-store")
	})

	// 设置跨域
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "auth-token"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	// v1 router
	versionForV1()
	// 错误信息返回
	beego.Router("/error/?:error_id", &controllers.ErrorsController{})
}

func versionForV1() {
	//初始化 namespace
	ns :=
		beego.NewNamespace("/v1",
			beego.NSCond(jwtAuth),
			// --get
			beego.NSRouter("/asset", &controllers.AssetController{}, "*:Get"),
			beego.NSRouter("/auth", &controllers.AuthController{}, "*:Get"),
			beego.NSRouter("/balance/?:coin_type", &controllers.BalanceController{}, "*:Get"),
			beego.NSRouter("/txs/?:coin_type", &controllers.TxrecordController{}, "*:Get"),
			beego.NSRouter("/txinfo/?:coin_type", &controllers.TxinfoController{}, "*:Get"),
			// --post
			beego.NSRouter("/register", &controllers.RegisterController{}, "*:Post"),
			beego.NSRouter("/import", &controllers.ImportWalletController{}, "*:Post"),
			beego.NSRouter("/rawtx/?:coin_type", &controllers.SendRawTx{}, "*:Post"),
		)
	//注册 namespace
	beego.AddNamespace(ns)
}

func jwtAuth(ctx *context.Context) bool {
	if strings.Compare(ctx.Request.URL.Path, "/v1/auth") == 0 {
		return true
	}
	if strings.Compare(ctx.Request.URL.Path, "/v1/register") == 0 {
		return true
	}
	if strings.Compare(ctx.Request.URL.Path, "/v1/import") == 0 {
		return true
	}
	token, check := ctx.Request.Header["Auth-Token"]
	if check && len(token) > 0 && hjwt.CheckToken(token[0]) {
		return true
	} else {
		// 跳转错误页面
		//ctx.Redirect(302, "/error/2002")
		return false
	}
	return false
}
