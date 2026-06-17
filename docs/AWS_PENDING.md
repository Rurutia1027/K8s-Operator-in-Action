
- [ ] Use dedicated AWS account or isolated project scope
- [ ] Add mandatory tags for cost tracking
- [ ] Set budget alarm in Billing
- [ ] Restrict instance type(s) for testing
- [ ] Clean up leftover instances daily
## 10) Common Failure Mapping
### `AuthFailure` / `InvalidClientTokenId`
- Access key is invalid or wrong profile/env is in use
### `UnauthorizedOperation`
- IAM policy is missing required EC2 actions
### `InvalidAMIID.NotFound`
- AMI is not valid in selected region
### `InvalidSubnetID.NotFound` / `InvalidGroup.NotFound`
- Subnet/SG mismatch with region or account
### `InsufficientInstanceCapacity`
- Change AZ or instance type
## 11) Done Criteria for M5
- [ ] Issue #12: real create works
- [ ] Issue #13: real delete works
- [ ] Issue #14: E2E checklist completed
- [ ] No orphan EC2 resources left
- [ ] Docs updated with final environment setup
