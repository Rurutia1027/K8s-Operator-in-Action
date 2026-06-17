# Learning CI Workflows

Copy these to `.github/workflows/` to enable CI for the learning path:

```bash
cp learning/ci/learning.yml .github/workflows/learning.yml
cp learning/ci/e2e-kind.yml .github/workflows/learning-e2e-kind.yml
cp learning/ci/aws-e2e.yml .github/workflows/learning-aws-e2e.yml
```

## Workflow summary

| File | Trigger | What it runs |
|------|---------|--------------|
| `learning.yml` | push, PR | fmt, vet, lint, all learning Go tests (Issues 1–11) |
| `e2e-kind.yml` | push, PR | Kind + deploy smoke (Issue 9) |
| `aws-e2e.yml` | **manual only** | Real AWS E2E (Issues 12–14) |

## Secrets for aws-e2e.yml

- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- Optional: `AWS_KEY_PAIR`, `AWS_SUBNET_ID`, `AWS_SECURITY_GROUP_ID`, `AWS_AMI_ID`
