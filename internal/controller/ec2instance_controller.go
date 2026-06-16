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
	computev1 "github.com/Rurutia1027/K8s-Operator-in-Action/api/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

const FinalizerName = "ec2instance.compute.cloud.com"

// Ec2InstanceReconciler reconciles a Ec2Instance object.
type Ec2InstanceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	EC2    EC2Client
}

// +kubebuilder:rbac:groups=compute.cloud.com,resources=ec2instances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=compute.cloud.com,resources=ec2instances/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=compute.cloud.com,resources=ec2instances/finalizers,verbs=update
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

	// Delete path: terminate cloud instance, then remove finalizer.
	if !ec2Instance.DeletionTimestamp.IsZero() {
		l.Info("deletion requested")
		if ec2Instance.Status.InstanceID != "" && r.EC2 != nil {
			if err := r.EC2.TerminateInstance(ctx, ec2Instance.Status.InstanceID, ec2Instance.Spec.Region); err != nil {
				return ctrl.Result{RequeueAfter: time.Second}, err
			}
		}
		controllerutil.RemoveFinalizer(ec2Instance, FinalizerName)
		if err := r.Update(ctx, ec2Instance); err != nil {
			return ctrl.Result{RequeueAfter: time.Second}, err
		}
		return ctrl.Result{}, nil
	}

	// Drift detection: instance already created.
	if ec2Instance.Status.InstanceID != "" {
		if r.EC2 == nil {
			return ctrl.Result{}, nil
		}

		exists, details, err := r.EC2.DescribeInstance(ctx, ec2Instance.Status.InstanceID, ec2Instance.Spec.Region)
		if err != nil {
			return ctrl.Result{}, err
		}

		drift := DetectDrift(ec2Instance, exists, details)
		if drift.Changed {
			ec2Instance.Status = drift.Status
			if err := r.Status().Update(ctx, ec2Instance); err != nil {
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil
	}

	// Create path: ensure finalizer, then RunInstance.
	if !controllerutil.ContainsFinalizer(ec2Instance, FinalizerName) {
		controllerutil.AddFinalizer(ec2Instance, FinalizerName)
		if err := r.Update(ctx, ec2Instance); err != nil {
			return ctrl.Result{RequeueAfter: time.Second}, err
		}
	}

	if r.EC2 == nil {
		return ctrl.Result{}, nil
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

	if err := r.Status().Update(ctx, ec2Instance); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *Ec2InstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&computev1.Ec2Instance{}).
		Named("ec2instance").
		Complete(r)
}
