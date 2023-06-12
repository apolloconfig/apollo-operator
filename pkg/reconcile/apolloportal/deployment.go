package apolloportal

import (
	apolloiov1alpha1 "apollo.io/apollo-operator/api/v1alpha1"
	"apollo.io/apollo-operator/pkg/reconcile"
	"apollo.io/apollo-operator/pkg/utils"
	"apollo.io/apollo-operator/pkg/utils/naming"
	"context"
	"fmt"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// +kubebuilder:rbac:groups="apps",resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Deployments reconciles the deployment(s) required for the instance in the current context.
func Deployments(ctx context.Context, params Params) error {
	desired := []appv1.Deployment{
		desiredDeployment(ctx, params),
	}

	// first, handle the create/update parts
	if err := expectedDeployments(ctx, params, desired); err != nil {
		return fmt.Errorf("failed to reconcile the expected deployments: %w", err)
	}

	// then, delete the extra objects
	if err := deleteDeployments(ctx, params, desired); err != nil {
		return fmt.Errorf("failed to reconcile the deployments to be deleted: %w", err)
	}

	return nil
}

func desiredDeployment(ctx context.Context, params Params) appv1.Deployment {
	name := naming.Apollo(&params.Instance)
	labels := reconcile.Labels(&params.Instance, name, []string{})

	container, _ := buildContainer(ctx, params.Instance)
	volume, _ := buildVolume(ctx, params.Instance)
	podTemplateSpec, _ := buildPodTemplate(ctx, params.Instance, container, volume)

	portalDepolyment := appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: params.Instance.Namespace,
			Labels:    labels,
		},
		Spec: appv1.DeploymentSpec{
			Replicas: &params.Instance.Spec.Replicas,
			Selector: &metav1.LabelSelector{MatchLabels: reconcile.SelectorLabels(&params.Instance)},
			Strategy: params.Instance.Spec.Strategy,
			Template: podTemplateSpec,
		},
	}
	return portalDepolyment
}

func expectedDeployments(ctx context.Context, params Params, expected []appv1.Deployment) error {
	for _, obj := range expected {
		desired := obj

		if err := controllerutil.SetControllerReference(&params.Instance, &desired, params.Scheme); err != nil {
			return fmt.Errorf("failed to set controller reference: %w", err)
		}

		existing := &appv1.Deployment{}
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

		updated.Spec = desired.Spec

		patch := client.MergeFrom(existing)
		if err := params.Client.Patch(ctx, updated, patch); err != nil {
			return fmt.Errorf("failed to apply changes: %w", err)
		}

		params.Log.V(2).Info("applied", "deployment.name", desired.Name, "deployment.namespace", desired.Namespace)
	}

	return nil
}

func deleteDeployments(ctx context.Context, params Params, expected []appv1.Deployment) error {
	opts := []client.ListOption{
		client.InNamespace(params.Instance.Namespace),
		client.MatchingLabels(map[string]string{
			"app.kubernetes.io/instance":   naming.Truncate("%s.%s", 63, params.Instance.Namespace, params.Instance.Name),
			"app.kubernetes.io/managed-by": "apollo-operator",
		}),
	}
	list := &appv1.DeploymentList{}
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
			params.Log.V(2).Info("deleted", "deployment.name", existing.Name, "deployment.namespace", existing.Namespace)
		}
	}

	return nil
}

func buildContainer(ctx context.Context, instance apolloiov1alpha1.ApolloPortal) (corev1.Container, error) {
	name := naming.Apollo(&instance)

	if instance.Spec.Env == nil {
		instance.Spec.Env = []corev1.EnvVar{}
	}

	// NOTE 和 volume 中内容保持一致
	volumeMounts := []corev1.VolumeMount{
		corev1.VolumeMount{
			Name:      naming.ConfigMap(&instance),
			MountPath: "/apollo-portal/config/application-github.properties",
			SubPath:   "application-github.properties",
		},
		corev1.VolumeMount{
			Name:      naming.ConfigMap(&instance),
			MountPath: "/apollo-portal/config/apollo-env.properties",
			SubPath:   "apollo-env.properties",
		},
	}
	for fileName := range instance.Spec.Config.File {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      naming.ConfigMap(&instance),
			MountPath: "/apollo-portal/config/" + fileName,
			SubPath:   fileName,
		})
	}

	livenessProbe, readinessProbe, _ := buildProbe(ctx, instance)

	container := corev1.Container{
		Name:            name,
		Image:           instance.Spec.Image,
		ImagePullPolicy: instance.Spec.ImagePullPolicy,
		Ports: []corev1.ContainerPort{
			corev1.ContainerPort{
				Name:          "http",
				ContainerPort: instance.Spec.ContainerPort,
				Protocol:      corev1.ProtocolTCP,
			},
		},
		Env: append(instance.Spec.Env, corev1.EnvVar{
			Name:  "SPRING_PROFILES_ACTIVE",
			Value: instance.Spec.Config.Profiles,
		}),
		VolumeMounts:   volumeMounts,
		LivenessProbe:  livenessProbe,
		ReadinessProbe: readinessProbe,
		Resources:      instance.Spec.Resources,
	}
	return container, nil
}

func buildVolume(ctx context.Context, instance apolloiov1alpha1.ApolloPortal) (corev1.Volume, error) {
	var defaultMode int32 = 420
	volume := corev1.Volume{
		Name: naming.ConfigMap(&instance),
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: naming.ConfigMap(&instance)}, // NOTE 和configmap的名字保持一致
				Items: []corev1.KeyToPath{
					corev1.KeyToPath{
						Key:  "application-github.properties",
						Path: "application-github.properties",
					},
					corev1.KeyToPath{
						Key:  "apollo-env.properties",
						Path: "apollo-env.properties",
					},
				},
				DefaultMode: &defaultMode,
			},
		},
	}

	for fileName := range instance.Spec.Config.File {
		volume.VolumeSource.ConfigMap.Items = append(volume.VolumeSource.ConfigMap.Items, corev1.KeyToPath{
			Key:  fileName,
			Path: fileName,
		})
	}
	return volume, nil
}

func buildPodTemplate(ctx context.Context, instance apolloiov1alpha1.ApolloPortal, container corev1.Container, volume corev1.Volume) (corev1.PodTemplateSpec, error) {
	template := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: reconcile.SelectorLabels(&instance),
		},
		Spec: corev1.PodSpec{
			Containers:       []corev1.Container{container},
			Volumes:          []corev1.Volume{volume},
			ImagePullSecrets: instance.Spec.ImagePullSecrets,
			NodeSelector:     instance.Spec.NodeSelector,
			Affinity:         &instance.Spec.Affinity,
			Tolerations:      instance.Spec.Tolerations,
		},
	}
	return template, nil
}

func buildProbe(ctx context.Context, instance apolloiov1alpha1.ApolloPortal) (livenessProbe, readinessProbe *corev1.Probe, err error) {
	livenessProbe = &instance.Spec.Probe.Liveness
	readinessProbe = &instance.Spec.Probe.Readineeds
	livenessProbe.ProbeHandler = corev1.ProbeHandler{
		TCPSocket: &corev1.TCPSocketAction{
			Port: intstr.IntOrString{Type: intstr.Int, IntVal: instance.Spec.ContainerPort},
		},
	}
	readinessProbe.ProbeHandler = corev1.ProbeHandler{
		HTTPGet: &corev1.HTTPGetAction{
			Port: intstr.IntOrString{Type: intstr.Int, IntVal: instance.Spec.ContainerPort},
			Path: instance.Spec.Config.ContextPath + "/health",
		},
	}
	return livenessProbe, readinessProbe, nil
}
