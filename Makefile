# Makefile for BoilerGo API

# Default target
.PHONY: all
all: run

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod tidy

# Run the server
.PHONY: run
run:
	@echo "Running server..."
	go run main.go

# Generate Swagger documentation
.PHONY: swagger
swagger:
	@echo "Generating Swagger documentation..."
	./scripts/generate-swagger.sh

# Build the application
.PHONY: build
build:
	@echo "Building application..."
	go build -o ./build/boilerGo main.go

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf ./build

# Help command
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make deps     - Install dependencies"
	@echo "  make run      - Run the server"
	@echo "  make swagger  - Generate Swagger documentation"
	@echo "  make build    - Build the application"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Clean build artifacts"
