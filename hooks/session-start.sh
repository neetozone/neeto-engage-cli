#!/bin/sh
# NeetoEngage CLI — session-start hook for Claude Code
# Lightweight auth liveness check. Always exits 0 (informational).

if ! command -v neetoengage >/dev/null 2>&1; then
  echo "NeetoEngage CLI is not installed or not on PATH."
  exit 0
fi

if neetoengage whoami >/dev/null 2>&1; then
  echo "NeetoEngage plugin active."
else
  echo "NeetoEngage CLI installed but not authenticated. Run 'neetoengage login' to authenticate."
fi

exit 0
