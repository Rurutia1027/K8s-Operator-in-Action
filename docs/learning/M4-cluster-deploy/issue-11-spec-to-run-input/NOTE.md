# Issue #11 — Map Spec to RunInstancesInput

**Milestone:** M4 Cluster Deploy  
**Needs AWS:** No (pure unit tests)

## Goal

Map CR spec fields to AWS SDK `RunInstancesInput` without calling AWS.

## Run

```bash
./run.sh
```

## Reference

See `kubernetes/ec2Instance.yaml` in repo root for full spec example.

## Acceptance criteria

- [ ] Security groups, tags, userData, storage mapped
- [ ] ≥5 unit tests pass
