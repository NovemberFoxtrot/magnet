Magnet
======

Magnet is a tiny self-hosted bookmarks management tool written in Go(lang). 

Works on Go 1.1 and Go 1.2.

![alt text](https://github.com/mvader/magnet/raw/master/magnet.png "Magnet screenshot")

Requisites
-------
* [Rethinkdb](http://rethinkdb.com)
* [Golang 1.1](http://golang.org/doc/install)

Setup
-------

```bash
git clone https://github.com/mvader/magnet magnet
cd magnet
go get .
go build
mv config.sample.json config.json
# Edit config.json
nano config.json
./magnet
```

Go dependencies 
-------
* [github.com/christopherhesse/rethinkgo](https://github.com/christopherhesse/rethinkgo)
* [github.com/gorilla/sessions](https://github.com/gorilla/sessions)
* [github.com/codegangsta/martini](https://github.com/codegangsta/martini)
* [github.com/hoisie/mustache](https://github.com/hoisie/mustache)
* [github.com/justinas/nosurf](https://github.com/justinas/nosurf)
