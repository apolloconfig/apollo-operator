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
	"apolloconfig.com/apollo-operator/pkg/reconcile"
	"apolloconfig.com/apollo-operator/pkg/reconcile/models"
	"context"
	"fmt"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	apolloiov1alpha1 "apolloconfig.com/apollo-operator/api/v1alpha1"
)

// ApolloReconciler reconciles a Apollo object
type ApolloReconciler struct {
	client.Client
	recorder record.EventRecorder
	scheme   *runtime.Scheme
	log      logr.Logger
	//config   config.Config

	tasks   []Task
	muTasks sync.RWMutex
}

// NewApolloAllInOneReconciler creates a new reconciler for ApolloAllInOne objects.
func NewApolloAllInOneReconciler(p ReconcilerParams) *ApolloReconciler {
	r := &ApolloReconciler{
		Client:   p.Client,
		log:      p.Log,
		scheme:   p.Scheme,
		tasks:    p.Tasks,
		recorder: p.Recorder,
	}
	if len(r.tasks) == 0 {
		r.tasks = []Task{
			{
				reconcile.ConfigMaps,
				"configmaps",
				true,
			},
			{
				reconcile.ServiceAccounts,
				"serviceaccounts",
				true,
			},
			{
				reconcile.Services,
				"services",
				true,
			},
			{
				reconcile.Deployments,
				"deployments",
				true,
			},
			{
				reconcile.StatefulSet,
				"statefulsets",
				true,
			},
			{
				reconcile.Ingresses,
				"ingresses",
				true,
			},
			{
				reconcile.Self,
				"apolloallinone",
				true,
			},
		}
	}
	return r
}

//+kubebuilder:rbac:groups=apolloconfig.com,resources=apolloes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apolloconfig.com,resources=apolloes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apolloconfig.com,resources=apolloes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Apollo object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *ApolloReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// _ = log.FromContext(ctx)
	log := r.log.WithValues("ApolloAllInOne", req.NamespacedName)
	log.Info("ApolloAllInOneReconciler Reconcile")
	// TODO(user): your logic here
	var instance apolloiov1alpha1.Apollo
	if err := r.Get(ctx, req.NamespacedName, &instance); err != nil {
		if !k8serrors.IsNotFound(err) {
			log.Error(err, "unable to fetch ApolloAllInOne")
		}

		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	params := models.Params{
		//Config:   r.config,
		Client:   r.Client,
		Log:      log,
		Scheme:   r.scheme,
		Recorder: r.recorder,
	}
	// TODO 为 instance 增加默认值

	if err := r.RunTasks(ctx, &instance, params); err != nil {
		//return ctrl.Result{}, err
		return ctrl.Result{RequeueAfter: time.Second * 5}, err
	}

	return ctrl.Result{RequeueAfter: time.Second * 3}, nil

}

// RunTasks runs all the tasks associated with this reconciler.
func (r *ApolloReconciler) RunTasks(ctx context.Context, instance client.Object, params models.Params) error {
	r.muTasks.RLock()
	defer r.muTasks.RUnlock()
	for _, task := range r.tasks {
		if err := task.Do(ctx, instance, params); err != nil {
			// If we get an error that occurs because a pod is being terminated, then exit this loop
			if k8serrors.IsForbidden(err) && k8serrors.HasStatusCause(err, corev1.NamespaceTerminatingCause) {
				r.log.V(2).Info("Exiting reconcile loop because namespace is being terminated", "namespace", instance.GetNamespace())
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
func (r *ApolloReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apolloiov1alpha1.Apollo{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&corev1.Service{}).
		Owns(&appsv1.Deployment{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&networkingv1.Ingress{}).
		Complete(r)
}
