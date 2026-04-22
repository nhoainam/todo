.PHONY: build test lint generate run

include .env
export

build:
	go build -o bin/server ./cmd/server

test:
	go test -race -v ./...

lint:
	golangci-lint run ./...

generate:
	go tool gqlgen generate .
	go generate ./...
run:
	go run ./cmd/server
