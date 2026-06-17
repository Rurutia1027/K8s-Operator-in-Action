# Issue #2 — Install CRD and Create CR Manually

**Milestone:** M1 API & CRD  
**Needs AWS:** No  
**Needs cluster:** Yes (Kind / minikube / any K8s)

## Goal

Confirm Kubernetes recognizes your CR **without** running the operator.

## Prerequisites

```bash
kubectl cluster-info
```

## Run

```bash
./run.sh
```

## What happens

1. `make install` — applies CRD to cluster
2. `kubectl apply -f sample.yaml` — creates Ec2Instance CR
3. Verifies CR exists and **status is empty**
4. Cleans up the sample CR (CRD stays installed)

## Copy into main project

- Use `sample.yaml` as template for `config/samples/compute_v1_ec2instance.yaml`

## Acceptance criteria

- [ ] `kubectl get crd ec2instances.compute.cloud.com` succeeds
- [ ] `kubectl get ec2instances` shows the sample
- [ ] `.status` is empty (no controller yet)

## CI

Optional cluster job; locally run `./run.sh` after `kind create cluster`.
