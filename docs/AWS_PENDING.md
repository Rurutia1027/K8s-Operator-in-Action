# AWS Pending Status and Resume Plan

This file records the current blocked state for AWS-related work and the exact resume workflow once AWS is ready.

## Current Status (Pending)

- AWS-dependent issues are pending because runtime AWS prerequisites are not fully ready yet (especially key pair and related network inputs).
- No blocker for non-AWS milestones (`M1` to `M4`).
- Pending scope is only `M5-real-aws`:
  - `issue-12-create-ec2`
  - `issue-13-delete-ec2`
  - `issue-14-e2e-validation`

## What Is Already Ready

- Structured learning path under `docs/learning/` by milestone and issue.
- Each issue folder includes runnable artifacts:
  - `NOTE.md` (context + expected outcome)
  - source/test files
  - `run.sh` (execution entry)
- CI workflow templates for learning:
  - `docs/learning/ci/learning.yml`
  - `docs/learning/ci/e2e-kind.yml`
  - `docs/learning/ci/aws-e2e.yml` (manual dispatch)

## AWS Prerequisites Checklist

Before resuming `M5`, ensure all items are available:

- [ ] `AWS_ACCESS_KEY_ID`
- [ ] `AWS_SECRET_ACCESS_KEY`
- [ ] `AWS_DEFAULT_REGION` (for example: `eu-central-1`)
- [ ] EC2 key pair name
- [ ] subnet id
- [ ] security group id
- [ ] AMI id compatible with chosen region
- [ ] cluster available (`kubectl cluster-info` succeeds)

## Resume From Context (When AWS Is Ready)

Follow the issue-local context in each `NOTE.md`, then run the matching script.

### 1) Issue 12: Create EC2

```bash
cd docs/learning/M5-real-aws/issue-12-create-ec2
# Edit sample-aws.yaml placeholders first
./run.sh
```

Expected:
- CR applied successfully
- `status.instanceId` populated with real EC2 id

### 2) Issue 13: Delete EC2

```bash
cd docs/learning/M5-real-aws/issue-13-delete-ec2
./run.sh
```

Expected:
- CR deletion triggers terminate flow
- instance is terminated (or no longer exists)

### 3) Issue 14: E2E Validation

```bash
cd docs/learning/M5-real-aws/issue-14-e2e-validation
./run.sh
```

Expected:
- create -> running -> delete lifecycle verified
- checklist completed in `checklist.md`

## Continue Developing Missing Parts

When resuming, use each issue folder as the source context and implement missing logic in project runtime code under:

- `internal/controller/`
- `cmd/main.go`
- `api/v1/` (if API updates are needed)

Recommended order:

1. Wire runtime client selection (`USE_FAKE_EC2` vs real AWS client).
2. Reuse builder mapping logic from `issue-11-spec-to-run-input`.
3. Validate reconcile create/delete/drift paths with issue scripts.
4. Run targeted tests for changed issue, then broader test pass.

## Quick Commands

Run all non-AWS learning tests:

```bash
make -C docs/learning test-all-no-aws
```

Run one issue:

```bash
make -C docs/learning test-issue-12
make -C docs/learning test-issue-13
make -C docs/learning test-issue-14
```

