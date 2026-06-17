/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	computev1 "github.com/Rurutia1027/K8s-Operator-in-Action/api/v1"
)

// DriftResult tells the reconciler whether status should be updated.
type DriftResult struct {
	Changed bool
	Status  computev1.Ec2InstanceStatus
}

// DetectDrift compares CR status with cloud state from EC2Client.
func DetectDrift(cr *computev1.Ec2Instance, exists bool, details *InstanceDetails) DriftResult {
	current := cr.Status
	out := current

	if cr.Status.InstanceID == "" {
		return DriftResult{Changed: false, Status: out}
	}

	if !exists {
		out.State = InstanceStateUnknown
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
