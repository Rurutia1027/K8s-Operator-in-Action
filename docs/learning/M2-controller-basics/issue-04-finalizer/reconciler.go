// Package issue04 demonstrates finalizer flow without AWS (Issue #4).
package issue04

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	computev1 "github.com/shkatara/ec2Operator/api/v1"
)

const FinalizerName = "ec2instance.compute.cloud.com"

// DeleteHook simulates external cleanup (AWS terminate). Tests use a fast stub.
type DeleteHook func(ctx context.Context, instance *computev1.Ec2Instance) error

// FinalizerReconciler adds/removes finalizers; uses DeleteHook instead of AWS.
type FinalizerReconciler struct {
	client.Client
	Delete DeleteHook
}

func (r *FinalizerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	ec2Instance := &computev1.Ec2Instance{}
	if err := r.Get(ctx, req.NamespacedName, ec2Instance); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if !ec2Instance.DeletionTimestamp.IsZero() {
		l.Info("deletion requested")
		if r.Delete != nil {
			if err := r.Delete(ctx, ec2Instance); err != nil {
				return ctrl.Result{Requeue: true}, err
			}
		}
		controllerutil.RemoveFinalizer(ec2Instance, FinalizerName)
		if err := r.Update(ctx, ec2Instance); err != nil {
			return ctrl.Result{Requeue: true}, err
		}
		return ctrl.Result{}, nil
	}

	if !controllerutil.ContainsFinalizer(ec2Instance, FinalizerName) {
		controllerutil.AddFinalizer(ec2Instance, FinalizerName)
		if err := r.Update(ctx, ec2Instance); err != nil {
			return ctrl.Result{Requeue: true}, err
		}
	}
	return ctrl.Result{}, nil
}

// StubDelete is a fast delete hook for tests.
func StubDelete(_ context.Context, _ *computev1.Ec2Instance) error {
	time.Sleep(10 * time.Millisecond)
	return nil
}
