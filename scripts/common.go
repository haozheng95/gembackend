package scripts

import (
	"github.com/gembackend/gembackendlog"
	"sync"
	"github.com/shopspring/decimal"
)

var (
	log              = gembackendlog.Log
	wg               sync.WaitGroup
	baseCoin         = "usdt"
	baseCoinFullName = "tether"
)

func DivString(m1, m2 string) (d3 decimal.Decimal) {
	//log.Debugf("check param m1=%s, m2=%s", m1, m2)
	defer func() {
		if err := recover(); err != nil {
			d3 ,_ = decimal.NewFromString("0")
			log.Errorf("DivString err: %s", err)
			log.Errorf("check param m1=%s, m2=%s", m1, m2)
		}
	}()

	d1, _ := decimal.NewFromString(m1)
	d2, _ := decimal.NewFromString(m2)
	//log.Debugf("check param d1=%s, d2=%s", d1.String(), d2.String())
	d3 = d1.Div(d2)
	return
}
