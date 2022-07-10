#!/bin/bash
os=$(uname)

if [ os = "Darwin" ]
then
	golangci-lint --version
	golangci-lint run ./cmd/...
else
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.45.2
	golangci-lint --version
	golangci-lint run ./cmd/...
fi