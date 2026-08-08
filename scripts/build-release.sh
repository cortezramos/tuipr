#!/usr/bin/env bash
set -euo pipefail

# build-release.sh - Build tuipr for all platforms
# Usage: ./scripts/build-release.sh [version]
# If no version is provided, it uses "dev"

VERSION="${1:-dev}"
APP_NAME="tuipr"
BUILD_DIR="dist"

# Clean previous builds
rm -rf "${BUILD_DIR}"
mkdir -p "${BUILD_DIR}"

echo "🔨 Building ${APP_NAME} v${VERSION}..."

# Build for macOS Intel (amd64)
echo "  → macOS Intel (amd64)..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.version=${VERSION}" \
  -o "${BUILD_DIR}/${APP_NAME}" ./cmd/tuipr/
tar -czf "${BUILD_DIR}/${APP_NAME}_${VERSION}_Darwin_x86_64.tar.gz" -C "${BUILD_DIR}" "${APP_NAME}"
rm "${BUILD_DIR}/${APP_NAME}"

# Build for macOS Apple Silicon (arm64)
echo "  → macOS Apple Silicon (arm64)..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X main.version=${VERSION}" \
  -o "${BUILD_DIR}/${APP_NAME}" ./cmd/tuipr/
tar -czf "${BUILD_DIR}/${APP_NAME}_${VERSION}_Darwin_arm64.tar.gz" -C "${BUILD_DIR}" "${APP_NAME}"
rm "${BUILD_DIR}/${APP_NAME}"

# Build for Linux Intel (amd64)
echo "  → Linux Intel (amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.version=${VERSION}" \
  -o "${BUILD_DIR}/${APP_NAME}" ./cmd/tuipr/
tar -czf "${BUILD_DIR}/${APP_NAME}_${VERSION}_Linux_x86_64.tar.gz" -C "${BUILD_DIR}" "${APP_NAME}"
rm "${BUILD_DIR}/${APP_NAME}"

echo ""
echo "✅ Build complete! Archives in ${BUILD_DIR}/"
echo ""
ls -lh "${BUILD_DIR}"/*.tar.gz
