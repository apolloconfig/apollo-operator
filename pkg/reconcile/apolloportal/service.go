package apolloportal

import (
	"apollo.io/apollo-operator/pkg/reconcile"
	"apollo.io/apollo-operator/pkg/utils/naming"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// headless label is to differentiate the headless service from the clusterIP service.
const (
	headlessLabel  = "apollo.io/apollo-headless-service"
	headlessExists = "Exists"
)

// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete

// Services reconciles the service(s) required for the instance in the current context.
func Services(ctx context.Context, params Params) error {
	desired := []corev1.Service{}
	type builder func(context.Context, Params) *corev1.Service
	for _, builder := range []builder{portaldbService, portalService} {
		svc := builder(ctx, params)
		// add only the non-nil to the list
		if svc != nil {
			desired = append(desired, *svc)
		}
	}

	// first, handle the create/update parts
	if err := expectedServices(ctx, params, desired); err != nil {
		return fmt.Errorf("failed to reconcile the expected services: %w", err)
	}

	// then, delete the extra objects
	if err := deleteServices(ctx, params, desired); err != nil {
		return fmt.Errorf("failed to reconcile the services to be deleted: %w", err)
	}

	return nil
}

func portaldbService(ctx context.Context, params Params) *corev1.Service {
	name := naming.ResourceNameWithSuffix(&params.Instance, "portaldb")
	labels := reconcile.Labels(&params.Instance, name, []string{})

	portaldbService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name, // TODO 暂时
			Namespace: params.Instance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type: params.Instance.Spec.PortalDB.Service.Type,
			Ports: []corev1.ServicePort{
				{ // portaldb 目前只需一个端口号即可
					Protocol:   corev1.ProtocolTCP,
					Port:       params.Instance.Spec.PortalDB.Service.Port,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: params.Instance.Spec.PortalDB.Port},
				},
			},
		},
	}
	if params.Instance.Spec.PortalDB.Service.Type == corev1.ServiceTypeExternalName {
		portaldbService.Spec.ExternalName = params.Instance.Spec.PortalDB.Host
	}

	// TODO 后端如果是statefulset的话需要在portaldbService中添加selector

	return portaldbService
}

func portalService(ctx context.Context, params Params) *corev1.Service {
	name := naming.ResourceNameWithSuffix(&params.Instance, "portal")
	labels := reconcile.Labels(&params.Instance, name, []string{})

	portalService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: params.Instance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type: params.Instance.Spec.Service.Type,
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       params.Instance.Spec.Service.Port,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: params.Instance.Spec.Service.TargetPort},
				},
			},
			Selector:        reconcile.SelectorLabels(&params.Instance),
			SessionAffinity: params.Instance.Spec.Service.SessionAffinity,
		},
	}
	return portalService
}

// TODO 多节点数据库一般需要使用，主节点读写，从节点读的情况，即自定义负载均衡
// headlessService 无头服务适用有状态应用部署,例如数据库
func headlessService(ctx context.Context, params Params) *corev1.Service {
	h := portaldbService(ctx, params)
	if h == nil {
		return nil
	}

	h.Name = naming.HeadlessService(&params.Instance)
	h.Labels[headlessLabel] = headlessExists
	h.Spec.ClusterIP = "None"
	return h
}

func expectedServices(ctx context.Context, params Params, expected []corev1.Service) error {
	for _, obj := range expected {
		desired := obj

		if err := controllerutil.SetControllerReference(&params.Instance, &desired, params.Scheme); err != nil {
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

func deleteServices(ctx context.Context, params Params, expected []corev1.Service) error {
	opts := []client.ListOption{
		client.InNamespace(params.Instance.Namespace),
		client.MatchingLabels(map[string]string{
			"app.kubernetes.io/instance":   naming.Truncate("%s.%s", 63, params.Instance.Namespace, params.Instance.Name),
			"app.kubernetes.io/managed-by": "apollo-operator",
		}),
	}
	list := &corev1.ServiceList{}
	if err := params.Client.List(ctx, list, opts...); err != nil {
		return fmt.Errorf("failed to list: %w", err)
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
