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
	"context"
	"errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ApolloEnvironmentReconciler reconciles a ApolloEnvironment object
type ApolloEnvironmentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	//K8sClient *k8sClient.K8sClient
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

	logger.Info("进入Reconcile", "ApolloEnvironment", "刚进入")
	// TODO(user): your logic here
	apolloEnvironmentInstance := &apolloiov1alpha1.ApolloEnvironment{}
	err := r.Client.Get(ctx, req.NamespacedName, apolloEnvironmentInstance)
	if err != nil {
		return reconcile.Result{}, errors.New(err.Error() + "get ApolloEnvironment error")
	}
	// 具体逻辑
	//ok := r.ReconcileWork(ctx, apolloEnvironmentInstance)
	//if !ok {
	//	return reconcile.Result{RequeueAfter: time.Second * 5}, nil
	//} else {
	//	return reconcile.Result{}, nil
	//}
	return ctrl.Result{}, nil
}

//func (r *ApolloEnvironmentReconciler) ReconcileWork(ctx context.Context, instance *apolloiov1alpha1.ApolloEnvironment) bool {
//	logger := log.FromContext(ctx)
//	logger.Info("进入ReconcileWork", "ApolloEnvironment", instance.Name)
//
//	deployment := &appv1.Deployment{
//		ObjectMeta: metav1.ObjectMeta{
//			Name: "nginx",
//			Labels: map[string]string{
//				"app": "nginx",
//				"env": "dev",
//			},
//		},
//		Spec: appv1.DeploymentSpec{
//			Replicas: &instance.Spec.AdminServiceCount,
//			Selector: &metav1.LabelSelector{
//				MatchLabels: map[string]string{
//					"app": "nginx",
//					"env": "dev",
//				},
//			},
//			Template: corev1.PodTemplateSpec{
//				ObjectMeta: metav1.ObjectMeta{
//					Name: "nginx",
//					Labels: map[string]string{
//						"app": "nginx",
//						"env": "dev",
//					},
//				},
//				Spec: corev1.PodSpec{
//					Containers: []corev1.Container{
//						{
//							Name:  "nginx",
//							Image: "nginx:1.16.1",
//							Ports: []corev1.ContainerPort{
//								{
//									Name:          "http",
//									Protocol:      corev1.ProtocolTCP,
//									ContainerPort: 80,
//								},
//							},
//						},
//					},
//				},
//			},
//		},
//	}
//
//	deploymentList, err := r.K8sClient.KubeClient.AppsV1().Deployments(instance.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
//	fmt.Println(err, deploymentList)
//	if err != nil {
//		logger.Info("Deployments创建失败", "ApolloEnvironment", instance.Name)
//		return false
//	}
//	logger.Info("Deployments创建成功", "ApolloEnvironment", instance.Name)
//	return true
//}

// SetupWithManager sets up the controller with the Manager.
func (r *ApolloEnvironmentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apolloiov1alpha1.ApolloEnvironment{}).
		Complete(r)
}
