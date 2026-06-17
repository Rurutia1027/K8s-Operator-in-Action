// Package issue08 is the consolidated envtest suite (Issue #8).
package issue08

import (
	"context"
	"errors"
	"testing"

	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	computev1 "github.com/shkatara/ec2Operator/api/v1"
	issue06 "github.com/shkatara/ec2Operator/learning/M3-fake-aws-client/issue-06-ec2-client-interface"
	issue07 "github.com/shkatara/ec2Operator/learning/M3-fake-aws-client/issue-07-full-reconcile-fake"
	"github.com/shkatara/ec2Operator/learning/pkg/testenv"
)

func TestSuite_CreateFinalizerAndStatus(t *testing.T) {
	g := NewWithT(t)
	env := testenv.Setup(t)
	r, req := newReconciler(env, issue06.NewFakeEC2Client(), "suite-create")

	cr := baseCR("suite-create")
	g.Expect(env.K8sClient.Create(env.Ctx, cr)).To(Succeed())

	for i := 0; i < 4; i++ {
		_, err := r.Reconcile(env.Ctx, req)
		g.Expect(err).NotTo(HaveOccurred())
	}

	got := &computev1.Ec2Instance{}
	g.Expect(env.K8sClient.Get(env.Ctx, req.NamespacedName, got)).To(Succeed())
	g.Expect(got.Status.InstanceID).NotTo(BeEmpty())
	g.Expect(controllerutil.ContainsFinalizer(got, issue07.FinalizerName)).To(BeTrue())
}

func TestSuite_DoesNotRecreateWhenInstanceIDExists(t *testing.T) {
	g := NewWithT(t)
	env := testenv.Setup(t)
	fake := issue06.NewFakeEC2Client()
	r, req := newReconciler(env, fake, "suite-no-recreate")

	cr := baseCR("suite-no-recreate")
	controllerutil.AddFinalizer(cr, issue07.FinalizerName)
	g.Expect(env.K8sClient.Create(env.Ctx, cr)).To(Succeed())
	cr.Status = computev1.Ec2InstanceStatus{InstanceID: "i-fake001", State: "running"}
	g.Expect(env.K8sClient.Status().Update(env.Ctx, cr)).To(Succeed())
	_, _ = fake.RunInstance(env.Ctx, cr)

	_, before, _ := fake.DescribeInstance(env.Ctx, "i-fake001", "eu-central-1")
	_, err := r.Reconcile(env.Ctx, req)
	g.Expect(err).NotTo(HaveOccurred())
	_, after, _ := fake.DescribeInstance(env.Ctx, "i-fake001", "eu-central-1")
	g.Expect(after).To(Equal(before))
}

func TestSuite_DeleteTerminatesAndRemovesFinalizer(t *testing.T) {
	g := NewWithT(t)
	env := testenv.Setup(t)
	fake := issue06.NewFakeEC2Client()
	r, req := newReconciler(env, fake, "suite-delete")

	cr := baseCR("suite-delete")
	controllerutil.AddFinalizer(cr, issue07.FinalizerName)
	g.Expect(env.K8sClient.Create(env.Ctx, cr)).To(Succeed())
	cr.Status = computev1.Ec2InstanceStatus{InstanceID: "i-fake001", State: "running"}
	g.Expect(env.K8sClient.Status().Update(env.Ctx, cr)).To(Succeed())
	_, _ = fake.RunInstance(env.Ctx, cr)

	g.Expect(env.K8sClient.Delete(env.Ctx, cr)).To(Succeed())
	_, err := r.Reconcile(env.Ctx, req)
	g.Expect(err).NotTo(HaveOccurred())

	got := &computev1.Ec2Instance{}
	err = env.K8sClient.Get(env.Ctx, req.NamespacedName, got)
	if apierrors.IsNotFound(err) {
		exists, _, _ := fake.DescribeInstance(env.Ctx, "i-fake001", "eu-central-1")
		g.Expect(exists).To(BeFalse())
		return
	}
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(controllerutil.ContainsFinalizer(got, issue07.FinalizerName)).To(BeFalse())

	exists, _, _ := fake.DescribeInstance(env.Ctx, "i-fake001", "eu-central-1")
	g.Expect(exists).To(BeFalse())
}

func TestSuite_ErrorFromFakeClientSurfaces(t *testing.T) {
	g := NewWithT(t)
	env := testenv.Setup(t)
	broken := &brokenClient{err: errors.New("simulated AWS error")}
	r, req := newReconciler(env, broken, "suite-error")

	cr := baseCR("suite-error")
	controllerutil.AddFinalizer(cr, issue07.FinalizerName)
	g.Expect(env.K8sClient.Create(env.Ctx, cr)).To(Succeed())

	_, err := r.Reconcile(env.Ctx, req)
	g.Expect(err).To(HaveOccurred())
}

type brokenClient struct{ err error }

func (b *brokenClient) RunInstance(context.Context, *computev1.Ec2Instance) (*computev1.CreatedInstanceInfo, error) {
	return nil, b.err
}
func (b *brokenClient) TerminateInstance(context.Context, string, string) error { return b.err }
func (b *brokenClient) DescribeInstance(context.Context, string, string) (bool, *issue06.InstanceDetails, error) {
	return false, nil, b.err
}

func newReconciler(env *testenv.Environment, ec2 issue06.EC2Client, name string) (*issue07.Reconciler, reconcile.Request) {
	return &issue07.Reconciler{Client: env.K8sClient, EC2: ec2},
		reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: "default"}}
}

func baseCR(name string) *computev1.Ec2Instance {
	return &computev1.Ec2Instance{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "default"},
		Spec: computev1.Ec2InstanceSpec{
			InstanceType: "t3.micro",
			AMIId:        "ami-placeholder",
			Region:       "eu-central-1",
		},
	}
}
