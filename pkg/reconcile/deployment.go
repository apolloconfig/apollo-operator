package reconcile

import (
	"apollo.io/apollo-operator/pkg/reconcile/models"
	"context"
	"fmt"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// +kubebuilder:rbac:groups="apps",resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Deployments reconciles the deployment(s) required for the instance in the current context.
func Deployments(ctx context.Context, instance client.Object, params models.Params) error {
	var obj ApolloObject

	// TODO switch 修改一下
	kind := instance.GetObjectKind().GroupVersionKind().Kind
	switch kind {
	case "ApolloPortal":
		obj = ApolloPortal()
	case "ApolloEnvironment":
		obj = ApolloEnvironment()
	case "Apollo":
		obj = ApolloAllInOne()
	}
	desired := obj.DesiredDeployments(ctx, instance, params)

	if err := obj.ExpectedDeployments(ctx, instance, params, desired); err != nil {
		return fmt.Errorf("failed to reconcile the expected deployments: %w", err)
	}

	// then, delete the extra objects
	if err := obj.DeleteDeployments(ctx, instance, params, desired); err != nil {
		return fmt.Errorf("failed to reconcile the deployments to be deleted: %w", err)
	}

	return nil
}
