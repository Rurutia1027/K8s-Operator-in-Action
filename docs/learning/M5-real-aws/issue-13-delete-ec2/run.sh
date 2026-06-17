#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SAMPLE="$SCRIPT_DIR/../issue-12-create-ec2/sample-aws.yaml"
CR_NAME="web-server-learning-12"
REGION="${AWS_DEFAULT_REGION:-eu-central-1}"

echo "== Issue #13: Real AWS delete =="

if ! kubectl get ec2instance "$CR_NAME" &>/dev/null; then
  echo "WARN: CR $CR_NAME not found — nothing to delete"
  exit 0
fi

INSTANCE_ID="$(kubectl get ec2instance "$CR_NAME" -o jsonpath='{.status.instanceId}' 2>/dev/null || true)"
echo ">> Instance ID: ${INSTANCE_ID:-<none>}"

echo ">> Deleting CR..."
kubectl delete -f "$SAMPLE" --ignore-not-found --wait=true

if [[ -n "$INSTANCE_ID" ]] && command -v aws &>/dev/null; then
  echo ">> Checking AWS instance state..."
  aws ec2 describe-instances --instance-ids "$INSTANCE_ID" --region "$REGION" \
    --query 'Reservations[0].Instances[0].State.Name' --output text || echo "Instance not found (OK)"
fi

echo ">> OK — Issue #13 delete flow completed"
