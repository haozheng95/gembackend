package exchange

import "github.com/astaxie/beego/orm"

func GetCoinFee(coin string) (fee string) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("fee").From("main_chain").Where("name=?")
	sql := qb.String()
	o := orm.NewOrm()
	o.Using(databases)
	o.Raw(sql, coin).QueryRow(&fee)
	return
}

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

type allTokenFullName struct {
	TokenFullName string
}

func GetAllFullTokenName() (r []allTokenFullName) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("token_full_name").From("eth_token")
	sql := qb.String()
	o := orm.NewOrm()
	o.Using(databases)
	o.Raw(sql).QueryRows(&r)
	return
}

// get token gas limit
func GetTokenGasLimit(contractAddr string) (r string) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("gas_limit").From("eth_token").Where("contract_addr=?")
	sql := qb.String()
	o := orm.NewOrm()
	o.Using(databases)
	o.Raw(sql, contractAddr).QueryRow(&r)
	return
}

// get main chain gas limit
func GetMainChainGasLimit(coin string) (r string) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("gas_limit").From("main_chain").Where("name=?")
	sql := qb.String()
	o := orm.NewOrm()
	o.Using(databases)
	o.Raw(sql, coin).QueryRow(&r)
	return
}

// get main chain gas price
func GetMainChainGasPrice(coin string) (r string) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("gas_price").From("main_chain").Where("name=?")
	sql := qb.String()
	o := orm.NewOrm()
	o.Using(databases)
	o.Raw(sql, coin).QueryRow(&r)
	return
}

// get main chain cny
func GetMainChainCnyByCoinName(coin string) (r string) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("cny").From("main_chain").Where("name=?")
	sql := qb.String()
	o := orm.NewOrm()
	o.Using(databases)
	o.Raw(sql, coin).QueryRow(&r)
	return
}

// get eth token cny
func GetEthTokenCnyByContractAddr(contractAddr string) (r string) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("cny").From("eth_token").Where("contract_addr=?")
	sql := qb.String()
	o := orm.NewOrm()
	o.Using(databases)
	o.Raw(sql, contractAddr).QueryRow(&r)
	return
}

// update main chain gasprice and gaslimit
func UpdateMainChainGasPriceAndGasLimit(coin, gaslimit, gasprice, fee string) {
	qb, _ := orm.NewQueryBuilder("mysql")
	//qb.Update("gas_limit=?", "gas_price=?", "fee=?").From("main_chain").Where("name=?")
	qb.Update("main_chain").Set("gas_limit=?", "gas_price=?", "fee=?").Where("name=?")
	sql := qb.String()
	o := orm.NewOrm()
	o.Using(databases)
	//log.Debug(sql)
	o.Raw(sql, gaslimit, gasprice, fee, coin).Exec()
}
