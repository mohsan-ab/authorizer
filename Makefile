BINARY=engine
VERSION = $(shell git describe --tags || echo "developer")

run:
	go run cmd/main.go

engine:
	go build -o ${BINARY} cmd/*.go

setup:
	export GOPRIVATE=github.com/mohsanabbas
	export GOSUMDB=off
	export GONOSUMDB=github.com/mohsanabbas
	export GONOPROXY=github.com/mohsanabbas
	go clean -modcache
	go get -u ./...
	go mod vendor

build:
	go mod vendor
	go build -o bin/authorizer -mod=vendor

clean:
	go clean
	-find . -name ".out" -exec rm -f {} \;