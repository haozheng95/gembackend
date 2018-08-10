package eth_query

import (
	"github.com/astaxie/beego/orm"
	"time"
)

func UpdateAddressAmount(unconfirm, addr, nonce, amount string) (err error) {
	qb, err := orm.NewQueryBuilder("mysql")
	qb.Update("address").Set("unconfirm_amount=?", "nonce=?", "amount=?").Where("addr=?")
	o := orm.NewOrm()
	o.Using(databases)
	_, err = o.Raw(qb.String(), unconfirm, nonce, amount, addr).Exec()
	if err != nil {
		log.Warning(err)
	}
	return
}

func UpdateTokenAddressAmount(amount, unconfirm, addr, contractaddr string) (err error) {
	qb, err := orm.NewQueryBuilder("mysql")
	qb.Update("token_address").Set("unconfirm_amount=?", "amount=?").
		Where("addr=?").And("contract_addr=?")
	o := orm.NewOrm()
	o.Using(databases)
	o.Raw(qb.String(), amount, unconfirm, addr, contractaddr)
	return
}

// delete token tx
func DeleteTokenTxbyHashWhere(hash string) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Delete("token_tx").Where("tx_hash=?").And("log_index=-1")
	sql := qb.String()
	o := orm.NewOrm()
	o.Using(databases)
	_, err := o.Raw(sql).Exec()
	if err != nil {
		log.Warning(err)
	}
}

// insert token tx
func InsertTokenTx(args ...string) (err error) {
	qb, _ := orm.NewQueryBuilder("mysql")
	for i := range args {
		args[i] = "'" + args[i] + "'"
	}
	qb.InsertInto("token_tx", "`from`", "`to`", "amount", "input_data",
		"nonce", "gas_limit", "gas_price", "gas_used", "fee", "tx_hash", "block_hash",
		"confirm_time", "created", "block_state", "is_token", "log_index", "contract_addr",
		"`decimal`").Values(args...)
	sql := qb.String()
	o := orm.NewOrm()
	o.Using(databases)
	_, err = o.Raw(sql).Exec()
	if err != nil {
		log.Warning(err)
	}
	return
}

// update one tx
func UpdateTxOneRawByHash(gasUsed, fee, status, blockHash, blockNumber, hash string) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Update("tx").Set("gas_used=?",
		"fee=?", "block_height=?", "confirm_time=?", "block_state=?", "tx_state=?",
		"block_hash=?").Where("tx_hash=?")
	sql := qb.String()
	o := orm.NewOrm()
	o.Using(databases)
	_, err := o.Raw(sql, gasUsed, fee, blockNumber, time.Now().Unix(), "1", status, blockHash, hash).Exec()
	if err != nil {
		log.Warning(err)
		log.Warning(sql)
	}
}

// select one tx
func GetTxOneRawByHash(hash string) (st struct {
	GasLimit, GasPrice string
}, err error) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("gas_limit", "gas_price").From("tx").Where("tx_hash=?")
	sql := qb.String()
	o := orm.NewOrm()
	o.Using(databases)
	err = o.Raw(sql, hash).QueryRow(&st)
	if err != nil {
		log.Error(err)
	}
	return
}

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
	Comment     string
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
func GetEthTxrecord(addr string, page uint64, size uint64) (txs []*ethtxrecordst, r int64) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("t1.id", "t1.from", "t1.to", "t1.amount", "t1.fee", "t1.nonce", "t1.tx_hash",
		"t1.block_height", "t1.confirm_time", "t1.created", "t1.tx_state", "t1.is_token", "t2.comment",
		"0 as collection").From("tx as t1").
		LeftJoin("tx_extra_info as t2").
		On("t1.tx_hash = t2.tx_hash").
		Where("t1.from = ?")
	sql1 := qb.String() + " union all "
	qb, _ = orm.NewQueryBuilder("mysql")
	qb.Select("t1.id", "t1.from", "t1.to", "t1.amount", "t1.fee", "t1.nonce", "t1.tx_hash",
		"t1.block_height", "t1.confirm_time", "t1.created", "t1.tx_state", "t1.is_token", "t2.comment",
		"1 as collection").From("tx as t1").
		LeftJoin("tx_extra_info as t2").
		On("t1.tx_hash = t2.tx_hash").
		Where("t1.to = ?").
		OrderBy("created DESC", "id").Desc().Limit(int(size)).Offset(int(page))
	sql2 := qb.String()
	sql := sql1 + sql2
	o := orm.NewOrm()
	o.Using(databases)
	r, err := o.Raw(sql, addr, addr).QueryRows(&txs)
	if err != nil {
		log.Errorf("GetEthTxrecord error %s", err)
		return
	}
	return
}

