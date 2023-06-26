package reconcile

import (
	"apollo.io/apollo-operator/pkg/reconcile/models"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// TODO 注释验证一下对不对
// +kubebuilder:rbac:groups="",resources=ingress,verbs=get;list;watch;create;update;patch;delete

// Ingresses reconciles the ingress(s) required for the instance in the current context.
func Ingresses(ctx context.Context, instance client.Object, params models.Params) error {

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

	nns := types.NamespacedName{Namespace: instance.GetNamespace(), Name: instance.GetName()}
	err := params.Client.Get(ctx, nns, &corev1.Service{}) // NOTE: check if service exists.
	serviceExists := err != nil

	var desired []networkingv1.Ingress
	if serviceExists {
		desired = obj.DesiredIngresses(ctx, instance, params)
	}

	// first, handle the create/update parts
	if err := obj.ExpectedIngresses(ctx, instance, params, desired); err != nil {
		return fmt.Errorf("failed to reconcile the expected ingresses: %w", err)
	}

	// then, delete the extra objects
	if err := obj.DeleteIngresses(ctx, instance, params, desired); err != nil {
		return fmt.Errorf("failed to reconcile the ingresses to be deleted: %w", err)
	}

	return nil
}
