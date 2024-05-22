#!/bin/sh
rm -rf ./output
mkdir -p ./output
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build --ldflags="-s -w" -o output/gpbuttond-armv6
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build --ldflags="-s -w" -o output/gpbuttond-armv7
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build --ldflags="-s -w" -o output/gpbuttond-armv8-aarch64
CGO_ENABLED=0 GOOS=linux GOARCH=riscv64 go build --ldflags="-s -w" -o output/gpbuttond-riscv64
