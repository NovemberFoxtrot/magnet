#!/bin/sh

go get .
go install github.com/mvader/magnet
cp $GOPATH/bin/magnet ./