// Package naming is for determining the names for components (containers, services, ...).
package naming

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
)

/* 所有在多个位置使用的名字，都要调用这里的函数，方便保持一致 */

/* Apollo Environment */

// AdminConfigMap builds the name for the admin configmap used in the apollo-operator.
func AdminConfigMap(obj client.Object) string {
	return DNSName(Truncate("%s-admin-configmap", 63, obj.GetName()))
}

// ConfigConfigMap builds the name for the config configmap used in the apollo-operator.
func ConfigConfigMap(obj client.Object) string {
	return DNSName(Truncate("%s-config-configmap", 63, obj.GetName()))
}

// AdminDeployment builds the name for the admin deployment used in the apollo-operator.
func AdminDeployment(obj client.Object) string {
	return DNSName(Truncate("%s-admin-deployment", 63, obj.GetName()))
}

// ConfigDeployment builds the name for the config deployment used in the apollo-operator.
func ConfigDeployment(obj client.Object) string {
	return DNSName(Truncate("%s-config-deployment", 63, obj.GetName()))
}

// AdminService builds the name for the admin service used in the apollo-operator.
func AdminService(obj client.Object) string {
	return DNSName(Truncate("%s-admin", 63, obj.GetName()))
}

// ConfigService builds the name for the config service used in the apollo-operator.
func ConfigService(obj client.Object) string {
	return DNSName(Truncate("%s-config", 63, obj.GetName()))
}

// ConfigDBService builds the name for the config db service used in the apollo-operator.
func ConfigDBService(obj client.Object) string {
	return DNSName(Truncate("%s-configdb", 63, obj.GetName()))
}

// AdminIngress builds the name for the admin ingress used in the apollo-operator.
func AdminIngress(obj client.Object) string {
	return DNSName(Truncate("%s-admin-ingress", 63, obj.GetName()))
}

// ConfigIngress builds the name for the config ingress used in the apollo-operator.
func ConfigIngress(obj client.Object) string {
	return DNSName(Truncate("%s-config-ingress", 63, obj.GetName()))
}

/* Apollo Portal */

// ConfigMap builds the name for the portal config map used in the apollo-operator.
func ConfigMap(obj client.Object) string {
	return DNSName(Truncate("%s-portal-configmap", 63, obj.GetName()))
}

// PortalDeployment builds the name for the portal deployment used in the apollo-operator.
func PortalDeployment(obj client.Object) string {
	return DNSName(Truncate("%s-portal-deployment", 63, obj.GetName()))
}

// PortalDBService builds the name for the portal db service used in the apollo-operator.
func PortalDBService(obj client.Object) string {
	return DNSName(Truncate("%s-portaldb", 63, obj.GetName()))
}

// PortalService builds the name for the portal service used in the apollo-operator.
func PortalService(obj client.Object) string {
	return DNSName(Truncate("%s-portal", 63, obj.GetName()))
}

// PortalIngress builds the name for the portal ingress used in the apollo-operator.
func PortalIngress(obj client.Object) string {
	return DNSName(Truncate("%s-portal-ingress", 63, obj.GetName()))
}

/* Apollo all in one */

// AllInOneAdminConfigMap builds the name for the admin configmap used in the apollo-operator.
func AllInOneAdminConfigMap(obj client.Object) string {
	return DNSName(Truncate("%s-admin-configmap-allinone", 63, obj.GetName()))
}

// AllInOneConfigConfigMap builds the name for the config configmap used in the apollo-operator.
func AllInOneConfigConfigMap(obj client.Object) string {
	return DNSName(Truncate("%s-config-configmap-allinone", 63, obj.GetName()))
}

// AllInOneAdminDeployment builds the name for the admin deployment used in the apollo-operator.
func AllInOneAdminDeployment(obj client.Object) string {
	return DNSName(Truncate("%s-admin-deployment-allinone", 63, obj.GetName()))
}

// AllInOneConfigDeployment builds the name for the config deployment used in the apollo-operator.
func AllInOneConfigDeployment(obj client.Object) string {
	return DNSName(Truncate("%s-config-deployment-allinone", 63, obj.GetName()))
}

// AllInOneAdminService builds the name for the admin service used in the apollo-operator.
func AllInOneAdminService(obj client.Object) string {
	return DNSName(Truncate("%s-admin-allinone", 63, obj.GetName()))
}

