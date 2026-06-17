# M2: Controller Basics

Learn the reconcile loop, finalizers, and status updates — **no AWS**.

| Issue | Topic | Test |
|-------|-------|------|
| #3 | Minimal reconcile (Get + log) | envtest |
| #4 | Finalizer flow | envtest |
| #5 | Status subresource | envtest |

```bash
make -C learning setup-envtest
make -C learning test-issue-03
make -C learning test-issue-04
make -C learning test-issue-05
```

Copy reference `reconciler.go` from each issue into `internal/controller/ec2instance_controller.go` as you progress.
