package issue03

import (
	"testing"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	computev1 "github.com/shkatara/ec2Operator/api/v1"
	"github.com/shkatara/ec2Operator/learning/pkg/testenv"
)

func TestMinimalReconcile_GetAndLog(t *testing.T) {
	g := NewWithT(t)
	env := testenv.Setup(t)

	cr := &computev1.Ec2Instance{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "minimal-test",
			Namespace: "default",
		},
		Spec: computev1.Ec2InstanceSpec{
			InstanceType: "t3.micro",
			AMIId:        "ami-placeholder",
			Region:       "eu-central-1",
		},
	}
	g.Expect(env.K8sClient.Create(env.Ctx, cr)).To(Succeed())

	r := &MinimalReconciler{Client: env.K8sClient}
	result, err := r.Reconcile(env.Ctx, reconcile.Request{
		NamespacedName: types.NamespacedName{Name: "minimal-test", Namespace: "default"},
	})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).To(Equal(ctrl.Result{}))
}

func TestMinimalReconcile_NotFound(t *testing.T) {
	g := NewWithT(t)
	env := testenv.Setup(t)

	r := &MinimalReconciler{Client: env.K8sClient}
	result, err := r.Reconcile(env.Ctx, reconcile.Request{
		NamespacedName: types.NamespacedName{Name: "does-not-exist", Namespace: "default"},
	})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).To(Equal(ctrl.Result{}))
}