// AllInOneConfigService builds the name for the config service used in the apollo-operator.
func AllInOneConfigService(obj client.Object) string {
	return DNSName(Truncate("%s-config-allinone", 63, obj.GetName()))
}

// AllInOneConfigDBService builds the name for the config db service used in the apollo-operator.
func AllInOneConfigDBService(obj client.Object) string {
	return DNSName(Truncate("%s-configdb-allinone", 63, obj.GetName()))
}

// AllInOneAdminIngress builds the name for the admin ingress used in the apollo-operator.
func AllInOneAdminIngress(obj client.Object) string {
	return DNSName(Truncate("%s-admin-ingress-allinone", 63, obj.GetName()))
}

// AllInOneConfigIngress builds the name for the config ingress used in the apollo-operator.
func AllInOneConfigIngress(obj client.Object) string {
	return DNSName(Truncate("%s-config-ingress-allinone", 63, obj.GetName()))
}

// AllInOneConfigMap builds the name for the portal config map used in the apollo-operator.
func AllInOneConfigMap(obj client.Object) string {
	return DNSName(Truncate("%s-portal-configmap-allinone", 63, obj.GetName()))
}

// AllInOnePortalDeployment builds the name for the portal deployment used in the apollo-operator.
func AllInOnePortalDeployment(obj client.Object) string {
	return DNSName(Truncate("%s-portal-deployment-allinone", 63, obj.GetName()))
}

// AllInOnePortalDBService builds the name for the portal db service used in the apollo-operator.
func AllInOnePortalDBService(obj client.Object) string {
	return DNSName(Truncate("%s-portaldb-allinone", 63, obj.GetName()))
}

// AllInOnePortalService builds the name for the portal service used in the apollo-operator.
func AllInOnePortalService(obj client.Object) string {
	return DNSName(Truncate("%s-portal-allinone", 63, obj.GetName()))
}

// AllInOneStatefulSet builds the name for the apollo allinone statefulset used in the apollo-operator.
func AllInOneStatefulSet(obj client.Object) string {
	return DNSName(Truncate("%s-statefulset-allinone", 63, obj.GetName()))
}

// AllInOneDBService builds the name for the allinone db service used in the apollo-operator.
func AllInOneDBService(obj client.Object) string {
	return DNSName(Truncate("%s-db-allinone", 63, obj.GetName()))
}

// AllInOnePVC builds the name for the apollo allinone pvc used in the apollo-operator.
func AllInOnePVC(obj client.Object) string {
	return DNSName(Truncate("%s-apolloDB-allinone", 63, obj.GetName()))
}

// AllInOneSqlScript builds the name for the apollo allinone sql script used in the apollo-operator.
func AllInOneSqlScript(obj client.Object) string {
	return DNSName(Truncate("%s-apolloDB-sqlscript", 63, obj.GetName()))
}

// AllInOneSqlScriptConfigmap builds the name for the apollo allinone sql script configmap used in the apollo-operator.
func AllInOneSqlScriptConfigmap(_ client.Object) string {
	return "mysql-initdb-config" // NOTE 包含初始化sql语句的configmap名字
}

/* Public name generation  */

// HeadlessService builds the name for the headless service used in the apollo-operator.
func HeadlessService(obj client.Object) string {
	return DNSName(Truncate("%s-headless", 63, Service(obj)))
}

// Service builds the name for the service used in the apollo-operator.
func Service(obj client.Object) string {
	return DNSName(Truncate("%s-service", 63, obj.GetName()))
}

// ServiceAccount builds the service account name based on the instance.
func ServiceAccount(obj client.Object) string {
	return DNSName(Truncate("%s-serviceaccount", 63, obj.GetName()))
}

// Apollo builds the apollo resource name used in the apollo-operator.
func Apollo(obj client.Object) string {
	return DNSName(Truncate("%s", 63, obj.GetName()))
}

// Container returns the name to use for the container in the pod.
func Container() string {
	return "apollo-container"
}

// InitContainer returns the name to use for the container in the pod.
func InitContainer() string {
	return "apollo-init-container"
}

// ResourceNameWithSuffix builds the resource name based on the instance and suffix.
func ResourceNameWithSuffix(obj client.Object, suffix string) string {
	if suffix == "" {
		return Apollo(obj)
	}
	return DNSName(Truncate("%s-%s", 63, obj.GetName(), suffix))
}
