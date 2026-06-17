package issue10

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"

	computev1 "github.com/shkatara/ec2Operator/api/v1"
	issue06 "github.com/shkatara/ec2Operator/learning/M3-fake-aws-client/issue-06-ec2-client-interface"
)

func TestDetectDrift_InstanceGone(t *testing.T) {
	g := NewWithT(t)
	cr := &computev1.Ec2Instance{
		Status: computev1.Ec2InstanceStatus{
			InstanceID: "i-fake001",
			State:      "running",
			PublicIP:   "203.0.113.1",
		},
	}

	result := DetectDrift(cr, false, nil)
	g.Expect(result.Changed).To(BeTrue())
	g.Expect(result.Status.State).To(Equal("Unknown"))
	g.Expect(result.Status.PublicIP).To(BeEmpty())
}

func TestDetectDrift_StateChanged(t *testing.T) {
	g := NewWithT(t)
	cr := &computev1.Ec2Instance{
		Status: computev1.Ec2InstanceStatus{InstanceID: "i-fake001", State: "running"},
	}
	details := &issue06.InstanceDetails{InstanceID: "i-fake001", State: "stopped"}

	result := DetectDrift(cr, true, details)
	g.Expect(result.Changed).To(BeTrue())
	g.Expect(result.Status.State).To(Equal("stopped"))
}

func TestDetectDrift_NoChange(t *testing.T) {
	g := NewWithT(t)
	cr := &computev1.Ec2Instance{
		Status: computev1.Ec2InstanceStatus{
			InstanceID: "i-fake001",
			State:      "running",
			PublicIP:   "203.0.113.1",
		},
	}
	details := &issue06.InstanceDetails{
		InstanceID: "i-fake001",
		State:      "running",
		PublicIP:   "203.0.113.1",
	}

	result := DetectDrift(cr, true, details)
	g.Expect(result.Changed).To(BeFalse())
}

func TestDetectDrift_WithFakeClient(t *testing.T) {
	g := NewWithT(t)
	fake := issue06.NewFakeEC2Client()
	ctx := context.Background()

	cr := &computev1.Ec2Instance{Spec: computev1.Ec2InstanceSpec{Region: "eu-central-1"}}
	info, err := fake.RunInstance(ctx, cr)
	g.Expect(err).NotTo(HaveOccurred())

	cr.Status = computev1.Ec2InstanceStatus{InstanceID: info.InstanceID, State: "running"}
	fake.SetInstanceState(info.InstanceID, "stopped")

	exists, details, err := fake.DescribeInstance(ctx, info.InstanceID, "eu-central-1")
	g.Expect(err).NotTo(HaveOccurred())
	result := DetectDrift(cr, exists, details)
	g.Expect(result.Status.State).To(Equal("stopped"))
}
