Magnet
======

Magnet is a tiny self-hosted bookmarks management tool written in Go(lang).

![alt text](https://github.com/mvader/magnet/raw/master/magnet.png "Magnet screenshot")

Requisites
-------
* [Rethinkdb](http://rethinkdb.com)
* [Golang 1.1](http://golang.org/doc/install)

Setup
-------

```bash
git clone git@github.com:mvader/magnet.git $GOPATH/github.com/mvader/magnet
cd $GOPATH/github.com/mvader/magnet
sh build.sh
# You must edit your config.sample.json and rename it to config.json
./magnet
```

Go dependencies 
-------
They're installed automatically when you run ```sh build.sh```, you don't need to install them yourself.
* [github.com/christopherhesse/rethinkgo](https://github.com/christopherhesse/rethinkgo)
* [github.com/gorilla/sessions](https://github.com/gorilla/sessions)
* [github.com/codegangsta/martini](https://github.com/codegangsta/martini)
* [github.com/hoisie/mustache](https://github.com/hoisie/mustache)
* [github.com/justinas/nosurf](https://github.com/justinas/nosurf)
