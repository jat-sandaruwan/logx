#!/bin/bash

# logx Installation Script
# This script will build and install logx on your system

set -e

echo "========================================="
echo "  logx - Remote Log Viewer Installation"
echo "========================================="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.21 or higher."
    echo "Visit: https://golang.org/dl/"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.21"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "Error: Go version $REQUIRED_VERSION or higher is required."
    echo "Current version: $GO_VERSION"
    exit 1
fi

echo "✓ Go version $GO_VERSION detected"
echo ""

# Download dependencies
echo "Downloading dependencies..."
go mod download
go mod tidy
echo "✓ Dependencies downloaded"
echo ""

# Build the binary
echo "Building logx..."
mkdir -p build
go build -o build/logx cmd/logx/main.go
echo "✓ Build complete"
echo ""

# Install system-wide (optional)
read -p "Install logx system-wide? (requires sudo) [y/N]: " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "Installing to /usr/local/bin..."
    sudo cp build/logx /usr/local/bin/
    echo "✓ Installed to /usr/local/bin/logx"
else
    echo "Binary available at: $(pwd)/build/logx"
    echo "To use it from anywhere, add to PATH or copy to /usr/local/bin/"
fi

echo ""
echo "========================================="
echo "  Installation Complete!"
echo "========================================="
echo ""
echo "Get started:"
echo "  1. Add a user:    logx user add"
echo "  2. Add an app:    logx app add"
echo "  3. View logs:     logx <appname>"
echo ""
echo "For more help:      logx help"
echo "View version:       logx version"
echo ""