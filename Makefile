<<<<<<< HEAD
# =========================
# Makefile — bloc (CLI)
# =========================

APP_NAME := bloc
CMD_PATH := ./cmd/bloc
BIN_DIR  := bin

GOOS ?= $(shell go env GOOS)

# -------------------------
# OS-specific settings
# -------------------------
ifeq ($(GOOS),windows)
	EXT := .exe
	INSTALL_SUPPORTED := false
else ifeq ($(GOOS),darwin)
	EXT :=
	INSTALL_SUPPORTED := false
else
	# linux
	EXT :=
	INSTALL_SUPPORTED := true
	INSTALL_DIR := /usr/local/bin
endif

BIN := $(BIN_DIR)/$(APP_NAME)$(EXT)

.PHONY: build install clean
=======
# Makefile for building and running Go server and client applications
>>>>>>> main

# -------------------------
# Build
# -------------------------
build:
<<<<<<< HEAD
	@echo "Building for $(GOOS)..."
	go build -o "$(BIN)" "$(CMD_PATH)"

# -------------------------
# Install (Linux only)
# -------------------------
install:
ifeq ($(INSTALL_SUPPORTED),true)
	@echo "Installing $(APP_NAME) to $(INSTALL_DIR)"
	sudo cp "$(BIN)" "$(INSTALL_DIR)/$(APP_NAME)"
else
	@echo "Install is not supported on $(GOOS). Build output is: $(BIN)"
endif

# -------------------------
# Clean
# -------------------------
=======
	go build -o bin/main ./main.go
	sudo cp bin/main /usr/local/bin/bloc

# Build server binary
build-server:
	go build -o bin/server ./server

# Build client binary
build-client:
	go build -o bin/client ./client

build-bloc:
	go build -o bin/bloc ./cmd

# Run server (auto-build)
run-server: build-server
	./bin/server

# Run client (auto-build)
run-client: build-client
	./bin/client

run-bloc: build-bloc
	./bin/bloc

run:
	make build
	./bin/main

# Clean binaries
>>>>>>> main
clean:
	rm -rf bin
