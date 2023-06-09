package reconcile

import (
	apolloiov1alpha1 "apollo.io/apollo-operator/api/v1alpha1"
	"apollo.io/apollo-operator/pkg/utils/naming"
	"regexp"
	"strings"
)

func isFilteredLabel(label string, filterLabels []string) bool {
	for _, pattern := range filterLabels {
		match, _ := regexp.MatchString(pattern, label)
		return match
	}
	return false
}

// Labels return the common labels to all objects that are part of a managed ApolloPortal.
func Labels(instance apolloiov1alpha1.ApolloPortal, name string, filterLabels []string) map[string]string {
	// new map every time, so that we don't touch the instance's label
	base := map[string]string{}
	if nil != instance.Labels {
		for k, v := range instance.Labels {
			if !isFilteredLabel(k, filterLabels) {
				base[k] = v
			}
		}
	}

	for k, v := range SelectorLabels(instance) {
		base[k] = v
	}

	version := strings.Split(instance.Spec.Image, ":")
	if len(version) > 1 {
		base["app.kubernetes.io/apollo-version"] = version[len(version)-1]
	} else {
		base["app.kubernetes.io/apollo-version"] = "latest"
	}

	// Don't override the app name if it already exists
	if _, ok := base["app.kubernetes.io/name"]; !ok {
		base["app.kubernetes.io/name"] = name
	}

	return base
}

// SelectorLabels return the common labels to all objects that are part of a managed OpenTelemetryCollector to use as selector.
// Selector labels are immutable for Deployment, StatefulSet and DaemonSet, therefore, no labels in selector should be
// expected to be modified for the lifetime of the object.
func SelectorLabels(instance apolloiov1alpha1.ApolloPortal) map[string]string {
	return map[string]string{
		"app.kubernetes.io/managed-by": "apollo-portal-operator",
		"app.kubernetes.io/instance":   naming.Truncate("%s.%s", 63, instance.Namespace, instance.Name),
		"app.kubernetes.io/part-of":    "apollo",
		"app.kubernetes.io/component":  "apollo-portal-operator",
	}
}
