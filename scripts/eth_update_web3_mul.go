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
	height uint64
}

func (updater *EthupdaterMul) Start(startHeight chan uint64) {
	for true {
		updater.height = <-startHeight
		updater.Forever()
		startHeight <- updater.height + 1
	}

}

func (updater *EthupdaterMul) Forever() {
	//updater.TableBlock = updater.TableBlock.SelectMaxHeight()
	//dbHehght := updater.TableBlock.BlockHeight
	//height := MaxIntByString(updater.StartHeight, dbHehght)
	rpcHeight, err := updater.connection.Eth.GetBlockNumber()
	height := updater.height
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

func StartEthupdaterMul(height uint64) {
	updater := NewEthupdaterMul()
	dbHeight := updater.TableBlock.SelectMaxHeight().BlockHeight
	log.Infof("db height = %s , input height = %d", dbHeight, height)
	height = MaxIntByString(height, dbHeight)
	log.Infof("update height = %d", height)
	c := make(chan uint64, 5)
	t := make(chan int)

	c <- height
	for i:=0;i<5 ;i++  {
		go updater.Start(c)
	}
	<-t
	log.Error("error exit")
}
