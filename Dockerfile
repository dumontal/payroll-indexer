FROM golang:1.7
MAINTAINER Alexandre Dumont <dmt.alexandre@gmail.com>

ADD . /go/src/github.com/dumontal/payroll-indexer

RUN go get github.com/aws/aws-sdk-go    \
  && go get gopkg.in/olivere/elastic.v5 \
  && go install github.com/dumontal/payroll-indexer

ENTRYPOINT /go/bin/payroll-indexer
