package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gembackend/models/eth_query"
	"github.com/gembackend/models/exchange"
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

	result := make([]assertControllerResponse, 0, 10)
	st := assertControllerResponse{}
	if beginInt == 0 {
		ethData := eth_query.GetEthInfoByWalletId(walletId)
		cny := exchange.GetMainChainCnyByCoinName("eth")
		amount := SubString(ethData.Amount, ethData.UnconfirmAmount)
		st.Coin = "eth"
		st.Amount = amount
		st.Dec = ethData.Decimal
		st.Istoken = "0"
		st.Price = MulString(cny, amount)
		result = append(result, st)
		sizeInt--
	} else {
		beginInt = beginInt*sizeInt - 1
	}
	ethTokenData := eth_query.GetAllTokenInfoWithUser(walletId, beginInt, sizeInt)
	for _, v := range ethTokenData {
		amount := SubString(v.Amount, v.UnconfirmAmount)
		cny := exchange.GetEthTokenCnyByContractAddr(v.ContractAddr)
		st.Istoken = "1"
		st.Dec = v.Decimal
		st.Price = MulString(cny, amount)
		st.Amount = amount
		st.Coin = v.TokenName
		st.ContractAddr = v.ContractAddr
		result = append(result, st)
	}

	a.Data["json"] = resultResponseMake(result)
	a.ServeJSON(true)
}

// result struct
type assertControllerResponse struct {
	Coin, Amount, Price, Dec, ContractAddr, Istoken string
}
