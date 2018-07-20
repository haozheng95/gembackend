# start eth updater for web3
FROM golang
MAINTAINER wind-yin yinhaozheng77625961@gmail.com

ADD src.tar.gz /go/
ADD config.ini /go/src/github.com/gembackend/conf/

WORKDIR /go/src/github.com/gembackend/startscript/
RUN go install
CMD ["/go/bin/startscript"]