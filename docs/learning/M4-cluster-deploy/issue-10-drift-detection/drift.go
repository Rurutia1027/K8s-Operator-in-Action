// Package issue10 implements drift detection logic (Issue #10).
package issue10

import (
	computev1 "github.com/shkatara/ec2Operator/api/v1"
	issue06 "github.com/shkatara/ec2Operator/learning/M3-fake-aws-client/issue-06-ec2-client-interface"
)

// DriftResult tells the reconciler whether status should be updated.
type DriftResult struct {
	Changed bool
	Status  computev1.Ec2InstanceStatus
}

// DetectDrift compares CR status with cloud state from EC2Client.
func DetectDrift(cr *computev1.Ec2Instance, exists bool, details *issue06.InstanceDetails) DriftResult {
	current := cr.Status
	out := current

	if cr.Status.InstanceID == "" {
		return DriftResult{Changed: false, Status: out}
	}

	if !exists {
		out.State = "Unknown"
		out.PublicIP = ""
		out.PrivateIP = ""
		return DriftResult{Changed: statusChanged(current, out), Status: out}
	}

	if details != nil {
		if out.State != details.State {
			out.State = details.State
		}
		if out.PublicIP != details.PublicIP {
			out.PublicIP = details.PublicIP
		}
		if out.PrivateIP != details.PrivateIP {
			out.PrivateIP = details.PrivateIP
		}
	}
	return DriftResult{Changed: statusChanged(current, out), Status: out}
}

func statusChanged(a, b computev1.Ec2InstanceStatus) bool {
	return a.State != b.State || a.PublicIP != b.PublicIP || a.PrivateIP != b.PrivateIP
}
