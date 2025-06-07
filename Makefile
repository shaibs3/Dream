.PHONY: build run test clean lint install-tools docker-up docker-down docker-up-build test-system

# Binary name
BINARY_NAME=dream

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run

# Docker Compose commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-up-build:
	docker-compose up --build

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/dream

# Run the application
run:
	$(GORUN) main.go

# Test the application
test:
	$(GOTEST) -v ./...

# Run system tests
test-system:
	$(GOTEST) -v system_test.go

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Install dependencies
deps:
	$(GOGET) -v ./...

# Install development tools
install-tools:
	@echo "Installing golangci-lint..."
	@if [ "$(shell uname)" = "Darwin" ]; then \
		brew install golangci-lint; \
	elif [ "$(shell uname)" = "Linux" ]; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2; \
	else \
		echo "Unsupported operating system"; \
		exit 1; \
	fi

# Lint the code
lint:
	golangci-lint run

# Help command
help:
	@echo "Available commands:"
	@echo "  make build         - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build files"
	@echo "  make deps         - Install dependencies"
	@echo "  make install-tools - Install development tools (golangci-lint)"
	@echo "  make lint         - Run linter"
	@echo "  make docker-up    - Start Docker Compose services"
	@echo "  make docker-down  - Stop Docker Compose services"
	@echo "  make docker-up-build - Start Docker Compose services with rebuild"
	@echo "  make test-system  - Run system tests" 