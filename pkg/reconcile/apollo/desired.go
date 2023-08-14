package apollo

import (
	apolloiov1alpha1 "apolloconfig.com/apollo-operator/api/v1alpha1"
	"apolloconfig.com/apollo-operator/pkg/reconcile/models"
	"apolloconfig.com/apollo-operator/pkg/utils"
	"apolloconfig.com/apollo-operator/pkg/utils/naming"
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

// DesiredConfigMaps 构建configmap对象
func (o ApolloAllInOne) DesiredConfigMaps(ctx context.Context, instance client.Object, params models.Params) []corev1.ConfigMap {
	desired := []corev1.ConfigMap{}
	type builder func(context.Context, client.Object, models.Params) *corev1.ConfigMap
	for _, builder := range []builder{adminServiceConfig, configServiceConfig, portalServiceConfig} {
		cm := builder(ctx, instance, params)
		// add only the non-nil to the list
		if cm != nil {
			desired = append(desired, *cm)
		}
	}
	return desired
}

func configServiceConfig(_ context.Context, obj client.Object, params models.Params) *corev1.ConfigMap {
	instance := obj.(*apolloiov1alpha1.Apollo)

	// NOTE 一定要和volume中使用的名字一致
	name := naming.ConfigConfigMap(instance)
	labels := utils.Labels(instance, name, []string{})

	// 从instance提取出configmap的data部分
	data := map[string]string{}

	// application-github.properties
	// TODO 多种数据库支持
	apolloGithubConfig := []string{
		fmt.Sprintf("spring.datasource.username = %s", "root"),
		fmt.Sprintf("spring.datasource.password = %s", "123456"),
		fmt.Sprintf("spring.datasource.url = jdbc:mysql://%s.%s:%d/%s?%s",
			naming.AllInOneDBService(instance), // NOTE 一定要确保和apollodbService服务名一致
			instance.Namespace,                 // NOTE 一定要确保和apollodbService服务的命名空间一致
			3306,
			"ApolloConfigDB",
			"characterEncoding=utf8"),

		// TODO 这里先默认k8s提供的服务发现地址
		fmt.Sprintf("apollo.config-service.url = http://%s.%s:%d%s",
			naming.ConfigService(instance), // NOTE 一定要确保和configService服务名一致
			instance.Namespace,             // NOTE 一定要确保和configService服务的命名空间一致
			8080,                           // instance.Spec.ConfigService.Service.Port
			""),                            // instance.Spec.ConfigService.Config.ContextPath
		fmt.Sprintf("apollo.admin-service.url = http://%s.%s:%d%s",
			naming.AdminService(instance), // NOTE 一定要确保和configService服务名一致
			instance.Namespace,            // NOTE 一定要确保和configService服务的命名空间一致
			8090,                          // instance.Spec.AdminService.Service.Port
			""),                           // instance.Spec.AdminService.Config.ContextPath
	}

	data["application-github.properties"] = strings.Join(apolloGithubConfig, "\n")

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
	instance := obj.(*apolloiov1alpha1.Apollo)

	// NOTE 一定要和volume中使用的名字一致
	name := naming.AdminConfigMap(instance)
	labels := utils.Labels(instance, name, []string{})

	// 从instance提取出configmap的data部分
	data := map[string]string{}

	// application-github.properties
	// TODO 多种数据库支持
	apolloGithubConfig := []string{
		fmt.Sprintf("spring.datasource.username = %s", "root"),
		fmt.Sprintf("spring.datasource.password = %s", "123456"),
		fmt.Sprintf("spring.datasource.url = jdbc:mysql://%s.%s:%d/%s?%s",
			naming.AllInOneDBService(instance), // NOTE 一定要确保和apollodbService服务名一致
			instance.Namespace,                 // NOTE 一定要确保和apollodbService服务的命名空间一致
			3306,
			"ApolloConfigDB",
			"characterEncoding=utf8"),
	}

	data["application-github.properties"] = strings.Join(apolloGithubConfig, "\n")

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

func portalServiceConfig(_ context.Context, obj client.Object, params models.Params) *corev1.ConfigMap {
	instance := obj.(*apolloiov1alpha1.Apollo)

	// NOTE 一定要和volume中使用的名字一致
	name := naming.ConfigMap(instance)
	labels := utils.Labels(instance, name, []string{})

	// 从instance提取出configmap的data部分
	data := map[string]string{}

	// apollo-env.properties
	// TODO 配置meta服务的地址, allinone模式下暂时不需要配置
	apolloEnvConfig := []string{
		fmt.Sprintf("%s.meta = http://%s.%s:%d",
			"dev",
			naming.ConfigService(instance), // NOTE 一定要确保和configService服务名一致
			instance.Namespace,             // NOTE 一定要确保和configService服务的命名空间一致
			8080,                           // NOTE instance.Spec.ConfigService.Service.Port
		),
	}
	//for env, address := range instance.Spec.PortalService.Config.MetaServers {
	//	apolloEnvConfig = append(apolloEnvConfig, fmt.Sprintf("%s.meta = %s", "dev", address))
	//}
	data["apollo-env.properties"] = strings.Join(apolloEnvConfig, "\n")

	// application-github.properties
	// TODO 多种数据库支持
	apolloGithubConfig := []string{
		fmt.Sprintf("spring.datasource.username = %s", "root"),
		fmt.Sprintf("spring.datasource.password = %s", "123456"),
		fmt.Sprintf("spring.datasource.url = jdbc:mysql://%s.%s:%d/%s?%s",
			naming.AllInOneDBService(instance), // NOTE 一定要确保和apollodbService服务名一致
			instance.Namespace,                 // NOTE 一定要确保和apollodbService服务的命名空间一致
			3306,
			"ApolloPortalDB",
			"characterEncoding=utf8"),
	}
	//if instance.Spec.PortalService.Config.Envs != "" {
	//	apolloGithubConfig = append(apolloGithubConfig, fmt.Sprintf("apollo.portal.envs = %s", instance.Spec.PortalService.Config.Envs))
	//}

	data["application-github.properties"] = strings.Join(apolloGithubConfig, "\n")

	// 其余配置文件
	for _, file := range instance.Spec.PortalService.Config.Files {
		data[file.Name] = file.Content
	}

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
func (o ApolloAllInOne) DesiredEndpoints(ctx context.Context, instance client.Object, params models.Params) []corev1.Endpoints {
	return []corev1.Endpoints{}
}

// DesiredServices 构建service对象
func (o ApolloAllInOne) DesiredServices(ctx context.Context, instance client.Object, params models.Params) []corev1.Service {
	desired := []corev1.Service{}
	type builder func(context.Context, client.Object, models.Params) *corev1.Service
	for _, builder := range []builder{apollodbService, configService, adminService, portalService} {
		svc := builder(ctx, instance, params)
		// add only the non-nil to the list
		if svc != nil {
			desired = append(desired, *svc)
		}
	}
	return desired
}

func apollodbService(ctx context.Context, obj client.Object, params models.Params) *corev1.Service {
	instance := obj.(*apolloiov1alpha1.Apollo)
	name := naming.AllInOneDBService(instance)
	labels := utils.Labels(instance, name, []string{})

	apollodbService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: instance.GetNamespace(),
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP, // TODO 稍后修改
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       3306,                                               // TODO 稍后修改
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 3306}, // TODO 稍后修改
				},
			},
			Selector: utils.SelectorLabelsWithCustom(instance, map[string]string{"app": "apollo-db"}),
		},
	}

	return apollodbService
}

