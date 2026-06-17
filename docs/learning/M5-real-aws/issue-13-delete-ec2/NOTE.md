# Issue #13 — Real AWS: Delete EC2

**Milestone:** M5 Real AWS  
**Needs AWS:** Yes

## Goal

`kubectl delete` terminates the EC2 instance and removes the finalizer.

## Run

```bash
./run.sh
```

## Verify

```bash
aws ec2 describe-instances --instance-ids <id> --region eu-central-1
# State: terminated
```

## Acceptance criteria

- [ ] EC2 terminated after CR delete
- [ ] No orphaned instances
