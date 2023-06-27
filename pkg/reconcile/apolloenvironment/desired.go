package apolloenvironment

import (
	apolloiov1alpha1 "apollo.io/apollo-operator/api/v1alpha1"
	"apollo.io/apollo-operator/pkg/reconcile/models"
	"apollo.io/apollo-operator/pkg/utils"
	"apollo.io/apollo-operator/pkg/utils/naming"
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

// DesiredConfigMaps 构建configmap对象
func (o ApolloEnvironment) DesiredConfigMaps(ctx context.Context, instance client.Object, params models.Params) []corev1.ConfigMap {
	desired := []corev1.ConfigMap{}
	type builder func(context.Context, client.Object, models.Params) *corev1.ConfigMap
	for _, builder := range []builder{adminServiceConfig, configServiceConfig} {
		cm := builder(ctx, instance, params)
		// add only the non-nil to the list
		if cm != nil {
			desired = append(desired, *cm)
		}
	}
	return desired
}

func configServiceConfig(_ context.Context, obj client.Object, params models.Params) *corev1.ConfigMap {
	instance := obj.(*apolloiov1alpha1.ApolloEnvironment)

	// NOTE 一定要和volume中使用的名字一致
	name := naming.ConfigConfigMap(instance)
	labels := utils.Labels(instance, name, []string{})

	// 从instance提取出configmap的data部分
	data := map[string]string{}

	// application-github.properties
	// TODO 多种数据库支持
	apolloGithubConfig := []string{
		fmt.Sprintf("spring.datasource.username = %s", instance.Spec.ConfigDB.Username),
		fmt.Sprintf("spring.datasource.password = %s", instance.Spec.ConfigDB.Password),
		fmt.Sprintf("spring.datasource.url = jdbc:mysql://%s.%s:%d/%s?%s",
			naming.ConfigDBService(instance), // NOTE 一定要确保和configdb服务名一致
			instance.Namespace,               // NOTE 一定要确保和configdb服务的命名空间一致
			instance.Spec.ConfigDB.Service.Port,
			instance.Spec.ConfigDB.DBName,
			instance.Spec.ConfigDB.ConnectionStringProperties),

		// TODO 这里先默认k8s提供的服务发现地址
		fmt.Sprintf("apollo.config-service.url = http://%s.%s:%d%s",
			naming.ConfigService(instance), // NOTE 一定要确保和configService服务名一致
			instance.Namespace,             // NOTE 一定要确保和configService服务的命名空间一致
			instance.Spec.ConfigService.Service.Port,
			instance.Spec.ConfigService.Config.ContextPath),
		fmt.Sprintf("apollo.admin-service.url = http://%s.%s:%d%s",
			naming.AdminService(instance), // NOTE 一定要确保和configService服务名一致
			instance.Namespace,            // NOTE 一定要确保和configService服务的命名空间一致
			instance.Spec.AdminService.Service.Port,
			instance.Spec.AdminService.Config.ContextPath),
	}

	if instance.Spec.ConfigService.Config.ContextPath != "" {
		apolloGithubConfig = append(apolloGithubConfig, fmt.Sprintf("server.servlet.context-path = %s", instance.Spec.ConfigService.Config.ContextPath))
	}

	data["application-github.properties"] = strings.Join(apolloGithubConfig, "\n")

	// 其余配置文件
	//for _, file := range instance.Spec.ConfigService.Config.Files {
	//	data[file.Name] = file.Content
	//}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   instance.GetNamespace(),
			Labels:      labels,
			Annotations: instance.GetAnnotations(),
		},
		Data: data,
	}

}

