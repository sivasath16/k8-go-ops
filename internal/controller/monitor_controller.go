/*
Copyright 2025.

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

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	monitorv1 "github.com/sivasath16/k8-go-ops/api/v1"
)

// MonitorReconciler reconciles a Monitor object
type MonitorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=monitor.sivasathwik.online,resources=monitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitor.sivasathwik.online,resources=monitors/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=monitor.sivasathwik.online,resources=monitors/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Monitor object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *MonitorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconcile called")

	// TODO(user): your logic here

	monitor := &monitorv1.Monitor{}
	if err := r.Get(ctx, req.NamespacedName, monitor); err != nil {
		log.Error(err, "unable to fetch Monitor")
		// Return error so it gets requeued
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	startTime := monitor.Spec.Start
	endTime := monitor.Spec.End
	replicas := monitor.Spec.Replicas

	currentHour := time.Now().UTC().Hour()
	if currentHour >= startTime && currentHour <= endTime {
		for _, deploy := range monitor.Spec.Deployments {
			deployment := &appsv1.Deployment{}
			if err := r.Get(ctx, types.NamespacedName{Namespace: deploy.Namespace, Name: deploy.Name}, deployment); err != nil {
				log.Error(err, "unable to fetch Deployment", "name", deploy.Name, "namespace", deploy.Namespace)
				return ctrl.Result{}, client.IgnoreNotFound(err)
			}
			if deployment.Spec.Replicas == nil || *deployment.Spec.Replicas != replicas {
				deployment.Spec.Replicas = &replicas
				if err := r.Update(ctx, deployment); err != nil {
					log.Error(err, "unable to update Deployment", "name", deploy.Name, "namespace", deploy.Namespace)
					return ctrl.Result{}, err
				}
				log.Info("Updated deployment replicas", "name", deploy.Name, "namespace", deploy.Namespace, "replicas", replicas)
			}
		}
	}
	return ctrl.Result{RequeueAfter: time.Duration(30 * time.Second)}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MonitorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitorv1.Monitor{}).
		Named("monitor").
		Complete(r)
}
