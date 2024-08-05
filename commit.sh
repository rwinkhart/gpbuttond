#!/bin/sh
gofmt -l -w -s .
git commit -am "$1"
git push