func GetEthTokenTxrecord(addr string, contract string, page uint64, size uint64) (txs []*ethtokentxrecordst, r int64) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("t1.id", "t1.tx_hash", "t1.from", "t1.to", "t1.amount", "t1.decimal", "t1.fee",
		"t1.block_height", "t1.confirm_time", "t1.created", "t1.log_index", "t1.tx_state", "t1.is_token",
		"t2.comment", "0 as collection").
		From("token_tx as t1").
		LeftJoin("tx_extra_info as t2").
		On("t1.tx_hash = t2.tx_hash").
		Where("t1.from = ?").And("t1.contract_addr = ?")
	sql1 := qb.String() + " union all "
	qb, _ = orm.NewQueryBuilder("mysql")
	qb.Select("t1.id", "t1.tx_hash", "t1.from", "t1.to", "t1.amount", "t1.decimal", "t1.fee",
		"t1.block_height", "t1.confirm_time", "t1.created", "t1.log_index", "t1.tx_state", "t1.is_token",
		"t2.comment", "1 as collection").
		From("token_tx as t1").
		LeftJoin("tx_extra_info as t2").
		On("t1.tx_hash = t2.tx_hash").
		Where("t1.to = ?").And("t1.contract_addr = ?").
		OrderBy("created DESC", "id").Desc().
		Limit(int(size)).Offset(int(page))

	sql2 := qb.String()
	sql := sql1 + sql2
	o := orm.NewOrm()
	o.Using(databases)
	r, err := o.Raw(sql, addr, contract, addr, contract).QueryRows(&txs)
	if err != nil {
		log.Errorf("GetEthTokenTxrecord error %s", err)
		return
	}
	return
}

func UpdateAddress(unconfirm, addr string) (err error) {
	qb, err := orm.NewQueryBuilder("mysql")
	qb.Update("address").Set("unconfirm_amount=?").Where("addr=?")
	o := orm.NewOrm()
	o.Using(databases)
	o.Raw(qb.String(), unconfirm, addr).Exec()
	return
}

func UpdateTokenAddress(unconfirm, addr, contractaddr string) (err error) {
	qb, err := orm.NewQueryBuilder("mysql")
	qb.Update("token_address").Set("unconfirm_amount=?").
		Where("addr=?").And("contract_addr=?")
	o := orm.NewOrm()
	o.Using(databases)
	o.Raw(qb.String(), unconfirm, addr, contractaddr)
	return
}

// 判断eth用户是否存在
func GetEthAddrExist(addr string) bool {
	o := orm.NewOrm()
	o.Using(databases)
	qs := o.QueryTable("address")
	return qs.Filter("addr", addr).Exist()
}

type tokentxinfores struct {
	TokenTx
	ConfirmNum string
}

// get token tx info
func GetTokenTxinfo(txhash string) (r []*tokentxinfores) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("`from`", "`to`", "amount", "nonce", "input_data", "gas_limit", "gas_price",
		"gas_used", "fee", "tx_hash", "block_hash", "block_height", "confirm_time", "tx_state",
		"is_token", "log_index", "contract_addr", "`decimal`", "created").
		From("token_tx").Where("tx_hash=?")
	o := orm.NewOrm()
	o.Using(databases)
	o.Raw(qb.String(), txhash).QueryRows(&r)
	//log.Debug(qb.String())
	//log.Debug(i, err)
	return
}

// get block height info
func GetBlockHeight() (height string) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("block_height").From("block").
		OrderBy("block_height").Desc().Limit(1)
	o := orm.NewOrm()
	o.Using(databases)
	o.Raw(qb.String()).QueryRow(&height)
	return
}

type txinfores struct {
	Tx
	ConfirmNum string
}

// get eth tx info
func GetTxInfo(txhash string) (r []*txinfores) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("`from`", "`to`", "amount", "nonce", "input_data", "gas_limit", "gas_price",
		"gas_used", "fee", "tx_hash", "block_hash", "block_height", "confirm_time", "tx_state",
		"is_token").From("tx").Where("tx_hash=?")
	o := orm.NewOrm()
	o.Using(databases)
	o.Raw(qb.String(), txhash).QueryRows(&r)
	return
}

// get eth info by wallet id
func GetEthInfoByWalletId(walletId string) (st struct{ Nonce, Amount, UnconfirmAmount, Addr, Decimal string }) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("nonce", "amount", "unconfirm_amount", "addr", "`decimal`").
		From("address").Where("wallet_id=?")
	o := orm.NewOrm()
	o.Using(databases)
	o.Raw(qb.String(), walletId).QueryRow(&st)
	return
}

// get token for eth info by wallet id
func GetEthTokenInfoByWalletId(walletId, contractAddr string) (st struct {
	Nonce,
	Amount,
	UnconfirmAmount,
	TokenAmount,
	TokenUnconfirmAmount,
	TokenName,
	Decimal string
}) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("t1.nonce", "t1.amount", "t1.unconfirm_amount", "t2.amount as token_amount",
		"t2.unconfirm_amount as token_unconfirm_amount", "t2.decimal", "t2.token_name").From("address as t1").
		LeftJoin("token_address as t2").On("t1.addr = t2.addr").
		Where("t1.wallet_id=?").And("t2.contract_addr=?")
	o := orm.NewOrm()
	o.Using(databases)
	o.Raw(qb.String(), walletId, contractAddr).QueryRow(&st)
	return
}

// get all token info
func GetAllTokenInfoWithUser(walletId string, begin, size int) (st []*struct {
	ContractAddr, Amount, UnconfirmAmount, TokenName, Decimal string
}) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("contract_addr", "amount", "unconfirm_amount", "token_name", "`decimal`").
		From("token_address").Where("wallet_id=?").And("added=?").Limit(size).Offset(begin)
	o := orm.NewOrm()
	o.Using(databases)
	//log.Debug(qb.String())
	_, err := o.Raw(qb.String(), walletId, 1).QueryRows(&st)
	if err != nil {
		log.Debug(err)
	}
	return
}
