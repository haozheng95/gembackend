package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/gembackend/conf"
	"github.com/gembackend/gembackendlog"
	"github.com/gembackend/models/btc_query"
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
		new(eth_query.TxExtraInfo),
		new(exchange.EthToken), new(exchange.MainChain),
		new(btc_query.AddressBtc), new(btc_query.TradeCollection), new(btc_query.TradingParticulars),
		new(btc_query.UnspentVout))
	orm.RegisterDriver("mysql", orm.DRMySQL)
	// eth query -----------------------
	orm.RegisterDataBase("default",
		"mysql", conf.MysqlUser+":"+conf.MysqlPasswd+
			"@tcp("+conf.MysqlHost+":"+conf.MysqlPort+")/eth_query?charset=utf8", maxIdle, maxConn)

	orm.RegisterDataBase("eth_query",
		"mysql", conf.MysqlUser+":"+conf.MysqlPasswd+
			"@tcp("+conf.MysqlHost+":"+conf.MysqlPort+")/eth_query?charset=utf8", maxIdle, maxConn)
	// exchange --------------------------
	orm.RegisterDataBase("exchange",
		"mysql", conf.MysqlUser+":"+conf.MysqlPasswd+
			"@tcp("+conf.MysqlHost+":"+conf.MysqlPort+")/exchange?charset=utf8", maxIdle, maxConn)
	// btc_query --------------------------
	orm.RegisterDataBase("btc_query",
		"mysql", conf.MysqlUser+":"+conf.MysqlPasswd+
			"@tcp("+conf.MysqlHost+":"+conf.MysqlPort+")/btc_query?charset=utf8", maxIdle, maxConn)
}

func CreateTable() {
	//orm.RunSyncdb("default", true, true)
	//orm.RunSyncdb("exchange", true, true)
	orm.RunSyncdb("btc_query", true, true)
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
			address(o)
		case "block":
		case "token_address":
		case "token_tx":
			token_tx(o)
		case "tx":
			tx(o)
		case "tx_extra_info":
			tx_extra_info(o)
		}
	}
}

func address(o orm.Ormer) {
	data := []eth_query.Address{
		{
			WalletId:        "d3ba134f262d6d197a93ade4a6c123ddb9122c5cc0ff666f5447639d36f5f155",
			Addr:            "0xd6cb6744b7f2da784c5afd6b023d957188522198",
			TypeId:          4,
			Nonce:           "4",
			Created:         time.Now(),
			Amount:          "4000",
			UnconfirmAmount: "300",
			Decimal:         18,
		},
	}

	if num, err := o.InsertMulti(len(data), data); err == nil {
		log.Infof("insert row : %d", num)
	} else {
		log.Errorf("insert error : %s", err)
	}
}

func tx_extra_info(o orm.Ormer) {
	data := []eth_query.TxExtraInfo{
		{
			From:        "0xd6cb6744b7f2da784c5afd6b023d957188522198",
			To:          "0x88a690553913a795c3c668275297635b903a29e5",
			TxHash:      "0x569c5b35f203ca6db6e2cec44bceba756fad513384e2bd79c06a8c0181273379",
			Nonce:       "4079",
			Amount:      "3.18096329",
			TokenAmount: "0",
			Comment:     "test1",
			Created:     time.Now(),
		},
		{
			From:        "0xd6cb6744b7f2da784c5afd6b023d957188522198",
			To:          "0x88a690553913a795c3c668275297635b903a29e5",
			TxHash:      "0x696a35492b283624ccf4ae9438ae2d5d5e84a4a00798155b568d1eb52606d829",
			Nonce:       "4079",
			Amount:      "3.18096329",
			TokenAmount: "0",
			Comment:     "test2",
			Created:     time.Now(),
		},
		{
			From:        "0xd6cb6744b7f2da784c5afd6b023d957188522198",
			To:          "0x88a690553913a795c3c668275297635b903a29e5",
			TxHash:      "0x696a35492b283624ccf4ae9438ae2d5d5e84a4a00798155b568d1eb52606d828",
			Nonce:       "4079",
			Amount:      "3.18096329",
			TokenAmount: "0",
			Comment:     "test3",
			Created:     time.Now(),
		},
	}
	if num, err := o.InsertMulti(len(data), data); err == nil {
		log.Infof("insert row : %d", num)
	} else {
		log.Errorf("insert error : %s", err)
	}
}

