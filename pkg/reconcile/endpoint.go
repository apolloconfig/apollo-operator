package reconcile

import (
	"apollo.io/apollo-operator/pkg/reconcile/models"
	"context"
	"fmt"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// +kubebuilder:rbac:groups="",resources=endpoints,verbs=get;list;watch;create;update;patch;delete

// Endpoints reconciles the endpoint(s) required for the instance in the current context.
func Endpoints(ctx context.Context, instance client.Object, params models.Params) error {

	var obj ApolloObject

	kind := instance.GetObjectKind().GroupVersionKind().Kind
	switch kind {
	case "ApolloPortal":
		obj = ApolloPortal()
	case "ApolloEnvironment":
		obj = ApolloEnvironment()
	case "Apollo":
		obj = ApolloAllInOne()
	}

	desired := obj.DesiredEndpoints(ctx, instance, params)

	// TODO 可以优化为先获取create、upodate、delete列表，然后再统一apply

	// first, handle the create/update parts
	if err := obj.ExpectedEndpoints(ctx, instance, params, desired, true); err != nil {
		return fmt.Errorf("failed to reconcile the expected Endpoints: %w", err)
	}

	// then, delete the extra objects
	if err := obj.DeleteEndpoints(ctx, instance, params, desired); err != nil {
		return fmt.Errorf("failed to reconcile the Endpoints to be deleted: %w", err)
	}

	return nil
}
