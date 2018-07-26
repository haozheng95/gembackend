package scripts

import (
	"net/url"
	"net/http"
	"io/ioutil"
	"regexp"
	"strings"
)

// tether -> usdt
// bitcoin -> btc
var (
	feixiaohaoBaseUrl = url.URL{Scheme: "https", Host: "www.feixiaohao.com", Path: "/currencies/"}
	cointype          = []string{"tether"}
	reg               = regexp.MustCompile(`coinprice(.*?)<`)
)

func FeixiaohaoStart() {
	feixiaohaoBaseUrl.Path = feixiaohaoBaseUrl.Path + cointype[0]
	log.Debugf("connect url = %s", feixiaohaoBaseUrl.String())
	original := feixiaohaoGetpage(feixiaohaoBaseUrl.String())
	price := feixiaohaoExtractPrice(original)
	log.Infof("the %s price : %s cny", cointype[0], price)
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
	s = strings.TrimSpace(s)
	return
}
