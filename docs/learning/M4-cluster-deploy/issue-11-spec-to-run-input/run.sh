#!/usr/bin/env bash
set -euo pipefail
REPO_ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
cd "$REPO_ROOT"
go test ./learning/M4-cluster-deploy/issue-11-spec-to-run-input/... -v
echo ">> OK — Issue #11 tests passed"
