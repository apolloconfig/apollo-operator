package utils

import "sigs.k8s.io/controller-runtime/pkg/client"

// Annotations return the annotations for Apollo operator.
func Annotations(instance client.Object) map[string]string {
	// new map every time, so that we don't touch the instance's annotations
	annotations := map[string]string{}

	// TODO 增加默认注解
	// set default apollo operator annotations
	annotations["apollo.io/apollo-portal/port"] = "8070"

	// allow override of prometheus annotations
	if instance.GetAnnotations() != nil {
		for k, v := range instance.GetAnnotations() {
			annotations[k] = v
		}
	}
	return annotations
}
