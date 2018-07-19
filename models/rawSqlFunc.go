package models

import "time"

type ethtxrecordst struct {
	Id          int64
	From        string
	To          string
	Amount      string
	Nonce       string
	Fee         string
	TxHash      string
	BlockHeight string
	ConfirmTime string
	Created     time.Time
	BlockState  int
	TxState     int
	IsToken     int
	Collection  int
}

type ethtokentxrecordst struct {
	ethtxrecordst
	LogIndex string
	Decimal  string
}

// eth-query 相关函数
func GetEthTxrecord(addr string, page uint64, size uint64) (txs *[]ethtxrecordst, r int64) {
	sql := "select id,`from`, `to`, amount, nonce,fee,tx_hash,block_height,confirm_time,created,tx_state,is_token, 0 as collection from tx as t1  where t1.from = ?  union all select id,`from`, `to`, amount, nonce,fee,tx_hash,block_height,confirm_time,created,tx_state,is_token, 1 as collection from tx as t2 where t2.to = ? order by created desc,id desc limit ?,?"
	o.Using("default")
	txs = new([]ethtxrecordst)

	r, err := o.Raw(sql, addr, addr, page, size).QueryRows(txs)
	if err != nil {
		log.Errorf("GetEthTxrecord error %s", err)
		return
	}
	log.Info(r)
	return
}

func GetEthTokenTxrecord(addr string, contract string, page uint64, size uint64) (txs *[]ethtokentxrecordst, r int64) {
	sql := "select id,tx_hash,`from`,`to`,amount,nonce,`decimal`,fee,block_height,confirm_time,created,log_index,tx_state,is_token, 0 as collection from token_tx as t1 where t1.from=? and t1.contract_addr=? union all select id,tx_hash,`from`,`to`,amount,nonce,`decimal`,fee,block_height,confirm_time,created,log_index,tx_state,is_token, 1 as collection from token_tx as t2 where t2.to=? and t2.contract_addr=? order by created desc,id desc limit ?,?"
	o.Using("default")
	txs = new([]ethtokentxrecordst)
	r, err := o.Raw(sql, addr, contract, addr, contract, page, size).QueryRows(txs)
	if err != nil {
		log.Errorf("GetEthTokenTxrecord error %s", err)
		return
	}
	log.Info(r)
	return
}

// 判断eth用户是否存在
func GetEthAddrExist(addr string) bool {
	o.Using("default")
	qs := o.QueryTable("address")
	return qs.Filter("addr", addr).Exist()
}
