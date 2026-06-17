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
	"context"
	"testing"

	computev1 "github.com/Rurutia1027/K8s-Operator-in-Action/api/v1"
	. "github.com/onsi/gomega"
)

func TestDetectDrift_InstanceGone(t *testing.T) {
	g := NewWithT(t)
	cr := &computev1.Ec2Instance{
		Status: computev1.Ec2InstanceStatus{
			InstanceID: FakeFirstInstanceID,
			State:      InstanceStateRunning,
			PublicIP:   "203.0.113.1",
		},
	}

	result := DetectDrift(cr, false, nil)
	g.Expect(result.Changed).To(BeTrue())
	g.Expect(result.Status.State).To(Equal(InstanceStateUnknown))
	g.Expect(result.Status.PublicIP).To(BeEmpty())
}

func TestDetectDrift_StateChanged(t *testing.T) {
	g := NewWithT(t)
	cr := &computev1.Ec2Instance{
		Status: computev1.Ec2InstanceStatus{InstanceID: FakeFirstInstanceID, State: InstanceStateRunning},
	}
	details := &InstanceDetails{InstanceID: FakeFirstInstanceID, State: InstanceStateStopped}
	result := DetectDrift(cr, true, details)
	g.Expect(result.Changed).To(BeTrue())
	g.Expect(result.Status.State).To(Equal(InstanceStateStopped))
}
func TestDetectDrift_NoChange(t *testing.T) {
	g := NewWithT(t)
	cr := &computev1.Ec2Instance{
		Status: computev1.Ec2InstanceStatus{
			InstanceID: FakeFirstInstanceID,
			State:      InstanceStateRunning,
			PublicIP:   FakeInstancePublicIP,
		},
	}
	details := &InstanceDetails{
		InstanceID: FakeFirstInstanceID,
		State:      InstanceStateRunning,
		PublicIP:   FakeInstancePublicIP,
	}
	result := DetectDrift(cr, true, details)
	g.Expect(result.Changed).To(BeFalse())
}
func TestDetectDrift_WithFakeClient(t *testing.T) {
	g := NewWithT(t)
	fake := NewFakeEC2Client()
	ctx := context.Background()
	cr := &computev1.Ec2Instance{Spec: computev1.Ec2InstanceSpec{Region: "us-east-1"}}
	info, err := fake.RunInstance(ctx, cr)
	g.Expect(err).NotTo(HaveOccurred())
	cr.Status = computev1.Ec2InstanceStatus{InstanceID: info.InstanceID, State: InstanceStateRunning}
	fake.SetInstanceState(info.InstanceID, InstanceStateStopped)
	exists, details, err := fake.DescribeInstance(ctx, info.InstanceID, "us-east-1")
	g.Expect(err).NotTo(HaveOccurred())
	result := DetectDrift(cr, exists, details)
	g.Expect(result.Status.State).To(Equal(InstanceStateStopped))
}
