FROM golang:1.9.0-alpine3.6 as compiler
RUN apk add --update git && apk add --update make && rm -rf /var/cache/apk/*
ADD . /go/src/github.com/${GITHUB_ORG:-ernestio}/datacenter-store
WORKDIR /go/src/github.com/${GITHUB_ORG:-ernestio}/datacenter-store
RUN make deps && CGO_ENABLED=0 go install -a -ldflags '-s' .

FROM scratch
COPY --from=compiler /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=compiler /go/bin/datacenter-store .
ENTRYPOINT ["./datacenter-store"]
