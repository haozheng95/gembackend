package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/gembackend/conf"
	"github.com/gembackend/gembackendlog"
	"github.com/gembackend/models/eth_query"
	"github.com/gembackend/models/exchange"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var log = gembackendlog.Log

func init() {
	maxIdle := 30
	maxConn := 30

	orm.RegisterModel(new(eth_query.Address), new(eth_query.Block),
		new(eth_query.TokenAddress), new(eth_query.TokenTx), new(eth_query.Tx),
		new(eth_query.TxExtraInfo), new(exchange.EthToken), new(exchange.MainChain))
	orm.RegisterDriver("mysql", orm.DRMySQL)
	// eth query -----------------------
	orm.RegisterDataBase("default",
		"mysql", conf.MysqlUser+":"+conf.MysqlPasswd+
			"@tcp("+conf.MysqlHost+":"+conf.MysqlPort+")/eth_query?charset=utf8", maxIdle, maxConn)

	// exchange --------------------------
	orm.RegisterDataBase("exchange",
		"mysql", conf.MysqlUser+":"+conf.MysqlPasswd+
			"@tcp("+conf.MysqlHost+":"+conf.MysqlPort+")/exchange?charset=utf8", maxIdle, maxConn)
}

func CreateTable() {
	orm.RunSyncdb("default", true, true)
	orm.RunSyncdb("exchange", true, true)
}

func AutoInsertData(dbname, tablename string) {
	o := orm.NewOrm()
	err := o.Using(dbname)
	if err != nil {
		log.Fatalf("error:%s", err)
		return
	}
	switch dbname {
	case "exchange":
		switch tablename {
		case "eth_token":
			eth_token(o)
		case "main_chain":
			main_chain(o)
		}
	case "eth_query":
		switch tablename {
		case "address":
		case "block":
		case "token_address":
		case "token_tx":
		case "tx":
		case "tx_extra_info":
		}
	}
}

func eth_token(o orm.Ormer) {
	data := []exchange.EthToken{
		{TokenName: "EOS", TokenFullName: "eos", ContractAddr: "0x86fa049857e0209aa7d9e616f7eb3b3b78ecfdb0", TokenDecimal: 18},
		{TokenName: "TRX", TokenFullName: "tron", ContractAddr: "0xf230b790e05390fc8295f4d3f60332c93bed42e2", TokenDecimal: 6},
		{TokenName: "OMG", TokenFullName: "omisego", ContractAddr: "0xd26114cd6ee289accf82350c8d8487fedb8a0c07", TokenDecimal: 18},
		{TokenName: "SNT", TokenFullName: "status", ContractAddr: "0x744d70fdbe2ba4cf95131626614a1763df805b9e", TokenDecimal: 18},
		{TokenName: "MKR", TokenFullName: "Maker", ContractAddr: "0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2", TokenDecimal: 18},
		{TokenName: "ZRX", TokenFullName: "0x", ContractAddr: "0xe41d2489571d322189246dafa5ebde1f4699f498", TokenDecimal: 18},
		{TokenName: "REP", TokenFullName: "augur", ContractAddr: "0xe94327d07fc17907b4db788e5adf2ed424addff6", TokenDecimal: 18},
		{TokenName: "BAT", TokenFullName: "basic-attention-token", ContractAddr: "0x0d8775f648430679a709e98d2b0cb6250d2887ef", TokenDecimal: 18},
		{TokenName: "SALT", TokenFullName: "salt", ContractAddr: "0x4156d3342d5c385a87d264f90653733592000581", TokenDecimal: 8},
		{TokenName: "FUN", TokenFullName: "funfair", ContractAddr: "0x419d0d8bdd9af5e606ae2232ed285aff190e711b", TokenDecimal: 8},
	}
	if num, err := o.InsertMulti(len(data), data); err == nil {
		log.Infof("insert row : %d", num)
	} else {
		log.Errorf("insert error : %s", err)
	}
}

func main_chain(o orm.Ormer) {
	data := []exchange.MainChain{
		{Name: "usdt", FullName: "tether", Decimal: 8, Created: time.Now(), Usdt: "1"},
		{Name: "btc", FullName: "bitcoin", Decimal: 8, Created: time.Now(), Usdt: "0"},
		{Name: "eth", FullName: "ethereum", Decimal: 18, Created: time.Now(), Usdt: "0"},
		{Name: "bch", FullName: "bitcoin-cash", Decimal: 8, Created: time.Now(), Usdt: "0"},
		{Name: "eos", FullName: "eos", Decimal: 18, Created: time.Now(), Usdt: "0"},
	}
	if num, err := o.InsertMulti(len(data), data); err == nil {
		log.Infof("insert row : %d", num)
	} else {
		log.Errorf("insert error : %s", err)
	}
}
