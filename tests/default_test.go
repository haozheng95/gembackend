package test

//software: GoLand
//file: test.go
//time: 2018/7/30 下午2:55
import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	_ "github.com/gembackend/models"
	_ "github.com/gembackend/routers"
	"github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	godKey       = "ASDFASDqwqfasvsfqioqjweamsdfmosejoqjma"
	walletId     = "d3ba134f262d6d197a93ade4a6c123ddb9122c5cc0ff666f5447639d36f5f155"
	ethAddr      = "0xd6cb6744b7f2da784c5afd6b023d957188522198"
	sign         = "6e904d69f5277bd863c4b09be37000cd4bf61b4a17f2a0099d5f1c5692e7402c"
	txHash       = "0x569c5b35f203ca6db6e2cec44bceba756fad513384e2bd79c06a8c0181273379"
	contractAddr = "0xd26114cd6ee289accf82350c8d8487fedb8a0c07"
	basePath     = "/v1"
)

func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, "../.."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

func decodeJson(s string) (r map[string]interface{}) {
	x := []byte(s)
	err := json.Unmarshal(x, &r)
	if err != nil {
		fmt.Println(err)
	}
	return
}

//Get ---
/**
@param://
@user_addr
@contract_addr
*/
func TestTxinfo(t *testing.T) {
	// eth
	param := "/eth"
	param += "?tx_hash=" + txHash

	r, _ := http.NewRequest("GET", basePath+"/txinfo"+param, nil)
	r.Header.Add("auth-token", godKey)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("txinfo", "TestTxinfo", fmt.Sprintf("Code[%d]\n%s", w.Code, w.Body.String()))
	z := decodeJson(w.Body.String())
	convey.Convey("txinfo eth", t, func() {
		convey.Convey("status code should be 0", func() {
			convey.So(z["status"], convey.ShouldEqual, 0)
		})
	})
	// contract
	param += "&contract_addr=" + contractAddr
	r, _ = http.NewRequest("GET", basePath+"/txinfo"+param, nil)
	r.Header.Add("auth-token", godKey)
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("txinfo", "TestTxinfo", fmt.Sprintf("Code[%d]\n%s", w.Code, w.Body.String()))
	z = decodeJson(w.Body.String())
	convey.Convey("txinfo  token", t, func() {
		convey.Convey("status code should be 0", func() {
			convey.So(z["status"], convey.ShouldEqual, 0)
		})
	})
}

/**
@param://
@wallet_id
@contract_addr
*/
func TestBalance(t *testing.T) {
	param := "/eth?wallet_id=" + walletId
	param += "&contract_addr="
	r, _ := http.NewRequest("GET", basePath+"/balance"+param, nil)
	r.Header.Add("auth-token", godKey)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("balance", "TestBalance", fmt.Sprintf("Code[%d]\n%s", w.Code, w.Body.String()))
	z := decodeJson(w.Body.String())
	convey.Convey("balance", t, func() {
		convey.Convey("status code should be 0", func() {
			convey.So(z["status"], convey.ShouldEqual, 0)
		})
	})

	contract := "0x86fa049857e0209aa7d9e616f7eb3b3b78ecfdb0"
	param += contract
	r, _ = http.NewRequest("GET", basePath+"/balance"+param, nil)
	r.Header.Add("auth-token", godKey)
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("balance token", "TestBalance", fmt.Sprintf("Code[%d]\n%s", w.Code, w.Body.String()))
	z = decodeJson(w.Body.String())
	convey.Convey("balance", t, func() {
		convey.Convey("status code should be 0", func() {
			convey.So(z["status"], convey.ShouldEqual, 0)
		})
	})
}

/**
@params://
@wallet_id
@begin
@size
*/
func TestAsset(t *testing.T) {
	begin := "0"
	size := "10"
	param := "?wallet_id=" + walletId
	param += "&begin=" + begin
	param += "&size=" + size
	r, _ := http.NewRequest("GET", basePath+"/asset"+param, nil)
	r.Header.Add("auth-token", godKey)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("asset", "TestAsset", fmt.Sprintf("Code[%d]\n%s", w.Code, w.Body.String()))
	z := decodeJson(w.Body.String())
	convey.Convey("Asset", t, func() {
		convey.Convey("status code should be 0", func() {
			convey.So(z["status"], convey.ShouldEqual, 0)
		})
	})
}

