package apolloportal

import (
	apolloiov1alpha1 "apollo.io/apollo-operator/api/v1alpha1"
	"apollo.io/apollo-operator/pkg/reconcile"
	"apollo.io/apollo-operator/pkg/utils"
	"apollo.io/apollo-operator/pkg/utils/naming"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// +kubebuilder:rbac:groups="",resources=endpoints,verbs=get;list;watch;create;update;patch;delete

// Endpoints reconciles the endpoint(s) required for the instance in the current context.
func Endpoints(ctx context.Context, params Params) error {
	desired := []corev1.Endpoints{
		desiredEndpoints(ctx, params),
	}

	// TODO 可以优化为先获取create、upodate、delete列表，然后再统一apply

	// first, handle the create/update parts
	if err := expectedEndpoints(ctx, params, desired, true); err != nil {
		return fmt.Errorf("failed to reconcile the expected Endpoints: %w", err)
	}

	// then, delete the extra objects
	if err := deleteEndpoints(ctx, params, desired); err != nil {
		return fmt.Errorf("failed to reconcile the Endpoints to be deleted: %w", err)
	}

	return nil
}

func desiredEndpoints(ctx context.Context, params Params) corev1.Endpoints {
	// TODO 目前需求只有一个subset，后续可以拓展为多个
	subset, _ := buildSubset(ctx, params.Instance)

	return corev1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Name:      naming.PortalDBService(&params.Instance), // NOTE endpoints 名字要和链接的service的服务名相同
			Namespace: params.Instance.Namespace,
			Labels:    reconcile.SelectorLabels(&params.Instance),
		},
		Subsets: []corev1.EndpointSubset{subset},
	}
}

func expectedEndpoints(ctx context.Context, params Params, expected []corev1.Endpoints, retry bool) error {
	for _, obj := range expected {
		desired := obj

		// 建立关联后，删除apollo-portal资源时就会将Endpoints也删除掉
		if err := controllerutil.SetControllerReference(&params.Instance, &desired, params.Scheme); err != nil {
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
					if err := expectedEndpoints(ctx, params, expected, false); err != nil {
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

func deleteEndpoints(ctx context.Context, params Params, expected []corev1.Endpoints) error {
	opts := []client.ListOption{
		client.InNamespace(params.Instance.Namespace),
		client.MatchingLabels(map[string]string{
			"app.kubernetes.io/instance":   naming.Truncate("%s.%s", 63, params.Instance.Namespace, params.Instance.Name),
			"app.kubernetes.io/managed-by": "apollo-operator",
		}),
	}
	Endpointslist := &corev1.EndpointsList{}
	if err := params.Client.List(ctx, Endpointslist, opts...); err != nil {
		return fmt.Errorf("failed to list endpoints : %w", err)
	}

	// 删除不属于expected的部分
	for i := range Endpointslist.Items {
		existing := Endpointslist.Items[i]
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

func buildSubset(_ context.Context, instance apolloiov1alpha1.ApolloPortal) (corev1.EndpointSubset, error) {
	return corev1.EndpointSubset{
		Addresses: []corev1.EndpointAddress{
			{
				IP: instance.Spec.PortalDB.Host,
			},
		},
		Ports: []corev1.EndpointPort{
			{
				Port:     instance.Spec.PortalDB.Port,
				Protocol: corev1.ProtocolTCP,
			},
		},
	}, nil
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
