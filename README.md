# K8s-Operator-in-Action | [![E2E Tests](https://github.com/Rurutia1027/K8s-Operator-in-Action/actions/workflows/test-e2e.yml/badge.svg)](https://github.com/Rurutia1027/K8s-Operator-in-Action/actions/workflows/test-e2e.yml)

## Overview 
This repository contains my learning notes and hands-on experiments while learning Kubernetes Operators Go and Kubebuilder. 

Goal: 
- Understand Kubernetes Controllers
- Build Custom Operators
- Learn Platform Engineering fundamentals
- Manage cloud resources through Kubernetes APIs

## Learning Roadmap 
### Fundamentals 
[ ] Controller Pattern 
[ ] Reconcile Loop 
[ ] Idempotency 
[ ] CRD & Custom Resources 
[ ] Kubebuilder 

### Core Internals 
[ ] Manager Architecture 
[ ] Informers
[ ] Cache 
[ ] WorkQueue
[ ] Finalizers 


### Cloud Integration 
[ ] AWS SDK 
[ ] EC2 Operator 
[ ] Resource Lifecycle Management 


### Deployment 
[ ] Helm Packaging 
[ ] RBAC 
[ ] Service Accounts 
[ ] Production Deployment 


## Local Environment 
Requirements: 
- Go 1.23+
- Docker
- Kind
- kubectl
- Kubebuilder

Create a local cluster 

```shell
kind create cluster --name operator-dev 
```

Run Operator: 

```shell
make install
make run 
```


## Learning Project 
Current project: **Kubernetes EC2 Operator**

Features: 
- Custom Resource Definitions (CRDs)
- EC2 Lifecycle Management
- Reconcile Loop Implementation
- Finalizer-based Cleanup
- AWS SDK Integration

## Long-Term Goal 

Became a Cloud Native Engineer capable of building: 
- Kubernetes Operators
- Internal Developer Platforms (IDP)
- Platform Engineering Solutions
- Cloud Infrastructure Controllers
