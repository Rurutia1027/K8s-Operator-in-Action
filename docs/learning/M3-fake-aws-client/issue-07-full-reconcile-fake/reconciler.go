// Package issue07 implements full reconcile with FakeEC2Client (Issue #7).
package issue07

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	computev1 "github.com/shkatara/ec2Operator/api/v1"
	issue06 "github.com/shkatara/ec2Operator/learning/M3-fake-aws-client/issue-06-ec2-client-interface"
)

const FinalizerName = "ec2instance.compute.cloud.com"

// Reconciler is the target shape for internal/controller after Issue #7.
type Reconciler struct {
	client.Client
	EC2 issue06.EC2Client
}

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	ec2Instance := &computev1.Ec2Instance{}
	if err := r.Get(ctx, req.NamespacedName, ec2Instance); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if !ec2Instance.DeletionTimestamp.IsZero() {
		if ec2Instance.Status.InstanceID != "" {
			if err := r.EC2.TerminateInstance(ctx, ec2Instance.Status.InstanceID, ec2Instance.Spec.Region); err != nil {
				return ctrl.Result{Requeue: true}, err
			}
		}
		controllerutil.RemoveFinalizer(ec2Instance, FinalizerName)
		return ctrl.Result{}, r.Update(ctx, ec2Instance)
	}

	if ec2Instance.Status.InstanceID != "" {
		exists, details, err := r.EC2.DescribeInstance(ctx, ec2Instance.Status.InstanceID, ec2Instance.Spec.Region)
		if err != nil {
			return ctrl.Result{}, err
		}
		if !exists {
			ec2Instance.Status.State = "Unknown"
			return ctrl.Result{}, r.Status().Update(ctx, ec2Instance)
		}
		if ec2Instance.Status.State != details.State {
			ec2Instance.Status.State = details.State
			return ctrl.Result{}, r.Status().Update(ctx, ec2Instance)
		}
		return ctrl.Result{}, nil
	}

	if !controllerutil.ContainsFinalizer(ec2Instance, FinalizerName) {
		controllerutil.AddFinalizer(ec2Instance, FinalizerName)
		if err := r.Update(ctx, ec2Instance); err != nil {
			return ctrl.Result{Requeue: true}, err
		}
	}

	l.Info("creating instance via EC2 client")
	info, err := r.EC2.RunInstance(ctx, ec2Instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	ec2Instance.Status.InstanceID = info.InstanceID
	ec2Instance.Status.State = info.State
	ec2Instance.Status.PublicIP = info.PublicIP
	ec2Instance.Status.PrivateIP = info.PrivateIP
	ec2Instance.Status.PublicDNS = info.PublicDNS
	ec2Instance.Status.PrivateDNS = info.PrivateDNS
	return ctrl.Result{}, r.Status().Update(ctx, ec2Instance)
}
