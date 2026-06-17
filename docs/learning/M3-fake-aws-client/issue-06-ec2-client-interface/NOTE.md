# Issue #6 — EC2Client Interface + Fake Implementation

**Milestone:** M3 Fake AWS Client  
**Needs AWS:** No  
**Needs cluster:** No (pure Go tests)

## Goal

**Critical fork:** controller depends on `EC2Client` interface, not AWS SDK directly.

## Files

- `ec2_client.go` — interface definition
- `fake_ec2_client.go` — in-memory implementation
- `ec2_client_test.go` — runnable tests

## Run

```bash
./run.sh
```

## Copy into main project

1. Move interface to `internal/controller/ec2_client.go`
2. Move fake to `internal/controller/fake_ec2_client.go`
3. In `cmd/main.go`:

```go
if os.Getenv("USE_FAKE_EC2") == "true" {
    reconciler.EC2 = controller.NewFakeEC2Client()
} else {
    reconciler.EC2 = controller.NewAWSEC2Client()
}
```

## Acceptance criteria

- [ ] All tests pass without network
- [ ] Fake returns deterministic `i-fake001`
