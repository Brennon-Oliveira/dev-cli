#!/bin/sh
set -eu

# Stop hooks must write valid JSON to stdout. Keep this hook advisory and
# non-blocking; detailed checks remain in the PostToolUse hook.
printf '{}\n'
