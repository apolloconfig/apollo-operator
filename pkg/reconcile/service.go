package reconcile

import (
	"apollo.io/apollo-operator/pkg/reconcile/models"
	"context"
	"fmt"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete

// Services reconciles the service(s) required for the instance in the current context.
func Services(ctx context.Context, instance client.Object, params models.Params) error {
	var obj ApolloObject

	// TODO switch 修改一下
	if instance.GetObjectKind().GroupVersionKind().Kind == "ApolloPortal" {
		obj = ApolloPortal()
	}

	desired := obj.DesiredServices(ctx, instance, params)

	// first, handle the create/update parts
	if err := obj.ExpectedServices(ctx, instance, params, desired); err != nil {
		return fmt.Errorf("failed to reconcile the expected services: %w", err)
	}

	// then, delete the extra objects
	if err := obj.DeleteServices(ctx, instance, params, desired); err != nil {
		return fmt.Errorf("failed to reconcile the services to be deleted: %w", err)
	}

	return nil
}
