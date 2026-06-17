# Issue #14 — Real AWS: End-to-End Validation

**Milestone:** M5 Real AWS  
**Needs AWS:** Yes

## Goal

Full lifecycle checklist: create → running → optional drift → delete.

## Run

```bash
./run.sh
```

## Manual drift test (optional)

1. Create instance via Issue #12
2. In AWS Console, manually terminate the EC2
3. Watch operator update `status.state` to `Unknown`

## Acceptance criteria

- [ ] All checklist items pass
- [ ] No leaked EC2 after delete
