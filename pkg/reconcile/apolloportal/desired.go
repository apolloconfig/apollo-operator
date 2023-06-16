package apolloportal

import (
	apolloiov1alpha1 "apollo.io/apollo-operator/api/v1alpha1"
	"apollo.io/apollo-operator/pkg/reconcile/models"
	"apollo.io/apollo-operator/pkg/utils"
	"apollo.io/apollo-operator/pkg/utils/naming"
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

func (o Apolloportal) DesiredConfigMaps(ctx context.Context, instance client.Object, params models.Params) []corev1.ConfigMap {
	// NOTE 一定要和volume中使用的名字一致
	name := naming.ConfigMap(instance)
	labels := utils.Labels(instance, name, []string{})

	data, _ := buildConfig(ctx, instance)

	configmap := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   instance.GetNamespace(),
			Labels:      labels,
			Annotations: instance.GetAnnotations(),
		},
		Data: data,
	}

	return []corev1.ConfigMap{configmap}
}

func buildConfig(_ context.Context, obj client.Object) (map[string]string, error) {
	instance := obj.(*apolloiov1alpha1.ApolloPortal)

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
			naming.PortalDBService(instance), // NOTE 一定要确保和portaldb服务名一致
			instance.Namespace,               // NOTE 一定要确保和portaldb服务的命名空间一致
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
	for _, file := range instance.Spec.Config.Files {
		data[file.Name] = file.Content
	}

	return data, nil

}

// 构建endpoints对象
func (o Apolloportal) DesiredEndpoints(ctx context.Context, instance client.Object, params models.Params) []corev1.Endpoints {
	// TODO 目前需求只有一个subset，后续可以拓展为多个
	subset, _ := buildSubset(ctx, instance)
	endpoints := corev1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Name:      naming.PortalDBService(instance), // NOTE endpoints 名字要和链接的service的服务名相同
			Namespace: instance.GetNamespace(),
			Labels:    utils.SelectorLabels(instance),
		},
		Subsets: []corev1.EndpointSubset{subset},
	}
	return []corev1.Endpoints{endpoints}
}

func buildSubset(_ context.Context, obj client.Object) (corev1.EndpointSubset, error) {
	instance := obj.(*apolloiov1alpha1.ApolloPortal)
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

// 构建service对象
func (o Apolloportal) DesiredServices(ctx context.Context, instance client.Object, params models.Params) []corev1.Service {
	desired := []corev1.Service{}
	type builder func(context.Context, client.Object, models.Params) *corev1.Service
	for _, builder := range []builder{portaldbService, portalService} {
		svc := builder(ctx, instance, params)
		// add only the non-nil to the list
		if svc != nil {
			desired = append(desired, *svc)
		}
	}
	return desired
}

func portaldbService(ctx context.Context, obj client.Object, params models.Params) *corev1.Service {
	instance := obj.(*apolloiov1alpha1.ApolloPortal)
	name := naming.PortalDBService(instance)
	labels := utils.Labels(instance, name, []string{})

	portaldbService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name, // TODO 暂时
			Namespace: instance.GetNamespace(),
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type: instance.Spec.PortalDB.Service.Type,
			Ports: []corev1.ServicePort{
				{ // portaldb 目前只需一个端口号即可
					Protocol:   corev1.ProtocolTCP,
					Port:       instance.Spec.PortalDB.Service.Port,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: instance.Spec.PortalDB.Port},
				},
			},
		},
	}
	if instance.Spec.PortalDB.Service.Type == corev1.ServiceTypeExternalName {
		portaldbService.Spec.ExternalName = instance.Spec.PortalDB.Host
	}

	// TODO 后端如果是statefulset的话需要在portaldbService中添加selector

	return portaldbService
}

