package issue04

import (
	"testing"

	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	computev1 "github.com/shkatara/ec2Operator/api/v1"
	"github.com/shkatara/ec2Operator/learning/pkg/testenv"
)

func TestFinalizer_AddedOnReconcile(t *testing.T) {
	g := NewWithT(t)
	env := testenv.Setup(t)

	cr := sampleCR("finalizer-add")
	g.Expect(env.K8sClient.Create(env.Ctx, cr)).To(Succeed())

	r := &FinalizerReconciler{Client: env.K8sClient, Delete: StubDelete}
	_, err := r.Reconcile(env.Ctx, nn("finalizer-add"))
	g.Expect(err).NotTo(HaveOccurred())

	updated := &computev1.Ec2Instance{}
	g.Expect(env.K8sClient.Get(env.Ctx, types.NamespacedName{Name: "finalizer-add", Namespace: "default"}, updated)).To(Succeed())
	g.Expect(controllerutil.ContainsFinalizer(updated, FinalizerName)).To(BeTrue())
}

func TestFinalizer_RemovedOnDelete(t *testing.T) {
	g := NewWithT(t)
	env := testenv.Setup(t)

	cr := sampleCR("finalizer-del")
	controllerutil.AddFinalizer(cr, FinalizerName)
	g.Expect(env.K8sClient.Create(env.Ctx, cr)).To(Succeed())

	r := &FinalizerReconciler{Client: env.K8sClient, Delete: StubDelete}
	_, err := r.Reconcile(env.Ctx, nn("finalizer-del"))
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(env.K8sClient.Delete(env.Ctx, cr)).To(Succeed())

	// Reconcile deletion
	_, err = r.Reconcile(env.Ctx, nn("finalizer-del"))
	g.Expect(err).NotTo(HaveOccurred())

	updated := &computev1.Ec2Instance{}
	err = env.K8sClient.Get(env.Ctx, types.NamespacedName{Name: "finalizer-del", Namespace: "default"}, updated)
	if apierrors.IsNotFound(err) {
		return // fully deleted
	}
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(controllerutil.ContainsFinalizer(updated, FinalizerName)).To(BeFalse())
}

func sampleCR(name string) *computev1.Ec2Instance {
	return &computev1.Ec2Instance{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "default"},
		Spec: computev1.Ec2InstanceSpec{
			InstanceType: "t3.micro",
			AMIId:        "ami-placeholder",
			Region:       "eu-central-1",
		},
	}
}

func nn(name string) reconcile.Request {
	return reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: "default"}}
}
