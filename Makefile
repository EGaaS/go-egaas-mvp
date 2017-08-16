.DEFAULT_GOAL := all


go-bindata:
	go get -u github.com/jteeuwen/go-bindata/...

static-files:
	rm -rf $GOPATH/src/github.com/EGaaS/go-egaas-mvp/packages/static/static.go
	${GOPATH}/bin/go-bindata -o="${GOPATH}/src/github.com/EGaaS/go-egaas-mvp/packages/static/static.go" -pkg="static" -prefix="${GOPATH}/src/github.com/EGaaS/go-egaas-mvp/" ${GOPATH}/src/github.com/EGaaS/go-egaas-mvp/static/...

build:
	go build github.com/EGaaS/go-egaas-mvp

install:
	go install github.com/EGaaS/go-egaas-mvp

all:
	make go-bindata
	make static-files
	make install
