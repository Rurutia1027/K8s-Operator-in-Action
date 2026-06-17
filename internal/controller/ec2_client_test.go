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

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	computev1 "github.com/Rurutia1027/K8s-Operator-in-Action/api/v1"
)

func TestFakeEC2Client_RunInstance(t *testing.T) {
	g := NewWithT(t)
	fake := NewFakeEC2Client()

	cr := &computev1.Ec2Instance{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec:       computev1.Ec2InstanceSpec{Region: "eu-central-1", InstanceType: testInstanceType, AMIId: testAMIX},
	}

	info, err := fake.RunInstance(context.Background(), cr)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(info.InstanceID).To(Equal(FakeFirstInstanceID))
	g.Expect(info.State).To(Equal(InstanceStateRunning))
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
	g.Expect(details.State).To(Equal(InstanceStateRunning))

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
