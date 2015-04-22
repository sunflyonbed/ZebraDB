#!/bin/bash

export CGO_CFLAGS=-I$GOPATH/deps/include
export CGO_LDFLAGS=-L$GOPATH/deps/libs
