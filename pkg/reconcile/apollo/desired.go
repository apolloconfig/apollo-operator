package apollo

import (
	"apollo.io/apollo-operator/pkg/reconcile/models"
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DesiredConfigMaps 构建configmap对象
func (o ApolloAllInOne) DesiredConfigMaps(ctx context.Context, instance client.Object, params models.Params) []corev1.ConfigMap {

	return []corev1.ConfigMap{}
}

// DesiredEndpoints 构建endpoints对象
func (o ApolloAllInOne) DesiredEndpoints(ctx context.Context, instance client.Object, params models.Params) []corev1.Endpoints {
	return []corev1.Endpoints{}
}

// DesiredServices 构建service对象
func (o ApolloAllInOne) DesiredServices(ctx context.Context, instance client.Object, params models.Params) []corev1.Service {

	return []corev1.Service{}
}

// DesiredDeployments 构建deployment对象
func (o ApolloAllInOne) DesiredDeployments(ctx context.Context, instance client.Object, params models.Params) []appsv1.Deployment {

	return []appsv1.Deployment{}
}

// DesiredIngresses 构建ingress对象
func (o ApolloAllInOne) DesiredIngresses(ctx context.Context, instance client.Object, params models.Params) []networkingv1.Ingress {

	return []networkingv1.Ingress{}
}
