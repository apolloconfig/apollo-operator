package reconcile

import (
	"apollo.io/apollo-operator/pkg/reconcile/models"
	"context"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete

// Secret reconciles the secret(s) required for the instance in the current context.
func Secret(ctx context.Context, instance client.Object, params models.Params) error {

	return nil
}
