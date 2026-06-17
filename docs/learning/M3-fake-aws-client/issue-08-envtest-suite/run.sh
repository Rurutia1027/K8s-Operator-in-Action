#!/usr/bin/env bash
set -euo pipefail
REPO_ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
cd "$REPO_ROOT"
make setup-envtest
go test ./learning/M3-fake-aws-client/issue-08-envtest-suite/... -v
echo ">> OK — Issue #8 tests passed"
