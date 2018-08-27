//software: GoLand
//file: block_btc.go
//time: 2018/8/27 下午3:51
package btc_query

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type BlockBtc struct {
	Id          int64
	BlockHeight int64  `orm:"index"`
	BlockHash   string `orm:"unique"`
	Previous    string `orm:"index"`
	ConfirmTime int64
	Updated     time.Time
	Nonce       uint32
}

func NewBlockBtc(height int64, blockhash string, prehash string, confirm int64, nonce uint32) (bc *BlockBtc) {
	bc = &BlockBtc{BlockHeight: height, BlockHash: blockhash, Previous: prehash, ConfirmTime: confirm, Nonce: nonce, Updated: time.Now()}
	return
}

func (b *BlockBtc) Insert() {
	o := orm.NewOrm()
	o.Using(databases)
	num, err := o.Insert(b)
	if err != nil {
		log.Warning(err)
	} else {
		log.Debug(num)
	}
}