func configService(ctx context.Context, obj client.Object, params models.Params) *corev1.Service {
	instance := obj.(*apolloiov1alpha1.Apollo)
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
	instance := obj.(*apolloiov1alpha1.Apollo)
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

func portalService(ctx context.Context, obj client.Object, params models.Params) *corev1.Service {
	instance := obj.(*apolloiov1alpha1.Apollo)
	name := naming.PortalService(instance)
	labels := utils.Labels(instance, name, []string{})

	portalService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type: instance.Spec.PortalService.Service.Type,
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       instance.Spec.PortalService.Service.Port,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: instance.Spec.PortalService.Service.TargetPort},
				},
			},
			Selector:        utils.SelectorLabelsWithCustom(instance, map[string]string{"app": "portalService"}),
			SessionAffinity: instance.Spec.PortalService.Service.SessionAffinity,
		},
	}
	return portalService
}

// DesiredStatefulSets 构建statefulset对象
func (o ApolloAllInOne) DesiredStatefulSets(ctx context.Context, instance client.Object, params models.Params) []appsv1.StatefulSet {

	desired := []appsv1.StatefulSet{}
	type builder func(context.Context, client.Object, models.Params) *appsv1.StatefulSet
	for _, builder := range []builder{apolloStatefulSet} {
		ss := builder(ctx, instance, params)
		// add only the non-nil to the list
		if ss != nil {
			desired = append(desired, *ss)
		}
	}
	return desired
}

