// Package naming is for determining the names for components (containers, services, ...).
package naming

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ConfigMap builds the name for the config map used in the Apollo.
func ConfigMap(obj client.Object) string {
	return DNSName(Truncate("%s-configmap", 63, obj.GetName()))
}

// Endpoints builds the name for the endpoints used in the Apollo.
func Endpoints(obj client.Object) string {
	return DNSName(Truncate("%s-endpoints", 63, obj.GetName()))
}

// Deployment builds the name for the deployment used in the Apollo.
func Deployment(obj client.Object) string {
	return DNSName(Truncate("%s-deployment", 63, obj.GetName()))
}

// Apollo builds the apollo resource name based on the instance.
func Apollo(obj client.Object) string {
	return DNSName(Truncate("%s", 63, obj.GetName()))
}

// Container returns the name to use for the container in the pod.
func Container() string {
	return "apollo-container"
}

// HeadlessService builds the name for the headless service based on the instance.
func HeadlessService(obj client.Object) string {
	return DNSName(Truncate("%s-headless", 63, Service(obj)))
}

// Service builds the name for the service used in the Apollo.
func Service(obj client.Object) string {
	return DNSName(Truncate("%s-service", 63, obj.GetName()))
}

// PortalDBService builds the name for the portal db service used in the Apollo.
func PortalDBService(obj client.Object) string {
	return DNSName(Truncate("%s-portaldb", 63, obj.GetName()))
}

// Ingress builds the ingress name based on the instance.
func Ingress(obj client.Object) string {
	return DNSName(Truncate("%s-ingress", 63, obj.GetName()))
}

// ServiceAccount builds the service account name based on the instance.
func ServiceAccount(obj client.Object) string {
	return DNSName(Truncate("%s-serviceaccount", 63, obj.GetName()))
}

// ResourceNameWithSuffix builds the resource name based on the instance.
func ResourceNameWithSuffix(obj client.Object, suffix string) string {
	if suffix == "" {
		return Apollo(obj)
	}
	return DNSName(Truncate("%s-%s", 63, obj.GetName(), suffix))
}
