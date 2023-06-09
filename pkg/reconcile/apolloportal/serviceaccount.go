package apolloportal

import (
	"context"
)

// +kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete

// ServiceAccounts reconciles the service account(s) required for the instance in the current context.
func ServiceAccounts(ctx context.Context, params Params) error {
	//desired := desiredServiceAccounts(params)
	//
	//// first, handle the create/update parts
	//if err := expectedServiceAccounts(ctx, params, desired); err != nil {
	//	return fmt.Errorf("failed to reconcile the expected service accounts: %w", err)
	//}
	//
	//// then, delete the extra objects
	//if err := deleteServiceAccounts(ctx, params, desired); err != nil {
	//	return fmt.Errorf("failed to reconcile the service accounts to be deleted: %w", err)
	//}

	return nil
}