func apolloStatefulSet(ctx context.Context, instance client.Object, params models.Params) *appsv1.StatefulSet {
	name := naming.AllInOneStatefulSet(instance)
	labels := utils.Labels(instance, name, []string{})

	spec, _ := buildApolloStatefulSetSpec(ctx, instance)

	apollodb := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: instance.GetNamespace(),
			Labels:    labels,
		},
		Spec: spec,
	}
	return apollodb
}

func buildApolloStatefulSetSpec(ctx context.Context, obj client.Object) (appsv1.StatefulSetSpec, error) {
	instance := obj.(*apolloiov1alpha1.Apollo)

	container, _ := buildMysqlContainer(ctx, instance)
	//initContainer, _ := buildMysqlInitContainer(ctx, instance)
	pvc, _ := buildMysqlPVC(ctx, instance)

	template := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: utils.SelectorLabelsWithCustom(instance, map[string]string{"app": "apollo-db"}),
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{container},
			Volumes: []corev1.Volume{
				{
					Name: naming.AllInOneSqlScript(instance),
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{Name: naming.AllInOneSqlScriptConfigmap(instance)}, // NOTE 和包含sql语句的configmap的名字保持一致
						},
					},
				},
			},
		},
	}
	var replicas int32 = 1 // TODO 如果不是1的话，需要使用 headless service
	return appsv1.StatefulSetSpec{
		Replicas:        &replicas, // TODO 修改
		ServiceName:     "apollo-db",
		MinReadySeconds: 10,
		Selector:        &metav1.LabelSelector{MatchLabels: utils.SelectorLabelsWithCustom(instance, map[string]string{"app": "apollo-db"})},
		Template:        template,

		VolumeClaimTemplates: []corev1.PersistentVolumeClaim{pvc},
	}, nil
}

func buildMysqlContainer(ctx context.Context, instance *apolloiov1alpha1.Apollo) (corev1.Container, error) {

	volumeMounts := []corev1.VolumeMount{
		corev1.VolumeMount{
			Name:      naming.AllInOnePVC(instance), // NOTE 和 volumeClaimTemplates 中内容保持一致
			MountPath: "/var/lib/mysql",
		},
		corev1.VolumeMount{
			Name:      naming.AllInOneSqlScript(instance),
			MountPath: "/mnt/sql-script",
		},
	}

	container := corev1.Container{
		Name:  naming.Container(),
		Image: "mysql:5.7", // TODO 修改
		Ports: []corev1.ContainerPort{
			corev1.ContainerPort{
				Name:          "mysql-port",
				ContainerPort: 3306, // TODO 修改
			},
		},
		Env: []corev1.EnvVar{
			{
				Name:  "MYSQL_ROOT_PASSWORD", // TODO 修改
				Value: "123456",              // TODO 修改
			},
		},
		Lifecycle: &corev1.Lifecycle{
			PostStart: &corev1.LifecycleHandler{
				Exec: &corev1.ExecAction{ // TODO 修改
					Command: []string{
						"bash", "-c",
						"set -ex\n# Copy the SQL script from the ConfigMap to a temporary location.\ncp /mnt/sql-script/initdb.sql /tmp/initdb.sql\n# Wait for the MySQL server to be ready.\nuntil mysql -h0.0.0.0 -uroot -p${MYSQL_ROOT_PASSWORD} -e \"SELECT 1\"; do sleep 1; done\n# Run the SQL script on the master node.\nmysql -h0.0.0.0 -uroot -p${MYSQL_ROOT_PASSWORD} < /tmp/initdb.sql",
					},
				},
			},
		},
		VolumeMounts: volumeMounts,
	}
	return container, nil
}

