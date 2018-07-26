package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gembackend/conf"
	"github.com/gembackend/models/eth_query"
	"github.com/gembackend/models/exchange"
)

func init() {
	maxIdle := 30
	maxConn := 30

	orm.RegisterModel(new(eth_query.Address), new(eth_query.Block),
		new(eth_query.TokenAddress), new(eth_query.TokenTx), new(eth_query.Tx),
		new(eth_query.TxExtraInfo), new(exchange.EthToken), new(exchange.MainChain))
	orm.RegisterDriver("mysql", orm.DRMySQL)
	// eth query -----------------------
	orm.RegisterDataBase("default",
		"mysql", conf.MysqlUser + ":" + conf.MysqlPasswd+
			"@tcp("+ conf.MysqlHost+ ":"+ conf.MysqlPort+ ")/eth_query?charset=utf8", maxIdle, maxConn)

	// exchange --------------------------
	orm.RegisterDataBase("exchange",
		"mysql", conf.MysqlUser + ":" + conf.MysqlPasswd+
			"@tcp("+ conf.MysqlHost+ ":"+ conf.MysqlPort+ ")/exchange?charset=utf8", maxIdle, maxConn)
}

func CreateTable() {
	orm.RunSyncdb("default", true, true)
	orm.RunSyncdb("exchange", true, true)
}
