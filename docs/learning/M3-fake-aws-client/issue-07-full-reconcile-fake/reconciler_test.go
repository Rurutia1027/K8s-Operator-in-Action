package issue07

import (
	"testing"

	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	computev1 "github.com/shkatara/ec2Operator/api/v1"
	issue06 "github.com/shkatara/ec2Operator/learning/M3-fake-aws-client/issue-06-ec2-client-interface"
	"github.com/shkatara/ec2Operator/learning/pkg/testenv"
)

func TestFullReconcile_CreateSetsStatus(t *testing.T) {
	g := NewWithT(t)
	env := testenv.Setup(t)
	fake := issue06.NewFakeEC2Client()

	cr := sampleCR("full-create")
	g.Expect(env.K8sClient.Create(env.Ctx, cr)).To(Succeed())

	r := &Reconciler{Client: env.K8sClient, EC2: fake}
	req := reconcile.Request{NamespacedName: types.NamespacedName{Name: "full-create", Namespace: "default"}}

	// First pass: add finalizer
	_, err := r.Reconcile(env.Ctx, req)
	g.Expect(err).NotTo(HaveOccurred())

	// Second pass: create + status (may need multiple passes depending on ordering)
	for i := 0; i < 3; i++ {
		_, err = r.Reconcile(env.Ctx, req)
		g.Expect(err).NotTo(HaveOccurred())
	}

	updated := &computev1.Ec2Instance{}
	g.Expect(env.K8sClient.Get(env.Ctx, req.NamespacedName, updated)).To(Succeed())
	g.Expect(updated.Status.InstanceID).To(Equal("i-fake001"))
	g.Expect(updated.Status.State).To(Equal("running"))
	g.Expect(controllerutil.ContainsFinalizer(updated, FinalizerName)).To(BeTrue())
}

func TestFullReconcile_DeleteRemovesFinalizer(t *testing.T) {
	g := NewWithT(t)
	env := testenv.Setup(t)
	fake := issue06.NewFakeEC2Client()

	cr := sampleCR("full-delete")
	controllerutil.AddFinalizer(cr, FinalizerName)
	g.Expect(env.K8sClient.Create(env.Ctx, cr)).To(Succeed())
	cr.Status = computev1.Ec2InstanceStatus{InstanceID: "i-fake001", State: "running"}
	g.Expect(env.K8sClient.Status().Update(env.Ctx, cr)).To(Succeed())
	_, _ = fake.RunInstance(env.Ctx, cr)

	r := &Reconciler{Client: env.K8sClient, EC2: fake}
	req := reconcile.Request{NamespacedName: types.NamespacedName{Name: "full-delete", Namespace: "default"}}

	g.Expect(env.K8sClient.Delete(env.Ctx, cr)).To(Succeed())
	_, err := r.Reconcile(env.Ctx, req)
	g.Expect(err).NotTo(HaveOccurred())

	updated := &computev1.Ec2Instance{}
	err = env.K8sClient.Get(env.Ctx, req.NamespacedName, updated)
	if errors.IsNotFound(err) {
		return // object fully removed — success
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