func buildMysqlInitContainer(ctx context.Context, instance *apolloiov1alpha1.Apollo) (corev1.Container, error) {

	// NOTE 和 volume 中内容保持一致
	volumeMounts := []corev1.VolumeMount{
		corev1.VolumeMount{
			Name:      naming.AllInOneSqlScript(instance),
			MountPath: "/mnt/sql-script",
		},
	}

	container := corev1.Container{
		Name:  naming.InitContainer(),
		Image: "mysql:5.7", // TODO 修改
		Env: []corev1.EnvVar{
			{
				Name:  "MYSQL_ROOT_PASSWORD", // TODO 修改
				Value: "123456",              // TODO 修改
			},
		},
		Command:      []string{"bash", "-c", "set -ex\n# Copy the SQL script from the ConfigMap to a temporary location.\ncp /mnt/sql-script/initdb.sql /tmp/initdb.sql\n# Wait for the MySQL server to be ready.\nuntil mysql -h mysql-0.mysql -uroot -p${MYSQL_ROOT_PASSWORD} -e \"SELECT 1\"; do sleep 1; done\n# Run the SQL script on the master node.\nif [[ `hostname` =~ -0$ ]]; then\n  mysql -h mysql-0.mysql -uroot -p${MYSQL_ROOT_PASSWORD} \u003c /tmp/initdb.sql\nfi\n"},
		VolumeMounts: volumeMounts,
	}
	return container, nil
}

func buildMysqlPVC(ctx context.Context, instance *apolloiov1alpha1.Apollo) (corev1.PersistentVolumeClaim, error) {
	name := naming.AllInOnePVC(instance) // NOTE 和 Container 中内容保持一致
	//labels := utils.Labels(instance, name, []string{})
	storageClassName := "standard"
	return corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: instance.GetNamespace(),
			//Labels:    labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}, // TODO 修改
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("2Gi"), // pv大小
				},
			}, // TODO 增加资源限制
			// TODO 指定存储类 不填使用默认存储类
			StorageClassName: &storageClassName,
		},
	}, nil
}

// DesiredDeployments 构建deployment对象
func (o ApolloAllInOne) DesiredDeployments(ctx context.Context, instance client.Object, params models.Params) []appsv1.Deployment {

	desired := []appsv1.Deployment{}
	type builder func(context.Context, client.Object, models.Params) *appsv1.Deployment
	for _, builder := range []builder{configDeployment, adminDeployment, portalDeployment} {
		deployment := builder(ctx, instance, params)
		// add only the non-nil to the list
		if deployment != nil {
			desired = append(desired, *deployment)
		}
	}
	return desired
}

