package apollo

import (
	"apolloconfig.com/apollo-operator/pkg/reconcile/models"
	"apolloconfig.com/apollo-operator/pkg/utils"
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ExpectedConfigMaps 创建或更新configmap
func (o ApolloAllInOne) ExpectedConfigMaps(ctx context.Context, instance client.Object, params models.Params, expected []corev1.ConfigMap, retry bool) error {
	for _, obj := range expected {
		desired := obj

		// After establishing the OwnerReference, deleting the apollo-portal resource will also delete the configmap
		if err := controllerutil.SetControllerReference(instance, &desired, params.Scheme); err != nil {
			return fmt.Errorf("failed to set controller reference: %w", err)
		}

		existing := &corev1.ConfigMap{}
		namespaceName := types.NamespacedName{Namespace: desired.Namespace, Name: desired.Name}
		getErr := params.Client.Get(ctx, namespaceName, existing)
		if getErr != nil && k8serrors.IsNotFound(getErr) {
			// 不存在则直接创建desired的资源
			if createErr := params.Client.Create(ctx, &desired); createErr != nil {
				if k8serrors.IsAlreadyExists(createErr) && retry {
					// let's try again? we probably had multiple updates at one, and now it exists already
					if err := o.ExpectedConfigMaps(ctx, instance, params, expected, false); err != nil {
						// somethin else happened now...
						return err
					}
					// we succeeded in the retry, exit this attempt
					return nil
				}
				return fmt.Errorf("failed to create: %w", createErr)
			}
			params.Log.V(2).Info("created", "configmap.name", desired.Name, "configmap.namespace", desired.Namespace)
			// Successfully created and entered the next cycle
			continue
		} else if getErr != nil {
			return fmt.Errorf("failed to get: %w", getErr)
		}

		// it exists already, merge the two if the end result isn't identical to the existing one
		updated := existing.DeepCopy()
		utils.InitObjectMeta(updated)

		updated.SetAnnotations(desired.GetAnnotations())
		updated.SetLabels(desired.GetLabels())
		updated.SetOwnerReferences(desired.GetOwnerReferences())

		updated.Data = desired.Data
		updated.BinaryData = desired.BinaryData

		// Modify the old configmap to the new configmap
		patch := client.MergeFrom(existing)
		if err := params.Client.Patch(ctx, updated, patch); err != nil {
			return fmt.Errorf("failed to apply changes: %w", err)
		}

		if !reflect.DeepEqual(desired.Data, existing.Data) {
			params.Recorder.Event(updated, "Normal", "ConfigUpdate ", fmt.Sprintf("ApolloPortal Config changed - %s/%s", desired.Namespace, desired.Name))
		}

		params.Log.V(2).Info("applied", "configmap.name", desired.Name, "configmap.namespace", desired.Namespace)
	}

	return nil
}

// ExpectedEndpoints Create or update endpoints
func (o ApolloAllInOne) ExpectedEndpoints(ctx context.Context, instance client.Object, params models.Params, expected []corev1.Endpoints, retry bool) error {
	for _, obj := range expected {
		desired := obj

		// After establishing the OwnerReference, deleting the apollo-portal resource will also delete the endpoints
		if err := controllerutil.SetControllerReference(instance, &desired, params.Scheme); err != nil {
			return fmt.Errorf("failed to set controller reference: %w", err)
		}

		existing := &corev1.Endpoints{}
		namespaceName := types.NamespacedName{Namespace: desired.Namespace, Name: desired.Name}
		getErr := params.Client.Get(ctx, namespaceName, existing)
		if getErr != nil && k8serrors.IsNotFound(getErr) {
			// If it does not exist, create the desired resource directly
			if createErr := params.Client.Create(ctx, &desired); createErr != nil {
				if k8serrors.IsAlreadyExists(createErr) && retry {
					// let's try again? we probably had multiple updates at one, and now it exists already
					if err := o.ExpectedEndpoints(ctx, instance, params, expected, false); err != nil {
						// somethin else happened now...
						return err
					}

					// we succeeded in the retry, exit this attempt
					return nil
				}
				return fmt.Errorf("failed to create: %w", createErr)
			}
			params.Log.V(2).Info("created", "endpoints.name", desired.Name, "endpoints.namespace", desired.Namespace)
			// Successfully created and entered the next cycle
			continue
		} else if getErr != nil {
			return fmt.Errorf("failed to get: %w", getErr)
		}

		// it exists already, merge the two if the end result isn't identical to the existing one
		updated := existing.DeepCopy()
		utils.InitObjectMeta(updated)

		if updated.Subsets == nil {
			updated.Subsets = []corev1.EndpointSubset{}
		}
		updated.SetAnnotations(desired.GetAnnotations())
		updated.SetLabels(desired.GetLabels())
		updated.SetOwnerReferences(desired.GetOwnerReferences())

		// Replace with desired
		for i, subset := range desired.Subsets {
			updated.Subsets[i] = subset
		}

		// Modify old Endpoints to new Endpoints
		patch := client.MergeFrom(existing)
		if err := params.Client.Patch(ctx, updated, patch); err != nil {
			return fmt.Errorf("failed to apply changes: %w", err)
		}

		if EndpointsChanged(&desired, existing) {
			params.Recorder.Event(updated, "Normal", "Endpoints Update ", fmt.Sprintf("ApolloPortal Endpoints changed - %s/%s", desired.Namespace, desired.Name))
		}

		params.Log.V(2).Info("applied", "endpoints.name", desired.Name, "endpoints.namespace", desired.Namespace)
	}

	return nil
}

func EndpointsChanged(desired *corev1.Endpoints, existing *corev1.Endpoints) bool {
	desSubsets := desired.Subsets
	for i, subset := range existing.Subsets {
		for j, address := range subset.Addresses {
			if address.IP != desSubsets[i].Addresses[j].IP {
				return true
			}
		}
		for j, port := range subset.Ports {
			if port.Port != desSubsets[i].Ports[j].Port {
				return true
			}
		}
	}
	return false
}

// ExpectedServices 创建或更新service
func (o ApolloAllInOne) ExpectedServices(ctx context.Context, instance client.Object, params models.Params, expected []corev1.Service) error {
	for _, obj := range expected {
		desired := obj

		if err := controllerutil.SetControllerReference(instance, &desired, params.Scheme); err != nil {
			return fmt.Errorf("failed to set controller reference: %w", err)
		}

		existing := &corev1.Service{}
		nns := types.NamespacedName{Namespace: desired.Namespace, Name: desired.Name}
		err := params.Client.Get(ctx, nns, existing)
		if err != nil && k8serrors.IsNotFound(err) {
			if clientErr := params.Client.Create(ctx, &desired); clientErr != nil {
				return fmt.Errorf("failed to create: %w", clientErr)
			}
			params.Log.V(2).Info("created", "service.name", desired.Name, "service.namespace", desired.Namespace)
			continue
		} else if err != nil {
			return fmt.Errorf("failed to get: %w", err)
		}

		// it exists already, merge the two if the end result isn't identical to the existing one
		updated := existing.DeepCopy()
		utils.InitObjectMeta(updated)
		//if updated.Annotations == nil {
		//	updated.Annotations = map[string]string{}
		//}
		//if updated.Labels == nil {
		//	updated.Labels = map[string]string{}
		//}
		//updated.ObjectMeta.OwnerReferences = desired.ObjectMeta.OwnerReferences
		updated.SetAnnotations(desired.GetAnnotations())
		updated.SetLabels(desired.GetLabels())
		updated.SetOwnerReferences(desired.GetOwnerReferences())

		updated.Spec.Type = desired.Spec.Type
		updated.Spec.Ports = desired.Spec.Ports
		updated.Spec.Selector = desired.Spec.Selector
		updated.Spec.SessionAffinity = desired.Spec.SessionAffinity

		patch := client.MergeFrom(existing)
		if err := params.Client.Patch(ctx, updated, patch); err != nil {
			return fmt.Errorf("failed to apply changes: %w", err)
		}

		params.Log.V(2).Info("applied", "service.name", desired.Name, "service.namespace", desired.Namespace)
	}

	return nil
}

// ExpectedDeployments Create or update deployment
func (o ApolloAllInOne) ExpectedDeployments(ctx context.Context, instance client.Object, params models.Params, expected []appsv1.Deployment) error {
	for _, obj := range expected {
		desired := obj

		if err := controllerutil.SetControllerReference(instance, &desired, params.Scheme); err != nil {
			return fmt.Errorf("failed to set controller reference: %w", err)
		}

		existing := &appsv1.Deployment{}
		nns := types.NamespacedName{Namespace: desired.Namespace, Name: desired.Name}
		err := params.Client.Get(ctx, nns, existing)
		if err != nil && k8serrors.IsNotFound(err) {
			if clientErr := params.Client.Create(ctx, &desired); clientErr != nil {
				return fmt.Errorf("failed to create: %w", clientErr)
			}
			params.Log.V(2).Info("created", "deployment.name", desired.Name, "deployment.namespace", desired.Namespace)
			continue
		} else if err != nil {
			return fmt.Errorf("failed to get: %w", err)
		}

		// Selector is an immutable field, if set, we cannot modify it otherwise we will have reconciliation error.
		if !apiequality.Semantic.DeepEqual(desired.Spec.Selector, existing.Spec.Selector) {
			params.Log.V(2).Info("Spec.Selector change detected, trying to delete, the new apollo-portal deployment will be created in the next reconcile cycle ", "deployment.name", existing.Name, "deployment.namespace", existing.Namespace)

			if err := params.Client.Delete(ctx, existing); err != nil {
				return fmt.Errorf("failed to delete deployment: %w", err)
			}
			continue
		}

		// it exists already, merge the two if the end result isn't identical to the existing one
		updated := existing.DeepCopy()
		utils.InitObjectMeta(updated)
		updated.SetAnnotations(desired.GetAnnotations())
		updated.SetLabels(desired.GetLabels())
		updated.SetOwnerReferences(desired.GetOwnerReferences())

		// Be sure to pay attention to the source of the slice in the spec, and it cannot be obtained by traversing the map, otherwise it will cause the pod to be recreated every time the tuning is performed
		updated.Spec = desired.Spec
		patch := client.MergeFrom(existing)
		if err := params.Client.Patch(ctx, updated, patch); err != nil {
			return fmt.Errorf("failed to apply changes: %w", err)
		}

		params.Log.V(2).Info("applied", "deployment.name", desired.Name, "deployment.namespace", desired.Namespace)
	}

	return nil
}

// ExpectedStatefulSets Create or update statefulset
func (o ApolloAllInOne) ExpectedStatefulSets(ctx context.Context, instance client.Object, params models.Params, expected []appsv1.StatefulSet) error {
	for _, obj := range expected {
		desired := obj

		if err := controllerutil.SetControllerReference(instance, &desired, params.Scheme); err != nil {
			return fmt.Errorf("failed to set controller reference: %w", err)
		}

		existing := &appsv1.StatefulSet{}
		nns := types.NamespacedName{Namespace: desired.Namespace, Name: desired.Name}
		err := params.Client.Get(ctx, nns, existing)
		if err != nil && k8serrors.IsNotFound(err) {
			if clientErr := params.Client.Create(ctx, &desired); clientErr != nil {
				return fmt.Errorf("failed to create: %w", clientErr)
			}
			params.Log.V(2).Info("created", "statefulset.name", desired.Name, "statefulset.namespace", desired.Namespace)
			continue
		} else if err != nil {
			return fmt.Errorf("failed to get: %w", err)
		}

		// Check for immutable fields. If set, we cannot modify the stateful set, otherwise we will face reconciliation error.
		if needsDeletion, fieldName := hasImmutableFieldChange(&desired, existing); needsDeletion {
			params.Log.V(2).Info("Immutable field change detected, trying to delete, the new collector statefulset will be created in the next reconcile cycle",
				"field", fieldName, "statefulset.name", existing.Name, "statefulset.namespace", existing.Namespace)

			if err := params.Client.Delete(ctx, existing); err != nil {
				return fmt.Errorf("failed to delete statefulset: %w", err)
			}
			continue
		}

		// it exists already, merge the two if the end result isn't identical to the existing one
		updated := existing.DeepCopy()
		utils.InitObjectMeta(updated)
		updated.SetAnnotations(desired.GetAnnotations())
		updated.SetLabels(desired.GetLabels())
		updated.SetOwnerReferences(desired.GetOwnerReferences())

		// Be sure to pay attention to the source of the slice in the spec, and it cannot be obtained by traversing the map, otherwise it will cause the pod to be recreated every time the tuning is performed
		updated.Spec = desired.Spec

		patch := client.MergeFrom(existing)
		if err := params.Client.Patch(ctx, updated, patch); err != nil {
			return fmt.Errorf("failed to apply changes: %w", err)
		}

		params.Log.V(2).Info("applied", "statefulset.name", desired.Name, "statefulset.namespace", desired.Namespace)
	}

	return nil
}

func hasImmutableFieldChange(desired, existing *appsv1.StatefulSet) (bool, string) {
	if !apiequality.Semantic.DeepEqual(desired.Spec.Selector, existing.Spec.Selector) {
		return true, "Spec.Selector"
	}

	if hasVolumeClaimsTemplatesChanged(desired, existing) {
		return true, "Spec.VolumeClaimTemplates"
	}

	return false, ""
}

// hasVolumeClaimsTemplatesChanged if volume claims template change has been detected.
// We need to do this manually due to some fields being automatically filled by the API server
// and these needs to be excluded from the comparison to prevent false positives.
func hasVolumeClaimsTemplatesChanged(desired, existing *appsv1.StatefulSet) bool {
	if len(desired.Spec.VolumeClaimTemplates) != len(existing.Spec.VolumeClaimTemplates) {
		return true
	}

	for i := range desired.Spec.VolumeClaimTemplates {
		// VolumeMode is automatically set by the API server, so if it is not set in the CR, assume it's the same as the existing one.
		if desired.Spec.VolumeClaimTemplates[i].Spec.VolumeMode == nil || *desired.Spec.VolumeClaimTemplates[i].Spec.VolumeMode == "" {
			desired.Spec.VolumeClaimTemplates[i].Spec.VolumeMode = existing.Spec.VolumeClaimTemplates[i].Spec.VolumeMode
		}

		if desired.Spec.VolumeClaimTemplates[i].Name != existing.Spec.VolumeClaimTemplates[i].Name {
			return true
		}
		if !apiequality.Semantic.DeepEqual(desired.Spec.VolumeClaimTemplates[i].Annotations, existing.Spec.VolumeClaimTemplates[i].Annotations) {
			return true
		}
		if !apiequality.Semantic.DeepEqual(desired.Spec.VolumeClaimTemplates[i].Spec, existing.Spec.VolumeClaimTemplates[i].Spec) {
			return true
		}
	}

	return false
}

// ExpectedIngresses Create or update ingresses
func (o ApolloAllInOne) ExpectedIngresses(ctx context.Context, instance client.Object, params models.Params, expected []networkingv1.Ingress) error {
	for _, obj := range expected {
		desired := obj

		if err := controllerutil.SetControllerReference(instance, &desired, params.Scheme); err != nil {
			return fmt.Errorf("failed to set controller reference: %w", err)
		}

		existing := &networkingv1.Ingress{}
		nns := types.NamespacedName{Namespace: desired.Namespace, Name: desired.Name}
		clientGetErr := params.Client.Get(ctx, nns, existing)
		if clientGetErr != nil && k8serrors.IsNotFound(clientGetErr) {
			if err := params.Client.Create(ctx, &desired); err != nil {
				return fmt.Errorf("failed to create: %w", err)
			}
			params.Log.V(2).Info("created", "ingress.name", desired.Name, "ingress.namespace", desired.Namespace)
			return nil
		} else if clientGetErr != nil {
			return fmt.Errorf("failed to get: %w", clientGetErr)
		}

		// it exists already, merge the two if the end result isn't identical to the existing one
		updated := existing.DeepCopy()
		utils.InitObjectMeta(updated)
		updated.SetAnnotations(desired.GetAnnotations())
		updated.SetLabels(desired.GetLabels())
		updated.SetOwnerReferences(desired.GetOwnerReferences())

		updated.Spec.Rules = desired.Spec.Rules
		updated.Spec.TLS = desired.Spec.TLS
		updated.Spec.IngressClassName = desired.Spec.IngressClassName

		patch := client.MergeFrom(existing)
		if err := params.Client.Patch(ctx, updated, patch); err != nil {
			return fmt.Errorf("failed to apply changes: %w", err)
		}

		params.Log.V(2).Info("applied", "ingress.name", desired.Name, "ingress.namespace", desired.Namespace)
	}
	return nil
}
