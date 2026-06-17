# Issue #1 — CRD API Types

**Milestone:** M1 API & CRD  
**Needs AWS:** No  
**Needs cluster:** No

## Goal

Understand **Spec** (desired state) vs **Status** (observed state). Learn how Go structs map to Kubernetes CRDs.

## Files in this folder

| File | Purpose |
|------|---------|
| `sample.yaml` | Example Ec2Instance CR with placeholder values |
| `types_test.go` | Runnable tests for serialization & spec/status separation |
| `reference/api_types.go` | Annotated copy of API types to study |
| `run.sh` | Run all checks |

## Run

```bash
./run.sh
# or from repo root:
make -C learning test-issue-01
```

## What you should learn

1. `Ec2InstanceSpec` — what the user wants (instance type, AMI, region)
2. `Ec2InstanceStatus` — what the operator reports back (instanceId, state, IPs)
3. Kubebuilder markers (`+kubebuilder:object:root=true`, `+kubebuilder:subresource:status`)
4. After editing types in `api/v1/`, run `make generate manifests`

## Copy into main project

Copy patterns from `reference/api_types.go` into:

- `api/v1/ec2instance_types.go`
- `config/samples/compute_v1_ec2instance.yaml` (use `sample.yaml` as template)

## Acceptance criteria

- [ ] `go test` in this folder passes
- [ ] `make generate manifests` succeeds at repo root
- [ ] You can explain Spec vs Status

## CI

Passes: `learning.yml` → `unit` job (`go test ./learning/M1-api-and-crd/...`)
