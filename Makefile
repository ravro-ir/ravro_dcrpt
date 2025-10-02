.PHONY: all build build-cli build-gui clean install-deps test help

# Variables
APP_NAME=ravro_dcrpt
CLI_BIN=build/$(APP_NAME)
GUI_BIN=build/$(APP_NAME)_gui
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION)"

##@ General

help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

all: clean build ## Clean and build everything

##@ Development

deps: ## Download dependencies
	$(GOMOD) download
	$(GOMOD) tidy

install-deps: ## Install system dependencies (Linux/macOS)
	@echo "📦 Installing system dependencies..."
	@if [ "$(shell uname)" = "Linux" ]; then \
		sudo apt-get update && \
		sudo apt-get install -y libssl-dev wkhtmltopdf; \
	elif [ "$(shell uname)" = "Darwin" ]; then \
		brew install openssl wkhtmltopdf; \
	fi

test: ## Run tests
	$(GOTEST) -v ./...

clean: ## Remove build artifacts
	$(GOCLEAN)
	rm -rf build/
	rm -rf template/
	rm -rf fyne-cross/

##@ Building

build: build-cli build-gui ## Build both CLI and GUI

build-cli: ## Build CLI application
	@echo "🔨 Building CLI..."
	@mkdir -p build
	$(GOBUILD) $(LDFLAGS) -o $(CLI_BIN) ./cmd/cli
	@echo "✅ CLI built: $(CLI_BIN)"

build-gui: ## Build GUI application
	@echo "🔨 Building GUI..."
	@mkdir -p build
ifeq ($(OS),Windows_NT)
	$(GOBUILD) $(LDFLAGS) -ldflags "-s -w -H windowsgui" -o $(GUI_BIN).exe ./cmd/gui
else
	$(GOBUILD) $(LDFLAGS) -o $(GUI_BIN) ./cmd/gui
endif
	@echo "✅ GUI built: $(GUI_BIN)"

##@ Running

run-cli: build-cli ## Build and run CLI
	$(CLI_BIN) --help

run-gui: build-gui ## Build and run GUI
	$(GUI_BIN)

##@ Cross-compilation (Advanced)

build-windows-cli: ## Build CLI for Windows (requires MinGW + OpenSSL)
	@echo "🔨 Building CLI for Windows..."
	@echo "⚠️  Requires: MinGW and OpenSSL for MinGW"
	@echo "⚠️  See CROSS_COMPILE.md for setup instructions"
	@mkdir -p build
	CGO_ENABLED=1 \
	GOOS=windows \
	GOARCH=amd64 \
	CC=x86_64-w64-mingw32-gcc \
	CXX=x86_64-w64-mingw32-g++ \
	$(GOBUILD) $(LDFLAGS) -o build/$(APP_NAME).exe ./cmd/cli
	@echo "✅ Windows CLI built: build/$(APP_NAME).exe"

build-windows-gui: ## Build GUI for Windows using fyne-cross
	@echo "🔨 Building GUI for Windows with fyne-cross..."
	@echo "⚠️  Requires: fyne-cross (go install github.com/fyne-io/fyne-cross@latest)"
	fyne-cross windows -arch=amd64 -app-id=ir.ravro.dcrpt ./cmd/gui
	@echo "✅ Windows GUI built in: fyne-cross/dist/windows-amd64/"

build-macos-cli: ## Build CLI for macOS
	@echo "🔨 Building CLI for macOS..."
	@mkdir -p build
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o build/$(APP_NAME)-darwin-arm64 ./cmd/cli
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o build/$(APP_NAME)-darwin-amd64 ./cmd/cli
	@echo "✅ macOS CLI built"

build-all-platforms: build build-macos-cli ## Build for all platforms
	@echo "✅ All platforms built!"

##@ Docker

docker-build: ## Build in Docker
	docker build -t $(APP_NAME):$(VERSION) .

##@ Installation

install: build ## Install binaries to /usr/local/bin
	@echo "📦 Installing..."
	sudo cp $(CLI_BIN) /usr/local/bin/
	sudo cp $(GUI_BIN) /usr/local/bin/
	@echo "✅ Installed to /usr/local/bin/"

uninstall: ## Uninstall binaries
	@echo "🗑️  Uninstalling..."
	sudo rm -f /usr/local/bin/$(APP_NAME)
	sudo rm -f /usr/local/bin/$(APP_NAME)_gui
	@echo "✅ Uninstalled"

##@ Release

release: clean ## Create release builds
	@echo "📦 Creating release builds..."
	@mkdir -p build/release
	# Linux
	$(GOBUILD) $(LDFLAGS) -o build/release/$(APP_NAME)-linux-amd64 ./cmd/cli
	$(GOBUILD) $(LDFLAGS) -o build/release/$(APP_NAME)_gui-linux-amd64 ./cmd/gui
	# macOS
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o build/release/$(APP_NAME)-darwin-arm64 ./cmd/cli
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o build/release/$(APP_NAME)-darwin-amd64 ./cmd/cli
	@echo "✅ Release builds created in build/release/"
	@ls -lh build/release/

##@ Utilities

fmt: ## Format code
	$(GOCMD) fmt ./...

lint: ## Run linters
	golangci-lint run ./...

vet: ## Run go vet
	$(GOCMD) vet ./...
