# Issue #14 E2E Checklist

- [ ] **Create:** `kubectl apply` → `status.instanceId` populated within 5 min
- [ ] **Running:** `status.state == running`
- [ ] **IPs:** public/private IP present (if subnet allows)
- [ ] **Logs:** reconcile timeline matches `reconcile_timeline.md`
- [ ] **Drift (optional):** manual terminate → `status.state == Unknown`
- [ ] **Delete:** `kubectl delete` → EC2 `terminated` in AWS
- [ ] **Finalizer:** CR removed from cluster
- [ ] **Cost:** instance type is `t3.micro` for learning

## Record your run

| Step | Time | Notes |
|------|------|-------|
| Apply CR | | |
| Instance running | | |
| Delete CR | | |
| AWS terminated | | |
