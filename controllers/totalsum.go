//software: GoLand
//file: totalsum.go
//time: 2018/9/10 下午2:56
package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gembackend/models/eth_query"
	"github.com/gembackend/models/exchange"
	"github.com/shopspring/decimal"
)

type Totalsum struct {
	beego.Controller
}

func (t *Totalsum) Get() {
	walletId := t.Input().Get("wallet_id")
	log.Debug(walletId)
	sum := t.getTotalSum(walletId)
	log.Debug("-----", sum)
	t.Data["json"] = resultResponseMake(sum)
	t.ServeJSON()

}
func (totalsum *Totalsum) getTotalSum(walletId string) interface{} {

	sum := decimal.New(0, 0)
	btcInfo := GetBtcInfo(walletId)
	btcSum := btcInfo.Amount
	btcPrice := btcInfo.Price
	sum = sum.Add(MulString2(btcSum, btcPrice))
	log.Debug("btc amount ====", btcInfo.Amount)
	log.Debug("btc price  ====", sum)

	ethInfo := eth_query.GetEthInfoByWalletId(walletId)
	ethPrice := exchange.GetMainChainCnyByCoinName("eth")
	ethSum := SubString(ethInfo.Amount, ethInfo.UnconfirmAmount)
	sum = sum.Add(MulString2(ethSum, ethPrice))

	tokenInfo := eth_query.GetAllTokenInfoWithUser(walletId, 0, 1000)
	tokenLastPrice := totalsum.getTokenPrice(tokenInfo)
	sum = sum.Add(tokenLastPrice)

	return sum

}
func (totalsum *Totalsum) getTokenPrice(tokenInfo []*struct {
	ContractAddr, Amount, UnconfirmAmount, TokenName, Decimal string
}) decimal.Decimal {
	sum := decimal.New(0, 0)
	for _, v := range tokenInfo {
		amount := SubStringDecimal(v.Amount, v.UnconfirmAmount)
		//log.Debug("contract addr ===", v.ContractAddr)
		//log.Debug("token name ===", v.TokenName)
		price := exchange.GetEthTokenCnyByContractAddr(v.ContractAddr)
		//log.Debug("price ===", price)
		sum = sum.Add(MulString2(price, amount.String()))
	}
	log.Debug("token sum ===", sum)
	return sum
}
