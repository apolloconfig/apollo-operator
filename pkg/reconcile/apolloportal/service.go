package apolloportal

import (
	"context"
)

// headless label is to differentiate the headless service from the clusterIP service.
const (
	headlessLabel  = "operator.opentelemetry.io/collector-headless-service"
	headlessExists = "Exists"
)

var portaldbCh = make(chan string, 1)

// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete

// Services reconciles the service(s) required for the instance in the current context.
func Services(ctx context.Context, params Params) error {
	//desired := []corev1.Service{}
	//if params.Instance.Spec.Mode != v1alpha1.ModeSidecar {
	//	type builder func(context.Context, Params) *corev1.Service
	//	for _, builder := range []builder{desiredService, headless, monitoringService} {
	//		svc := builder(ctx, params)
	//		// add only the non-nil to the list
	//		if svc != nil {
	//			desired = append(desired, *svc)
	//		}
	//	}
	//}
	//
	//if params.Instance.Spec.TargetAllocator.Enabled {
	//	desired = append(desired, desiredTAService(params))
	//}
	//
	//// first, handle the create/update parts
	//if err := expectedServices(ctx, params, desired); err != nil {
	//	return fmt.Errorf("failed to reconcile the expected services: %w", err)
	//}
	//
	//// then, delete the extra objects
	//if err := deleteServices(ctx, params, desired); err != nil {
	//	return fmt.Errorf("failed to reconcile the services to be deleted: %w", err)
	//}

	return nil
}
