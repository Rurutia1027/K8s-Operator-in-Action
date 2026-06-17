#!/usr/bin/env bash
set -euo pipefail
REPO_ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
cd "$REPO_ROOT"
go test ./learning/M3-fake-aws-client/issue-06-ec2-client-interface/... -v
echo ">> OK — Issue #6 tests passed"
