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

// isWithinTimeRange checks if current hour is within the specified time range
// Handles cases where the time range crosses midnight (e.g., 17:00 to 04:00)
func isWithinTimeRange(currentHour, startHour, endHour int) bool {
	if startHour <= endHour {
		// Normal range (e.g., 9 to 17)
		return currentHour >= startHour && currentHour <= endHour
	} else {
		// Range crosses midnight (e.g., 17 to 4)
		return currentHour >= startHour || currentHour <= endHour
	}
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *MonitorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconcile called")

	monitor := &monitorv1.Monitor{}
	if err := r.Get(ctx, req.NamespacedName, monitor); err != nil {
		log.Error(err, "unable to fetch Monitor")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	startTime := monitor.Spec.Start
	endTime := monitor.Spec.End
	replicas := monitor.Spec.Replicas

	currentHour := time.Now().UTC().Hour()

	log.Info("Checking time range",
		"currentHour", currentHour,
		"startTime", startTime,
		"endTime", endTime,
		"targetReplicas", replicas)

	if isWithinTimeRange(currentHour, startTime, endTime) {
		log.Info("Within monitoring time range, checking deployments")

		// Track if we updated any deployments
		updated := false

		for _, deploy := range monitor.Spec.Deployments {
			deployment := &appsv1.Deployment{}
			if err := r.Get(ctx, types.NamespacedName{Namespace: deploy.Namespace, Name: deploy.Name}, deployment); err != nil {
				log.Error(err, "unable to fetch Deployment", "name", deploy.Name, "namespace", deploy.Namespace)
				monitor.Status.Status = monitorv1.FAILED
				if statusErr := r.Status().Update(ctx, monitor); statusErr != nil {
					log.Error(statusErr, "unable to update Monitor status")
				}
				return ctrl.Result{}, client.IgnoreNotFound(err)
			}

			currentReplicas := int32(0)
			if deployment.Spec.Replicas != nil {
				currentReplicas = *deployment.Spec.Replicas
			}

			if currentReplicas != replicas {
				log.Info("Scaling deployment",
					"name", deploy.Name,
					"namespace", deploy.Namespace,
					"from", currentReplicas,
					"to", replicas)

				deployment.Spec.Replicas = &replicas
				if err := r.Update(ctx, deployment); err != nil {
					monitor.Status.Status = monitorv1.FAILED
					if statusErr := r.Status().Update(ctx, monitor); statusErr != nil {
						log.Error(statusErr, "unable to update Monitor status")
					}
					log.Error(err, "unable to update Deployment", "name", deploy.Name, "namespace", deploy.Namespace)
					return ctrl.Result{}, err
				}
				updated = true
				log.Info("Successfully scaled deployment", "name", deploy.Name, "namespace", deploy.Namespace, "replicas", replicas)
			} else {
				log.Info("Deployment already at target replicas", "name", deploy.Name, "replicas", currentReplicas)
			}
		}

		if updated {
			monitor.Status.Status = monitorv1.SUCCESS
		} else {
			// All deployments already at target replicas
			monitor.Status.Status = monitorv1.SUCCESS
		}
	} else {
		log.Info("Outside monitoring time range", "currentHour", currentHour)
		monitor.Status.Status = "Outside time range"
	}

	// Update status
	if statusErr := r.Status().Update(ctx, monitor); statusErr != nil {
		log.Error(statusErr, "unable to update Monitor status")
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
