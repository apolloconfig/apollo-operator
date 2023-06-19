// Package naming is for determining the names for components (containers, services, ...).
package naming

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
)

/* 所有在多个位置使用的名字，都要调用这里的函数，方便保持一致 */

// ConfigMap builds the name for the portal config map used in the apollo-operator.
func ConfigMap(obj client.Object) string {
	return DNSName(Truncate("%s-portal-configmap", 63, obj.GetName()))
}

// AdminConfigMap builds the name for the admin configmap used in the apollo-operator.
func AdminConfigMap(obj client.Object) string {
	return DNSName(Truncate("%s-admin-configmap", 63, obj.GetName()))
}

// ConfigConfigMap builds the name for the config configmap used in the apollo-operator.
func ConfigConfigMap(obj client.Object) string {
	return DNSName(Truncate("%s-config-configmap", 63, obj.GetName()))
}

// PortalDeployment builds the name for the portal deployment used in the apollo-operator.
func PortalDeployment(obj client.Object) string {
	return DNSName(Truncate("%s-portal-deployment", 63, obj.GetName()))
}

// AdminDeployment builds the name for the admin deployment used in the apollo-operator.
func AdminDeployment(obj client.Object) string {
	return DNSName(Truncate("%s-admin-deployment", 63, obj.GetName()))
}

// ConfigDeployment builds the name for the config deployment used in the apollo-operator.
func ConfigDeployment(obj client.Object) string {
	return DNSName(Truncate("%s-config-deployment", 63, obj.GetName()))
}

// Apollo builds the apollo resource name used in the apollo-operator.
func Apollo(obj client.Object) string {
	return DNSName(Truncate("%s", 63, obj.GetName()))
}

// Container returns the name to use for the container in the pod.
func Container() string {
	return "apollo-container"
}

// HeadlessService builds the name for the headless service used in the apollo-operator.
func HeadlessService(obj client.Object) string {
	return DNSName(Truncate("%s-headless", 63, Service(obj)))
}

// Service builds the name for the service used in the apollo-operator.
func Service(obj client.Object) string {
	return DNSName(Truncate("%s-service", 63, obj.GetName()))
}

// PortalService builds the name for the portal service used in the apollo-operator.
func PortalService(obj client.Object) string {
	return DNSName(Truncate("%s-portal", 63, obj.GetName()))
}

// AdminService builds the name for the admin service used in the apollo-operator.
func AdminService(obj client.Object) string {
	return DNSName(Truncate("%s-admin", 63, obj.GetName()))
}

// ConfigService builds the name for the config service used in the apollo-operator.
func ConfigService(obj client.Object) string {
	return DNSName(Truncate("%s-config", 63, obj.GetName()))
}

// PortalDBService builds the name for the portal db service used in the apollo-operator.
func PortalDBService(obj client.Object) string {
	return DNSName(Truncate("%s-portaldb", 63, obj.GetName()))
}

// ConfigDBService builds the name for the config db service used in the apollo-operator.
func ConfigDBService(obj client.Object) string {
	return DNSName(Truncate("%s-configdb", 63, obj.GetName()))
}

// Ingress builds the ingress name based on the instance.
func Ingress(obj client.Object) string {
	return DNSName(Truncate("%s-ingress", 63, obj.GetName()))
}

// ServiceAccount builds the service account name based on the instance.
func ServiceAccount(obj client.Object) string {
	return DNSName(Truncate("%s-serviceaccount", 63, obj.GetName()))
}

// ResourceNameWithSuffix builds the resource name based on the instance and suffix.
func ResourceNameWithSuffix(obj client.Object, suffix string) string {
	if suffix == "" {
		return Apollo(obj)
	}
	return DNSName(Truncate("%s-%s", 63, obj.GetName(), suffix))
}
