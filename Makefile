# Default binary name
BINARY ?= taxirestapi

GO ?= go
BIN ?= bin
BIN_DIR ?= ./$(BIN)

check-bin:
	@mkdir -p $(BIN_DIR)
.PHONY: check-bin

generate:
	mkdir -p gen
	swagger generate server -t gen -f ./swagger.yml --exclude-main
.PHONY: generate

start:
	@$(BIN_DIR)/$(BINARY)
.PHONY: start

build: check-bin
	@echo "==> building $(BINARY) binary"
	@$(GO) build -o $(BIN_DIR)/$(BINARY)
.PHONY: build