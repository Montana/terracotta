#!/bin/bash
# scripts/cross-compile.sh - Cross-compilation script
set -e
BINARY_NAME="terracotta"
VERSION="1.0.0"
BUILD_DIR="build"
BUILD_FLAGS="-ldflags=-X main.version=${VERSION}"
echo "Cross-compiling Terracotta v${VERSION}"
mkdir -p ${BUILD_DIR}
declare -a platforms=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)
for platform in "${platforms[@]}"; do
    IFS='/' read -r -a platform_split <<< "$platform"
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name="${BINARY_NAME}-${GOOS}-${GOARCH}"
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi
    echo "Building for ${GOOS}/${GOARCH}..."
    GOOS=$GOOS GOARCH=$GOARCH go build $BUILD_FLAGS -o ${BUILD_DIR}/${output_name}
    if [ $? -ne 0 ]; then
        echo "Failed to build for ${GOOS}/${GOARCH}"
        exit 1
    fi
done
echo "Cross-compilation complete!"
echo "Binaries available in ${BUILD_DIR}/"
ls -la ${BUILD_DIR}/
