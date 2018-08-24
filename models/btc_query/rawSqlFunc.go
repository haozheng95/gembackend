//software: GoLand
//file: rawSqlFunc.go
//time: 2018/8/23 下午6:01
package btc_query

import "github.com/astaxie/beego/orm"

const insertnumber = 500

func InsertMulTradeCollection(data []*TradeCollection) (err error) {
	o := orm.NewOrm()
	o.Using(databases)
	start := 0
	end := insertnumber
	long := len(data)
	log.Debug("total long ====", long)
	for start < long {
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
		if end > long {
			end = long
		}
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
		log.Debug("insert UnspentVout start === ", start)
		log.Debug("insert UnspentVout end   === ", end)
		if num, err := o.InsertMulti(end-start, data[start:end]); err == nil {
			log.Infof("insert row : %d", num)
		} else {
			log.Errorf("insert error : %s", err)
			break
		}
		start = end
		end += insertnumber
		if end > long {
			end = long
		}
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
