#!/bin/bash

# Ensure swag is installed
if ! command -v swag &> /dev/null; then
    echo "Installing swag..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Clean any existing generated docs
rm -rf ./internal/docs

# Make sure the output directory exists
mkdir -p ./internal/docs

# Generate Swagger documentation
echo "Generating Swagger documentation..."
swag init \
    --output ./internal/docs \
    --parseDepth 10

echo "Swagger documentation generated in internal/docs/" 