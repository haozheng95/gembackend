package rpc

import (
	"strings"
	"encoding/hex"
)

// 解析十六进制数据
func DecodeHexResponseToString(hexString string) (string, error){
	var str string
	if strings.HasPrefix(hexString, "0x") || strings.HasPrefix(hexString, "0X") {
		str = hexString[2:]
	} else {
		str = hexString
	}
	b, err := hex.DecodeString(str)
	if err != nil{
		return str, err
	}
	b = RemoveZero(b)
	str = strings.Replace(string(b), " ", "", -1)
	return str, err
}

func RemoveZero(slice []byte) []byte {
	if len(slice) == 0{
		return slice
	}
	c := make([]byte, len(slice))
	z := 0
	for _, v := range slice{
		if v != 0{
			c[z] = v
			z ++
		}
	}
	result := make([]byte, z, z)
	result = c[:z]
	return result
}