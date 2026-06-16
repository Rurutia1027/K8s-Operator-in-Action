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

func reconcilerWithStub() *Ec2InstanceReconciler {
	return &Ec2InstanceReconciler{
		Client: k8sClient,
		Scheme: k8sClient.Scheme(),
		Delete: StubDelete,
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

var _ = Describe("Ec2Instance Controller", func() {
	ctx := context.Background()

	Context("minimal reconcile", func() {
		const name = "minimal-test"

		BeforeEach(func() {
			cr := sampleCR(name)
			_ = k8sClient.Create(ctx, cr)
		})

		AfterEach(func() {
			cleanupCR(ctx, name)
		})

		It("reconciles an existing CR without error", func() {
			result, err := reconcilerWithStub().Reconcile(ctx, reconcile.Request{
				NamespacedName: nn(name),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(ctrl.Result{}))
		})
	})

	Context("When the resource does not exist", func() {
		It("returns without error", func() {
			result, err := reconcilerWithStub().Reconcile(ctx, reconcile.Request{
				NamespacedName: nn("does-not-exist"),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(ctrl.Result{}))
		})
	})

	Context("finalizer added on first reconcile", func() {
		const name = "finalizer-add"

		AfterEach(func() {
			cleanupCR(ctx, name)
		})

		It("adds finalizer ec2instance.compute.cloud.com", func() {
			Expect(k8sClient.Create(ctx, sampleCR(name))).To(Succeed())

			_, err := reconcilerWithStub().Reconcile(ctx, reconcile.Request{
				NamespacedName: nn(name),
			})
			Expect(err).NotTo(HaveOccurred())

			updated := &computev1.Ec2Instance{}
			Expect(k8sClient.Get(ctx, nn(name), updated)).To(Succeed())
			Expect(controllerutil.ContainsFinalizer(updated, FinalizerName)).To(BeTrue())
		})
	})

	Context("finalizer idempotent on second reconcile", func() {
		const name = "finalizer-twice"

		AfterEach(func() {
			cleanupCR(ctx, name)
		})

		It("does not duplicate finalizer", func() {
			Expect(k8sClient.Create(ctx, sampleCR(name))).To(Succeed())

			r := reconcilerWithStub()
			_, err := r.Reconcile(ctx, reconcile.Request{NamespacedName: nn(name)})
			Expect(err).NotTo(HaveOccurred())

			_, err = r.Reconcile(ctx, reconcile.Request{NamespacedName: nn(name)})
			Expect(err).NotTo(HaveOccurred())

			updated := &computev1.Ec2Instance{}
			Expect(k8sClient.Get(ctx, nn(name), updated)).To(Succeed())
			Expect(updated.Finalizers).To(ConsistOf(FinalizerName))
		})
	})

	Context("finalizer removed on delete", func() {
		const name = "finalizer-del"

		It("runs stub cleanup and deletes CR", func() {
			cr := sampleCR(name)
			controllerutil.AddFinalizer(cr, FinalizerName)
			Expect(k8sClient.Create(ctx, cr)).To(Succeed())

			r := reconcilerWithStub()
			_, err := r.Reconcile(ctx, reconcile.Request{NamespacedName: nn(name)})
			Expect(err).NotTo(HaveOccurred())

			Expect(k8sClient.Delete(ctx, cr)).To(Succeed())

			_, err = r.Reconcile(ctx, reconcile.Request{NamespacedName: nn(name)})
			Expect(err).NotTo(HaveOccurred())

			updated := &computev1.Ec2Instance{}
			err = k8sClient.Get(ctx, nn(name), updated)
			if apierrors.IsNotFound(err) {
				return
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(controllerutil.ContainsFinalizer(updated, FinalizerName)).To(BeFalse())
		})
	})

	Context("delete with nil DeleteHook", func() {
		const name = "finalizer-nil-hook"

		It("still removes finalizer without external cleanup", func() {
			cr := sampleCR(name)
			controllerutil.AddFinalizer(cr, FinalizerName)
			Expect(k8sClient.Create(ctx, cr)).To(Succeed())

			r := &Ec2InstanceReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
				Delete: nil,
			}

			Expect(k8sClient.Delete(ctx, cr)).To(Succeed())

			result, err := r.Reconcile(ctx, reconcile.Request{NamespacedName: nn(name)})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(ctrl.Result{}))

			updated := &computev1.Ec2Instance{}
			err = k8sClient.Get(ctx, nn(name), updated)
			Expect(apierrors.IsNotFound(err)).To(BeTrue())
		})
	})

	Context("delete hook failure", func() {
		const name = "finalizer-hook-fail"

		AfterEach(func() {
			cleanupCR(ctx, name)
		})

		It("requeues and keeps finalizer when cleanup fails", func() {
			cr := sampleCR(name)
			controllerutil.AddFinalizer(cr, FinalizerName)
			Expect(k8sClient.Create(ctx, cr)).To(Succeed())
			Expect(k8sClient.Delete(ctx, cr)).To(Succeed())

			r := &Ec2InstanceReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
				Delete: func(context.Context, *computev1.Ec2Instance) error {
					return errors.New("cleanup failed")
				},
			}

			result, err := r.Reconcile(ctx, reconcile.Request{NamespacedName: nn(name)})
			Expect(err).To(HaveOccurred())
			Expect(result.RequeueAfter).To(BeNumerically(">", time.Duration(0)))

			updated := &computev1.Ec2Instance{}
			Expect(k8sClient.Get(ctx, nn(name), updated)).To(Succeed())
			Expect(controllerutil.ContainsFinalizer(updated, FinalizerName)).To(BeTrue())
		})
	})

	Context("status updates", func() {
		const name = "status-update"
		AfterEach(func() {
			cleanupCR(ctx, name)
		})

		It("writes fake status when status is empty", func() {
			Expect(k8sClient.Create(ctx, sampleCR(name))).To(Succeed())

			r := reconcilerWithStub()
			_, err := r.Reconcile(ctx, reconcile.Request{NamespacedName: nn(name)})
			Expect(err).NotTo(HaveOccurred())

			// First reconcile may only add finalizer; second writes status.
			_, err = r.Reconcile(ctx, reconcile.Request{NamespacedName: nn(name)})
			Expect(err).NotTo(HaveOccurred())

			updated := &computev1.Ec2Instance{}
			Expect(k8sClient.Get(ctx, nn(name), updated)).To(Succeed())
			Expect(updated.Status.InstanceID).To(Equal("i-fake123"))
			Expect(updated.Status.State).To(Equal(InstanceStateRunning))
			Expect(updated.Status.PublicIP).To(Equal("203.0.113.1"))
			Expect(updated.Status.PrivateIP).To(Equal("10.0.0.1"))
		})

		It("does not overwrite existing status", func() {
			cr := sampleCR("status-idempotent")
			Expect(k8sClient.Create(ctx, cr)).To(Succeed())
			defer cleanupCR(ctx, "status-idempotent")
			// add finalizer first to pass finalizer gate
			r := reconcilerWithStub()
			_, err := r.Reconcile(ctx, reconcile.Request{NamespacedName: nn("status-idempotent")})
			Expect(err).NotTo(HaveOccurred())
			current := &computev1.Ec2Instance{}
			Expect(k8sClient.Get(ctx, nn("status-idempotent"), current)).To(Succeed())
			current.Status.InstanceID = "i-existing"
			current.Status.State = InstanceStateRunning
			Expect(k8sClient.Status().Update(ctx, current)).To(Succeed())
			_, err = r.Reconcile(ctx, reconcile.Request{NamespacedName: nn("status-idempotent")})
			Expect(err).NotTo(HaveOccurred())
			updated := &computev1.Ec2Instance{}
			Expect(k8sClient.Get(ctx, nn("status-idempotent"), updated)).To(Succeed())
			Expect(updated.Status.InstanceID).To(Equal("i-existing"))
		})
	})
})
