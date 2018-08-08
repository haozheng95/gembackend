//software: GoLand
//file: eth_update_api_mul.go
//time: 2018/8/8 上午10:45
package eth

import (
	"github.com/gembackend/models/eth_query"
	"time"
)

type EthUpdateApiMul struct {
	EthUpdater
}

func (updater *EthUpdateApiMul) Forever(height uint64, tolerance *chan uint64) {
	defer func() {
		if err := recover(); err != nil {
			*tolerance <- height
		}
	}()
keep:
	blockInfo, err := updater.getBlockInfo(height)
	log.Warningf("\n The height %d run", height)
	if err != nil {
		log.Warningf("get block info error, will retry； %s", err)
		goto keep
	}

	if blockInfo.Result != nil {
		updater.BeginUpdateBlockInfo(blockInfo)
	} else {
		time.Sleep(time.Second * 5)
		goto keep
	}
}

func NewEthupdaterApiMul() *EthUpdateApiMul {
	u := new(EthUpdateApiMul)
	u.TableBlock = new(eth_query.Block)
	u.TableTx = new(eth_query.Tx)
	u.TableTokenTx = new(eth_query.TokenTx)
	u.TableAddress = new(eth_query.Address)
	u.TableTokenAddress = new(eth_query.TokenAddress)
	return u
}

func startEthUpdateApiMul(c chan uint64, e *EthUpdateApiMul, tolerance *chan uint64) {
	for true {
		height := <-c
		c <- height + 1
		log.Infof("the api update height == %d", height)
		e.Forever(height, tolerance)
	}
}

func StartEthApiMul(height uint64) {
	updater := NewEthupdaterApiMul()
	dbHeight := updater.TableBlock.SelectMaxHeight().BlockHeight
	log.Infof("db height = %s , input height = %d", dbHeight, height)
	height = MaxIntByString(height, dbHeight)
	log.Infof("api update height = %d", height)

	time.Sleep(time.Second)
	wg.Add(1)
	c := make(chan uint64, 5)
	tolerance := make(chan uint64)
	c <- height
	defer close(c)
	for i := 0; i < 10; i++ {
		go startEthUpdateApiMul(c, NewEthupdaterApiMul(), &tolerance)
	}
	i := 0
	select {
	case h := <-tolerance:
		i++
		log.Debugf("rescue count %d the hright == %d", i, h)
		c <- h
		go startEthUpdateApiMul(c, NewEthupdaterApiMul(), &tolerance)
	}
	wg.Wait()
}
