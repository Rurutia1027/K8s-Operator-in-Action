#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

echo "== Issue #1: CRD API Types =="
cd "$REPO_ROOT"

echo ">> Running Go tests..."
go test ./learning/M1-api-and-crd/issue-01-crd-api-types/... -v

echo ">> Verifying code generation (dry run)..."
make generate manifests
test -f config/crd/bases/compute.cloud.com_ec2instances.yaml

echo ">> OK — Issue #1 checks passed"
