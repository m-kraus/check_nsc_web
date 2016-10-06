#/usr/bin/env bash

set -x

mkdir -p build/linux/amd64
GOOS=linux GOARCH=amd64 go build -o build/linux/amd64/nscrestc nscrestc.go
mkdir -p build/linux/386
GOOS=linux GOARCH=386 go build -o build/linux/386/nscrestc nscrestc.go
mkdir -p build/darwin/amd64
GOOS=darwin GOARCH=amd64 go build -o build/darwin/amd64/nscrestc nscrestc.go
mkdir -p build/windows/amd64
GOOS=windows GOARCH=amd64 go build -o build/windows/amd64/nscrestc.exe nscrestc.go
mkdir -p build/windows/386
GOOS=windows GOARCH=386 go build -o build/windows/386/nscrestc.exe nscrestc.go