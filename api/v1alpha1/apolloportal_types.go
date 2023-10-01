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
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ApolloPortalSpec defines the desired state of ApolloPortal
type ApolloPortalSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Image string `json:"image,omitempty" default:"apolloconfig/apollo-portal:2.1.0"`

	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	Replicas int32 `json:"replicas,omitempty" default:"1"`

	ContainerPort int32 `json:"containerPort,omitempty"`

	Strategy appv1.DeploymentStrategy `json:"strategy,omitempty"`

	Env []corev1.EnvVar `json:"env,omitempty"`

	Service Service `json:"service,omitempty"`

	Config PortalConfig `json:"config,omitempty"`

	PortalDB PortalDB `json:"portaldb,omitempty"`

	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	Probe Probe `json:"probe,omitempty"`

	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	Affinity corev1.Affinity `json:"affinity,omitempty"`

	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// Ingress is used to specify how ApolloPortal is exposed.
	// +optional
	Ingress Ingress `json:"ingress,omitempty"`
}

type Service struct {
	Port       int32              `json:"port,omitempty"`
	TargetPort int32              `json:"targetPort,omitempty"`
	Type       corev1.ServiceType `json:"type,omitempty"`
	// TODO Follow up to see if necessary, delete if not necessary
	SessionAffinity corev1.ServiceAffinity `json:"sessionAffinity,omitempty"`
}

type PortalConfig struct {
	Envs        string            `json:"envs,omitempty"`
	MetaServers map[string]string `json:"metaServers,omitempty"`
	Profiles    string            `json:"profiles,omitempty"`
	ContextPath string            `json:"contextPath,omitempty"`
	Files       []File            `json:"file,omitempty"`
}

type PortalDB struct {
	Username                   string          `json:"username,omitempty"`
	Password                   string          `json:"password,omitempty"`
	Host                       string          `json:"host,omitempty"`
	Port                       int32           `json:"port,omitempty"`
	DBName                     string          `json:"dbName,omitempty"`
	ConnectionStringProperties string          `json:"connectionStringProperties,omitempty"`
	Service                    PortalDBService `json:"service,omitempty"`
}

type PortalDBService struct {
	Name string             `json:"name,omitempty"`
	Port int32              `json:"port,omitempty"`
	Type corev1.ServiceType `json:"type,omitempty"`
}

type File struct {
	Name    string `json:"name,omitempty"`
	Content string `json:"content,omitempty"`
}

type Probe struct {
	Liveness   corev1.Probe `json:"livenessProbe,omitempty"`
	Readineeds corev1.Probe `json:"readinessProbe,omitempty"`
}

type Ingress struct {

	// IngressClassName is the name of an IngressClass cluster resource. Ingress
	// controller implementations use this field to know whether they should be
	// serving this Ingress resource.
	// +optional
	IngressClassName *string `json:"ingressClassName,omitempty"`

	// Annotations to add to ingress.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// HTTP Host URL
	// +optional
	Hosts []string `json:"hosts,omitempty"`

	// TLS configuration.
	// +optional
	TLS []networkingv1.IngressTLS `json:"tls,omitempty"`
}

// ApolloPortalStatus defines the observed state of ApolloPortal
type ApolloPortalStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ApolloPortal is the Schema for the apolloportals API
type ApolloPortal struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApolloPortalSpec   `json:"spec,omitempty"`
	Status ApolloPortalStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ApolloPortalList contains a list of ApolloPortal
type ApolloPortalList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApolloPortal `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ApolloPortal{}, &ApolloPortalList{})
}
