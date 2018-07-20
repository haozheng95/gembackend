package models

import "time"

type Erc20Token struct {
	Id            int64
	TokenName     string
	TokenFullName string
	ContractAddr  string    `orm:"unique"`
	TokenDecimal  int64
	Created       time.Time `orm:"auto_now_add;type(datetime)"`
	ContractAbi   string    `orm:"type(text)"`
}

//func (u *Erc20Token) TableEngine() string {
//	return "MYISAM"
//}
