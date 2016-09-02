install:
	go install -v

build:
	go build -v ./...

deps: dev-deps
	go get github.com/jinzhu/gorm
	go get github.com/nats-io/nats
	go get github.com/lib/pq
	go get github.com/r3labs/natsdb
	go get github.com/ernestio/ernest-config-client

dev-deps:
	go get github.com/golang/lint/golint
	go get github.com/smartystreets/goconvey/convey

test:
	go test -v ./...

lint:
	golint ./...
	go vet ./...
