#!/usr/bin/env bash
set -uo pipefail

# test-brew.sh - Validates the tuipr Homebrew formula
# Runs inside the Docker container with Homebrew installed

eval "$(/home/brewuser/.linuxbrew/bin/brew shellenv)"

PASS=0
FAIL=0

pass() {
  echo "  ✓ $1"
  ((PASS++))
}

fail() {
  echo "  ✗ $1"
  ((FAIL++))
}

echo "=== tuipr Homebrew Formula Tests ==="
echo ""

# ---------------------------------------------------------------------------
# T1: brew install --build-from-source (from local tap)
# ---------------------------------------------------------------------------
echo "[T1] brew install --build-from-source local/tuipr/tuipr"
BREW_OUTPUT=$(brew install --build-from-source local/tuipr/tuipr 2>&1) && BREW_EXIT=0 || BREW_EXIT=$?
echo "$BREW_OUTPUT"
if echo "$BREW_OUTPUT" | grep -q "Cellar/tuipr"; then
  pass "Formula builds and installs successfully"
else
  fail "Formula failed to install (exit code: ${BREW_EXIT})"
fi
echo ""

# ---------------------------------------------------------------------------
# T2: brew audit (installed formula)
# ---------------------------------------------------------------------------
echo "[T2] brew audit tuipr"
AUDIT_OUTPUT=$(brew audit tuipr 2>&1) && AUDIT_EXIT=0 || AUDIT_EXIT=$?
if [ "${AUDIT_EXIT}" -eq 0 ]; then
  pass "Formula passes audit"
else
  echo "$AUDIT_OUTPUT"
  fail "Formula fails audit (exit code: ${AUDIT_EXIT})"
fi
echo ""

# ---------------------------------------------------------------------------
# T3: tuipr --help
# ---------------------------------------------------------------------------
echo "[T3] tuipr --help"
OUTPUT=$(tuipr --help 2>&1 || true)
if echo "$OUTPUT" | grep -q "tuipr"; then
  pass "--help prints usage information"
else
  fail "--help did not contain expected output"
fi
echo ""

# ---------------------------------------------------------------------------
# T4: tuipr --version
# ---------------------------------------------------------------------------
echo "[T4] tuipr --version"
OUTPUT=$(tuipr --version 2>&1 || true)
if echo "$OUTPUT" | grep -q "tuipr"; then
  pass "--version prints version info"
else
  fail "--version did not contain expected output"
fi
echo ""

# ---------------------------------------------------------------------------
# T5: brew test tuipr
# ---------------------------------------------------------------------------
echo "[T5] brew test tuipr"
TEST_OUTPUT=$(brew test tuipr 2>&1) && TEST_EXIT=0 || TEST_EXIT=$?
if [ "${TEST_EXIT}" -eq 0 ]; then
  pass "brew test passes"
else
  echo "$TEST_OUTPUT"
  fail "brew test failed (exit code: ${TEST_EXIT})"
fi
echo ""

# ---------------------------------------------------------------------------
# T6: Verify binary is executable
# ---------------------------------------------------------------------------
echo "[T6] Verify binary exists and is executable"
BREW_PREFIX=$(brew --prefix tuipr 2>/dev/null || echo "")
BINARY="${BREW_PREFIX}/bin/tuipr"
if [ -x "${BINARY}" ]; then
  pass "Binary is executable at ${BINARY}"
else
  fail "Binary not found or not executable at ${BINARY}"
fi
echo ""

# ---------------------------------------------------------------------------
# T7: brew info tuipr
# ---------------------------------------------------------------------------
echo "[T7] brew info tuipr"
OUTPUT=$(brew info tuipr 2>&1 || true)
if echo "$OUTPUT" | grep -q "tuipr"; then
  pass "brew info shows formula info"
else
  fail "brew info did not show expected output"
fi
echo ""

# ---------------------------------------------------------------------------
# T8: brew uninstall tuipr
# ---------------------------------------------------------------------------
echo "[T8] brew uninstall tuipr"
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