func tx(o orm.Ormer) {
	data := []eth_query.Tx{
		{
			From:        "0xd6cb6744b7f2da784c5afd6b023d957188522198",
			To:          "0x88a690553913a795c3c668275297635b903a29e5",
			Amount:      "3.18096329",
			InputData:   "0x",
			Nonce:       "4079",
			GasLimit:    "121000",
			GasPrice:    "134000000000",
			GasUsed:     "21000",
			Fee:         "0.002814",
			TxHash:      "0x569c5b35f203ca6db6e2cec44bceba756fad513384e2bd79c06a8c0181273379",
			BlockHash:   "0x7d5a4369273c723454ac137f48a4f142b097aa2779464e6505f1b1c5e37b5382",
			BlockHeight: "5000000",
			ConfirmTime: "1517319693",
			Created:     time.Now(),
			BlockState:  1,
			TxState:     1,
			IsToken:     0,
		},
		{
			From:        "0x88a690553913a795c3c668275297635b903a29e5",
			To:          "0xd6cb6744b7f2da784c5afd6b023d957188522198",
			Amount:      "3.18096329",
			InputData:   "0xa9059cbb000000000000000000000000f53354a8dc35416d28ab2523589d1b44843e025c00000000000000000000000000000000000000000000009a41e07a74a99ec000",
			Nonce:       "221655",
			GasLimit:    "79670",
			GasPrice:    "100000000000",
			GasUsed:     "39835",
			Fee:         "0.0039835",
			TxHash:      "0x696a35492b283624ccf4ae9438ae2d5d5e84a4a00798155b568d1eb52606d829",
			BlockHash:   "0x7d5a4369273c723454ac137f48a4f142b097aa2779464e6505f1b1c5e37b5382",
			BlockHeight: "5000000",
			ConfirmTime: "1517319693",
			Created:     time.Now(),
			BlockState:  1,
			TxState:     1,
			IsToken:     1,
		},
		{
			From:      "0x88a690553913a795c3c668275297635b903a29e5",
			To:        "0xd6cb6744b7f2da784c5afd6b023d957188522198",
			Amount:    "31.18096329",
			InputData: "0xa9059cbb000000000000000000000000f53354a8dc35416d28ab2523589d1b44843e025c00000000000000000000000000000000000000000000009a41e07a74a99ec000",
			Nonce:     "221655",
			GasLimit:  "79670",
			GasPrice:  "100000000000",
			Fee:       "0.0039835",
			TxHash:    "0x696a35492b283624ccf4ae9438ae2d5d5e84a4a00798155b568d1eb52606d828",
			Created:   time.Now(),
			TxState:   -1,
			IsToken:   1,
		},
	}
	if num, err := o.InsertMulti(len(data), data); err == nil {
		log.Infof("insert row : %d", num)
	} else {
		log.Errorf("insert error : %s", err)
	}
}

