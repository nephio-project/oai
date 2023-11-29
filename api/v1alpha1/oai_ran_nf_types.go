/*
Copyright 2023 The Nephio Authors.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OAIConfigSpec defines the desired state of OAIConfig
type OAIConfigSpec struct {
	//image defines the image location for the OAI NF
	Image string `json:"image"`
}

// OAIConfigStatus defines the observed state of OAIConfig
type OAIConfigStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// OAIConfig is the Schema for the OAIConfigs API
type OAIConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OAIConfigSpec   `json:"spec,omitempty"`
	Status OAIConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OAIConfigList contains a list of OAIConfig
type OAIConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OAIConfig `json:"items"`
}
