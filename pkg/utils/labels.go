package utils

import (
	"apollo.io/apollo-operator/pkg/utils/naming"
	"regexp"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
func Labels(instance client.Object, name string, filterLabels []string) map[string]string {
	// new map every time, so that we don't touch the instance's label
	base := map[string]string{}
	if nil != instance.GetLabels() {
		for k, v := range instance.GetLabels() {
			if !isFilteredLabel(k, filterLabels) {
				base[k] = v
			}
		}
	}

	for k, v := range SelectorLabels(instance) {
		base[k] = v
	}

	// Don't override the app name if it already exists
	if _, ok := base["app.kubernetes.io/name"]; !ok {
		base["app.kubernetes.io/name"] = name
	}

	return base
}

// SelectorLabels return the common labels to all objects that are part of a managed Apollo Operator to use as selector.
// Selector labels are immutable for Deployment, StatefulSet and DaemonSet, therefore, no labels in selector should be
// expected to be modified for the lifetime of the object.
func SelectorLabels(instance client.Object) map[string]string {
	// 如果你修改这里，那么所有资源的delete部分中，ListOption都要同步修改
	return map[string]string{
		"app.kubernetes.io/managed-by": "apollo-operator",
		"app.kubernetes.io/instance":   naming.Truncate("%s.%s", 63, instance.GetNamespace(), instance.GetName()),
		"app.kubernetes.io/part-of":    "apollo-operator",
		"app.kubernetes.io/component":  strings.ToLower(instance.GetObjectKind().GroupVersionKind().Kind), // eg. apolloportal
	}
}

func SelectorLabelsWithCustom(instance client.Object, custom map[string]string) map[string]string {
	commonLabels := SelectorLabels(instance)
	return MergeTwoMap(commonLabels, custom)
}
