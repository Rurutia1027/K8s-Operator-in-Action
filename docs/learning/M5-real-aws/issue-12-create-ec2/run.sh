#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
SAMPLE="$SCRIPT_DIR/sample-aws.yaml"

echo "== Issue #12: Real AWS create =="

if [[ "${USE_FAKE_EC2:-}" == "true" ]]; then
  echo "ERROR: USE_FAKE_EC2=true — set USE_FAKE_EC2=false for real AWS"
  exit 1
fi
if [[ -z "${AWS_ACCESS_KEY_ID:-}" || -z "${AWS_SECRET_ACCESS_KEY:-}" ]]; then
  echo "ERROR: Set AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY"
  exit 1
fi
if grep -q "REPLACE" "$SAMPLE"; then
  echo "ERROR: Edit sample-aws.yaml — replace REPLACE_ME placeholders (keypair, ami, subnet, sg)"
  exit 1
fi
if ! kubectl cluster-info &>/dev/null; then
  echo "ERROR: Kubernetes cluster required"
  exit 1
fi

cd "$REPO_ROOT"
export USE_FAKE_EC2=false

echo ">> Deploy operator (ensure main project wires real AWS client)..."
make install
make deploy IMG="${IMG:-ec2operator:aws-learning}"

kubectl wait --for=condition=Available deployment/ec2operator-controller-manager \
  -n ec2operator-system --timeout=120s

echo ">> Applying AWS sample CR..."
kubectl apply -f "$SAMPLE"

echo ">> Watch status (Ctrl+C to stop watching)..."
kubectl get ec2instance web-server-learning-12 -w
