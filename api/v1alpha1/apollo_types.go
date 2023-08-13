/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ApolloSpec defines the desired state of Apollo
type ApolloSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	ConfigService ConfigService `json:"configService,omitempty"`

	AdminService AdminService `json:"adminService,omitempty"`

	PortalService PortalService `json:"portalService,omitempty"`
}

type PortalService struct {
	Image string `json:"image,omitempty"`

	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	Replicas int32 `json:"replicas,omitempty"`

	ContainerPort int32 `json:"containerPort,omitempty"`

	Strategy appv1.DeploymentStrategy `json:"strategy,omitempty"`

	Env []corev1.EnvVar `json:"env,omitempty"`

	Service Service `json:"service,omitempty"`

	Config PortalServiceConfig `json:"config,omitempty"`

	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	Probe Probe `json:"probe,omitempty"`

	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	Affinity corev1.Affinity `json:"affinity,omitempty"`

	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// Ingress is used to specify how ApolloAdmin is exposed.
	// +optional
	Ingress Ingress `json:"ingress,omitempty"`
}

type PortalServiceConfig struct {
	Envs        string `json:"envs,omitempty"`
	Profiles    string `json:"profiles,omitempty"`
	ContextPath string `json:"contextPath,omitempty"`
	Files       []File `json:"file,omitempty"`
}

// ApolloStatus defines the observed state of Apollo
type ApolloStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Apollo is the Schema for the apolloes API
type Apollo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApolloSpec   `json:"spec,omitempty"`
	Status ApolloStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ApolloList contains a list of Apollo
type ApolloList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Apollo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Apollo{}, &ApolloList{})
}
