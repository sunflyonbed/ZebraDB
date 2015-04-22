#!/bin/bash

export GOPATH=`pwd`/ZebraDB
export CGO_CFLAGS=-I$GOPATH/deps/include
export CGO_LDFLAGS=-L$GOPATH/deps/libs