func adminServiceConfig(_ context.Context, obj client.Object, params models.Params) *corev1.ConfigMap {
	instance := obj.(*apolloiov1alpha1.ApolloEnvironment)

	// NOTE 一定要和volume中使用的名字一致
	name := naming.AdminConfigMap(instance)
	labels := utils.Labels(instance, name, []string{})

	// 从instance提取出configmap的data部分
	data := map[string]string{}

	// application-github.properties
	// TODO 多种数据库支持
	apolloGithubConfig := []string{
		fmt.Sprintf("spring.datasource.username = %s", instance.Spec.ConfigDB.Username),
		fmt.Sprintf("spring.datasource.password = %s", instance.Spec.ConfigDB.Password),
		fmt.Sprintf("spring.datasource.url = jdbc:mysql://%s.%s:%d/%s?%s",
			naming.ConfigDBService(instance), // NOTE 一定要确保和configdb服务名一致
			instance.Namespace,               // NOTE 一定要确保和configdb服务的命名空间一致
			instance.Spec.ConfigDB.Service.Port,
			instance.Spec.ConfigDB.DBName,
			instance.Spec.ConfigDB.ConnectionStringProperties),
	}

	if instance.Spec.AdminService.Config.ContextPath != "" {
		apolloGithubConfig = append(apolloGithubConfig, fmt.Sprintf("server.servlet.context-path = %s", instance.Spec.ConfigService.Config.ContextPath))
	}

	data["application-github.properties"] = strings.Join(apolloGithubConfig, "\n")

	// 其余配置文件
	//for _, file := range instance.Spec.ConfigService.Config.Files {
	//	data[file.Name] = file.Content
	//}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   instance.GetNamespace(),
			Labels:      labels,
			Annotations: instance.GetAnnotations(),
		},
		Data: data,
	}

}

// DesiredEndpoints 构建endpoints对象
func (o ApolloEnvironment) DesiredEndpoints(ctx context.Context, instance client.Object, params models.Params) []corev1.Endpoints {
	// TODO 目前需求只有一个subset，后续可以拓展为多个
	subset, _ := buildSubset(ctx, instance)
	endpoints := corev1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Name:      naming.ConfigDBService(instance), // NOTE endpoints 名字要和链接的service的服务名相同
			Namespace: instance.GetNamespace(),
			Labels:    utils.SelectorLabels(instance),
		},
		Subsets: []corev1.EndpointSubset{subset},
	}
	return []corev1.Endpoints{endpoints}
}

func buildSubset(_ context.Context, obj client.Object) (corev1.EndpointSubset, error) {
	instance := obj.(*apolloiov1alpha1.ApolloEnvironment)
	return corev1.EndpointSubset{
		Addresses: []corev1.EndpointAddress{
			{
				IP: instance.Spec.ConfigDB.Host,
			},
		},
		Ports: []corev1.EndpointPort{
			{
				Port:     instance.Spec.ConfigDB.Port,
				Protocol: corev1.ProtocolTCP,
			},
		},
	}, nil
}

// DesiredServices 构建service对象
func (o ApolloEnvironment) DesiredServices(ctx context.Context, instance client.Object, params models.Params) []corev1.Service {
	desired := []corev1.Service{}
	type builder func(context.Context, client.Object, models.Params) *corev1.Service
	for _, builder := range []builder{configdbService, configService, adminService} {
		svc := builder(ctx, instance, params)
		// add only the non-nil to the list
		if svc != nil {
			desired = append(desired, *svc)
		}
	}
	return desired
}

func configdbService(ctx context.Context, obj client.Object, params models.Params) *corev1.Service {
	instance := obj.(*apolloiov1alpha1.ApolloEnvironment)
	name := naming.ConfigDBService(instance)
	labels := utils.Labels(instance, name, []string{})

	configdbService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: instance.GetNamespace(),
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type: instance.Spec.ConfigDB.Service.Type,
			Ports: []corev1.ServicePort{
				{ // configdb 目前只需一个端口号即可
					Protocol:   corev1.ProtocolTCP,
					Port:       instance.Spec.ConfigDB.Service.Port,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: instance.Spec.ConfigDB.Port},
				},
			},
		},
	}
	if instance.Spec.ConfigDB.Service.Type == corev1.ServiceTypeExternalName {
		configdbService.Spec.ExternalName = instance.Spec.ConfigDB.Host
	}

	// TODO 后端如果是statefulset的话需要在configdbService中添加selector

	return configdbService
}

