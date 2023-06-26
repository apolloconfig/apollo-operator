package apollo

import (
	"apollo.io/apollo-operator/pkg/reconcile/models"
	"apollo.io/apollo-operator/pkg/utils"
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

func (o ApolloAllInOne) ExpectedConfigMaps(ctx context.Context, instance client.Object, params models.Params, expected []corev1.ConfigMap, retry bool) error {
	for _, obj := range expected {
		desired := obj

		// 建立关联后，删除apollo-portal资源时就会将configmap也删除掉
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
			// 创建成功进入下次循环
			continue
		} else if getErr != nil {
			return fmt.Errorf("failed to get: %w", getErr)
		}

		// it exists already, merge the two if the end result isn't identical to the existing one
		updated := existing.DeepCopy()
		utils.InitObjectMeta(updated)
		// TODO 删除该日志
		params.Log.V(2).Info("查看existing和updated", "existing configmap：", existing, "updated configmap：", updated)

		updated.SetAnnotations(desired.GetAnnotations())
		updated.SetLabels(desired.GetLabels())
		updated.SetOwnerReferences(desired.GetOwnerReferences())

		updated.Data = desired.Data
		updated.BinaryData = desired.BinaryData

		// 将旧的configmap修改为新的configmap
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

// 创建或更新endpoints
func (o ApolloAllInOne) ExpectedEndpoints(ctx context.Context, instance client.Object, params models.Params, expected []corev1.Endpoints, retry bool) error {
	for _, obj := range expected {
		desired := obj

		// 建立关联后，删除apollo-portal资源时就会将Endpoints也删除掉
		if err := controllerutil.SetControllerReference(instance, &desired, params.Scheme); err != nil {
			return fmt.Errorf("failed to set controller reference: %w", err)
		}

		existing := &corev1.Endpoints{}
		namespaceName := types.NamespacedName{Namespace: desired.Namespace, Name: desired.Name}
		getErr := params.Client.Get(ctx, namespaceName, existing)
		if getErr != nil && k8serrors.IsNotFound(getErr) {
			// 不存在则直接创建desired的资源
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
			// 创建成功进入下次循环
			continue
		} else if getErr != nil {
			return fmt.Errorf("failed to get: %w", getErr)
		}

		// it exists already, merge the two if the end result isn't identical to the existing one
		updated := existing.DeepCopy()
		utils.InitObjectMeta(updated)
		// TODO 删除该日志
		params.Log.V(2).Info("查看existing和updated", "existing endpoints：", existing, "updated endpoints：", updated)

		if updated.Subsets == nil {
			updated.Subsets = []corev1.EndpointSubset{}
		}
		updated.SetAnnotations(desired.GetAnnotations())
		updated.SetLabels(desired.GetLabels())
		updated.SetOwnerReferences(desired.GetOwnerReferences())

		// 用 desired 替换
		for i, subset := range desired.Subsets {
			updated.Subsets[i] = subset
		}

		// 将旧的Endpoints修改为新的Endpoints
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

// 创建或更新service
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

// 创建或更新deployment
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

		updated.Spec = desired.Spec // 一定要注意spec中的slice的来源，不能用遍历map的方式获取，否则会导致每次调谐重新创建pod

		patch := client.MergeFrom(existing)
		if err := params.Client.Patch(ctx, updated, patch); err != nil {
			return fmt.Errorf("failed to apply changes: %w", err)
		}

		params.Log.V(2).Info("applied", "deployment.name", desired.Name, "deployment.namespace", desired.Namespace)
	}

	return nil
}

// ExpectedIngresses 创建或更新ingresses
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
