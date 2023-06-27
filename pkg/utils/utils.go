package utils

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

func MergeTwoMap(m1 map[string]string, m2 map[string]string) (res map[string]string) {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}

// InitObjectMeta will set the required default settings to
// kubernetes objects metadata if is required.
func InitObjectMeta(obj metav1.Object) {
	if obj.GetLabels() == nil {
		obj.SetLabels(map[string]string{})
	}

	if obj.GetAnnotations() == nil {
		obj.SetAnnotations(map[string]string{})
	}
}
