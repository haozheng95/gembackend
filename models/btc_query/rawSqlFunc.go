//software: GoLand
//file: rawSqlFunc.go
//time: 2018/8/23 下午6:01
package btc_query

import "github.com/astaxie/beego/orm"

// insert unspent tx
func InsertUnspentVout(args ...string) (err error) {
	qb, err := orm.NewQueryBuilder("mysql")
	qb.InsertInto("unspent_vout", "txid", "index", "value", "address").Values(args...)
	o := orm.NewOrm()
	o.Using(databases)
	_, err = o.Raw(qb.String()).Exec()
	if err != nil {
		log.Warning(err)
	}
	return
}

// update state
func UpdateUnspentVout(txid, n, spenttx string) (err error) {
	qb, err := orm.NewQueryBuilder("mysql")
	qb.Update("unspent_vout").Set("spent=1", "spent_txid=?").Where("txid=?").And("index=?")
	o := orm.NewOrm()
	o.Using(databases)
	_, err = o.Raw(qb.String(), spenttx, txid, n).Exec()
	if err != nil {
		log.Warning(err)
	}
	return
}
