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

	// result ----
	//result := make([]assertControllerResponse, sizeInt, sizeInt)
	result := make([]assertControllerResponse, 0, sizeInt)
	st := assertControllerResponse{}
	if beginInt == 0 {
		ethData := eth_query.GetEthInfoByWalletId(walletId)
		cny := exchange.GetMainChainCnyByCoinName("eth")
		amount := SubString(ethData.Amount, ethData.UnconfirmAmount)
		//result[0].Coin = "eth"
		//result[0].Amount = amount
		//result[0].Dec = ethData.Decimal
		//result[0].Istoken = "0"
		//result[0].Price = MulString(cny, amount)
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
	//log.Debug(walletId)//d3ba134f262d6d197a93ade4a6c123ddb9122c5cc0ff666f5447639d36f5f155
	//if beginInt != 0 {
	//	sizeInt--
	//}
	for _, v := range ethTokenData {
		amount := SubString(v.Amount, v.UnconfirmAmount)
		cny := exchange.GetEthTokenCnyByContractAddr(v.ContractAddr)
		//result[sizeInt-i].Istoken = "1"
		//result[sizeInt-i].Dec = v.Decimal
		//result[sizeInt-i].Price = MulString(cny, amount)
		//result[sizeInt-i].Amount = amount
		//result[sizeInt-i].Coin = v.TokenName
		st.Istoken = "1"
		st.Dec = v.Decimal
		st.Price = MulString(cny, amount)
		st.Amount = amount
		st.Coin = v.TokenName
		st.ContractAddr = v.ContractAddr
		result = append(result, st)
	}
	//log.Debug(cap(result))
	//log.Debug(len(result))
	//res := map[string]interface{}{
	//	"data":result,
	//	"last":sizeInt < len(result),
	//}
	a.Data["json"] = resultResponseMake(result)
	a.ServeJSON(true)
}

// result struct
type assertControllerResponse struct {
	Coin, Amount, Price, Dec, ContractAddr, Istoken string
}
