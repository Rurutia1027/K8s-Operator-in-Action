# Issue #10 — Drift Detection (Fake Client)

**Milestone:** M4 Cluster Deploy  
**Needs AWS:** No

## Goal

When `status.instanceId` is set, call `DescribeInstance` and update status if the instance disappeared or changed state.

## Run

```bash
./run.sh
```

## Copy into main project

Uncomment drift block in `internal/controller/ec2instance_controller.go` and use `EC2Client.DescribeInstance`.

## Acceptance criteria

- [ ] Instance removed from fake → `status.state = Unknown`
- [ ] State change `running` → `stopped` updates status
