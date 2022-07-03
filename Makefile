.PHONY: generate docker-build docker-netrc clean download

CURRENT_DIR = $(notdir $(shell pwd))
IMAGE_BASE=github.com/anicoll/unicom/${CURRENT_DIR}
GIT_REV=$(shell git rev-parse --short HEAD)

build:
	CGO_ENABLED=0 go build -o main -ldflags="-X 'main.version=${GIT_REV}'" ./cmd/

lint:
	golangci-lint run ./...

format:
	gofmt -l -w .

test:
	go test ./...

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main -ldflags="-X 'main.version=${GIT_REV}'" ./cmd/

download:
	go mod download

docker-netrc:
	@touch ~/.netrc
	@echo "machine gitlab.com" > ~/.netrc
	@echo login gitlab-ci-token >> ~/.netrc
	@echo password ${CI_JOB_TOKEN} >> ~/.netrc

clean:
	rm main

docker-build:
	docker build -f Dockerfile . --build-arg CI_JOB_TOKEN=${CI_JOB_TOKEN} -t ${IMAGE_BASE}:${GIT_REV} -t ${IMAGE_BASE}:latest

docker-build-local: build-linux
	docker build -f Dockerfile.local . -t ${IMAGE_BASE}:${GIT_REV} -t ${IMAGE_BASE}:latest

docker-push:
	docker push ${IMAGE_BASE}:latest

generate:
	rm -rf gen/
	buf generate