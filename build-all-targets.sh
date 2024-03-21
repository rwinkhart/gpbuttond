#!/bin/sh
rm -rf ./output
mkdir -p ./output
go mod tidy
GOOS=linux GOARCH=arm GOARM=6 go build -o output/gpbuttond-armv6
GOOS=linux GOARCH=arm GOARM=7 go build -o output/gpbuttond-armv7
GOOS=linux GOARCH=arm64 go build -o output/gpbuttond-armv8-aarch64
GOOS=linux GOARCH=riscv64 go build -o output/gpbuttond-riscv64