func configService(ctx context.Context, obj client.Object, params models.Params) *corev1.Service {
	instance := obj.(*apolloiov1alpha1.ApolloEnvironment)
	name := naming.ConfigService(instance)
	labels := utils.Labels(instance, name, []string{})

	configService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type: instance.Spec.ConfigService.Service.Type,
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       instance.Spec.ConfigService.Service.Port,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: instance.Spec.ConfigService.Service.TargetPort},
				},
			},
			Selector: utils.SelectorLabelsWithCustom(instance, map[string]string{"app": "configService"}),
			//SessionAffinity: instance.Spec.Service.SessionAffinity,
		},
	}
	return configService
}

func adminService(ctx context.Context, obj client.Object, params models.Params) *corev1.Service {
	instance := obj.(*apolloiov1alpha1.ApolloEnvironment)
	name := naming.AdminService(instance)
	labels := utils.Labels(instance, name, []string{})

	adminService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type: instance.Spec.AdminService.Service.Type,
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       instance.Spec.AdminService.Service.Port,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: instance.Spec.AdminService.Service.TargetPort},
				},
			},
			Selector: utils.SelectorLabelsWithCustom(instance, map[string]string{"app": "adminService"}),
			//SessionAffinity: instance.Spec.Service.SessionAffinity,
		},
	}
	return adminService
}

// TODO 多节点数据库一般需要使用，主节点读写，从节点读的情况，即自定义负载均衡
// headlessService 无头服务适用有状态应用部署,例如数据库
func headlessService(ctx context.Context, instance client.Object, params models.Params) *corev1.Service {
	h := configdbService(ctx, instance, params)
	if h == nil {
		return nil
	}

	h.Name = naming.HeadlessService(instance)
	h.Labels[utils.HeadlessLabel] = utils.HeadlessExists
	h.Spec.ClusterIP = "None"
	return h
}

// DesiredDeployments 构建deployment对象
func (o ApolloEnvironment) DesiredDeployments(ctx context.Context, instance client.Object, params models.Params) []appsv1.Deployment {

	desired := []appsv1.Deployment{}
	type builder func(context.Context, client.Object, models.Params) *appsv1.Deployment
	for _, builder := range []builder{configDeployment, adminDeployment} {
		deployment := builder(ctx, instance, params)
		// add only the non-nil to the list
		if deployment != nil {
			desired = append(desired, *deployment)
		}
	}
	return desired
}

func configDeployment(ctx context.Context, instance client.Object, params models.Params) *appsv1.Deployment {
	name := naming.ConfigDeployment(instance)
	labels := utils.Labels(instance, name, []string{})

	spec, _ := buildConfigDepolymentSpec(ctx, instance)

	configDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: instance.GetNamespace(),
			Labels:    labels,
		},
		Spec: spec,
	}
	return configDeployment
}

func buildConfigDepolymentSpec(ctx context.Context, obj client.Object) (appsv1.DeploymentSpec, error) {
	instance := obj.(*apolloiov1alpha1.ApolloEnvironment)

	container, _ := buildConfigContainer(ctx, instance)
	volume, _ := buildConfigVolume(ctx, instance)

	template := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: utils.SelectorLabelsWithCustom(instance, map[string]string{"app": "configService"}),
		},
		Spec: corev1.PodSpec{
			Containers:       []corev1.Container{container},
			Volumes:          []corev1.Volume{volume},
			ImagePullSecrets: instance.Spec.ConfigService.ImagePullSecrets,
			NodeSelector:     instance.Spec.ConfigService.NodeSelector,
			Affinity:         &instance.Spec.ConfigService.Affinity,
			Tolerations:      instance.Spec.ConfigService.Tolerations,
		},
	}
	return appsv1.DeploymentSpec{
		Replicas: &instance.Spec.ConfigService.Replicas,
		Selector: &metav1.LabelSelector{MatchLabels: utils.SelectorLabelsWithCustom(instance, map[string]string{"app": "configService"})},
		Strategy: instance.Spec.ConfigService.Strategy,
		Template: template,
	}, nil
}

