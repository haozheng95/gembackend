package eth

import (
	"strings"
	"github.com/shopspring/decimal"
	"math"
	"github.com/regcostajr/go-web3/dto"
	"strconv"
	"fmt"
)

func HexDec(h string) (n string) {
	//log.Debugf("------- %s", h)
	if len(h) > 2 && strings.HasPrefix(strings.ToLower(h[:2]), "0x") {
		h = h[2:]
	} else if strings.Compare(h, "0x") == 0 {
		h = "0"
	}

	s := strings.Split(strings.ToUpper(h), "")
	l := len(s)
	i := 0
	d := decimal.NewFromFloat(0)
	hex := map[string]string{"A": "10", "B": "11", "C": "12", "D": "13", "E": "14", "F": "15"}
	for i = 0; i < l; i++ {
		c := s[i]
		if v, ok := hex[c]; ok {
			c = v
		}
		f, err := strconv.ParseFloat(c, 10)
		if err != nil {
			fmt.Println(h)
			log.Error(err)
			return decimal.NewFromFloat(-1).String()
		}
		d = d.Add(decimal.NewFromFloat(f * math.Pow(16, float64(l-i-1))))
	}
	return d.String()
}
func MaxInt(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func MaxIntByString(a uint64, b string) uint64 {
	c, _ := strconv.ParseUint(b, 10, 64)
	if a > c {
		return a
	}
	return c + 1
}

func FormatAddress(addr string) string {
	k := false
	s1 := "0x"
	for _, v := range addr[2:] {
		if v != 48 {
			k = true
		}
		if k {
			s1 += string(v)
		}
	}

	l := len(s1)
	if l < 42 {
		t := ""
		for i := 0; i < 42-l; i++ {
			t += "0"
		}
		s1 = s1[:2] + t + s1[2:]
	}

	if len(s1) > 42 {
		log.Errorf("format addr error %s", addr)
		panic("format addr error %s")
	}
	return s1
}

func AnalysisTokenLog(logs dto.TransactionLogs) (from, to, amount, logindex string) {
	if len(logs.Topics) > 2 && strings.Compare(logs.Topics[0], _TRANSACTION_TOPIC) == 0 {
		to = FormatAddress(logs.Topics[2])
		from = FormatAddress(logs.Topics[1])

		amount = logs.Data
		logindex = logs.LogIndex.String()
	}
	return
}

func formatAmount(s string, d int) string {
	var str string
	if len(s) > 2 && strings.HasPrefix(strings.ToLower(s[:2]), "0x") {
		str = s[2:]
	} else {
		str = s
	}
	amount := HexDec(str)
	return format10Decimals(amount, d)
}



func format10Decimals(amount string, i int) string {
	tempFloat, _ := decimal.NewFromString(amount)
	t := decimal.NewFromFloat(1.0 * math.Pow(10, float64(i)))
	tempFloat = tempFloat.Div(t)
	return tempFloat.String()
}

func makeFee(gasprice, gasused string) string {
	a, _ := decimal.NewFromString(gasprice)
	b, _ := decimal.NewFromString(gasused)
	a = a.Div(decimal.NewFromFloat(1.0 * math.Pow(10, 18)))
	c := a.Mul(b)
	return c.String()
}

func formatAmountString(s string, d string) string {
	var str string
	if len(s) > 2 && strings.HasPrefix(strings.ToLower(s[:2]), "0x") {
		str = s[2:]
	} else {
		str = s
	}
	amount := HexDec(str)
	return format10DecimalsString(amount, d)
}

func format10DecimalsString(amount string, i string) string {
	tempFloat, _ := decimal.NewFromString(amount)
	d,_ := decimal.NewFromString(i)
	p ,_ := d.Float64()
	t := decimal.NewFromFloat(1.0 * math.Pow(10, p))
	tempFloat = tempFloat.Div(t)
	return tempFloat.String()
}
