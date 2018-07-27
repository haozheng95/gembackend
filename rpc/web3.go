package rpc

import (
	"github.com/regcostajr/go-web3/providers"
	"github.com/regcostajr/go-web3"
	"github.com/regcostajr/go-web3/net"
	"github.com/regcostajr/go-web3/personal"
	"github.com/regcostajr/go-web3/utils"
	"github.com/regcostajr/go-web3/dto"
)

// Web3 extension

type Web3 struct {
	web3.Web3
	Eth *Eth
}

type Eth struct {
	provider providers.ProviderInterface
}

func NewEth(provider providers.ProviderInterface) *Eth {
	go_eth := new(Eth)
	go_eth.provider = provider
	return go_eth
}

func NewWeb3(provider providers.ProviderInterface) *Web3 {
	go_web3 := new(Web3)
	go_web3.Provider = provider
	go_web3.Eth = NewEth(provider)
	go_web3.Net = net.NewNet(provider)
	go_web3.Personal = personal.NewPersonal(provider)
	go_web3.Utils = utils.NewUtils(provider)
	return go_web3
}

// add eth_sendRawTransaction
func (eth *Eth) SendRawTransaction(params string) (string, error) {

	pointer := &dto.RequestResult{}

	err := eth.provider.SendRequest(&pointer, "eth_sendRawTransaction", params)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}