func buildConfigContainer(ctx context.Context, instance *apolloiov1alpha1.ApolloEnvironment) (corev1.Container, error) {

	if instance.Spec.ConfigService.Env == nil {
		instance.Spec.ConfigService.Env = []corev1.EnvVar{}
	}

	// NOTE 和 volume 中内容保持一致
	volumeMounts := []corev1.VolumeMount{
		corev1.VolumeMount{
			Name:      naming.ConfigConfigMap(instance),
			MountPath: "/apollo-configservice/config/application-github.properties",
			SubPath:   "application-github.properties",
		},
	}

	livenessProbe, readinessProbe, _ := buildConfigProbe(ctx, instance)

	container := corev1.Container{
		Name:            naming.Container(),
		Image:           instance.Spec.ConfigService.Image,
		ImagePullPolicy: instance.Spec.ConfigService.ImagePullPolicy,
		Ports: []corev1.ContainerPort{
			corev1.ContainerPort{
				Name:          "http",
				ContainerPort: instance.Spec.ConfigService.ContainerPort,
				Protocol:      corev1.ProtocolTCP,
			},
		},
		Env: append(instance.Spec.ConfigService.Env, corev1.EnvVar{
			Name:  "SPRING_PROFILES_ACTIVE",
			Value: instance.Spec.ConfigService.Config.Profiles,
		}),
		VolumeMounts:   volumeMounts,
		LivenessProbe:  livenessProbe,
		ReadinessProbe: readinessProbe,
		Resources:      instance.Spec.ConfigService.Resources,
	}
	return container, nil
}

func buildConfigVolume(ctx context.Context, instance *apolloiov1alpha1.ApolloEnvironment) (corev1.Volume, error) {
	var defaultMode int32 = 420
	volume := corev1.Volume{
		Name: naming.ConfigConfigMap(instance),
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: naming.ConfigConfigMap(instance)}, // NOTE 和configmap的名字保持一致
				Items: []corev1.KeyToPath{
					corev1.KeyToPath{
						Key:  "application-github.properties",
						Path: "application-github.properties",
					},
				},
				DefaultMode: &defaultMode,
			},
		},
	}

	return volume, nil
}

func buildConfigProbe(ctx context.Context, instance *apolloiov1alpha1.ApolloEnvironment) (livenessProbe, readinessProbe *corev1.Probe, err error) {
	livenessProbe = &instance.Spec.ConfigService.Probe.Liveness
	readinessProbe = &instance.Spec.ConfigService.Probe.Readineeds
	livenessProbe.ProbeHandler = corev1.ProbeHandler{
		TCPSocket: &corev1.TCPSocketAction{
			Port: intstr.IntOrString{Type: intstr.Int, IntVal: instance.Spec.ConfigService.ContainerPort},
		},
	}
	readinessProbe.ProbeHandler = corev1.ProbeHandler{
		HTTPGet: &corev1.HTTPGetAction{
			Port: intstr.IntOrString{Type: intstr.Int, IntVal: instance.Spec.ConfigService.ContainerPort},
			Path: instance.Spec.ConfigService.Config.ContextPath + "/health",
		},
	}
	return livenessProbe, readinessProbe, nil
}

func adminDeployment(ctx context.Context, instance client.Object, params models.Params) *appsv1.Deployment {
	name := naming.AdminDeployment(instance)
	labels := utils.Labels(instance, name, []string{})

	spec, _ := buildAdminDepolymentSpec(ctx, instance)

	adminDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: instance.GetNamespace(),
			Labels:    labels,
		},
		Spec: spec,
	}
	return adminDeployment
}

func buildAdminDepolymentSpec(ctx context.Context, obj client.Object) (appsv1.DeploymentSpec, error) {
	instance := obj.(*apolloiov1alpha1.ApolloEnvironment)

	container, _ := buildAdminContainer(ctx, instance)
	volume, _ := buildAdminVolume(ctx, instance)

	template := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: utils.SelectorLabelsWithCustom(instance, map[string]string{"app": "adminService"}),
		},
		Spec: corev1.PodSpec{
			Containers:       []corev1.Container{container},
			Volumes:          []corev1.Volume{volume},
			ImagePullSecrets: instance.Spec.AdminService.ImagePullSecrets,
			NodeSelector:     instance.Spec.AdminService.NodeSelector,
			Affinity:         &instance.Spec.AdminService.Affinity,
			Tolerations:      instance.Spec.AdminService.Tolerations,
		},
	}
	return appsv1.DeploymentSpec{
		Replicas: &instance.Spec.AdminService.Replicas,
		Selector: &metav1.LabelSelector{MatchLabels: utils.SelectorLabelsWithCustom(instance, map[string]string{"app": "adminService"})},
		Strategy: instance.Spec.AdminService.Strategy,
		Template: template,
	}, nil
}

