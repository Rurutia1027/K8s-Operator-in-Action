#!/usr/bin/env bash
set -euo pipefail
REPO_ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
cd "$REPO_ROOT"
make setup-envtest
go test ./learning/M3-fake-aws-client/issue-07-full-reconcile-fake/... -v
echo ">> OK — Issue #7 tests passed"
