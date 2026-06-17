# Issue #5 — Status Subresource Updates

**Milestone:** M2 Controller Basics  
**Needs AWS:** No

## Goal

Learn `r.Status().Update()` and observe `kubectl get` print columns.

## Run

```bash
./run.sh
```

## Key concept

Updating status triggers **another reconcile** (see `reconcile_timeline.md` at repo root).

## Acceptance criteria

- [ ] Fake status written on first reconcile when `status.instanceId` is empty
- [ ] `status.instanceId == i-fake123` after reconcile
