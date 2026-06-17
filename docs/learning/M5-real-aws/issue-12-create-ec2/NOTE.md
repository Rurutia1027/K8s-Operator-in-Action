# Issue #12 — Real AWS: Create EC2

**Milestone:** M5 Real AWS  
**Needs AWS:** Yes  
**Needs keypair:** Yes

## Prerequisites

```bash
export AWS_ACCESS_KEY_ID=...
export AWS_SECRET_ACCESS_KEY=...
export AWS_DEFAULT_REGION=eu-central-1
export USE_FAKE_EC2=false

# Edit sample-aws.yaml with your real:
# - amiId, keyPair, subnet, securityGroups
```

## Run

```bash
./run.sh
```

## Copy into main project

- Wire `NewAWSEC2Client()` in `cmd/main.go`
- Use `BuildRunInstancesInput` from Issue #11 in `createInstance.go`

## Acceptance criteria

- [ ] `status.instanceId` is a real `i-...` ID
- [ ] Instance visible in AWS console

## CI

Manual only: `learning-aws-e2e.yml` workflow_dispatch
