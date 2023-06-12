package apolloportal

import (
	apolloiov1alpha1 "apollo.io/apollo-operator/api/v1alpha1"
	"apollo.io/apollo-operator/pkg/reconcile"
	"apollo.io/apollo-operator/pkg/utils"
	"apollo.io/apollo-operator/pkg/utils/naming"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strings"
)

// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete

// ConfigMaps reconciles the configmap(s) required for the instance in the current context.
func ConfigMaps(ctx context.Context, params Params) error {
	desired := []corev1.ConfigMap{
		desiredConfigMap(ctx, params),
	}

	// TODO 可以优化为先获取create、upodate、delete列表，然后再统一apply

	// first, handle the create/update parts
	if err := expectedConfigMaps(ctx, params, desired, true); err != nil {
		return fmt.Errorf("failed to reconcile the expected configmaps: %w", err)
	}

	// then, delete the extra objects
	if err := deleteConfigMaps(ctx, params, desired); err != nil {
		return fmt.Errorf("failed to reconcile the configmaps to be deleted: %w", err)
	}

	return nil
}

func desiredConfigMap(ctx context.Context, params Params) corev1.ConfigMap {
	// NOTE 一定要和volume中使用的名字一致
	name := naming.ConfigMap(&params.Instance)
	labels := reconcile.Labels(&params.Instance, name, []string{})

	config, _ := buildConfig(ctx, params.Instance)

	return corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   params.Instance.Namespace,
			Labels:      labels,
			Annotations: params.Instance.Annotations,
		},
		Data: config,
	}
}

func expectedConfigMaps(ctx context.Context, params Params, expected []corev1.ConfigMap, retry bool) error {
	for _, obj := range expected {
		desired := obj

		// 建立关联后，删除apollo-portal资源时就会将configmap也删除掉
		if err := controllerutil.SetControllerReference(&params.Instance, &desired, params.Scheme); err != nil {
			return fmt.Errorf("failed to set controller reference: %w", err)
		}

		existing := &corev1.ConfigMap{}
		namespaceName := types.NamespacedName{Namespace: desired.Namespace, Name: desired.Name}
		getErr := params.Client.Get(ctx, namespaceName, existing)
		if getErr != nil && apierrors.IsNotFound(getErr) {
			// 不存在则直接创建desired的资源
			if createErr := params.Client.Create(ctx, &desired); createErr != nil {
				if apierrors.IsAlreadyExists(createErr) && retry {
					// let's try again? we probably had multiple updates at one, and now it exists already
					if err := expectedConfigMaps(ctx, params, expected, false); err != nil {
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

		if configMapChanged(&desired, existing) {
			params.Recorder.Event(updated, "Normal", "ConfigUpdate ", fmt.Sprintf("ApolloPortal Config changed - %s/%s", desired.Namespace, desired.Name))
		}

		params.Log.V(2).Info("applied", "configmap.name", desired.Name, "configmap.namespace", desired.Namespace)
	}

	return nil
}

func deleteConfigMaps(ctx context.Context, params Params, expected []corev1.ConfigMap) error {
	opts := []client.ListOption{
		client.InNamespace(params.Instance.Namespace),
		client.MatchingLabels(map[string]string{
			"app.kubernetes.io/instance":   naming.Truncate("%s.%s", 63, params.Instance.Namespace, params.Instance.Name),
			"app.kubernetes.io/managed-by": "apollo-operator",
		}),
	}
	configmaplist := &corev1.ConfigMapList{}
	if err := params.Client.List(ctx, configmaplist, opts...); err != nil {
		return fmt.Errorf("failed to list configmap : %w", err)
	}

	// 删除不属于expected的部分
	for i := range configmaplist.Items {
		existing := configmaplist.Items[i]
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

func buildConfig(_ context.Context, instance apolloiov1alpha1.ApolloPortal) (map[string]string, error) {
	// 从instance提取出configmap的data部分
	data := map[string]string{}

	// apollo-env.properties
	// TODO 配置meta服务的地址，比如CR中指定configservice的namespace和name
	var apolloEnvConfig []string
	for env, address := range instance.Spec.Config.MetaServers {
		apolloEnvConfig = append(apolloEnvConfig, fmt.Sprintf("%s.meta = %s", env, address))
	}
	data["apollo-env.properties"] = strings.Join(apolloEnvConfig, "\n")

	// application-github.properties
	// TODO 多种数据库支持
	apolloGithubConfig := []string{
		fmt.Sprintf("spring.datasource.username = %s", instance.Spec.PortalDB.Username),
		fmt.Sprintf("spring.datasource.password = %s", instance.Spec.PortalDB.Password),
		fmt.Sprintf("spring.datasource.url = jdbc:mysql://%s.%s:%d/%s?%s",
			naming.PortalDBService(&instance), // NOTE 一定要确保和portaldb服务名一致
			instance.Namespace,                // NOTE 一定要确保和portaldb服务的命名空间一致
			instance.Spec.PortalDB.Service.Port,
			instance.Spec.PortalDB.DBName,
			instance.Spec.PortalDB.ConnectionStringProperties),
	}
	if instance.Spec.Config.Envs != "" {
		apolloGithubConfig = append(apolloGithubConfig, fmt.Sprintf("apollo.portal.envs = %s", instance.Spec.Config.Envs))
	}
	if instance.Spec.Config.ContextPath != "" {
		apolloGithubConfig = append(apolloGithubConfig, fmt.Sprintf("server.servlet.context-path = %s", instance.Spec.Config.ContextPath))
	}
	// TODO config.contextPath

	data["application-github.properties"] = strings.Join(apolloGithubConfig, "\n")

	// 其余配置文件
	for fileName, content := range instance.Spec.Config.File {
		data[fileName] = content
	}

	return data, nil

}

func configMapChanged(desired *corev1.ConfigMap, existing *corev1.ConfigMap) bool {
	return !reflect.DeepEqual(desired.Data, existing.Data)
}
