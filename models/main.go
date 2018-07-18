package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gembackend/conf"
)

func init() {
	orm.RegisterModel(new(Address), new(Block),
		new(TokenAddress), new(TokenTx), new(Tx),
		new(TxExtraInfo), new(Erc20Token))
	orm.RegisterDriver("mysql", orm.DRMySQL)
	maxIdle := 30
	maxConn := 30
	orm.RegisterDataBase("default",
		"mysql", conf.MysqlUser + ":" + conf.MysqlPasswd+
			"@tcp("+ conf.MysqlHost+ ":"+ conf.EthRpcPort+ ")/eth_query?charset=utf8", maxIdle, maxConn)
	orm.RunSyncdb("default", true, true)
}