/**
@param://
@wallet_id
*/
func TestAuth(t *testing.T) {
	param := "?wallet_id=" + walletId
	r, _ := http.NewRequest("GET", "/v1/auth"+param, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("auth", "TestAuth", fmt.Sprintf("Code[%d]\n%s", w.Code, w.Body.String()))
	z := decodeJson(w.Body.String())
	convey.Convey("Auth", t, func() {
		convey.Convey("status code should be 0", func() {
			convey.So(z["status"], convey.ShouldEqual, 0)
		})
	})
}

/**
@param://
@wallet_id
@contract_addr
@begin_page
@size
*/
func TestTxs(t *testing.T) {
	param := "/eth"
	param += "?wallet_id=" + walletId
	param += "&begin_page=0"
	param += "&size=10"
	param += "&contract_addr="
	r, _ := http.NewRequest("GET", basePath+"/txs"+param, nil)
	r.Header.Add("auth-token", godKey)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("txs", "TestTxs", fmt.Sprintf("Code[%d]\n%s", w.Code, w.Body.String()))
	z := decodeJson(w.Body.String())
	convey.Convey("txs", t, func() {
		convey.Convey("status code should be 0", func() {
			convey.So(z["status"], convey.ShouldEqual, 0)
		})
	})

	param += contractAddr
	r, _ = http.NewRequest("GET", basePath+"/txs"+param, nil)
	r.Header.Add("auth-token", godKey)
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("txs", "TestTxs", fmt.Sprintf("Code[%d]\n%s", w.Code, w.Body.String()))
	z = decodeJson(w.Body.String())
	convey.Convey("txs", t, func() {
		convey.Convey("status code should be 0", func() {
			convey.So(z["status"], convey.ShouldEqual, 0)
		})
	})
}

//Post
/**
@json params://
@wallet_id
@sign
@eth_addr
*/
func TestRegister(t *testing.T) {

	param := struct {
		WalletId string `json:"wallet_id"`
		Sign     string `json:"sign"`
		EthAddr  string `json:"eth_addr"`
	}{
		walletId, sign, ethAddr,
	}
	jsons, _ := json.Marshal(param)
	r, _ := http.NewRequest("POST", basePath+"/register", bytes.NewBuffer(jsons))
	r.Header.Add("auth-token", godKey)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("register", "TestRegister", fmt.Sprintf("Code[%d]\n%s", w.Code, w.Body.String()))
	z := decodeJson(w.Body.String())
	convey.Convey("register", t, func() {
		convey.Convey("status code should be 0", func() {
			convey.So(z["status"], convey.ShouldEqual, 0)
		})
	})
}

/**
@json params://
@wallet_id
@sign
@eth_addr
*/
func TestImport(t *testing.T) {

	param := struct {
		WalletId string `json:"wallet_id"`
		Sign     string `json:"sign"`
		EthAddr  string `json:"eth_addr"`
	}{
		walletId, sign, ethAddr,
	}
	jsons, _ := json.Marshal(param)
	r, _ := http.NewRequest("POST", basePath+"/import", bytes.NewBuffer(jsons))
	r.Header.Add("auth-token", godKey)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("import", "TestImport", fmt.Sprintf("Code[%d]\n%s", w.Code, w.Body.String()))
	z := decodeJson(w.Body.String())
	convey.Convey("import", t, func() {
		convey.Convey("status code should be 0", func() {
			convey.So(z["status"], convey.ShouldEqual, 0)
		})
	})
}

/**
@coin_type in path
@json params://
@amount
@raw
@fee
@from
@to
@note
@gaslimit
@gasprice
@dec
@contract_addr
*/
func TestSendRaw(t *testing.T) {
	coin := "/eth"
	param := struct {
		Amount       string `json:"amount"`
		Raw          string `json:"raw"`
		Fee          string `json:"fee"`
		From         string `json:"from"`
		To           string `json:"to"`
		Note         string `json:"note"`
		Nonce        string `json:"nonce"`
		Gaslimit     string `json:"gaslimit"`
		Gasprice     string `json:"gasprice"`
		Dec          string `json:"dec"`
		ContractAddr string `json:"contract_addr"`
	}{
		"1",
		"0xf86902847735940082520894ff1ef64ea9fdddb1a1e17bb44c5ca1ddc508cd8d870874e5c4192400801ca0fb08865641002e949d11a00590e10c7309cbbdf2345e68bdc88cf3ae428b20f59f456366f06cc372f16b407e80b6b210000bf30d92e705668feaae56c2e27415",
		"0.4",
		"0xd6cb6744b7f2da784c5afd6b023d957188522198",
		"0x4ac751f0152b6799a5acfc25614072fbb06dca06",
		"sup no token",
		"22",
		"150000",
		"1000000",
		"18",
		"",
	}
	jsons, _ := json.Marshal(param)
	r, _ := http.NewRequest("POST", basePath+"/rawtx"+coin, bytes.NewBuffer(jsons))
	r.Header.Add("auth-token", godKey)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("SendRaw", "TestSendRaw", fmt.Sprintf("Code[%d]\n%s", w.Code, w.Body.String()))
	z := decodeJson(w.Body.String())
	convey.Convey("SendRaw", t, func() {
		convey.Convey("eth ==== status code should be 0", func() {
			convey.So(z["status"], convey.ShouldEqual, 0)
		})
	})

	// token ----------
	param.ContractAddr = contractAddr
	jsons, _ = json.Marshal(param)
	r, _ = http.NewRequest("POST", basePath+"/rawtx"+coin, bytes.NewBuffer(jsons))
	r.Header.Add("auth-token", godKey)
	r.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("SendRaw", "TestSendRaw", fmt.Sprintf("Code[%d]\n%s", w.Code, w.Body.String()))
	z = decodeJson(w.Body.String())
	convey.Convey("SendRaw", t, func() {
		convey.Convey("eth ==== status code should be 0", func() {
			convey.So(z["status"], convey.ShouldEqual, 0)
		})
	})
}
