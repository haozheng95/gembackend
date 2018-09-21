package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gembackend/models/btc_query"
	"github.com/gembackend/models/eth_query"
	"github.com/gembackend/models/exchange"
	"github.com/shopspring/decimal"
	"strconv"
)

type AssetController struct {
	beego.Controller
}

func (a *AssetController) Get() {
	walletId := a.Input().Get("wallet_id")
	begin := a.Input().Get("begin")
	size := a.Input().Get("size")

	if len(begin) < 1 {
		begin = "0"
	}
	if len(size) < 1 {
		size = "10"
	}
	beginInt, _ := strconv.Atoi(begin)
	sizeInt, _ := strconv.Atoi(size)

	result := make([]*assertControllerResponse, 0, 10)

	if beginInt == 0 {
		st := new(assertControllerResponse)
		ethData := eth_query.GetEthInfoByWalletId(walletId)
		cny := exchange.GetMainChainCnyByCoinName("eth")
		amount := SubString(ethData.Amount, ethData.UnconfirmAmount)
		st.Coin = "eth"
		st.Amount = amount
		st.Dec = ethData.Decimal
		st.Istoken = "0"
		st.Price = MulString(cny, amount)
		result = append(result, st)
		result = append(result, GetBtcInfo(walletId))
		sizeInt -= 2
	} else {
		beginInt = beginInt*sizeInt - 2
	}
	ethTokenData := eth_query.GetAllTokenInfoWithUser(walletId, beginInt, sizeInt)
	for _, v := range ethTokenData {
		st := new(assertControllerResponse)
		amount := SubString(v.Amount, v.UnconfirmAmount)
		cny := exchange.GetEthTokenCnyByContractAddr(v.ContractAddr)
		st.Istoken = "1"
		st.Dec = v.Decimal
		st.Price = MulString(cny, amount)
		st.Coin = v.TokenName
		//log.Debug(st.Coin)
		st.ContractAddr = v.ContractAddr
		result = append(result, st)
	}

	a.Data["json"] = resultResponseMake(result)
	a.ServeJSON()
}

func GetBtcInfo(walletId string) *assertControllerResponse {
	allAddr := btc_query.GetUserInfo(walletId)
	amount := decimal.New(0, 8)
	unconfirmAmount := decimal.New(0, 8)
	for _, value := range allAddr {
		tempAmount, _ := decimal.NewFromString(value.Amount)
		tempUnconfirmAmount, _ := decimal.NewFromString(value.UnconfirmAmount)
		amount = amount.Add(tempAmount)
		unconfirmAmount = unconfirmAmount.Add(tempUnconfirmAmount)
	}
	resultAmount := amount.Sub(unconfirmAmount)
	log.Debug("amount   ===", amount)
	//log.Debug("unamount ===", unconfirmAmount)
	//log.Debug("result   ===", resultAmount)
	cny := exchange.GetMainChainCnyByCoinName("btc")
	log.Debug("cny   ===", cny)
	return NewassertControllerResponse("btc", resultAmount.String(), cny, "8", "", "0")
}

// result struct
type assertControllerResponse struct {
	Coin, Amount, Price, Dec, ContractAddr, Istoken string
}

func NewassertControllerResponse(Coin, Amount, Price, Dec, ContractAddr, Istoken string) (res *assertControllerResponse) {
	res = &assertControllerResponse{Coin, Amount, Price, Dec, ContractAddr, Istoken}
	return
}
