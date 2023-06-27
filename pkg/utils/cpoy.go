package utils

import "sigs.k8s.io/controller-runtime/pkg/client"

func CopyAnnotationsLabels(from, to client.Object) {
	to.SetAnnotations(from.GetAnnotations())
	to.SetLabels(from.GetLabels())
}