func buildAdminContainer(ctx context.Context, instance *apolloiov1alpha1.ApolloEnvironment) (corev1.Container, error) {

	if instance.Spec.AdminService.Env == nil {
		instance.Spec.AdminService.Env = []corev1.EnvVar{}
	}

	// NOTE 和 volume 中内容保持一致
	volumeMounts := []corev1.VolumeMount{
		corev1.VolumeMount{
			Name:      naming.AdminConfigMap(instance),
			MountPath: "/apollo-adminservice/config/application-github.properties",
			SubPath:   "application-github.properties",
		},
	}

	livenessProbe, readinessProbe, _ := buildAdminProbe(ctx, instance)

	container := corev1.Container{
		Name:            naming.Container(),
		Image:           instance.Spec.AdminService.Image,
		ImagePullPolicy: instance.Spec.AdminService.ImagePullPolicy,
		Ports: []corev1.ContainerPort{
			corev1.ContainerPort{
				Name:          "http",
				ContainerPort: instance.Spec.AdminService.ContainerPort,
				Protocol:      corev1.ProtocolTCP,
			},
		},
		Env: append(instance.Spec.AdminService.Env, corev1.EnvVar{
			Name:  "SPRING_PROFILES_ACTIVE",
			Value: instance.Spec.AdminService.Config.Profiles,
		}),
		VolumeMounts:   volumeMounts,
		LivenessProbe:  livenessProbe,
		ReadinessProbe: readinessProbe,
		Resources:      instance.Spec.AdminService.Resources,
	}
	return container, nil
}

func buildAdminVolume(ctx context.Context, instance *apolloiov1alpha1.ApolloEnvironment) (corev1.Volume, error) {
	var defaultMode int32 = 420
	volume := corev1.Volume{
		Name: naming.AdminConfigMap(instance),
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: naming.AdminConfigMap(instance)}, // NOTE 和configmap的名字保持一致
				Items: []corev1.KeyToPath{
					corev1.KeyToPath{
						Key:  "application-github.properties",
						Path: "application-github.properties",
					},
				},
				DefaultMode: &defaultMode,
			},
		},
	}

	return volume, nil
}

func buildAdminProbe(ctx context.Context, instance *apolloiov1alpha1.ApolloEnvironment) (livenessProbe, readinessProbe *corev1.Probe, err error) {
	livenessProbe = &instance.Spec.AdminService.Probe.Liveness
	readinessProbe = &instance.Spec.AdminService.Probe.Readineeds

	// TODO 删除 ProbeHandler，因为已完全开放probe的字段
	livenessProbe.ProbeHandler = corev1.ProbeHandler{
		TCPSocket: &corev1.TCPSocketAction{
			Port: intstr.IntOrString{Type: intstr.Int, IntVal: instance.Spec.AdminService.ContainerPort},
		},
	}
	readinessProbe.ProbeHandler = corev1.ProbeHandler{
		HTTPGet: &corev1.HTTPGetAction{
			Port: intstr.IntOrString{Type: intstr.Int, IntVal: instance.Spec.AdminService.ContainerPort},
			Path: instance.Spec.AdminService.Config.ContextPath + "/health",
		},
	}
	return livenessProbe, readinessProbe, nil
}

// DesiredIngresses 构建ingress对象
func (o ApolloEnvironment) DesiredIngresses(ctx context.Context, instance client.Object, params models.Params) []networkingv1.Ingress {
	desired := []networkingv1.Ingress{}
	type builder func(context.Context, client.Object, models.Params) *networkingv1.Ingress
	for _, builder := range []builder{configIngress, adminIngress} {
		ingress := builder(ctx, instance, params)
		// add only the non-nil to the list
		if ingress != nil {
			desired = append(desired, *ingress)
		}
	}
	return desired
}

