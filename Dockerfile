FROM golang:1.6.2-alpine

RUN apk add --update git && apk add --update make && rm -rf /var/cache/apk/*

ADD . /go/src/github.com/ernestio/datacenter-store
WORKDIR /go/src/github.com/ernestio/datacenter-store

RUN make deps && go install

ENTRYPOINT ./entrypoint.sh
