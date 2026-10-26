BINARY := httpserver
MAIN_PACKAGE := ./cmd/httpserver

.DEFAULT_GOAL := build

fmt:
	go fmt ./...
.PHONY: fmt

vet: fmt
	go vet ./...
.PHONY: vet

build: vet
	go build -o bin/$(BINARY) $(MAIN_PACKAGE)
.PHONY: build

run:
	go run $(MAIN_PACKAGE)
.PHONY: run

test:
	go test ./...
.PHONY: test

clean:
	rm -f $(BINARY)
.PHONY: clean

fix:
	go fix ./...
.PHONY: fix
