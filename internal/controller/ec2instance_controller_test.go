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
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	computev1 "github.com/Rurutia1027/K8s-Operator-in-Action/api/v1"
)

func sampleCR(name string) *computev1.Ec2Instance {
	return &computev1.Ec2Instance{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Spec: computev1.Ec2InstanceSpec{
			InstanceType: "t3.micro",
			AMIId:        "ami-placeholder",
			Region:       "us-east-1",
		},
	}
}

func reconcilerWithFake(ec2 EC2Client) *Ec2InstanceReconciler {
	if ec2 == nil {
		ec2 = NewFakeEC2Client()
	}
	return &Ec2InstanceReconciler{
		Client: k8sClient,
		Scheme: k8sClient.Scheme(),
		EC2:    ec2,
	}
}

func nn(name string) types.NamespacedName {
	return types.NamespacedName{Name: name, Namespace: "default"}
}

func cleanupCR(ctx context.Context, name string) {
	cr := &computev1.Ec2Instance{}
	if err := k8sClient.Get(ctx, nn(name), cr); err != nil {
		return
	}
	controllerutil.RemoveFinalizer(cr, FinalizerName)
	_ = k8sClient.Update(ctx, cr)
	_ = k8sClient.Delete(ctx, cr)
}

func reconcileUntilStatus(ctx context.Context, r *Ec2InstanceReconciler, req reconcile.Request) *computev1.Ec2Instance {
	for range 4 {
		_, err := r.Reconcile(ctx, req)
		Expect(err).NotTo(HaveOccurred())
	}
	updated := &computev1.Ec2Instance{}
	Expect(k8sClient.Get(ctx, req.NamespacedName, updated)).To(Succeed())
	return updated
}

// brokenEC2Client simulates AWS errors for Issue #8.
type brokenEC2Client struct{ err error }

func (b *brokenEC2Client) RunInstance(context.Context, *computev1.Ec2Instance) (*computev1.CreatedInstanceInfo, error) {
	return nil, b.err
}
func (b *brokenEC2Client) TerminateInstance(context.Context, string, string) error { return b.err }
func (b *brokenEC2Client) DescribeInstance(context.Context, string, string) (bool, *InstanceDetails, error) {
	return false, nil, b.err
}

