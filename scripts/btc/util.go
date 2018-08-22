//software: GoLand
//file: util.go
//time: 2018/8/22 上午10:43
package btc

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

func makeChHash(str string) (hash chainhash.Hash) {
	h, _ := hex.DecodeString(str)
	for j, v3 := range h {
		hash[chainhash.HashSize-1-j] = v3
	}
	return
}
