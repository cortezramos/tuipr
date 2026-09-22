#!/usr/bin/env bash
set -uo pipefail

# test-brew-install.sh - Simulates real brew install flow
# 1. Creates local tap with formula pointing to local HTTP server
# 2. Installs tuipr via brew (downloads from local server)
# 3. Verifies the installed binary works

eval "$(/home/brewuser/.linuxbrew/bin/brew shellenv)"

PASS=0
FAIL=0
SERVER_PID=""

pass() {
  echo "  ✓ $1"
  ((PASS++))
}

fail() {
  echo "  ✗ $1"
  ((FAIL++))
}

cleanup() {
  if [ -n "${SERVER_PID}" ]; then
    kill "${SERVER_PID}" 2>/dev/null || true
  fi
}
trap cleanup EXIT

echo "=== tuipr Brew Install Test (Real Flow) ==="
echo ""

# ---------------------------------------------------------------------------
# T1: Start local HTTP server with release binary
# ---------------------------------------------------------------------------
echo "[T1] Start local HTTP server for release binary"
cd /tmp
python3 -m http.server 8888 &>/dev/null &
SERVER_PID=$!
sleep 1
if curl -s -o /dev/null http://localhost:8888/tuipr_0.1.0_Linux_x86_64.tar.gz; then
  pass "HTTP server running, binary accessible"
else
  fail "HTTP server not accessible"
fi
echo ""

# ---------------------------------------------------------------------------
# T2: Create local tap with formula pointing to local server
# ---------------------------------------------------------------------------
echo "[T2] Create local tap with formula"
TAP_DIR="/home/brewuser/.linuxbrew/Homebrew/Library/Taps/local/homebrew-tuipr"
mkdir -p "${TAP_DIR}/Formula"

SHA256=$(sha256sum /tmp/tuipr_0.1.0_Linux_x86_64.tar.gz | cut -d' ' -f1)

cat > "${TAP_DIR}/Formula/tuipr.rb" <<RUBY
class Tuipr < Formula
  desc "Keyboard-driven Pull Request Lifecycle Manager for your terminal"
  homepage "https://github.com/cortezramos/tuipr"
  version "0.1.0"
  license "MIT"

  depends_on "gh"

  on_linux do
    if Hardware::CPU.intel?
      url "http://localhost:8888/tuipr_0.1.0_Linux_x86_64.tar.gz"
      sha256 "${SHA256}"
    end
  end

  def install
    bin.install "tuipr"
  end

  test do
    assert_path_exists bin/"tuipr"
    assert_predicate bin/"tuipr", :executable?
  end
end
RUBY
pass "Tap created at ${TAP_DIR}, SHA256=${SHA256}"
echo ""

# ---------------------------------------------------------------------------
# T3: brew install local/tuipr/tuipr (downloads from local server)
# ---------------------------------------------------------------------------
echo "[T3] brew install local/tuipr/tuipr (real download flow)"
INSTALL_OUTPUT=$(brew install local/tuipr/tuipr 2>&1) && INSTALL_EXIT=0 || INSTALL_EXIT=$?
echo "$INSTALL_OUTPUT" | tail -5
if echo "$INSTALL_OUTPUT" | grep -q "Cellar/tuipr"; then
  pass "tuipr installed successfully via brew install"
else
  echo "$INSTALL_OUTPUT"
  fail "brew install failed (exit code: ${INSTALL_EXIT})"
fi
echo ""

# ---------------------------------------------------------------------------
# T4: Verify installed binary works
# ---------------------------------------------------------------------------
echo "[T4] Verify tuipr --help works"
HELP_OUTPUT=$(tuipr --help 2>&1 || true)
if echo "$HELP_OUTPUT" | grep -q "tuipr"; then
  pass "--help outputs usage information"
else
  fail "--help did not produce expected output"
fi
echo ""

# ---------------------------------------------------------------------------
# T5: Verify tuipr --version works
# ---------------------------------------------------------------------------
echo "[T5] Verify tuipr --version works"
VERSION_OUTPUT=$(tuipr --version 2>&1 || true)
if echo "$VERSION_OUTPUT" | grep -q "tuipr"; then
  pass "--version outputs version info: ${VERSION_OUTPUT}"
else
  fail "--version did not produce expected output"
fi
echo ""

# ---------------------------------------------------------------------------
# T6: Verify binary location
# ---------------------------------------------------------------------------
echo "[T6] Verify binary is in Homebrew bin"
BREW_BIN=$(brew --prefix tuipr)/bin/tuipr
if [ -x "${BREW_BIN}" ]; then
  pass "Binary at ${BREW_BIN} is executable"
else
  fail "Binary not found at ${BREW_BIN}"
fi
echo ""

# ---------------------------------------------------------------------------
# T7: brew info tuipr
# ---------------------------------------------------------------------------
echo "[T7] brew info tuipr"
INFO_OUTPUT=$(brew info tuipr 2>&1 || true)
if echo "$INFO_OUTPUT" | grep -q "tuipr"; then
  pass "brew info shows formula details"
else
  fail "brew info did not show expected output"
fi
echo ""

# ---------------------------------------------------------------------------
# T8: brew test tuipr
# ---------------------------------------------------------------------------
echo "[T8] brew test tuipr"
TEST_OUTPUT=$(brew test tuipr 2>&1) && TEST_EXIT=0 || TEST_EXIT=$?
if [ "${TEST_EXIT}" -eq 0 ]; then
  pass "brew test passes"
else
  echo "$TEST_OUTPUT"
  fail "brew test failed (exit code: ${TEST_EXIT})"
fi
echo ""

# ---------------------------------------------------------------------------
# T9: brew uninstall tuipr
# ---------------------------------------------------------------------------
echo "[T9] brew uninstall tuipr"
UNINSTALL_OUTPUT=$(brew uninstall tuipr 2>&1) && UNINSTALL_EXIT=0 || UNINSTALL_EXIT=$?
if [ "${UNINSTALL_EXIT}" -eq 0 ]; then
  pass "Uninstall succeeds cleanly"
else
  echo "$UNINSTALL_OUTPUT"
  fail "Uninstall failed (exit code: ${UNINSTALL_EXIT})"
fi
echo ""

# ---------------------------------------------------------------------------
# Summary
# ---------------------------------------------------------------------------
echo "=== Summary ==="
echo "  Passed: ${PASS}"
echo "  Failed: ${FAIL}"
echo ""

if [ "${FAIL}" -gt 0 ]; then
  echo "RESULT: SOME TESTS FAILED"
  exit 1
else
  echo "RESULT: ALL TESTS PASSED"
  exit 0
fi
