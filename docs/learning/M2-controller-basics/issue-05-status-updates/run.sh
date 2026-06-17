#!/usr/bin/env bash
set -euo pipefail
REPO_ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
cd "$REPO_ROOT"
make setup-envtest
go test ./learning/M2-controller-basics/issue-05-status-updates/... -v
echo ">> OK — Issue #5 tests passed"
