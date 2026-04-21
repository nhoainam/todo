.PHONY: build test lint generate run

# Lấy đường dẫn GOPATH từ môi trường Go hiện tại
GOPATH := $(shell go env GOPATH)

# Thêm thư mục bin của Go vào biến PATH dùng chung cho toàn bộ Makefile
export PATH := $(PATH):$(GOPATH)/bin
include .env
export

build:
	go build -o bin/server ./cmd/server

test:
	go test -race -v ./...

lint:
	golangci-lint run ./...

generate:
	go generate ./...

run:
	go run ./cmd/server

# Apply schema changes to the database
migrate:
	ridgepole -c "mysql2://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)" \
		--apply -f database/users.schema

# Dry-run — preview SQL without executing
migrate-dry-run:
	ridgepole -c "mysql2://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)" \
		--apply --dry-run -f database/users.schema
