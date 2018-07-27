package exchange

import "github.com/astaxie/beego/orm"

type mainchainFullname struct {
	FullName string
}

func GetAllCoinFullNaml() (r []mainchainFullname) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("full_name").From("main_chain")
	sql := qb.String()
	o := orm.NewOrm()
	o.Using(databases)
	o.Raw(sql).QueryRows(&r)
	return
}

func GetMainChainCny(coin string) (cny string) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("cny").From("main_chain").Where("full_name=?")
	sql := qb.String()
	o := orm.NewOrm()
	o.Using(databases)
	o.Raw(sql, coin).QueryRow(&cny)
	return
}

type allTokenName struct {
	TokenName string
}

func GetAllTokenName() (r []allTokenName) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("token_name").From("eth_token")
	sql := qb.String()
	o := orm.NewOrm()
	o.Using(databases)
	o.Raw(sql).QueryRows(&r)
	return
}
