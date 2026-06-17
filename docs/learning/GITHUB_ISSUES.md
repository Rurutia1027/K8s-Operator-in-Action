# GitHub Issues — copy-paste bodies

Create milestones first: `M1: API & CRD` … `M5: Real AWS`  
Labels: `learning`, `no-aws`, `needs-aws`, `test`, `ci`

Use: `gh issue create --title "..." --label "learning,no-aws" --body-file learning/issues/issue-01.md`

See each `learning/**/issue-*/NOTE.md` for the full lesson.  
Runnable code is in the matching issue folder.

| # | Title | Labels | Folder |
|---|-------|--------|--------|
| 1 | CRD API types | learning, no-aws | `M1-api-and-crd/issue-01-crd-api-types` |
| 2 | Install CRD manually | learning, no-aws | `M1-api-and-crd/issue-02-install-crd` |
| 3 | Minimal reconcile | learning, no-aws, test | `M2-controller-basics/issue-03-minimal-reconcile` |
| 4 | Finalizer flow | learning, no-aws, test | `M2-controller-basics/issue-04-finalizer` |
| 5 | Status updates | learning, no-aws, test | `M2-controller-basics/issue-05-status-updates` |
| 6 | EC2Client + Fake | learning, no-aws, test | `M3-fake-aws-client/issue-06-ec2-client-interface` |
| 7 | Full reconcile fake | learning, no-aws, test | `M3-fake-aws-client/issue-07-full-reconcile-fake` |
| 8 | envtest suite | learning, no-aws, test, ci | `M3-fake-aws-client/issue-08-envtest-suite` |
| 9 | Kind deploy | learning, no-aws, ci | `M4-cluster-deploy/issue-09-kind-deploy` |
| 10 | Drift detection | learning, no-aws, test | `M4-cluster-deploy/issue-10-drift-detection` |
| 11 | Spec → RunInput | learning, no-aws, test | `M4-cluster-deploy/issue-11-spec-to-run-input` |
| 12 | Real AWS create | learning, needs-aws | `M5-real-aws/issue-12-create-ec2` |
| 13 | Real AWS delete | learning, needs-aws | `M5-real-aws/issue-13-delete-ec2` |
| 14 | AWS E2E validation | learning, needs-aws, ci | `M5-real-aws/issue-14-e2e-validation` |

## Test command per issue

```bash
make -C learning test-issue-01   # through test-issue-14
```

## Batch create (example)

```bash
gh issue create --title "Issue 1: CRD API types" --label "learning,no-aws" \
  --body "See learning/M1-api-and-crd/issue-01-crd-api-types/NOTE.md. Run: make -C learning test-issue-01"
```
