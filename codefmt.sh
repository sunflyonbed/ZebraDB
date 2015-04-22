#!/bin/bash

find ./src -iname "*.go" | xargs gofmt -w -s
