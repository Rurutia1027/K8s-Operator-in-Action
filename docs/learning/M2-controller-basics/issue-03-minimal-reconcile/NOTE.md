# Issue #3 — Minimal Reconcile (Get + Log)

**Milestone:** M2 Controller Basics  
**Needs AWS:** No  
**Needs cluster:** envtest (local binary, no Kind required for tests)

## Goal

Run the smallest possible reconciler: fetch the CR, log fields, return success.

## Reference file

`reconciler.go` — copy this pattern into `internal/controller/ec2instance_controller.go`

## Run tests

```bash
./run.sh
```

## Deploy to Kind (optional)

```bash
# After copying into internal/controller/
make docker-build IMG=ec2operator:dev
make deploy IMG=ec2operator:dev
kubectl apply -f ../../M1-api-and-crd/issue-01-crd-api-types/sample.yaml
kubectl logs -n ec2operator-system deployment/ec2operator-controller-manager -f
```

## Acceptance criteria

- [ ] envtest test passes
- [ ] Reconcile returns no error for existing CR
- [ ] Reconcile returns nil for deleted CR (NotFound)

## CI

`go test ./learning/M2-controller-basics/issue-03-minimal-reconcile/...`
