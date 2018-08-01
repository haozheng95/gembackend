package scripts

import (
	"github.com/gembackend/gembackendlog"
	"github.com/shopspring/decimal"
	"sync"
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
			d3, _ = decimal.NewFromString("0")
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

func MulString(d1, d2 string) (r string) {
	m1, _ := decimal.NewFromString(d1)
	m2, _ := decimal.NewFromString(d2)
	m3 := m1.Mul(m2)
	r = m3.String()
	return
}
