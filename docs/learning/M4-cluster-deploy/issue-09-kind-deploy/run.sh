#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
IMG="${IMG:-ec2operator:learning}"
SAMPLE="$SCRIPT_DIR/sample.yaml"
CR_NAME="ec2instance-kind-09"

echo "== Issue #9: Kind deploy smoke test =="

if ! command -v kind &>/dev/null; then
  echo "ERROR: kind not installed. See https://kind.sigs.k8s.io/"
  exit 1
fi
if ! kubectl cluster-info &>/dev/null; then
  echo "ERROR: no cluster. Run: kind create cluster"
  exit 1
fi

cd "$REPO_ROOT"
export USE_FAKE_EC2=true

echo ">> Building image $IMG ..."
make docker-build IMG="$IMG"

echo ">> Loading image into kind..."
kind load docker-image "$IMG"

echo ">> Installing CRDs and deploying operator..."
make install
make deploy IMG="$IMG"

echo ">> Waiting for controller pod..."
kubectl wait --for=condition=Available deployment/ec2operator-controller-manager \
  -n ec2operator-system --timeout=120s

echo ">> Applying sample CR..."
kubectl apply -f "$SAMPLE"

echo ">> Checking controller logs (last 20 lines)..."
kubectl logs -n ec2operator-system deployment/ec2operator-controller-manager --tail=20 || true

echo ">> Note: status will populate after you wire USE_FAKE_EC2 in cmd/main.go (Issues #6-#7)"
kubectl get ec2instance "$CR_NAME" -o wide || true

echo ">> Cleanup sample CR..."
kubectl delete -f "$SAMPLE" --ignore-not-found

echo ">> OK — Issue #9 deploy smoke passed (operator is running)"
