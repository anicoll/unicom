.PHONY: generate docker-build docker-netrc clean download lint lint-ci

CURRENT_DIR = $(notdir $(shell pwd))
IMAGE_BASE=ghcr.io/anicoll/${CURRENT_DIR}
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

build-linux-arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o main -ldflags="-X 'main.version=${GIT_REV}'" -ldflags="-X 'main.author=${CURRENT_DIR}'" ./cmd/

download:
	go mod download

clean:
	rm main

docker-build: docker-build-amd64

docker-build-amd64: build-linux
	docker build --platform=linux/amd64 -f Dockerfile . -t ${IMAGE_BASE}:${GIT_REV} -t ${IMAGE_BASE}:latest

docker-build-arm64: build-linux-arm
	docker build --platform=linux/arm64 -f Dockerfile . -t ${IMAGE_BASE}:${GIT_REV} -t ${IMAGE_BASE}:latest

docker-push-arm: docker-build-arm64 docker-push

docker-push-amd: docker-build-amd64 docker-push

docker-push:
	docker push -a ${IMAGE_BASE}

generate:
	rm -rf gen/
	buf generate

generate-docker:
	rm -rf gen/
	docker compose up
