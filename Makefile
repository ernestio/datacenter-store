install:
	go install -v

build:
	go build -v ./...

deps:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

dev-deps: deps
	go get github.com/smartystreets/goconvey
	go get github.com/alecthomas/gometalinter
	gometalinter --install

test:
	go test --cover -v $(go list ./... | grep -v /vendor/)

lint:
	gometalinter --config .linter.conf
