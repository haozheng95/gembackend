//software: GoLand
//file: rawSqlFunc.go
//time: 2018/8/23 下午6:01
package btc_query

import "github.com/astaxie/beego/orm"

func InsertAddress(data []*AddressBtc) {
	o := orm.NewOrm()
	o.Using(databases)
	if num, err := o.InsertMulti(len(data), data); err == nil {
		log.Infof("insert row : %d", num)
	} else {
		log.Errorf("insert error : %s", err)
		log.Debug(err == nil)
	}
}

type Txs struct {
	TradeCollection
	Comment string
}

func GetTxs(walletId string, size, index int) (res []*Txs, r int64) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("t1.*,t3.comment").From("trade_collection as t1").LeftJoin("address_btc as t2").
		On("t1.addr=t2.addr").LeftJoin("tx_extra_info_btc as t3").On("t1.txid=t3.tx_hash").
		Where("t2.wallet_id=?").OrderBy("id").Desc().Limit(size).Offset(index)
	o := orm.NewOrm()
	o.Using(databases)
	//log.Debug(qb.String())
	r, err := o.Raw(qb.String(), walletId).QueryRows(&res)
	if err != nil {
		log.Warning("error ==== ", err)
	}
	return
}

func GetUnspent(walletId string) (res []*UnspentVout) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("t1.*").From("unspent_vout as t1").LeftJoin("address_btc as t2").On("t1.address=t2.addr").
		Where("t1.spent=0").And("t2.wallet_id=?")
	o := orm.NewOrm()
	o.Using(databases)
	//log.Debug(qb.String())
	_, err := o.Raw(qb.String(), walletId).QueryRows(&res)
	if err != nil {
		log.Warning("error ==== ", err)
	}
	return
}

func GetTxInfo(txhash string) (res []*TradingParticulars) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("*").From("trading_particulars").Where("txid=?")
	o := orm.NewOrm()
	o.Using(databases)
	//log.Debug(qb.String())
	_, err := o.Raw(qb.String(), txhash).QueryRows(&res)
	if err != nil {
		log.Warning("GetTxInfo error ==== ", err)
	}
	return
}

func GetUserInfo(walletId string) (res []*AddressBtc) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("*").From("address_btc").Where("wallet_id=?")
	o := orm.NewOrm()
	o.Using(databases)
	_, err := o.Raw(qb.String(), walletId).QueryRows(&res)
	if err != nil {
		log.Warning("GetUserInfo error ==== ", err)
	}
	return
}

func UpdateAddr(addr, amount string) (err error) {
	//addr = addr
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Update("address_btc").Set("amount=?", "unconfirm_amount=0").Where("addr=?")
	o := orm.NewOrm()
	o.Using(databases)
	log.Debug(amount, addr)
	_, err = o.Raw(qb.String(), amount, addr).Exec()
	if err != nil {
		log.Warning("update error ==== ", err)
	}
	return
}

func GetAllUnspent(addr string) (res []*struct{ Value string }) {
	addr = "`" + addr + "`"
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("value").From("unspent_vout").Where("spent=0").And("address=?")
	o := orm.NewOrm()
	o.Using(databases)
	_, err := o.Raw(qb.String(), addr).QueryRows(&res)
	if err != nil {
		log.Warning("get unspents error ==== ", err)
	}
	return
}

// btc address exist
func GetBtcAddrExist(addr string) bool {
	o := orm.NewOrm()
	o.Using(databases)
	qs := o.QueryTable("address_btc")
	res := qs.Filter("addr", addr).Exist()
	if res {
		log.Debug("addr =====", addr)
	}
	return res
}

func CurrBlockNum() (num int64) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("block_height").From("block_btc").OrderBy("block_height").Desc().Limit(1)
	o := orm.NewOrm()
	o.Using(databases)
	err := o.Raw(qb.String()).QueryRow(&num)
	if err != nil {
		log.Warning(err)
	}
	return
}

func Getblockhash(height int64) (blockhash string) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("block_hash").From("block_btc").Where("block_height=?").Limit(1)
	o := orm.NewOrm()
	o.Using(databases)
	err := o.Raw(qb.String(), height).QueryRow(&blockhash)
	if err != nil {
		log.Warning(err)
		blockhash = ""
	}
	return
}

