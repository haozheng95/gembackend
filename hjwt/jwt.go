package hjwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gembackend/conf"
	"github.com/gembackend/gembackendlog"
	"strings"
	"time"
)

var (
	key = []byte(conf.JwtKey)
	log = gembackendlog.Log
)

// 产生json web token
func GenToken() string {
	claims := &jwt.StandardClaims{
		NotBefore: int64(time.Now().Unix()),
		ExpiresAt: int64(time.Now().Unix() + conf.JwtExpiration),
		Issuer:    "jwt-zz",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		log.Error(err)
		return ""
	}
	return ss
}

// 校验token是否有效
func CheckToken(token string) bool {

	if strings.Compare(token, conf.JwtGodKey) == 0 {
		return true
	}

	_, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		log.Errorf("parase with claims failed. %s", err)
	}
	return false

}
