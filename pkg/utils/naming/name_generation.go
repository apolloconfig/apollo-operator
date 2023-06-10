// Package naming is for determining the names for components (containers, services, ...).
package naming

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ConfigMap builds the name for the config map used in the ApolloPortal containers.
func ConfigMap(obj client.Object) string {
	return DNSName(Truncate("%s-config", 63, obj.GetName()))
}

// Endpoints builds the name for the endpoints used in the ApolloPortalDB Service.
func Endpoints(obj client.Object) string {
	return DNSName(Truncate("%s-endpoints", 63, obj.GetName()))
}

// Apollo builds the collector (deployment/daemonset) name based on the instance.
func Apollo(obj client.Object) string {
	return DNSName(Truncate("%s-apollo", 63, obj.GetName()))
}

// HeadlessService builds the name for the headless service based on the instance.
func HeadlessService(obj client.Object) string {
	return DNSName(Truncate("%s-headless", 63, Service(obj)))
}

// Service builds the service name based on the instance.
func Service(obj client.Object) string {
	return DNSName(Truncate("%s-service", 63, obj.GetName()))
}

// Ingress builds the ingress name based on the instance.
func Ingress(obj client.Object) string {
	return DNSName(Truncate("%s-ingress", 63, obj.GetName()))
}

// ServiceAccount builds the service account name based on the instance.
func ServiceAccount(obj client.Object) string {
	return DNSName(Truncate("%s-serviceaccount", 63, obj.GetName()))
}

// ResourceNameWithSuffix builds the service name based on the instance.
func ResourceNameWithSuffix(obj client.Object, suffix string) string {
	return DNSName(Truncate("%s-%s", 63, obj.GetName(), suffix))
}
