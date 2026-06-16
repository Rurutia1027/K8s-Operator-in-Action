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
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	computev1 "github.com/Rurutia1027/K8s-Operator-in-Action/api/v1"
)

const FinalizerName = "ec2instance.compute.cloud.com"

// DeleteHook simulates external cleanup (AWS terminate). Tests use StubDelete.
// Wire real AWS delete in production later; nil skips external cleanup.
type DeleteHook func(ctx context.Context, instance *computev1.Ec2Instance) error

// Ec2InstanceReconciler reconciles a Ec2Instance object.
type Ec2InstanceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Delete DeleteHook
}

// +kubebuilder:rbac:groups=compute.cloud.com,resources=ec2instances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=compute.cloud.com,resources=ec2instances/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=compute.cloud.com,resources=ec2instances/finalizers,verbs=update

// Reconcile adds/removes finalizers and runs delete cleanup when a CR is deleted.
func (r *Ec2InstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
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

	// Delete path: user deleted the CR → run cleanup hook → remove finalizer.
	if !ec2Instance.DeletionTimestamp.IsZero() {
		l.Info("deletion requested")
		if r.Delete != nil {
			if err := r.Delete(ctx, ec2Instance); err != nil {
				return ctrl.Result{RequeueAfter: time.Second}, err
			}
		}
		controllerutil.RemoveFinalizer(ec2Instance, FinalizerName)
		if err := r.Update(ctx, ec2Instance); err != nil {
			return ctrl.Result{RequeueAfter: time.Second}, err
		}
		return ctrl.Result{}, nil
	}

	// Create/update path: ensure finalizer is registered on the CR.
	if !controllerutil.ContainsFinalizer(ec2Instance, FinalizerName) {
		controllerutil.AddFinalizer(ec2Instance, FinalizerName)
		if err := r.Update(ctx, ec2Instance); err != nil {
			return ctrl.Result{RequeueAfter: time.Second}, err
		}
		return ctrl.Result{}, nil
	}

	l.Info("got resource", "instanceType", ec2Instance.Spec.InstanceType, "region", ec2Instance.Spec.Region)
	return ctrl.Result{}, nil
}

// StubDelete is a fast delete hook for tests (no AWS).
func StubDelete(_ context.Context, _ *computev1.Ec2Instance) error {
	time.Sleep(10 * time.Millisecond)
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Ec2InstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&computev1.Ec2Instance{}).
		Named("ec2instance").
		Complete(r)
}
