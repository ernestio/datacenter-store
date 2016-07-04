install:
	go install -v

build:
	go build -v ./...

deps: dev-deps
	go get -u github.com/jinzhu/gorm
	go get -u github.com/nats-io/nats
	go get -u github.com/lib/pq
	go get -u github.com/r3labs/natsdb

dev-deps:
	go get -u github.com/golang/lint/golint
	go get -u github.com/smartystreets/goconvey/convey

test:
	go test -v ./...

lint:
	golint ./...
	go vet ./...
