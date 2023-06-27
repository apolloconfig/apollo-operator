package reconcile

import (
	"apollo.io/apollo-operator/pkg/reconcile/models"
	"context"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// +kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete

// ServiceAccounts reconciles the service account(s) required for the instance in the current context.
func ServiceAccounts(ctx context.Context, instance client.Object, params models.Params) error {

	return nil
}
