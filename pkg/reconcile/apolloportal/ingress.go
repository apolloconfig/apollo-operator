package apolloportal

import (
	"context"
)

// Ingresses reconciles the ingress(s) required for the instance in the current context.
func Ingresses(ctx context.Context, params Params) error {
	//isSupportedMode := true
	//if params.Instance.Spec.Mode == v1alpha1.ModeSidecar {
	//	params.Log.V(3).Info("ingress settings are not supported in sidecar mode")
	//	isSupportedMode = false
	//}
	//
	//nns := types.NamespacedName{Namespace: params.Instance.Namespace, Name: params.Instance.Name}
	//err := params.Client.Get(ctx, nns, &corev1.Service{}) // NOTE: check if service exists.
	//serviceExists := err != nil
	//
	//var desired []networkingv1.Ingress
	//if isSupportedMode && serviceExists {
	//	if d := desiredIngresses(ctx, params); d != nil {
	//		desired = append(desired, *d)
	//	}
	//}
	//
	//// first, handle the create/update parts
	//if err := expectedIngresses(ctx, params, desired); err != nil {
	//	return fmt.Errorf("failed to reconcile the expected ingresses: %w", err)
	//}
	//
	//// then, delete the extra objects
	//if err := deleteIngresses(ctx, params, desired); err != nil {
	//	return fmt.Errorf("failed to reconcile the ingresses to be deleted: %w", err)
	//}

	return nil
}