func configDeployment(ctx context.Context, instance client.Object, params models.Params) *appsv1.Deployment {
	name := naming.ConfigDeployment(instance) // TODO 调用allinone专门的 名字服务
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
	instance := obj.(*apolloiov1alpha1.Apollo)

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

func buildConfigContainer(ctx context.Context, instance *apolloiov1alpha1.Apollo) (corev1.Container, error) {

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

func buildConfigVolume(ctx context.Context, instance *apolloiov1alpha1.Apollo) (corev1.Volume, error) {
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

func buildConfigProbe(ctx context.Context, instance *apolloiov1alpha1.Apollo) (livenessProbe, readinessProbe *corev1.Probe, err error) {
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
	instance := obj.(*apolloiov1alpha1.Apollo)

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

func buildAdminContainer(ctx context.Context, instance *apolloiov1alpha1.Apollo) (corev1.Container, error) {

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

func buildAdminVolume(ctx context.Context, instance *apolloiov1alpha1.Apollo) (corev1.Volume, error) {
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

func buildAdminProbe(ctx context.Context, instance *apolloiov1alpha1.Apollo) (livenessProbe, readinessProbe *corev1.Probe, err error) {
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

func portalDeployment(ctx context.Context, instance client.Object, params models.Params) *appsv1.Deployment {
	name := naming.PortalDeployment(instance)
	labels := utils.Labels(instance, name, []string{})

	spec, _ := buildPortalDepolymentSpec(ctx, instance)

	portalDepolyment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: instance.GetNamespace(),
			Labels:    labels,
		},
		Spec: spec,
	}
	return portalDepolyment
}

func buildPortalDepolymentSpec(ctx context.Context, obj client.Object) (appsv1.DeploymentSpec, error) {
	instance := obj.(*apolloiov1alpha1.Apollo)

	container, _ := buildPortalContainer(ctx, instance)
	volume, _ := buildPortalVolume(ctx, instance)

	template := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: utils.SelectorLabelsWithCustom(instance, map[string]string{"app": "portalService"}),
		},
		Spec: corev1.PodSpec{
			Containers:       []corev1.Container{container},
			Volumes:          []corev1.Volume{volume},
			ImagePullSecrets: instance.Spec.PortalService.ImagePullSecrets,
			NodeSelector:     instance.Spec.PortalService.NodeSelector,
			Affinity:         &instance.Spec.PortalService.Affinity,
			Tolerations:      instance.Spec.PortalService.Tolerations,
		},
	}
	return appsv1.DeploymentSpec{
		Replicas: &instance.Spec.PortalService.Replicas,
		Selector: &metav1.LabelSelector{MatchLabels: utils.SelectorLabelsWithCustom(instance, map[string]string{"app": "portalService"})},
		Strategy: instance.Spec.PortalService.Strategy,
		Template: template,
	}, nil
}

func buildPortalContainer(ctx context.Context, instance *apolloiov1alpha1.Apollo) (corev1.Container, error) {

	if instance.Spec.PortalService.Env == nil {
		instance.Spec.PortalService.Env = []corev1.EnvVar{}
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

	for _, file := range instance.Spec.PortalService.Config.Files {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      naming.ConfigMap(instance),
			MountPath: "/apollo-portal/config/" + file.Name,
			SubPath:   file.Name,
		})
	}

	livenessProbe, readinessProbe, _ := buildPortalProbe(ctx, instance)

	container := corev1.Container{
		Name:            naming.Container(),
		Image:           instance.Spec.PortalService.Image,
		ImagePullPolicy: instance.Spec.PortalService.ImagePullPolicy,
		Ports: []corev1.ContainerPort{
			corev1.ContainerPort{
				Name:          "http",
				ContainerPort: instance.Spec.PortalService.ContainerPort,
				Protocol:      corev1.ProtocolTCP,
			},
		},
		Env: append(instance.Spec.PortalService.Env, corev1.EnvVar{
			Name:  "SPRING_PROFILES_ACTIVE",
			Value: instance.Spec.PortalService.Config.Profiles,
		}),
		VolumeMounts:   volumeMounts,
		LivenessProbe:  livenessProbe,
		ReadinessProbe: readinessProbe,
		Resources:      instance.Spec.PortalService.Resources,
	}
	return container, nil
}

func buildPortalVolume(ctx context.Context, instance *apolloiov1alpha1.Apollo) (corev1.Volume, error) {
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

	for _, file := range instance.Spec.PortalService.Config.Files {
		volume.VolumeSource.ConfigMap.Items = append(volume.VolumeSource.ConfigMap.Items, corev1.KeyToPath{
			Key:  file.Name,
			Path: file.Name,
		})
	}
	return volume, nil
}

func buildPortalProbe(ctx context.Context, instance *apolloiov1alpha1.Apollo) (livenessProbe, readinessProbe *corev1.Probe, err error) {
	livenessProbe = &instance.Spec.PortalService.Probe.Liveness
	readinessProbe = &instance.Spec.PortalService.Probe.Readineeds

	// TODO 删除 ProbeHandler，因为已完全开放probe的字段
	livenessProbe.ProbeHandler = corev1.ProbeHandler{
		TCPSocket: &corev1.TCPSocketAction{
			Port: intstr.IntOrString{Type: intstr.Int, IntVal: instance.Spec.PortalService.ContainerPort},
		},
	}
	readinessProbe.ProbeHandler = corev1.ProbeHandler{
		HTTPGet: &corev1.HTTPGetAction{
			Port: intstr.IntOrString{Type: intstr.Int, IntVal: instance.Spec.PortalService.ContainerPort},
			Path: instance.Spec.PortalService.Config.ContextPath + "/health",
		},
	}
	return livenessProbe, readinessProbe, nil
}

// DesiredIngresses 构建ingress对象
func (o ApolloAllInOne) DesiredIngresses(ctx context.Context, instance client.Object, params models.Params) []networkingv1.Ingress {
	desired := []networkingv1.Ingress{}
	type builder func(context.Context, client.Object, models.Params) *networkingv1.Ingress
	for _, builder := range []builder{configIngress, adminIngress, portalIngress} {
		ingress := builder(ctx, instance, params)
		// add only the non-nil to the list
		if ingress != nil {
			desired = append(desired, *ingress)
		}
	}
	return desired
}

func configIngress(ctx context.Context, obj client.Object, params models.Params) *networkingv1.Ingress {
	instance := obj.(*apolloiov1alpha1.Apollo)
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

func buildConfigIngressSpec(ctx context.Context, instance *apolloiov1alpha1.Apollo) (networkingv1.IngressSpec, error) {
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

func buildConfigRule(instance *apolloiov1alpha1.Apollo, host string) networkingv1.IngressRule {
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
	instance := obj.(*apolloiov1alpha1.Apollo)
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

func buildAdminIngressSpec(ctx context.Context, instance *apolloiov1alpha1.Apollo) (networkingv1.IngressSpec, error) {
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

func buildAdminRule(instance *apolloiov1alpha1.Apollo, host string) networkingv1.IngressRule {
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

func portalIngress(ctx context.Context, obj client.Object, params models.Params) *networkingv1.Ingress {
	instance := obj.(*apolloiov1alpha1.Apollo)
	name := naming.PortalIngress(instance)
	labels := utils.Labels(instance, name, []string{})

	spec, _ := buildPortalIngressSpec(ctx, instance)

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

func buildPortalIngressSpec(ctx context.Context, instance *apolloiov1alpha1.Apollo) (networkingv1.IngressSpec, error) {
	rules := make([]networkingv1.IngressRule, 0, len(instance.Spec.PortalService.Ingress.Hosts))
	for _, host := range instance.Spec.PortalService.Ingress.Hosts {
		rules = append(rules, buildPortalRule(instance, host))
	}

	return networkingv1.IngressSpec{
		TLS:              instance.Spec.PortalService.Ingress.TLS,
		Rules:            rules,
		IngressClassName: instance.Spec.PortalService.Ingress.IngressClassName,
	}, nil
}

func buildPortalRule(instance *apolloiov1alpha1.Apollo, host string) networkingv1.IngressRule {
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
								Name: naming.PortalService(instance),
								Port: networkingv1.ServiceBackendPort{
									Number: instance.Spec.PortalService.Service.Port,
								},
							},
						},
					},
				},
			},
		},
	}
}
