.PHONY: generate docker-build docker-netrc clean download

CURRENT_DIR = $(notdir $(shell pwd))
GIT_REV=$(shell git rev-parse --short HEAD)

build:
	CGO_ENABLED=0 go build -o main -ldflags="-X 'main.version=${GIT_REV}'" -ldflags="-X 'main.author=${CURRENT_DIR}'" ./cmd/

lint:
	./lint.sh

format:
	gofmt -l -w .

test:
	go test ./...

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