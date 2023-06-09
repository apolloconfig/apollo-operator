package k8sClient

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

type k8sClient struct {
	KubeClient kubernetes.Interface
	scheme     *runtime.Scheme
}

func NewK8sClient(kubeClient kubernetes.Interface, scheme *runtime.Scheme) *k8sClient {
	return &k8sClient{
		KubeClient: kubeClient,
		scheme:     scheme,
	}
}
