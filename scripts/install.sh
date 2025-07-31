#!/bin/bash
# scripts/install.sh - Installation script for Terracotta
set -e
BINARY_NAME="terracotta"
INSTALL_DIR="/usr/local/bin"
BUILD_DIR="build"
echo "Installing Terracotta"
if [ ! -f "${BUILD_DIR}/${BINARY_NAME}" ]; then
    echo "Binary not found. Please run 'make build' first."
    exit 1
fi
if [ "$EUID" -ne 0 ]; then
    echo "This script requires sudo permissions to install to ${INSTALL_DIR}"
    sudo cp ${BUILD_DIR}/${BINARY_NAME} ${INSTALL_DIR}/
else
    cp ${BUILD_DIR}/${BINARY_NAME} ${INSTALL_DIR}/
fi
chmod +x ${INSTALL_DIR}/${BINARY_NAME}
echo "Terracotta installed to ${INSTALL_DIR}/${BINARY_NAME}"
echo "Run 'terracotta -help' to get started"
