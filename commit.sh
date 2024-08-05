#!/bin/sh
gofmt -l -w -s ./main.go
git commit -am "$1"
git push
