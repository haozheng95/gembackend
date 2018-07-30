package scripts

import (
	"fmt"
	"github.com/gembackend/models/exchange"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// tether -> usdt
// bitcoin -> btc
var (
	cointype []string
	reg      = regexp.MustCompile(`coinprice(.*?)<`)
	baseHost = "www.feixiaohao.com"
	interval = time.Duration(10)
)

func FeixiaohaoStart() {
	fullnames := exchange.GetAllCoinFullNaml()
	cointype = make([]string, len(fullnames))
	for i := range fullnames {
		cointype[i] = fullnames[i].FullName
	}
	log.Infof("coin num:%d", len(cointype))
	wg.Add(1)

	for _, coin := range cointype {

		go func(coin string) {
			ticker := time.NewTicker(time.Second * interval)
			defer ticker.Stop()

			for range ticker.C {
				feixiaohaoBaseUrl := url.URL{Scheme: "https", Host: baseHost, Path: "/currencies/%s"}
				feixiaohaoBaseUrl.Path = fmt.Sprintf(feixiaohaoBaseUrl.Path, coin)
				//log.Infof("connect url = %s", feixiaohaoBaseUrl.String())
				original := feixiaohaoGetpage(feixiaohaoBaseUrl.String())
				price := feixiaohaoExtractPrice(original)
				//log.Infof("the %s price : %s cny", coin, price)
				updatemainchain(coin, price)
			}
		}(coin)
	}

	// start token ================

	for _, v := range exchange.GetAllFullTokenName() {
		go func(coin string) {
			ticker1 := time.NewTicker(time.Second * interval)
			defer ticker1.Stop()
			for range ticker1.C {
				//log.Debug(coin)
				feixiaohaoBaseUrl := url.URL{Scheme: "https", Host: baseHost, Path: "/currencies/%s"}
				feixiaohaoBaseUrl.Path = fmt.Sprintf(feixiaohaoBaseUrl.Path, coin)
				original := feixiaohaoGetpage(feixiaohaoBaseUrl.String())
				price := feixiaohaoExtractPrice(original)
				feixiaohaoupdateEthtoken(coin, price)
			}
		}(v.TokenFullName)

	}

	wg.Wait()

}

func updatemainchain(coin, cny string) {
	st := new(exchange.MainChain)
	st.FullName = baseCoinFullName
	st.SelectCny()
	usdtCny := st.Cny
	usdtNum := DivString(cny, usdtCny)

	st.FullName = coin
	st.Cny = cny
	st.Usdt = usdtNum.String()
	st.UpdateCnyAndUsdt()
	log.Infof("cny :%s, coin :%s, usdt_num :%s", cny, coin, usdtNum.String())
}

func feixiaohaoGetpage(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Error(err)
		return ""
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	s1 := string(body)
	s2 := reg.FindAllString(s1, -1)
	if len(s2) != 1 {
		log.Warning("find error the len : %d", len(s2))
		log.Warning("the original : %s", s1)
		return ""
	}
	return s2[0]
}

func feixiaohaoExtractPrice(original string) (s string) {
	s = strings.Replace(original, "coinprice>￥", "", -1)
	s = strings.Replace(s, "<", "", -1)
	s = strings.Replace(s, ",", "", -1)
	s = strings.TrimSpace(s)
	return
}

//func FeixiaohaoTokenStart() {
//
//	wg.Add(1)
//
//	ticker := time.NewTicker(time.Second * 1)
//	defer ticker.Stop()
//
//	for _, v := range exchange.GetAllFullTokenName() {
//		go func(coin string) {
//			for range ticker.C {
//				log.Debug(coin)
//				feixiaohaoBaseUrl := url.URL{Scheme: "https", Host: baseHost, Path: "/currencies/%s"}
//				feixiaohaoBaseUrl.Path = fmt.Sprintf(feixiaohaoBaseUrl.Path, coin)
//				original := feixiaohaoGetpage(feixiaohaoBaseUrl.String())
//				price := feixiaohaoExtractPrice(original)
//				feixiaohaoupdateEthtoken(coin, price)
//			}
//		}(v.TokenFullName)
//
//	}
//	wg.Wait()
//}

func feixiaohaoupdateEthtoken(coin, cny string) {
	st := new(exchange.MainChain)
	st.FullName = baseCoinFullName
	st.SelectCny()
	usdtCny := st.Cny
	usdtNum := DivString(cny, usdtCny)

	stToken := new(exchange.EthToken)
	stToken.TokenFullName = coin
	stToken.Cny = cny
	stToken.Usdt = usdtNum.String()
	stToken.UpdateCnyAndUsdtByFullName()
	log.Infof("token === :cny :%s, coin :%s, usdt_num :%s", cny, coin, usdtNum.String())
}
