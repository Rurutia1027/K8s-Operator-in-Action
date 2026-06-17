# Issue #4 — Finalizer Flow (No AWS)

**Milestone:** M2 Controller Basics  
**Needs AWS:** No

## Goal

Learn deletion protection: when user deletes a CR, Kubernetes waits until the finalizer is removed.

## Behavior

1. **Create path:** append `ec2instance.compute.cloud.com` finalizer
2. **Delete path:** stub cleanup (sleep 10ms in tests) → remove finalizer

## Run

```bash
./run.sh
```

## Copy into main project

Replace AWS `deleteEc2Instance` with stub first, then wire real delete in Issue #13.

## Acceptance criteria

- [ ] Finalizer added after first reconcile on new CR
- [ ] Delete completes after stub cleanup removes finalizer
