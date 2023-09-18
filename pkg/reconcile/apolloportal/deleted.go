package apolloportal

import (
	"apolloconfig.com/apollo-operator/pkg/reconcile/models"
	"apolloconfig.com/apollo-operator/pkg/utils/naming"
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (o ApolloPortal) DeleteConfigMaps(ctx context.Context, instance client.Object, params models.Params, expected []corev1.ConfigMap) error {
	opts := []client.ListOption{
		client.InNamespace(instance.GetNamespace()),
		client.MatchingLabels(map[string]string{
			"app.kubernetes.io/instance":   naming.Truncate("%s.%s", 63, instance.GetNamespace(), instance.GetName()),
			"app.kubernetes.io/managed-by": "apollo-operator",
		}),
	}
	list := &corev1.ConfigMapList{}
	if err := params.Client.List(ctx, list, opts...); err != nil {
		return fmt.Errorf("failed to list configmap : %w", err)
	}

	// Delete parts that are not expected
	for i := range list.Items {
		existing := list.Items[i]
		del := true
		for _, keep := range expected {
			if keep.Name == existing.Name && keep.Namespace == existing.Namespace {
				del = false
				break
			}
		}

		if del {
			if err := params.Client.Delete(ctx, &existing); err != nil {
				return fmt.Errorf("failed to delete: %w", err)
			}
			params.Log.V(2).Info("deleted", "configmap.name", existing.Name, "configmap.namespace", existing.Namespace)
		}
	}

	return nil
}

func (o ApolloPortal) DeleteEndpoints(ctx context.Context, instance client.Object, params models.Params, expected []corev1.Endpoints) error {
	opts := []client.ListOption{
		client.InNamespace(instance.GetNamespace()),
		client.MatchingLabels(map[string]string{
			"app.kubernetes.io/instance":   naming.Truncate("%s.%s", 63, instance.GetNamespace(), instance.GetName()),
			"app.kubernetes.io/managed-by": "apollo-operator",
		}),
	}
	list := &corev1.EndpointsList{}
	if err := params.Client.List(ctx, list, opts...); err != nil {
		return fmt.Errorf("failed to list endpoints : %w", err)
	}

	// Delete parts that are not expected
	for i := range list.Items {
		existing := list.Items[i]
		del := true
		for _, keep := range expected {
			if keep.Name == existing.Name && keep.Namespace == existing.Namespace {
				del = false
				break
			}
		}

		if del {
			if err := params.Client.Delete(ctx, &existing); err != nil {
				return fmt.Errorf("failed to delete: %w", err)
			}
			params.Log.V(2).Info("deleted", "endpoints.name", existing.Name, "endpoints.namespace", existing.Namespace)
		}
	}

	return nil
}

func (o ApolloPortal) DeleteServices(ctx context.Context, instance client.Object, params models.Params, expected []corev1.Service) error {
	opts := []client.ListOption{
		client.InNamespace(instance.GetNamespace()),
		client.MatchingLabels(map[string]string{
			"app.kubernetes.io/instance":   naming.Truncate("%s.%s", 63, instance.GetNamespace(), instance.GetName()),
			"app.kubernetes.io/managed-by": "apollo-operator",
		}),
	}
	list := &corev1.ServiceList{}
	if err := params.Client.List(ctx, list, opts...); err != nil {
		return fmt.Errorf("failed to list service: %w", err)
	}

	for i := range list.Items {
		existing := list.Items[i]
		del := true
		for _, keep := range expected {
			if keep.Name == existing.Name && keep.Namespace == existing.Namespace {
				del = false
				break
			}
		}

		if del {
			if err := params.Client.Delete(ctx, &existing); err != nil {
				return fmt.Errorf("failed to delete: %w", err)
			}
			params.Log.V(2).Info("deleted", "service.name", existing.Name, "service.namespace", existing.Namespace)
		}
	}

	return nil
}

func (o ApolloPortal) DeleteDeployments(ctx context.Context, instance client.Object, params models.Params, expected []appsv1.Deployment) error {
	opts := []client.ListOption{
		client.InNamespace(instance.GetNamespace()),
		client.MatchingLabels(map[string]string{
			"app.kubernetes.io/instance":   naming.Truncate("%s.%s", 63, instance.GetNamespace(), instance.GetName()),
			"app.kubernetes.io/managed-by": "apollo-operator",
		}),
	}
	list := &appsv1.DeploymentList{}
	if err := params.Client.List(ctx, list, opts...); err != nil {
		return fmt.Errorf("failed to list deployment: %w", err)
	}

	for i := range list.Items {
		existing := list.Items[i]
		del := true
		for _, keep := range expected {
			if keep.Name == existing.Name && keep.Namespace == existing.Namespace {
				del = false
				break
			}
		}

		if del {
			if err := params.Client.Delete(ctx, &existing); err != nil {
				return fmt.Errorf("failed to delete: %w", err)
			}
			params.Log.V(2).Info("deleted", "deployment.name", existing.Name, "deployment.namespace", existing.Namespace)
		}
	}

	return nil
}

func (o ApolloPortal) DeleteIngresses(ctx context.Context, instance client.Object, params models.Params, expected []networkingv1.Ingress) error {
	opts := []client.ListOption{
		client.InNamespace(instance.GetNamespace()),
		client.MatchingLabels(map[string]string{
			"app.kubernetes.io/instance":   naming.Truncate("%s.%s", 63, instance.GetNamespace(), instance.GetName()),
			"app.kubernetes.io/managed-by": "apollo-operator",
		}),
	}
	list := &networkingv1.IngressList{}
	if err := params.Client.List(ctx, list, opts...); err != nil {
		return fmt.Errorf("failed to list ingress : %w", err)
	}

	// Delete parts that are not expected
	for i := range list.Items {
		existing := list.Items[i]
		del := true
		for _, keep := range expected {
			if keep.Name == existing.Name && keep.Namespace == existing.Namespace {
				del = false
				break
			}
		}

		if del {
			if err := params.Client.Delete(ctx, &existing); err != nil {
				return fmt.Errorf("failed to delete: %w", err)
			}
			params.Log.V(2).Info("deleted", "ingress.name", existing.Name, "ingress.namespace", existing.Namespace)
		}
	}

	return nil
}
