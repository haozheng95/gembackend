package scripts

import (
	"fmt"
	"github.com/gembackend/rpc"
	"time"
	"github.com/gembackend/models"
	"github.com/gembackend/conf"
	"github.com/regcostajr/go-web3"
	"github.com/regcostajr/go-web3/providers"
)

type EthupdaterMul struct {
	EthUpdaterWeb3
}

func (updater *EthupdaterMul) Forever(height uint64) {
	rpcHeight, err := updater.connection.Eth.GetBlockNumber()
	if err != nil {
		log.Errorf("web3 rpcHeight error %s", err)
		return
	}
Again:
	log.Infof("db height %d ==== rpc height %s", height, rpcHeight.String())
	if rpcHeight.Uint64() > height {
		hexHeight := fmt.Sprintf("0x%x", height)
		updater.parityParam["method"] = "eth_getBlockByNumber"
		updater.parityParam["params"] = []interface{}{hexHeight, true}
		blockInfoString := rpc.HttpPost(updater.parityParam)
		updater.rpcRes, err = rpc.FormatResponse(&blockInfoString)
		if err != nil {
			log.Errorf("FormatResponse error %s", err)
			panic(err)
		}
		updater.disposeBlockInfo()
		updater.disposeTransactions()
		log.Infof("block update success %d", height)
	} else {
		//Have a rest
		log.Info("block pending")
		time.Sleep(time.Second * 5)
		log.Info("again")
		goto Again
	}
	log.Infof("------------------ %d", height)

}

func NewEthupdaterMul() *EthupdaterMul {
	u := new(EthupdaterMul)
	u.TableBlock = new(models.Block)
	u.TableTx = new(models.Tx)
	u.TableTokenTx = new(models.TokenTx)
	u.TableAddress = new(models.Address)
	u.TableTokenAddress = new(models.TokenAddress)
	timeOut := conf.EthRpcTimeOut
	source := conf.EthRpcSecure
	url := conf.EthRpcHost + ":" + conf.EthRpcPort
	u.connection = web3.NewWeb3(providers.NewHTTPProvider(url, int32(timeOut), source))
	u.parityParam = map[string]interface{}{
		"id":      1,
		"jsonrpc": "2.0",
	}
	return u
}

func Start(startHeight chan uint64, updater *EthupdaterMul) {
	for true {
		t := <-startHeight
		log.Infof("height %d", t)
		startHeight <- t + 1
		updater.Forever(t)
	}
}

func StartEthupdaterMul(height uint64) {
	updater := NewEthupdaterMul()
	dbHeight := updater.TableBlock.SelectMaxHeight().BlockHeight
	log.Infof("db height = %s , input height = %d", dbHeight, height)
	height = MaxIntByString(height, dbHeight)
	log.Infof("update height = %d", height)

	time.Sleep(time.Second * 2)
	c := make(chan uint64, 5)
	wg.Add(1)
	c <- height
	for i := 0; i < 5; i++ {
		u := NewEthupdaterMul()
		go Start(c, u)
	}
	wg.Wait()
	log.Error("error exit")
}
