package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gembackend/models/eth_query"
	"github.com/gembackend/models/exchange"
)

type BalanceController struct {
	beego.Controller
}

// 获取eth余额 和合约地址余额
func (b *BalanceController) Get() {
	cointype := b.Ctx.Input.Param(":coin_type")
	walletid := b.Input().Get("wallet_id")

	res := new(balanceControllerResponse)
	switch cointype {
	case "eth":
		contractaddr := b.Input().Get("contract_addr")
		if len(contractaddr) > 0 {
			// contract dispose
			st := eth_query.GetEthTokenInfoByWalletId(walletid, contractaddr)
			res.Amount = SubString(st.Amount, st.UnconfirmAmount)
			res.TokenAmount = SubString(st.TokenAmount, st.TokenUnconfirmAmount)
			res.TokenName = st.TokenName
			res.Nonce = st.Nonce
			res.Dec = st.Decimal
			res.Istoken = "1"
			tokenGasLimit := exchange.GetTokenGasLimit(contractaddr)
			if len(tokenGasLimit) == 0 {
				res.Gaslimit = defaultTokenGasLimit
			} else {
				res.Gaslimit = tokenGasLimit
			}
		} else {
			// eth dispose
			st := eth_query.GetEthInfoByWalletId(walletid)
			res.Amount = SubString(st.Amount, st.UnconfirmAmount)
			res.Nonce = st.Nonce
			res.Dec = st.Decimal
			res.Istoken = "0"
			gaslimit := exchange.GetMainChainGasLimit(cointype)
			if len(gaslimit) > 0 {
				res.Gaslimit = gaslimit
			} else {
				res.Gaslimit = defaultGasLimit
			}
		}
		res.Coin = cointype
		res.ContractAddr = contractaddr
		gasprice := exchange.GetMainChainGasPrice(cointype)
		if len(gasprice) > 0 {
			res.Gasprice = gasprice
		} else {
			res.Gasprice = defaultGasPrice
		}
	case "btc":
		tempst := GetBtcInfo(walletid)
		res.Coin = tempst.Coin
		res.Amount = tempst.Amount
		res.Dec = "8"
		res.Istoken = "0"
		res.Fee = exchange.GetCoinFee("btc")

	default:
		// error
		b.Data["json"] = resultResponseErrorMake(2010, nil)
		b.ServeJSON()
		return
	}
	//log.Debug(res)
	b.Data["json"] = resultResponseMake(res)
	b.ServeJSON()
}

type balanceControllerResponse struct {
	Coin, Amount, TokenAmount, Nonce, Gasprice, Gaslimit, Dec, Istoken, ContractAddr, TokenName, Fee string
}
