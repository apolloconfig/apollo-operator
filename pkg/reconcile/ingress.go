package reconcile

import (
	"apollo.io/apollo-operator/pkg/reconcile/models"
	"context"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// TODO 注释验证一下对不对
// +kubebuilder:rbac:groups="",resources=ingress,verbs=get;list;watch;create;update;patch;delete

// Ingresses reconciles the ingress(s) required for the instance in the current context.
func Ingresses(ctx context.Context, instance client.Object, params models.Params) error {

	return nil
}
