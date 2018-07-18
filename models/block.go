package models

import (
	"time"
	"github.com/astaxie/beego/orm"
)

type Block struct {
	Id          int64
	BlockHeight string    `orm:"index"`
	BlockHash   string    `orm:"unique"`
	TimeStamp   string    `orm:"index"`
	ParentHash  string    `orm:"index"`
	Miner       string    `orm:"index"`
	MixHash     string
	Nonce       string
	ExtraData   string
	GasLimit    string
	GasUsed     string
	Size        string
	Created     time.Time `orm:"auto_now_add;type(datetime)"`
}

func (block *Block) SelectMaxHeight() *Block {
	o := orm.NewOrm()
	qs := o.QueryTable(block)
	qs.OrderBy("-id").Limit(1).One(block)
	return block
}

func (block *Block) SelectRawByHeight(height uint64) *Block {
	o := orm.NewOrm()
	qs := o.QueryTable(block)
	qs.Filter("block_height", height).Limit(1).One(block)
	return block
}

func (block *Block) DeleteOneRaw(blockHash string) *Block {
	o := orm.NewOrm()
	qs := o.QueryTable(block)
	num, err := qs.Filter("block_hash", blockHash).Delete()
	if err != nil {
		log.Errorf("block delete error %s", err)
	}
	log.Debugf("block delete num = %startscript", num)
	return block
}

func (block *Block) InsertOneRaw(data *Block) *Block {
	o := orm.NewOrm()
	data.Id = 0
	data.Created = time.Now()
	id, err := o.Insert(data)
	if err != nil {
		log.Errorf("block insert error %s", err)
	}
	log.Debugf("block insert id %d startscript", id)
	return block
}
