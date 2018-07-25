package scripts

import (
	"github.com/gembackend/gembackendlog"
	"sync"
)

type Updater interface {
	Forever()
}

const (
	_TRANSFER          = "0xa9059cbb"
	_TRANSFER_FROM     = "0x23b872dd"
	_TRANSACTION_TOPIC = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	_tokenName         = "0x06fdde03"
	_tokenSymbol       = "0x95d89b41"
	_tokenDecimals     = "0x313ce567"
	_tokenBalance      = "0x70a08231000000000000000000000000"
	_tag               = "latest"
)

var
(
	log = gembackendlog.Log
	wg sync.WaitGroup
)
