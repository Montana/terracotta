#!/bin/bash
# scripts/build.sh - Build script for Terracotta
set -e
BINARY_NAME="terracotta"
VERSION="1.0.0"
BUILD_DIR="build"
echo "Building Terracotta v${VERSION}"
mkdir -p ${BUILD_DIR}
echo "Building for current platform..."
go build -ldflags="-X main.version=${VERSION}" -o ${BUILD_DIR}/${BINARY_NAME}
echo "Build complete: ${BUILD_DIR}/${BINARY_NAME}"
chmod +x ${BUILD_DIR}/${BINARY_NAME}
echo "Ready to run: ./${BUILD_DIR}/${BINARY_NAME} -help"
