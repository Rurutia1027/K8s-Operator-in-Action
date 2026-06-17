#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
SAMPLE="$SCRIPT_DIR/sample.yaml"
CR_NAME="ec2instance-learning-02"

echo "== Issue #2: Install CRD + create CR =="

if ! kubectl cluster-info &>/dev/null; then
  echo "ERROR: No Kubernetes cluster. Run: kind create cluster"
  exit 1
fi

cd "$REPO_ROOT"

echo ">> Installing CRDs..."
make install

echo ">> Applying sample CR..."
kubectl apply -f "$SAMPLE"

echo ">> Waiting for CR..."
kubectl wait --for=condition=Established crd/ec2instances.compute.cloud.com --timeout=60s 2>/dev/null || true
kubectl get ec2instance "$CR_NAME" -o wide

STATUS="$(kubectl get ec2instance "$CR_NAME" -o jsonpath='{.status}' 2>/dev/null || echo "")"
if [[ -n "$STATUS" && "$STATUS" != "{}" ]]; then
  echo "WARN: status is not empty (controller may be running): $STATUS"
else
  echo ">> OK: status is empty (expected without controller)"
fi

echo ">> API resource registered:"
kubectl api-resources | grep -i ec2 || true

echo ">> Cleanup sample CR..."
kubectl delete -f "$SAMPLE" --ignore-not-found

echo ">> OK — Issue #2 checks passed"
