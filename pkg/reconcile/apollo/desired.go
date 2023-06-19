package apollo

import (
	"apollo.io/apollo-operator/pkg/reconcile/models"
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (o ApolloAllInOne) DesiredConfigMaps(ctx context.Context, instance client.Object, params models.Params) []corev1.ConfigMap {

	return []corev1.ConfigMap{}
}

// 构建endpoints对象
func (o ApolloAllInOne) DesiredEndpoints(ctx context.Context, instance client.Object, params models.Params) []corev1.Endpoints {
	return []corev1.Endpoints{}
}

// 构建service对象
func (o ApolloAllInOne) DesiredServices(ctx context.Context, instance client.Object, params models.Params) []corev1.Service {

	return []corev1.Service{}
}

// 构建deployment对象
func (o ApolloAllInOne) DesiredDeployments(ctx context.Context, instance client.Object, params models.Params) []appsv1.Deployment {

	return []appsv1.Deployment{}
}
