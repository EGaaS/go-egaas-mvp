[![Go Report Card](https://goreportcard.com/badge/github.com/EGaaS/go-egaas-mvp)](https://goreportcard.com/report/github.com/EGaaS/go-egaas-mvp) [![Build Status](https://travis-ci.org/EGaaS/go-egaas-mvp.svg?branch=master)](https://travis-ci.org/EGaaS/go-egaas-mvp) [![API Reference](
https://camo.githubusercontent.com/915b7be44ada53c290eb157634330494ebe3e30a/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f676f6c616e672f6764646f3f7374617475732e737667
)](https://godoc.org/github.com/EGaaS/go-egaas-mvp)
[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/EGaaS/go-egaas-mvp?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)


### Installation v1.x - only egs-wallet

[![Join the chat at https://gitter.im/go-egaas-mvp/Lobby](https://badges.gitter.im/go-egaas-mvp/Lobby.svg)](https://gitter.im/go-egaas-mvp/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Install golang >=1.6 https://golang.org/dl/<br>
Set GOPATH<br>
Install git https://git-scm.com/
```
go get -u github.com/jteeuwen/go-bindata/...
git clone -b 1.0 https://github.com/EGaaS/go-egaas-mvp.git
cd go-egaas-mvp
rm -rf packages/static/static.go
$GOPATH/bin/go-bindata -o="packages/static/static.go" -pkg="static" static/..
go build
./go-egaas-mvp
```

### Installation v0.x - full egaas (private blockchain)


Install golang >=1.6 https://golang.org/dl/<br>
Set GOPATH<br>
Install git https://git-scm.com/
```
go get -u github.com/jteeuwen/go-bindata/...
go get -u github.com/EGaaS/go-egaas-mvp
cd $GOPATH/src/github.com/EGaaS/go-egaas-mvp
$GOPATH/bin/go-bindata -o="packages/static/static.go" -pkg="static" static/..
go build
./go-egaas-mvp
```


### Questions?
email: hello@egaas.org
