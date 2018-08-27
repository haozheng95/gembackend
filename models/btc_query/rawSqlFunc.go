//software: GoLand
//file: rawSqlFunc.go
//time: 2018/8/23 下午6:01
package btc_query

import "github.com/astaxie/beego/orm"

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
		log.Debug("insert TradingParticulars start === ", start)
		log.Debug("insert TradingParticulars end   === ", end)

		if num, err := o.InsertMulti(end-start, data[start:end]); err == nil {
			log.Infof("TradingParticulars insert row : %d", num)
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
		log.Debug("insert TradeCollection start === ", start)
		log.Debug("insert TradeCollection end   === ", end)

		if num, err := o.InsertMulti(end-start, data[start:end]); err == nil {
			log.Infof("insert row : %d", num)
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
		log.Debug("insert UnspentVout start === ", start)
		log.Debug("insert UnspentVout end   === ", end)

		if num, err := o.InsertMulti(end-start, data[start:end]); err == nil {
			log.Infof("insert row : %d", num)
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
