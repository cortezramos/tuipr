#!/usr/bin/env bash
set -euo pipefail

# test-brew-local.sh - Builds and runs Homebrew formula tests in Docker
# Usage: ./scripts/test-brew-local.sh

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
IMAGE_NAME="tuipr-brew-test"

echo "=== tuipr Homebrew Formula Test Runner ==="
echo "Project root: ${PROJECT_ROOT}"
echo ""

# Build Docker image
echo "[1/2] Building Docker image..."
docker build \
  -t "${IMAGE_NAME}" \
  -f "${PROJECT_ROOT}/docker/tests/Dockerfile" \
  "${PROJECT_ROOT}"
echo ""

# Run tests
echo "[2/2] Running formula tests in Docker..."
docker run --rm "${IMAGE_NAME}"
EXIT_CODE=$?

echo ""
if [ "${EXIT_CODE}" -eq 0 ]; then
  echo "=== All Homebrew formula tests passed ==="
else
  echo "=== Some Homebrew formula tests failed ==="
fi

exit "${EXIT_CODE}"