func portalService(ctx context.Context, obj client.Object, params models.Params) *corev1.Service {
	instance := obj.(*apolloiov1alpha1.ApolloPortal)
	name := naming.ResourceNameWithSuffix(instance, "portal")
	labels := utils.Labels(instance, name, []string{})

	portalService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type: instance.Spec.Service.Type,
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       instance.Spec.Service.Port,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: instance.Spec.Service.TargetPort},
				},
			},
			Selector:        utils.SelectorLabels(instance),
			SessionAffinity: instance.Spec.Service.SessionAffinity,
		},
	}
	return portalService
}

// TODO 多节点数据库一般需要使用，主节点读写，从节点读的情况，即自定义负载均衡
// headlessService 无头服务适用有状态应用部署,例如数据库
func headlessService(ctx context.Context, instance client.Object, params models.Params) *corev1.Service {
	h := portaldbService(ctx, instance, params)
	if h == nil {
		return nil
	}

	h.Name = naming.HeadlessService(instance)
	h.Labels[utils.HeadlessLabel] = utils.HeadlessExists
	h.Spec.ClusterIP = "None"
	return h
}

// 构建deployment对象
func (o Apolloportal) DesiredDeployments(ctx context.Context, instance client.Object, params models.Params) []appsv1.Deployment {
	name := naming.Apollo(instance)
	labels := utils.Labels(instance, name, []string{})

	spec, _ := buildDepolymentSpec(ctx, instance)

	portalDepolyment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: instance.GetNamespace(),
			Labels:    labels,
		},
		Spec: spec,
	}
	return []appsv1.Deployment{portalDepolyment}
}

func buildContainer(ctx context.Context, instance *apolloiov1alpha1.ApolloPortal) (corev1.Container, error) {

	if instance.Spec.Env == nil {
		instance.Spec.Env = []corev1.EnvVar{}
	}

	// NOTE 和 volume 中内容保持一致
	volumeMounts := []corev1.VolumeMount{
		corev1.VolumeMount{
			Name:      naming.ConfigMap(instance),
			MountPath: "/apollo-portal/config/application-github.properties",
			SubPath:   "application-github.properties",
		},
		corev1.VolumeMount{
			Name:      naming.ConfigMap(instance),
			MountPath: "/apollo-portal/config/apollo-env.properties",
			SubPath:   "apollo-env.properties",
		},
	}
	// TODO map无序，导致调谐前后不一致
	for _, file := range instance.Spec.Config.Files {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      naming.ConfigMap(instance),
			MountPath: "/apollo-portal/config/" + file.Name,
			SubPath:   file.Name,
		})
	}

	livenessProbe, readinessProbe, _ := buildProbe(ctx, instance)

	container := corev1.Container{
		Name:            naming.Container(),
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

func buildVolume(ctx context.Context, instance *apolloiov1alpha1.ApolloPortal) (corev1.Volume, error) {
	var defaultMode int32 = 420
	volume := corev1.Volume{
		Name: naming.ConfigMap(instance),
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: naming.ConfigMap(instance)}, // NOTE 和configmap的名字保持一致
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

	for _, file := range instance.Spec.Config.Files {
		volume.VolumeSource.ConfigMap.Items = append(volume.VolumeSource.ConfigMap.Items, corev1.KeyToPath{
			Key:  file.Name,
			Path: file.Name,
		})
	}
	return volume, nil
}

func buildDepolymentSpec(ctx context.Context, obj client.Object) (appsv1.DeploymentSpec, error) {
	instance := obj.(*apolloiov1alpha1.ApolloPortal)

	container, _ := buildContainer(ctx, instance)
	volume, _ := buildVolume(ctx, instance)

	template := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: utils.SelectorLabels(instance),
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
	return appsv1.DeploymentSpec{
		Replicas: &instance.Spec.Replicas,
		Selector: &metav1.LabelSelector{MatchLabels: utils.SelectorLabels(instance)},
		Strategy: instance.Spec.Strategy,
		Template: template,
	}, nil
}

func buildProbe(ctx context.Context, instance *apolloiov1alpha1.ApolloPortal) (livenessProbe, readinessProbe *corev1.Probe, err error) {
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
