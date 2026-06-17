# Issue #7 — Full Reconcile with Fake EC2 Client

**Milestone:** M3 Fake AWS Client  
**Needs AWS:** No

## Goal

Wire create / exists / delete paths through `FakeEC2Client`.

## State machine

| Condition | Action |
|-----------|--------|
| `DeletionTimestamp` set | Terminate → remove finalizer |
| `status.instanceId` empty | Add finalizer → RunInstance → update status |
| `status.instanceId` set | Describe (no-op if healthy) |

## Run

```bash
./run.sh
```

## Deploy with fake (optional)

```bash
USE_FAKE_EC2=true make docker-build IMG=ec2operator:dev
USE_FAKE_EC2=true make deploy IMG=ec2operator:dev
```

## Acceptance criteria

- [ ] Create sets `i-fake001` in status
- [ ] Delete removes finalizer after fake terminate
