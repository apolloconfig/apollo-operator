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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ApolloPortalSpec defines the desired state of ApolloPortal
type ApolloPortalSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Image string `json:"image,omitempty"`

	Service Service `json:"service,omitempty"`

	Config Config `json:"config,omitempty"`

	PortalDB PortalDB `json:"portaldb,omitempty"`
}

type Service struct {
	Port            int32                  `json:"port"`
	TargetPort      int32                  `json:"targetPort"`
	Type            corev1.ServiceType     `json:"type"`
	SessionAffinity corev1.ServiceAffinity `json:"sessionAffinity"`
}

type Config struct {
	Envs        string            `json:"envs,omitempty"`
	MetaServers map[string]string `json:"metaServers,omitempty"`
	File        map[string]string `json:"file,omitempty"`
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
