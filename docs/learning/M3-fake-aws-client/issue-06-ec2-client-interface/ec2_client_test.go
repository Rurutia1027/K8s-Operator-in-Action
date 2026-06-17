package issue06

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	computev1 "github.com/shkatara/ec2Operator/api/v1"
)

func TestFakeEC2Client_RunInstance(t *testing.T) {
	g := NewWithT(t)
	fake := NewFakeEC2Client()

	cr := &computev1.Ec2Instance{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec:       computev1.Ec2InstanceSpec{Region: "eu-central-1", InstanceType: "t3.micro", AMIId: "ami-x"},
	}

	info, err := fake.RunInstance(context.Background(), cr)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(info.InstanceID).To(Equal("i-fake001"))
	g.Expect(info.State).To(Equal("running"))
}

func TestFakeEC2Client_TerminateAndDescribe(t *testing.T) {
	g := NewWithT(t)
	fake := NewFakeEC2Client()
	cr := &computev1.Ec2Instance{Spec: computev1.Ec2InstanceSpec{Region: "eu-central-1"}}

	info, err := fake.RunInstance(context.Background(), cr)
	g.Expect(err).NotTo(HaveOccurred())

	exists, details, err := fake.DescribeInstance(context.Background(), info.InstanceID, "eu-central-1")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(exists).To(BeTrue())
	g.Expect(details.State).To(Equal("running"))

	g.Expect(fake.TerminateInstance(context.Background(), info.InstanceID, "eu-central-1")).To(Succeed())

	exists, _, err = fake.DescribeInstance(context.Background(), info.InstanceID, "eu-central-1")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(exists).To(BeFalse())
}

func TestFakeEC2Client_SetInstanceState(t *testing.T) {
	g := NewWithT(t)
	fake := NewFakeEC2Client()
	info, _ := fake.RunInstance(context.Background(), &computev1.Ec2Instance{Spec: computev1.Ec2InstanceSpec{Region: "r"}})

	fake.SetInstanceState(info.InstanceID, "stopped")
	_, details, _ := fake.DescribeInstance(context.Background(), info.InstanceID, "r")
	g.Expect(details.State).To(Equal("stopped"))
}
