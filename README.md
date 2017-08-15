[![Build Status](https://travis-ci.org/EGaaS/go-egaas-mvp.svg?branch=1.0)](https://travis-ci.org/EGaaS/go-egaas-mvp)

# Installation

## Requirements

* Go >=1.6
* git

## Build

Clone release beanch:
```
git clone -b 1.0 https://github.com/EGaaS/go-egaas-mvp.git $GOPATH/src/github.com/EGaaS/go-egaas-mvp
```

Build EgaaS:
```
go get -u github.com/jteeuwen/go-bindata/...
rm -rf $GOPATH/src/github.com/EGaaS/go-egaas-mvp/packages/static/static.go
$GOPATH/bin/go-bindata -o="$GOPATH/src/github.com/EGaaS/go-egaas-mvp/packages/static/static.go" -pkg="static" $GOPATH/src/github.com/EGaaS/go-egaas-mvp/static/...
go install github.com/EGaaS/go-egaas-mvp
```

## Running

Create EGaaS directory and copy binary:
```
mkdir ~/egaas
cp $GOPATH/bin/go-egaas-mvp ~/egaas
```

Run EGaaS:
```
~/egaas/go-egaas-mvp
```
