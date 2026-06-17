// Package issue03 demonstrates a minimal reconciler (Issue #3).
// Copy Reconcile logic into internal/controller/ec2instance_controller.go.
package issue03

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	computev1 "github.com/shkatara/ec2Operator/api/v1"
)

// MinimalReconciler only fetches the resource and logs — no AWS, no finalizer.
type MinimalReconciler struct {
	client.Client
}

func (r *MinimalReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	l.Info("reconcile started", "namespace", req.Namespace, "name", req.Name)

	ec2Instance := &computev1.Ec2Instance{}
	if err := r.Get(ctx, req.NamespacedName, ec2Instance); err != nil {
		if errors.IsNotFound(err) {
			l.Info("resource deleted, nothing to do")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	l.Info("got resource", "instanceType", ec2Instance.Spec.InstanceType, "region", ec2Instance.Spec.Region)
	return ctrl.Result{}, nil
}
