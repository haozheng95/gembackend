# start eth updater for web3
FROM golang
MAINTAINER wind-yin

RUN go get github.com/Shopify/sarama
RUN mkdir -p /go/src/golang.org/x
RUN git clone https://github.com/golang/crypto.git /go/src/golang.org/x/crypto
RUN go get github.com/astaxie/beego
RUN go get github.com/astaxie/beego/config/xml
RUN go get github.com/dgrijalva/jwt-go
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/op/go-logging
RUN go get github.com/regcostajr/go-web3
RUN go get github.com/shopspring/decimal
RUN go get github.com/beego/bee
RUN git clone https://github.com/explosive77/gembackend.git /go/src/github.com/gembackend
