.DEFAULT_GOAL := all


go-bindata:
	go get -u github.com/jteeuwen/go-bindata/...

static:
	rm -rf packages/static/static.go
	$GOPATH/bin/go-bindata -o="packages/static/static.go" -pkg="static" static/..

build:
	go build

all:
	make go-bindata
	make static
	make build

