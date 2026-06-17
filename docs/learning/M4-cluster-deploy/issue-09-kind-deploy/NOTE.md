# Issue #9 — Deploy Operator on Kind (Fake AWS)

**Milestone:** M4 Cluster Deploy  
**Needs AWS:** No  
**Needs Kind:** Yes

## Goal

Full pipeline: build image → load Kind → deploy → apply CR → wait for fake status.

## Prerequisites

```bash
kind create cluster
```

## Run

```bash
./run.sh
```

## Environment

Uses `USE_FAKE_EC2=true` — implement this env check in `cmd/main.go` when copying Issue #6–#7.

## Acceptance criteria

- [ ] Controller pod Running
- [ ] Sample CR gets `status.instanceId` (when fake wired in main project)
- [ ] Script completes without error

## CI

`learning-e2e-kind.yml` runs a subset of this on every PR.