var _ = Describe("Ec2Instance Controller", func() {
	ctx := context.Background()

	Context("When the resource does not exist", func() {
		It("returns without error", func() {
			result, err := reconcilerWithFake(nil).Reconcile(ctx, reconcile.Request{
				NamespacedName: nn("does-not-exist"),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(ctrl.Result{}))
		})
	})

	Context("full reconcile create", func() {
		const name = "full-create"

		AfterEach(func() { cleanupCR(ctx, name) })

		It("adds finalizer and sets i-fake001 in status", func() {
			Expect(k8sClient.Create(ctx, sampleCR(name))).To(Succeed())
			req := reconcile.Request{NamespacedName: nn(name)}
			updated := reconcileUntilStatus(ctx, reconcilerWithFake(nil), req)

			Expect(updated.Status.InstanceID).To(Equal(FakeFirstInstanceID))
			Expect(updated.Status.State).To(Equal(InstanceStateRunning))
			Expect(controllerutil.ContainsFinalizer(updated, FinalizerName)).To(BeTrue())
		})
	})

	Context("does not recreate when instance ID exists", func() {
		const name = "no-recreate"

		AfterEach(func() { cleanupCR(ctx, name) })

		It("keeps the same fake instance", func() {
			fake := NewFakeEC2Client()
			cr := sampleCR(name)
			controllerutil.AddFinalizer(cr, FinalizerName)
			Expect(k8sClient.Create(ctx, cr)).To(Succeed())
			cr.Status = computev1.Ec2InstanceStatus{InstanceID: FakeFirstInstanceID, State: InstanceStateRunning}
			Expect(k8sClient.Status().Update(ctx, cr)).To(Succeed())
			_, _ = fake.RunInstance(ctx, cr)

			req := reconcile.Request{NamespacedName: nn(name)}
			_, before, _ := fake.DescribeInstance(ctx, FakeFirstInstanceID, "us-east-1")
			_, err := reconcilerWithFake(fake).Reconcile(ctx, req)
			Expect(err).NotTo(HaveOccurred())
			_, after, _ := fake.DescribeInstance(ctx, FakeFirstInstanceID, "us-east-1")
			Expect(after).To(Equal(before))
		})
	})

	Context("full reconcile delete", func() {
		const name = "full-delete"

		It("terminates fake instance and removes finalizer", func() {
			fake := NewFakeEC2Client()
			cr := sampleCR(name)
			controllerutil.AddFinalizer(cr, FinalizerName)
			Expect(k8sClient.Create(ctx, cr)).To(Succeed())
			cr.Status = computev1.Ec2InstanceStatus{InstanceID: FakeFirstInstanceID, State: InstanceStateRunning}
			Expect(k8sClient.Status().Update(ctx, cr)).To(Succeed())
			_, _ = fake.RunInstance(ctx, cr)

			req := reconcile.Request{NamespacedName: nn(name)}
			Expect(k8sClient.Delete(ctx, cr)).To(Succeed())
			_, err := reconcilerWithFake(fake).Reconcile(ctx, req)
			Expect(err).NotTo(HaveOccurred())

			exists, _, _ := fake.DescribeInstance(ctx, FakeFirstInstanceID, "us-east-1")
			Expect(exists).To(BeFalse())

			got := &computev1.Ec2Instance{}
			err = k8sClient.Get(ctx, req.NamespacedName, got)
			if apierrors.IsNotFound(err) {
				return
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(controllerutil.ContainsFinalizer(got, FinalizerName)).To(BeFalse())
		})
	})

	Context("terminate failure", func() {
		const name = "terminate-fail"

		AfterEach(func() { cleanupCR(ctx, name) })

		It("requeues and keeps finalizer", func() {
			fake := NewFakeEC2Client()
			cr := sampleCR(name)
			controllerutil.AddFinalizer(cr, FinalizerName)
			Expect(k8sClient.Create(ctx, cr)).To(Succeed())
			cr.Status = computev1.Ec2InstanceStatus{InstanceID: FakeFirstInstanceID, State: InstanceStateRunning}
			Expect(k8sClient.Status().Update(ctx, cr)).To(Succeed())
			_, _ = fake.RunInstance(ctx, cr)
			Expect(k8sClient.Delete(ctx, cr)).To(Succeed())

			r := reconcilerWithFake(&brokenTerminateClient{fake: fake})
			result, err := r.Reconcile(ctx, reconcile.Request{NamespacedName: nn(name)})
			Expect(err).To(HaveOccurred())
			Expect(result.RequeueAfter).To(BeNumerically(">", time.Duration(0)))

			got := &computev1.Ec2Instance{}
			Expect(k8sClient.Get(ctx, nn(name), got)).To(Succeed())
			Expect(controllerutil.ContainsFinalizer(got, FinalizerName)).To(BeTrue())
		})
	})

	Context("drift detection", func() {
		const name = "drift-state"

		AfterEach(func() { cleanupCR(ctx, name) })

		It("updates status when fake instance state changes", func() {
			fake := NewFakeEC2Client()
			cr := sampleCR(name)
			controllerutil.AddFinalizer(cr, FinalizerName)
			Expect(k8sClient.Create(ctx, cr)).To(Succeed())
			info, _ := fake.RunInstance(ctx, cr)
			cr.Status = computev1.Ec2InstanceStatus{
				InstanceID: info.InstanceID,
				State:      InstanceStateRunning,
				PublicIP:   info.PublicIP,
				PrivateIP:  info.PrivateIP,
			}
			Expect(k8sClient.Status().Update(ctx, cr)).To(Succeed())
			fake.SetInstanceState(info.InstanceID, InstanceStateStopped)

			req := reconcile.Request{NamespacedName: nn(name)}
			_, err := reconcilerWithFake(fake).Reconcile(ctx, req)
			Expect(err).NotTo(HaveOccurred())

			updated := &computev1.Ec2Instance{}
			Expect(k8sClient.Get(ctx, nn(name), updated)).To(Succeed())
			Expect(updated.Status.State).To(Equal(InstanceStateStopped))
		})

		It("sets Unknown when instance disappears from fake", func() {
			fake := NewFakeEC2Client()
			cr := sampleCR(name)
			controllerutil.AddFinalizer(cr, FinalizerName)
			Expect(k8sClient.Create(ctx, cr)).To(Succeed())
			info, _ := fake.RunInstance(ctx, cr)
			cr.Status = computev1.Ec2InstanceStatus{
				InstanceID: info.InstanceID,
				State:      InstanceStateRunning,
				PublicIP:   info.PublicIP,
			}
			Expect(k8sClient.Status().Update(ctx, cr)).To(Succeed())
			fake.DeleteInstance(info.InstanceID)

			req := reconcile.Request{NamespacedName: nn(name)}
			_, err := reconcilerWithFake(fake).Reconcile(ctx, req)
			Expect(err).NotTo(HaveOccurred())

			updated := &computev1.Ec2Instance{}
			Expect(k8sClient.Get(ctx, nn(name), updated)).To(Succeed())
			Expect(updated.Status.State).To(Equal(InstanceStateUnknown))
			Expect(updated.Status.PublicIP).To(BeEmpty())
		})
	})

	Context("fake client error", func() {
		const name = "client-error"

		AfterEach(func() { cleanupCR(ctx, name) })

		It("surfaces RunInstance errors", func() {
			cr := sampleCR(name)
			controllerutil.AddFinalizer(cr, FinalizerName)
			Expect(k8sClient.Create(ctx, cr)).To(Succeed())

			r := reconcilerWithFake(&brokenEC2Client{err: errors.New("simulated AWS error")})
			_, err := r.Reconcile(ctx, reconcile.Request{NamespacedName: nn(name)})
			Expect(err).To(HaveOccurred())
		})
	})
})

type brokenTerminateClient struct {
	fake *FakeEC2Client
}

func (b *brokenTerminateClient) RunInstance(ctx context.Context, instance *computev1.Ec2Instance) (*computev1.CreatedInstanceInfo, error) {
	return b.fake.RunInstance(ctx, instance)
}
func (b *brokenTerminateClient) DescribeInstance(ctx context.Context, instanceID, region string) (bool, *InstanceDetails, error) {
	return b.fake.DescribeInstance(ctx, instanceID, region)
}
func (b *brokenTerminateClient) TerminateInstance(context.Context, string, string) error {
	return errors.New("terminate failed")
}
