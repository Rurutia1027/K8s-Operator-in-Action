#!/usr/bin/env bash
set -euo pipefail
REPO_ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
cd "$REPO_ROOT"
go test ./learning/M4-cluster-deploy/issue-10-drift-detection/... -v
echo ">> OK — Issue #10 tests passed"
