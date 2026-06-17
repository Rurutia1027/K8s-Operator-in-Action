# EC2 Operator — Copy-Along Learning Path

Self-contained, runnable reference code for each GitHub issue.  
Work through milestones in order; AWS issues (#12–#14) come last.

## Current Progress

- Non-AWS tracks are ready for execution.
- AWS track is currently pending environment readiness.
- Resume notes are recorded in `AWS_PENDING.md`.

## Structure

```
learning/
├── Makefile                 # Run any issue: make test-issue-01
├── ci/                      # GitHub Actions workflows (copy to .github/workflows/)
├── M1-api-and-crd/
├── M2-controller-basics/
├── M3-fake-aws-client/
├── M4-cluster-deploy/
└── M5-real-aws/
```

## Quick start

```bash
# From repo root
make -C learning help
make -C learning test-all-no-aws    # Issues 1–11 (no cluster/AWS for most)
```

## Per-issue workflow

```bash
cd learning/M1-api-and-crd/issue-01-crd-api-types
cat NOTE.md          # Read the lesson
./run.sh             # Run tests / verification
```

## Milestones

| Milestone | Issues | Needs AWS | Needs K8s cluster |
|-----------|--------|-----------|-------------------|
| M1: API & CRD | #1–#2 | No | #2 only |
| M2: Controller Basics | #3–#5 | No | #3–#5 (Kind) |
| M3: Fake AWS Client | #6–#8 | No | #8 envtest only |
| M4: Cluster Deploy | #9–#11 | No | #9 Kind |
| M5: Real AWS | #12–#14 | **Yes** | Yes |

## CI

Copy workflows from `learning/ci/` to `.github/workflows/`:

```bash
cp learning/ci/learning.yml .github/workflows/learning.yml
cp learning/ci/e2e-kind.yml .github/workflows/learning-e2e-kind.yml
cp learning/ci/aws-e2e.yml .github/workflows/learning-aws-e2e.yml
```

## How to use with the main project

Each issue folder contains **reference code**. To implement in the real operator:

1. Read `NOTE.md`
2. Run `./run.sh` to see expected behavior
3. Copy patterns from `reference/` or top-level `.go` files into `internal/controller/`, `api/v1/`, etc.
4. Open a PR that closes the GitHub issue
