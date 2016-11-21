#/usr/bin/env bash

set -x

mkdir -p build/linux/amd64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/linux/amd64/check_nsc_web check_nsc_web.go
mkdir -p build/linux/386
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o build/linux/386/check_nsc_web check_nsc_web.go
mkdir -p build/darwin/amd64
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/darwin/amd64/check_nsc_web check_nsc_web.go
mkdir -p build/windows/amd64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/windows/amd64/check_nsc_web.exe check_nsc_web.go
mkdir -p build/windows/386
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o build/windows/386/check_nsc_web.exe check_nsc_web.go
