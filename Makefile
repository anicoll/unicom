.PHONY: generate docker-build docker-netrc clean download lint lint-ci

CURRENT_DIR = $(notdir $(shell pwd))
IMAGE_BASE=github.com/anicoll/${CURRENT_DIR}
GIT_REV=$(shell git rev-parse --short HEAD)

build:
	CGO_ENABLED=0 go build -o main -ldflags="-X 'main.version=${GIT_REV}'" -ldflags="-X 'main.author=${CURRENT_DIR}'" ./cmd/

lint:
	golangci-lint run ./...

lint-ci:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.45.2
	golangci-lint run ./cmd/...

format:
	gofmt -l -w .

test:
	gotestsum --junitfile junit.xml ./...

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main -ldflags="-X 'main.version=${GIT_REV}'" -ldflags="-X 'main.author=${CURRENT_DIR}'" ./cmd/

download:
	go mod download

clean:
	rm main

docker-build:
	docker build -f Dockerfile . -t ${IMAGE_BASE}:${GIT_REV} -t ${IMAGE_BASE}:latest

docker-build-local: build-linux
	docker build -f Dockerfile.local . -t ${IMAGE_BASE}:${GIT_REV} -t ${IMAGE_BASE}:latest

docker-push:
	docker push ${IMAGE_BASE}:latest

generate:
	rm -rf gen/
	buf generate