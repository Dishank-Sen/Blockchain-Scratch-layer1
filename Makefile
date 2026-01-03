# =========================
# Makefile — bloc (CLI)
# =========================

APP_NAME=bloc
CMD_PATH=./cmd/bloc
BIN_DIR=bin

GOOS ?= $(shell go env GOOS)

ifeq ($(GOOS),windows)
	EXT=.exe
	INSTALL_DIR=C:/Program Files/bloc
else
	EXT=
	INSTALL_DIR=/usr/local/bin
endif

BIN=$(BIN_DIR)/$(APP_NAME)$(EXT)

.PHONY: build install clean

build:
	go build -o "$(BIN)" "$(CMD_PATH)"

install: build
ifeq ($(GOOS),windows)
	@echo "Installing $(APP_NAME) (Administrator required)"
	powershell -Command " \
	  if (-not ([Security.Principal.WindowsPrincipal] \
	    [Security.Principal.WindowsIdentity]::GetCurrent() \
	    ).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) { \
	      Write-Error 'Run make install as Administrator'; exit 1 \
	  }"
	powershell -Command "New-Item -ItemType Directory -Force '$(INSTALL_DIR)'"
	powershell -Command "Copy-Item '$(BIN)' '$(INSTALL_DIR)/$(APP_NAME)$(EXT)' -Force"
	powershell -Command " \
	  $$path = [Environment]::GetEnvironmentVariable('Path','Machine'); \
	  if ($$path -notlike '*$(INSTALL_DIR)*') { \
	    [Environment]::SetEnvironmentVariable('Path', $$path + ';$(INSTALL_DIR)', 'Machine') \
	  }"
	@echo ""
	@echo "SUCCESS: Installed to $(INSTALL_DIR)"
	@echo "Restart all terminals to use 'bloc'"
else
	sudo cp "$(BIN)" "$(INSTALL_DIR)/$(APP_NAME)"
endif

clean:
	rm -rf "$(BIN_DIR)"
