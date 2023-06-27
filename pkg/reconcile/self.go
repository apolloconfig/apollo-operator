package reconcile

import (
	"apollo.io/apollo-operator/pkg/reconcile/models"
	"context"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Self updates this instance's self data. This should be the last item in the reconciliation, as it causes changes
// making params.Instance obsolete. Default values should be set in the Defaulter webhook, this should only be used
// for the Status, which can't be set by the defaulter.
func Self(ctx context.Context, instance client.Object, params models.Params) error {

	return nil
}
