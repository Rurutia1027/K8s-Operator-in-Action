# M3: Fake AWS Client

Decouple from AWS SDK; test everything locally with `FakeEC2Client`.

| Issue | Topic |
|-------|-------|
| #6 | `EC2Client` interface + fake |
| #7 | Full reconcile state machine |
| #8 | envtest test suite |

```bash
make -C learning test-issue-06
make -C learning test-issue-07
make -C learning test-issue-08
```

Set `USE_FAKE_EC2=true` when deploying (Issue #9).
