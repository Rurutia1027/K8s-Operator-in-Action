# Issue #8 — envtest Unit Test Suite

**Milestone:** M3 Fake AWS Client  
**Needs AWS:** No

## Goal

Consolidate controller tests: create, idempotent reconcile, delete, fake client error.

## Run

```bash
./run.sh
```

## Copy into main project

Merge patterns into `internal/controller/ec2instance_controller_test.go`.

## Acceptance criteria

- [ ] ≥4 test cases pass
- [ ] Coverage of create + delete + no-recreate paths
