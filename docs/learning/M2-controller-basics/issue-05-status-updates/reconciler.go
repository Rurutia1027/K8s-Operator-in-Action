// Package issue05 demonstrates status subresource updates (Issue #5).
package issue05

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	computev1 "github.com/shkatara/ec2Operator/api/v1"
)

// StatusReconciler writes fake status when instance ID is empty (no AWS).
type StatusReconciler struct {
	client.Client
}

func (r *StatusReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	ec2Instance := &computev1.Ec2Instance{}
	if err := r.Get(ctx, req.NamespacedName, ec2Instance); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if ec2Instance.Status.InstanceID != "" {
		l.Info("status already set", "instanceID", ec2Instance.Status.InstanceID)
		return ctrl.Result{}, nil
	}

	ec2Instance.Status.InstanceID = "i-fake123"
	ec2Instance.Status.State = "running"
	ec2Instance.Status.PublicIP = "203.0.113.1"
	ec2Instance.Status.PrivateIP = "10.0.0.1"

	if err := r.Status().Update(ctx, ec2Instance); err != nil {
		return ctrl.Result{}, err
	}
	l.Info("status updated", "instanceID", ec2Instance.Status.InstanceID)
	return ctrl.Result{}, nil
}
