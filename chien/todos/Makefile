.PHONY: build test lint generate run

build:
	go build -o bin/server ./cmd/server

test:
	go test -race -v ./...

lint:
	golangci-lint run ./...

generate:
	go generate ./...
	@# Fix wire_gen.go if needed
	@if [ -f internal/registry/wire_gen.go ]; then \
		sed -i '' 's/return nil, nil/return nil/g' internal/registry/wire_gen.go 2>/dev/null || true; \
	fi

run:
	go run ./cmd/server
