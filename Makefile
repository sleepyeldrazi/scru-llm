.PHONY: all build test clean lint run install help

# Variables
BINARY_NAME=scru-llm
MAIN_PACKAGE=./cmd/scru-llm
BUILD_DIR=bin
GO=go
GOFLAGS=-v

# Default target
all: clean build

## help: Show this help message
help:
	@echo "Available targets:"
	@awk '/^##/{c=substr($$0,3);next}c&&/^[[:alpha:]][[:alnum:]_-]+:/{print "  ",substr($$1,1,index($$1,":")-1),c}' $(MAKEFILE_LIST) | column -t -s '  '

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@if [ ! -d cmd/scru-llm ]; then \
		echo "cmd/scru-llm is not implemented yet. This repository is currently a planning skeleton."; \
		exit 1; \
	fi
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## test: Run all tests
test:
	@echo "Running tests..."
	$(GO) test $(GOFLAGS) ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@$(GO) clean
	@echo "Clean complete"

## lint: Run linter
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, running go vet..."; \
		$(GO) vet ./...; \
	fi

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

## vet: Run go vet
vet:
	@echo "Running go vet..."
	$(GO) vet ./...

## run: Build and run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

## install: Install the binary to GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install $(MAIN_PACKAGE)
	@echo "Installed to $(GOPATH)/bin/$(BINARY_NAME)"

## deps: Download and verify dependencies
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod verify

## tidy: Clean up go.mod and go.sum
tidy:
	@echo "Tidying modules..."
	$(GO) mod tidy

## generate: Run go generate
generate:
	@echo "Running go generate..."
	$(GO) generate ./...

## check: Run all checks (fmt, vet, lint, test)
check: fmt vet lint test
	@echo "All checks passed!"

## dev-setup: Set up development environment
dev-setup:
	@echo "Setting up development environment..."
	@mkdir -p ~/.scru-llm/workspaces
	@mkdir -p ~/.scru-llm/logs
	@if [ ! -f config/scru-llm.yaml ]; then \
		echo "Creating config from example..."; \
		cp config/scru-llm.example.yaml config/scru-llm.yaml; \
		echo "Please edit config/scru-llm.yaml with your settings"; \
	fi
	@echo "Development environment ready!"

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):latest .

## docker-run: Run in Docker container
docker-run: docker-build
	@echo "Running in Docker..."
	docker run --rm -it $(BINARY_NAME):latest

## benchmark: Run benchmarks
benchmark:
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem ./...

## version: Show version
version:
	@$(GO) version
