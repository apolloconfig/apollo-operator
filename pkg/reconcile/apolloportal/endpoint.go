package apolloportal

import (
	"context"
)

// +kubebuilder:rbac:groups="apps",resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Deployments reconciles the deployment(s) required for the instance in the current context.
func Endpoints(ctx context.Context, params Params) error {
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

	return nil
}
