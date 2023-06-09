package reconcile

import (
	apolloiov1alpha1 "apollo.io/apollo-operator/api/v1alpha1"
	"context"
	"errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// +kubebuilder:rbac:groups="apps",resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Deployments reconciles the deployment(s) required for the instance in the current context.
//func Endpoints(ctx context.Context, params Params) error {
//desired := []appsv1.Deployment{}
//if params.Instance.Spec.Mode == "deployment" {
//	desired = append(desired, collector.Deployment(params.Config, params.Log, params.Instance))
//}
//
//if params.Instance.Spec.TargetAllocator.Enabled {
//	desired = append(desired, targetallocator.Deployment(params.Config, params.Log, params.Instance))
//}
//
//// first, handle the create/update parts
//if err := expectedDeployments(ctx, params, desired); err != nil {
//	return fmt.Errorf("failed to reconcile the expected deployments: %w", err)
//}
//
//// then, delete the extra objects
//if err := deleteDeployments(ctx, params, desired); err != nil {
//	return fmt.Errorf("failed to reconcile the deployments to be deleted: %w", err)
//}

//	return nil
//}

type ParamsAny struct {
	Instance client.Object
}

type EndpointsInterface interface {

	// 定义策略执行方法
	ApplyEndpoints(ctx context.Context, params ParamsAny) error
}

type EndpointsApolloPortal struct {
}

func (e EndpointsApolloPortal) ApplyEndpoints(ctx context.Context, params ParamsAny) error {
	return nil
}

type EndpointsApollo struct {
}

func (e EndpointsApollo) ApplyEndpoints(ctx context.Context, params ParamsAny) error {
	_ = params.Instance.(*apolloiov1alpha1.Apollo)
	return nil
}

// NewEndpoints test
func NewApply(kind string) (EndpointsInterface, error) {
	if kind == "apollo" {
		return EndpointsApollo{}, nil
	} else if kind == "apolloportal" {
		return EndpointsApolloPortal{}, nil
	}
	return nil, errors.New("无对象")
}

func ApplyEndpoints(ctx context.Context, params ParamsAny) error {
	kind := params.Instance.GetObjectKind()
	if kind.GroupVersionKind().Kind == "apollo" {

	}
	return nil
}
