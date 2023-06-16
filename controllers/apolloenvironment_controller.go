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
	"apollo.io/apollo-operator/pkg/reconcile"
	"context"
	"errors"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sync"
)

// ApolloEnvironmentReconciler reconciles a ApolloEnvironment object
type ApolloEnvironmentReconciler struct {
	client.Client
	recorder record.EventRecorder
	scheme   *runtime.Scheme
	log      logr.Logger
	//config   config.Config

	tasks   []Task
	muTasks sync.RWMutex
}

// NewApolloPortalReconciler creates a new reconciler for ApolloPortal objects.
func NewApolloEnvironmentReconciler(p ReconcilerParams) *ApolloEnvironmentReconciler {
	r := &ApolloEnvironmentReconciler{
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
				reconcile.Endpoints,
				"endpoints",
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
				reconcile.Ingresses,
				"ingresses",
				true,
			},
			{
				reconcile.Self,
				"apolloenvironment",
				true,
			},
		}
	}
	return r
}

//+kubebuilder:rbac:groups=apollo.io,resources=apolloenvironments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apollo.io,resources=apolloenvironments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apollo.io,resources=apolloenvironments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ApolloEnvironment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *ApolloEnvironmentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("进入ApolloEnvironmentReconciler Reconcile", "ApolloEnvironment", "刚进入")
	// TODO(user): your logic here
	apolloEnvironmentInstance := &apolloiov1alpha1.ApolloEnvironment{}
	err := r.Client.Get(ctx, req.NamespacedName, apolloEnvironmentInstance)
	if err != nil {
		return ctrl.Result{}, errors.New(err.Error() + "get ApolloEnvironment error")
	}
	// TODO 具体逻辑
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApolloEnvironmentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apolloiov1alpha1.ApolloEnvironment{}).
		Complete(r)
}
