package reconcile

import (
	"apolloconfig.com/apollo-operator/pkg/reconcile/models"
	"context"
	"fmt"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// +kubebuilder:rbac:groups="",resources=statefulsets,verbs=get;list;watch;create;update;patch;delete

// StatefulSet reconciles the statefulset required for the instance in the current context.
func StatefulSet(ctx context.Context, instance client.Object, params models.Params) error {

	obj := ApolloAllInOne()

	desired := obj.DesiredStatefulSets(ctx, instance, params)

	if err := obj.ExpectedStatefulSets(ctx, instance, params, desired); err != nil {
		return fmt.Errorf("failed to reconcile the expected statefulset: %w", err)
	}

	// then, delete the extra objects
	if err := obj.DeleteStatefulSets(ctx, instance, params, desired); err != nil {
		return fmt.Errorf("failed to reconcile the statefulsets to be deleted: %w", err)
	}

	return nil
}
