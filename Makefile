SHELL:=/bin/bash

GO_BIN:=go1.20

.PHONY:$(MAKECMDGOALS)

all: test build-all
	./build/ttt

run:
	$(GO_BIN) run cmd/ttt/main.go

build: test
	$(GO_BIN) build -o build/ttt cmd/ttt/main.go

build-all: test
	$(GO_BIN) build -o build/ ./...

clean:
	rm -r build/*

test:
	$(GO_BIN) test -v ./...

dev:
	air
