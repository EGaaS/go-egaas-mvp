<h1>Installation</h1>

* Install golang >=1.6 https://golang.org/dl/
* Set GOPATH
* Install git https://git-scm.com/
* go get -u github.com/jteeuwen/go-bindata/...
* git clone -b 1.0 https://github.com/EGaaS/go-egaas-mvp.git
* cd go-egaas-mvp
* rm -rf packages/static/static.go
* $GOPATH/bin/go-bindata -o="packages/static/static.go" -pkg="static" static/...
* go build
* ./go-egaas-mvp
