[![Go Report Card](https://goreportcard.com/badge/github.com/EGaaS/go-egaas-mvp)](https://goreportcard.com/report/github.com/EGaaS/go-egaas-mvp) 
[![Build Status](https://travis-ci.org/EGaaS/go-egaas-mvp.svg?branch=master)](https://travis-ci.org/EGaaS/go-egaas-mvp) 
[![Documentation](https://img.shields.io/badge/docs-latest-brightgreen.svg?style=flat)](http://egaas-en.readthedocs.io/en/latest/)
[![API Reference](
https://camo.githubusercontent.com/915b7be44ada53c290eb157634330494ebe3e30a/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f676f6c616e672f6764646f3f7374617475732e737667
)](https://godoc.org/github.com/EGaaS/go-egaas-mvp)
[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/go-egaas-mvp?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)


# Installation

## Requirements

* Go >=1.6
* git

## Build

Clone:
```
git clone https://github.com/EGaaS/go-egaas-mvp.git $GOPATH/src/github.com/EGaaS/go-egaas-mvp
```

Build EgaaS:
```
go get -u github.com/jteeuwen/go-bindata/...
$GOPATH/bin/go-bindata -o="$GOPATH/src/github.com/EGaaS/go-egaas-mvp/packages/static/static.go" -pkg="static" -prefix="$GOPATH/src/github.com/EGaaS/go-egaas-mvp/" $GOPATH/src/github.com/EGaaS/go-egaas-mvp/static/...
go install github.com/EGaaS/go-egaas-mvp
```

# Running

Create EGaaS directory and copy binary:
```
mkdir ~/egaas
cp $GOPATH/bin/go-egaas-mvp ~/egaas
```

Run EGaaS:
```
~/egaas/go-egaas-mvp
```
Open EGaaS: http://localhost:7079/


----------


### Questions?
email: hello@egaas.org