func Deleteblockhash(blockhash string) (err error) {
	qb1, _ := orm.NewQueryBuilder("mysql")
	qb2, _ := orm.NewQueryBuilder("mysql")
	qb3, _ := orm.NewQueryBuilder("mysql")
	qb4, _ := orm.NewQueryBuilder("mysql")
	qb1.Delete().From("trade_collection").Where("block_hash=?")
	qb2.Delete().From("trading_particulars").Where("block_hash=?")
	qb3.Delete().From("unspent_vout").Where("block_hash=?")
	qb4.Delete().From("block_btc").Where("block_hash=?")

	o := orm.NewOrm()
	o.Using(databases)

	_, err = o.Raw(qb1.String(), blockhash).Exec()
	if err != nil {
		log.Warning("qb1 ", err)
	}
	_, err = o.Raw(qb2.String(), blockhash).Exec()
	if err != nil {
		log.Warning("qb2 ", err)
	}
	_, err = o.Raw(qb3.String(), blockhash).Exec()
	if err != nil {
		log.Warning("qb3 ", err)
	}
	_, err = o.Raw(qb4.String(), blockhash).Exec()
	if err != nil {
		log.Warning("qb4 ", err)
	}
	return
}

func Deletetxrecord(txid string) (err error) {
	qb1, _ := orm.NewQueryBuilder("mysql")
	qb2, _ := orm.NewQueryBuilder("mysql")
	qb3, _ := orm.NewQueryBuilder("mysql")
	qb1.Delete().From("trading_particulars").Where("txid=?")
	qb2.Delete().From("trade_collection").Where("txid=?")
	qb3.Delete().From("unspent_vout").Where("txid = ?")
	o := orm.NewOrm()
	o.Using(databases)
	//log.Debug(qb1.String())
	_, err = o.Raw(qb1.String(), txid).Exec()
	if err != nil {
		log.Warning("qb1 ", err)
	}
	_, err = o.Raw(qb2.String(), txid).Exec()
	if err != nil {
		log.Warning("qb2 ", err)
	}
	_, err = o.Raw(qb3.String(), txid).Exec()
	if err != nil {
		log.Warning("qb3 ", err)
	}
	return
}

const insertnumber = 500

func InsertMulTradingParticulars(data []*TradingParticulars) (err error) {
	o := orm.NewOrm()
	o.Using(databases)
	start := 0
	end := insertnumber
	long := len(data)
	log.Debug("total long ====", long)
	for start < long {

		if end > long {
			end = long
		}
		//log.Debug("insert TradingParticulars start === ", start)
		//log.Debug("insert TradingParticulars end   === ", end)

		if _, err := o.InsertMulti(end-start, data[start:end]); err == nil {
			//log.Infof("TradingParticulars insert row : %d", num)
		} else {
			log.Errorf("TradingParticulars insert error : %s", err)
			break
		}
		start += insertnumber
		end += insertnumber
	}
	return
}

func InsertMulTradeCollection(data []*TradeCollection) (err error) {
	o := orm.NewOrm()
	o.Using(databases)
	start := 0
	end := insertnumber
	long := len(data)
	log.Debug("total long ====", long)
	for start < long {

		if end > long {
			end = long
		}
		//log.Debug("insert TradeCollection start === ", start)
		//log.Debug("insert TradeCollection end   === ", end)

		if _, err := o.InsertMulti(end-start, data[start:end]); err == nil {
			//log.Infof("insert row : %d", num)
		} else {
			log.Errorf("insert error : %s", err)
			break
		}
		start += insertnumber
		end += insertnumber
	}
	return
}

// insert mul umspent txs
func InsertMulUnspentVout(data []*UnspentVout) (err error) {
	o := orm.NewOrm()
	o.Using(databases)
	start := 0
	end := insertnumber
	long := len(data)
	log.Debug("total long ====", long)

	for start < long {
		if end > long {
			end = long
		}
		//log.Debug("insert UnspentVout start === ", start)
		//log.Debug("insert UnspentVout end   === ", end)

		if _, err := o.InsertMulti(end-start, data[start:end]); err == nil {
			//log.Infof("insert row : %d", num)
		} else {
			log.Errorf("insert error : %s", err)
			log.Debug(err == nil)
			break
		}
		start = end
		end += insertnumber

	}

	return
}

// insert unspent tx
func InsertUnspentVout(args ...string) (err error) {
	qb, err := orm.NewQueryBuilder("mysql")
	for i := range args {
		args[i] = "'" + args[i] + "'"
	}
	args = append(args, "now()")

	qb.InsertInto("unspent_vout", "txid", "`index`", "value", "address", "updated").Values(args...)
	o := orm.NewOrm()
	o.Using(databases)
	sql := qb.String()
	//log.Debug(sql)
	_, err = o.Raw(sql).Exec()
	if err != nil {
		log.Warning(err)
	}
	return
}

// update state
func UpdateUnspentVout(txid, n, spenttx string) (err error) {
	qb, err := orm.NewQueryBuilder("mysql")
	qb.Update("unspent_vout").Set("spent=1", "spent_txid=?").Where("txid=?").And("`index`=?")
	o := orm.NewOrm()
	o.Using(databases)
	_, err = o.Raw(qb.String(), spenttx, txid, n).Exec()
	if err != nil {
		log.Warning(err)
	}
	return
}