func configIngress(ctx context.Context, obj client.Object, params models.Params) *networkingv1.Ingress {
	instance := obj.(*apolloiov1alpha1.ApolloEnvironment)
	name := naming.ConfigIngress(instance)
	labels := utils.Labels(instance, name, []string{})

	spec, _ := buildConfigIngressSpec(ctx, instance)

	return &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   instance.GetNamespace(),
			Annotations: instance.Spec.ConfigService.Ingress.Annotations,
			Labels:      labels,
		},
		Spec: spec,
	}
}

func buildConfigIngressSpec(ctx context.Context, instance *apolloiov1alpha1.ApolloEnvironment) (networkingv1.IngressSpec, error) {
	rules := make([]networkingv1.IngressRule, 0, len(instance.Spec.ConfigService.Ingress.Hosts))
	for _, host := range instance.Spec.ConfigService.Ingress.Hosts {
		rules = append(rules, buildConfigRule(instance, host))
	}

	return networkingv1.IngressSpec{
		TLS:              instance.Spec.ConfigService.Ingress.TLS,
		Rules:            rules,
		IngressClassName: instance.Spec.ConfigService.Ingress.IngressClassName,
	}, nil
}

func buildConfigRule(instance *apolloiov1alpha1.ApolloEnvironment, host string) networkingv1.IngressRule {
	pathType := networkingv1.PathTypePrefix // NOTE: 先默认 Prefix
	return networkingv1.IngressRule{
		Host: host,
		IngressRuleValue: networkingv1.IngressRuleValue{
			HTTP: &networkingv1.HTTPIngressRuleValue{
				Paths: []networkingv1.HTTPIngressPath{
					{
						PathType: &pathType,
						Path:     "/",
						Backend: networkingv1.IngressBackend{
							Service: &networkingv1.IngressServiceBackend{
								Name: naming.ConfigService(instance),
								Port: networkingv1.ServiceBackendPort{
									Number: instance.Spec.ConfigService.Service.Port,
								},
							},
						},
					},
				},
			},
		},
	}
}

func adminIngress(ctx context.Context, obj client.Object, params models.Params) *networkingv1.Ingress {
	instance := obj.(*apolloiov1alpha1.ApolloEnvironment)
	name := naming.AdminIngress(instance)
	labels := utils.Labels(instance, name, []string{})

	spec, _ := buildAdminIngressSpec(ctx, instance)

	return &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   instance.GetNamespace(),
			Annotations: instance.Spec.AdminService.Ingress.Annotations,
			Labels:      labels,
		},
		Spec: spec,
	}
}

func buildAdminIngressSpec(ctx context.Context, instance *apolloiov1alpha1.ApolloEnvironment) (networkingv1.IngressSpec, error) {
	rules := make([]networkingv1.IngressRule, 0, len(instance.Spec.AdminService.Ingress.Hosts))
	for _, host := range instance.Spec.AdminService.Ingress.Hosts {
		rules = append(rules, buildAdminRule(instance, host))
	}

	return networkingv1.IngressSpec{
		TLS:              instance.Spec.AdminService.Ingress.TLS,
		Rules:            rules,
		IngressClassName: instance.Spec.AdminService.Ingress.IngressClassName,
	}, nil
}

func buildAdminRule(instance *apolloiov1alpha1.ApolloEnvironment, host string) networkingv1.IngressRule {
	pathType := networkingv1.PathTypePrefix // NOTE: 先默认 Prefix
	return networkingv1.IngressRule{
		Host: host,
		IngressRuleValue: networkingv1.IngressRuleValue{
			HTTP: &networkingv1.HTTPIngressRuleValue{
				Paths: []networkingv1.HTTPIngressPath{
					{
						PathType: &pathType,
						Path:     "/",
						Backend: networkingv1.IngressBackend{
							Service: &networkingv1.IngressServiceBackend{
								Name: naming.AdminService(instance),
								Port: networkingv1.ServiceBackendPort{
									Number: instance.Spec.AdminService.Service.Port,
								},
							},
						},
					},
				},
			},
		},
	}
}