func token_tx(o orm.Ormer) {
	data := []eth_query.TokenTx{
		{
			From:         "0xd6cb6744b7f2da784c5afd6b023d957188522198",
			To:           "0x4ac751f0152b6799a5acfc25614072fbb06dca06",
			Amount:       "43",
			InputData:    "0xa9059cbb0000000000000000000000004ac751f0152b6799a5acfc25614072fbb06dca06000000000000000000000000000000000000000000000000000000000000a7f8",
			Nonce:        "4577088",
			GasLimit:     "150000",
			GasPrice:     "90000000000",
			GasUsed:      "37175",
			Fee:          "0.00334575",
			TxHash:       "0x06303859b5a5e00a72e6d020ea4479f14fe6ebd51e6fdd07d255f0c8a75608b8",
			BlockHash:    "0x02d10d64c97fc42e71f4bb2ac896470a66cb4ec9b7236334ab9288fcb81c77ae",
			BlockHeight:  "5000149",
			ConfirmTime:  "1517321714",
			Created:      time.Now(),
			BlockState:   1,
			TxState:      1,
			IsToken:      1,
			LogIndex:     "0",
			ContractAddr: "0xd26114cd6ee289accf82350c8d8487fedb8a0c07",
			Decimal:      "18",
		},
		{
			From:         "0xd6cb6744b7f2da784c5afd6b023d957188522198",
			To:           "0x9174f429ba9cab8b654d6da6e00ad043982b49f0",
			Amount:       "10509.71557555",
			InputData:    "0xa9059cbb0000000000000000000000009174f429ba9cab8b654d6da6e00ad043982b49f0000000000000000000000000000000000000000000000239bb9a491a67c86c00",
			Nonce:        "4576829",
			GasLimit:     "150000",
			GasPrice:     "90000000000",
			GasUsed:      "37286",
			Fee:          "0.00335574",
			TxHash:       "0xf3af7032b99f1c7b72f8ca79be607b8f9565d85cc5277a59982e6442d0467666",
			BlockHash:    "0x5ca74f2fe6cba2615054713fe003e8eb0a8ea4e470990a7e96820f6f714ee0c0",
			BlockHeight:  "5000056",
			ConfirmTime:  "1517320567",
			Created:      time.Now(),
			BlockState:   1,
			TxState:      1,
			IsToken:      1,
			LogIndex:     "8",
			ContractAddr: "0xd26114cd6ee289accf82350c8d8487fedb8a0c07",
			Decimal:      "18",
		},
		{
			From:         "0x24dd7159c188b399dba59ecc65196f9cc1476cce",
			To:           "0xd6cb6744b7f2da784c5afd6b023d957188522198",
			Amount:       "10509.71557555",
			InputData:    "0xa9059cbb0000000000000000000000009174f429ba9cab8b654d6da6e00ad043982b49f0000000000000000000000000000000000000000000000239bb9a491a67c86c00",
			Nonce:        "4576830",
			GasLimit:     "150000",
			GasPrice:     "90000000000",
			GasUsed:      "37286",
			Fee:          "0.00335574",
			TxHash:       "0x1b5df18dc595e70d67b0abca83aa95fe4bb57ffb117c053f3b36501870202b25",
			BlockHash:    "0x5ca74f2fe6cba2615054713fe003e8eb0a8ea4e470990a7e96820f6f714ee0c0",
			BlockHeight:  "5000056",
			ConfirmTime:  "1517320567",
			Created:      time.Now(),
			BlockState:   1,
			TxState:      1,
			IsToken:      1,
			LogIndex:     "8",
			ContractAddr: "0xd26114cd6ee289accf82350c8d8487fedb8a0c07",
			Decimal:      "18",
		},
		{
			From:         "0xd6cb6744b7f2da784c5afd6b023d957188522198",
			To:           "0x24dd7159c188b399dba59ecc65196f9cc1476cce",
			Amount:       "10509.71557555",
			InputData:    "0xa9059cbb0000000000000000000000009174f429ba9cab8b654d6da6e00ad043982b49f0000000000000000000000000000000000000000000000239bb9a491a67c86c00",
			Nonce:        "4576831",
			GasLimit:     "150000",
			GasPrice:     "90000000000",
			Fee:          "0.00335574",
			TxHash:       "0x1b5df18dc595e70d67b0abca83aa95fe4bb57ffb117c053f3b36501870202b25",
			Created:      time.Now(),
			TxState:      -1,
			IsToken:      1,
			ContractAddr: "0xd26114cd6ee289accf82350c8d8487fedb8a0c07",
			Decimal:      "18",
		},
	}
	if num, err := o.InsertMulti(len(data), data); err == nil {
		log.Infof("insert row : %d", num)
	} else {
		log.Errorf("insert error : %s", err)
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
