#!/bin/bash
set -e
echo "Running Terracotta tests"
echo "Running unit tests..."
go test -v ./...
echo "Testing build..."
go build -o /tmp/terracotta-test
rm -f /tmp/terracotta-test
echo "All tests passed!"
