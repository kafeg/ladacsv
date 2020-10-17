#!/bin/bash

ls -alh
export GOPATH=~/.gopath
go get -v
go build
