#!/bin/bash

# This script generates Swagger documentation for the BoilerGo API

# Set error reporting
set -e

echo "Updating Swagger dependencies..."
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/echo-swagger
go get -u github.com/swaggo/files

# Install swag command
echo "Installing swag command..."
go install github.com/swaggo/swag/cmd/swag@latest

# Add the Go bin directory to PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Make sure swag is in the path
which swag || { echo "Error: swag command not found in PATH"; exit 1; }

# Generate Swagger docs
echo "Generating Swagger documentation..."
swag init -g main.go -d ./ --parseDependency --parseDepth 2 --output ./docs

echo "Swagger documentation generated successfully!"
echo "You can access the Swagger UI at: http://localhost:8080/swagger/index.html"
