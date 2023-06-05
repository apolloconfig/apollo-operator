package k8sClient

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

type IK8sClient interface {
}

type K8sClient struct {
	KubeClient kubernetes.Interface
	scheme     *runtime.Scheme
}

func NewK8sClient(kubeClient kubernetes.Interface, scheme *runtime.Scheme) *K8sClient {
	return &K8sClient{
		KubeClient: kubeClient,
		scheme:     scheme,
	}
}
