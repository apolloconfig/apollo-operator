package apolloportal

import (
	"context"
)

// Self updates this instance's self data. This should be the last item in the reconciliation, as it causes changes
// making params.Instance obsolete. Default values should be set in the Defaulter webhook, this should only be used
// for the Status, which can't be set by the defaulter.
func Self(ctx context.Context, params Params) error {
	//changed := params.Instance
	//
	//// this field is only changed for new instances: on existing instances this
	//// field is reconciled when the operator is first started, i.e. during
	//// the upgrade mechanism
	//if params.Instance.Status.Version == "" {
	//	// a version is not set, otherwise let the upgrade mechanism take care of it!
	//	changed.Status.Version = version.OpenTelemetryCollector()
	//}
	//
	//if err := updateScaleSubResourceStatus(ctx, params.Client, &changed); err != nil {
	//	return fmt.Errorf("failed to update the scale subresource status for the OpenTelemetry CR: %w", err)
	//}
	//
	//statusPatch := client.MergeFrom(&params.Instance)
	//if err := params.Client.Status().Patch(ctx, &changed, statusPatch); err != nil {
	//	return fmt.Errorf("failed to apply status changes to the OpenTelemetry CR: %w", err)
	//}

	return nil
}
