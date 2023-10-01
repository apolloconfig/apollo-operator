package reconcile

import (
	"apolloconfig.com/apollo-operator/pkg/reconcile/apollo"
	"apolloconfig.com/apollo-operator/pkg/reconcile/apolloenvironment"
	"apolloconfig.com/apollo-operator/pkg/reconcile/apolloportal"
	"apolloconfig.com/apollo-operator/pkg/reconcile/models"
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ApolloObject interface {
	// configmap
	DesiredConfigMaps(ctx context.Context, instance client.Object, params models.Params) []corev1.ConfigMap                              // 构建configmap对象
	ExpectedConfigMaps(ctx context.Context, instance client.Object, params models.Params, expected []corev1.ConfigMap, retry bool) error // 创建或更新configmap
	DeleteConfigMaps(ctx context.Context, instance client.Object, params models.Params, expected []corev1.ConfigMap) error               // 删除configmap

	// endpoints
	DesiredEndpoints(ctx context.Context, instance client.Object, params models.Params) []corev1.Endpoints                              // 构建endpoints对象
	ExpectedEndpoints(ctx context.Context, instance client.Object, params models.Params, expected []corev1.Endpoints, retry bool) error // 创建或更新endpoints
	DeleteEndpoints(ctx context.Context, instance client.Object, params models.Params, expected []corev1.Endpoints) error               // 删除endpoints

	// service
	DesiredServices(ctx context.Context, instance client.Object, params models.Params) []corev1.Service                  // 构建service对象
	ExpectedServices(ctx context.Context, instance client.Object, params models.Params, expected []corev1.Service) error // 创建或更新service
	DeleteServices(ctx context.Context, instance client.Object, params models.Params, expected []corev1.Service) error   // 删除service

	// deployment
	DesiredDeployments(ctx context.Context, instance client.Object, params models.Params) []appsv1.Deployment                  // 构建deployment对象
	ExpectedDeployments(ctx context.Context, instance client.Object, params models.Params, expected []appsv1.Deployment) error // 创建或更新deployment
	DeleteDeployments(ctx context.Context, instance client.Object, params models.Params, expected []appsv1.Deployment) error   // 删除deployment

	// ingerss
	DesiredIngresses(ctx context.Context, instance client.Object, params models.Params) []networkingv1.Ingress                  // 构建ingress对象
	ExpectedIngresses(ctx context.Context, instance client.Object, params models.Params, expected []networkingv1.Ingress) error // 创建或更新ingress
	DeleteIngresses(ctx context.Context, instance client.Object, params models.Params, expected []networkingv1.Ingress) error   // 删除ingress

}

var (
	apolloPortal      apolloportal.ApolloPortal
	apolloEnvironment apolloenvironment.ApolloEnvironment
	apolloAllInOne    apollo.ApolloAllInOne
)

func init() {
	apolloPortal = apolloportal.NewApolloPortal()
	apolloEnvironment = apolloenvironment.NewApolloEnvironment()
	apolloAllInOne = apollo.NewApolloAllInOne()
}

func ApolloPortal() apolloportal.ApolloPortal {
	return apolloPortal
}

func ApolloEnvironment() apolloenvironment.ApolloEnvironment {
	return apolloEnvironment
}

func ApolloAllInOne() apollo.ApolloAllInOne {
	return apolloAllInOne
}
