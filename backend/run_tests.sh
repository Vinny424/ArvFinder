#!/bin/bash

# ARV Finder Backend Test Runner
# This script runs all unit tests for the ARV service

set -e  # Exit on any error

echo "ARV Finder Backend Test Suite"
echo "============================="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "ERROR: Go is not installed or not in PATH"
    exit 1
fi

# Display Go version
echo "Go version: $(go version)"
echo ""

# Run tests with verbose output
echo "Running ARV Service Tests..."
echo "----------------------------"

# Run specific service tests
go test ./services -v -cover

echo ""
echo "All tests completed!"
echo ""

# Optional: Run tests for all packages
echo "Running all package tests..."
echo "----------------------------"
go test ./... -v

echo ""
echo "Test suite finished successfully!"