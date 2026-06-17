#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "== Issue #14: Full AWS E2E =="
echo ">> Step 1/2: Create (Issue #12)"
"$SCRIPT_DIR/../issue-12-create-ec2/run.sh" &
CREATE_PID=$!

# Wait up to 5 minutes for instance ID
for i in $(seq 1 30); do
  ID="$(kubectl get ec2instance web-server-learning-12 -o jsonpath='{.status.instanceId}' 2>/dev/null || true)"
  if [[ -n "$ID" ]]; then
    echo ">> Instance ready: $ID"
    kill $CREATE_PID 2>/dev/null || true
    break
  fi
  sleep 10
done

kill $CREATE_PID 2>/dev/null || true

STATE="$(kubectl get ec2instance web-server-learning-12 -o jsonpath='{.status.state}' 2>/dev/null || true)"
echo ">> State: ${STATE:-unknown}"

echo ">> Step 2/2: Delete (Issue #13)"
"$SCRIPT_DIR/../issue-13-delete-ec2/run.sh"

echo ">> Review checklist.md and fill in timings"
cat "$SCRIPT_DIR/checklist.md"

echo ">> OK — Issue #14 E2E script finished"
