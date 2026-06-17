package issue05

import (
	"testing"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	computev1 "github.com/shkatara/ec2Operator/api/v1"
	"github.com/shkatara/ec2Operator/learning/pkg/testenv"
)

func TestStatusReconciler_PopulatesFakeStatus(t *testing.T) {
	g := NewWithT(t)
	env := testenv.Setup(t)

	cr := &computev1.Ec2Instance{
		ObjectMeta: metav1.ObjectMeta{Name: "status-test", Namespace: "default"},
		Spec: computev1.Ec2InstanceSpec{
			InstanceType: "t3.micro",
			AMIId:        "ami-placeholder",
			Region:       "eu-central-1",
		},
	}
	g.Expect(env.K8sClient.Create(env.Ctx, cr)).To(Succeed())

	r := &StatusReconciler{Client: env.K8sClient}
	_, err := r.Reconcile(env.Ctx, reconcile.Request{
		NamespacedName: types.NamespacedName{Name: "status-test", Namespace: "default"},
	})
	g.Expect(err).NotTo(HaveOccurred())

	updated := &computev1.Ec2Instance{}
	g.Expect(env.K8sClient.Get(env.Ctx, types.NamespacedName{Name: "status-test", Namespace: "default"}, updated)).To(Succeed())
	g.Expect(updated.Status.InstanceID).To(Equal("i-fake123"))
	g.Expect(updated.Status.State).To(Equal("running"))
	g.Expect(updated.Status.PublicIP).To(Equal("203.0.113.1"))
}

func TestStatusReconciler_Idempotent(t *testing.T) {
	g := NewWithT(t)
	env := testenv.Setup(t)

	cr := &computev1.Ec2Instance{
		ObjectMeta: metav1.ObjectMeta{Name: "status-idem", Namespace: "default"},
		Spec:       computev1.Ec2InstanceSpec{InstanceType: "t3.micro", AMIId: "ami-x", Region: "eu-central-1"},
	}
	g.Expect(env.K8sClient.Create(env.Ctx, cr)).To(Succeed())

	cr.Status = computev1.Ec2InstanceStatus{InstanceID: "i-existing", State: "running"}
	g.Expect(env.K8sClient.Status().Update(env.Ctx, cr)).To(Succeed())

	r := &StatusReconciler{Client: env.K8sClient}
	_, err := r.Reconcile(env.Ctx, reconcile.Request{
		NamespacedName: types.NamespacedName{Name: "status-idem", Namespace: "default"},
	})
	g.Expect(err).NotTo(HaveOccurred())

	updated := &computev1.Ec2Instance{}
	g.Expect(env.K8sClient.Get(env.Ctx, types.NamespacedName{Name: "status-idem", Namespace: "default"}, updated)).To(Succeed())
	g.Expect(updated.Status.InstanceID).To(Equal("i-existing"))
}
