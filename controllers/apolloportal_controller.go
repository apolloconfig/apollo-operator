/*
Copyright 2023.

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

package controllers

import (
	apolloiov1alpha1 "apollo.io/apollo-operator/api/v1alpha1"
	"apollo.io/apollo-operator/pkg/reconcile/apolloportal"
	"context"
	"fmt"
	"github.com/go-logr/logr"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sync"
)

// ApolloPortalReconciler reconciles a ApolloPortal object
type ApolloPortalReconciler struct {
	client.Client
	recorder record.EventRecorder
	scheme   *runtime.Scheme
	log      logr.Logger
	//config   config.Config

	tasks   []Task
	muTasks sync.RWMutex
}

// Task represents a reconciliation task to be executed by the reconciler.
type Task struct {
	Do          func(context.Context, apolloportal.Params) error
	Name        string
	BailOnError bool
}

// Params is the set of options to build a new ApolloPortalReconciler.
type Params struct { // 感觉可以改个名字，不叫 Params
	client.Client
	Recorder record.EventRecorder
	Scheme   *runtime.Scheme
	Log      logr.Logger
	Tasks    []Task
	//Config   config.Config
}

// NewApolloPortalReconciler creates a new reconciler for ApolloPortal objects.
func NewApolloPortalReconciler(p Params) *ApolloPortalReconciler {
	r := &ApolloPortalReconciler{
		Client: p.Client,
		log:    p.Log,
		scheme: p.Scheme,
		//config:   p.Config,
		tasks:    p.Tasks,
		recorder: p.Recorder,
	}
	if len(r.tasks) == 0 {
		r.tasks = []Task{
			{
				apolloportal.ConfigMaps,
				"configmaps",
				true,
			},
			{
				apolloportal.ServiceAccounts,
				"serviceaccounts",
				true,
			},
			{
				apolloportal.Endpoints,
				"endpoints",
				true,
			},
			{
				apolloportal.Services,
				"services",
				true,
			},
			{
				apolloportal.Deployments,
				"deployments",
				true,
			},
			//{
			//	apolloportal.HorizontalPodAutoscalers,
			//	"horizontal pod autoscalers",
			//	true,
			//},
			//{
			//	apolloportal.DaemonSets,
			//	"daemon sets",
			//	true,
			//},
			{
				apolloportal.Ingresses,
				"ingresses",
				true,
			},
			{
				apolloportal.Self,
				"opentelemetry",
				true,
			},
		}
	}
	return r
}

//+kubebuilder:rbac:groups=apollo.io,resources=apolloportals,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apollo.io,resources=apolloportals/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apollo.io,resources=apolloportals/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ApolloPortal object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *ApolloPortalReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.log.WithValues("ApolloPortal", req.NamespacedName)

	// TODO(user): your logic here
	var instance apolloiov1alpha1.ApolloPortal
	if err := r.Get(ctx, req.NamespacedName, &instance); err != nil {
		if !k8serrors.IsNotFound(err) {
			log.Error(err, "unable to fetch OpenTelemetryCollector")
		}

		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	params := apolloportal.Params{
		//Config:   r.config,
		Client:   r.Client,
		Instance: instance,
		Log:      log,
		Scheme:   r.scheme,
		Recorder: r.recorder,
	}

	if err := r.RunTasks(ctx, params); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// RunTasks runs all the tasks associated with this reconciler.
func (r *ApolloPortalReconciler) RunTasks(ctx context.Context, params apolloportal.Params) error {
	r.muTasks.RLock()
	defer r.muTasks.RUnlock()
	for _, task := range r.tasks {
		if err := task.Do(ctx, params); err != nil {
			// If we get an error that occurs because a pod is being terminated, then exit this loop
			if k8serrors.IsForbidden(err) && k8serrors.HasStatusCause(err, corev1.NamespaceTerminatingCause) {
				r.log.V(2).Info("Exiting reconcile loop because namespace is being terminated", "namespace", params.Instance.Namespace)
				return nil
			}
			r.log.Error(err, fmt.Sprintf("failed to reconcile %s", task.Name))
			if task.BailOnError {
				return err
			}
		}
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApolloPortalReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apolloiov1alpha1.ApolloPortal{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&corev1.Service{}).
		Owns(&appv1.Deployment{}).
		Owns(&appv1.DaemonSet{}).
		Owns(&appv1.StatefulSet{}).
		Complete(r)
}
