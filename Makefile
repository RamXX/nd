MODULE := github.com/RamXX/nd
BIN    := nd
PREFIX := $(HOME)/.local/bin

.PHONY: help build test vet install clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

VERSION := $(shell git describe --tags --always 2>/dev/null || echo dev)

build: ## Build nd binary
	go build -ldflags "-X github.com/RamXX/nd/cmd.version=$(VERSION)" -o $(BIN) .

test: ## Run tests
	go test -v ./...

vet: ## Run go vet
	go vet ./...

install: build ## Install nd to ~/.local/bin
	mkdir -p $(PREFIX)
	cp $(BIN) $(PREFIX)/$(BIN)

clean: ## Remove build artifacts
	rm -f $(BIN)
