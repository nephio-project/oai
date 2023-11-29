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

// RANNfConfigSpec defines the desired state of RANNfConfig
type RANConfigSpec struct {
	//cellIdentity defines the cell identity of a cell
	CellIdentity string `json:"cellIdentity"`
	//physicalCellId defines the physical cell identity of a cell
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=503
	PhysicalCellID            uint32 `json:"physicalCellID"`
	DownlinkFrequencyBand     uint32 `json:"downlinkFrequencyBand"`
	DownlinkSubCarrierSpacing uint16 `json:"downlinkSubCarrierSpacing"`
	DownlinkCarrierBandwidth  uint32 `json:"downlinkCarrierBandwidth"`
	UplinkFrequencyBand       uint32 `json:"uplinkFrequencyBand"`
	UplinkSubCarrierSpacing   uint16 `json:"uplinkSubCarrierSpacing"`
	UplinkCarrierBandwidth    uint32 `json:"uplinkCarrierBandwidth"`
}

// RANConfigStatus defines the observed state of RANConfig
type RANConfigStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RANConfig is the Schema for the RANconfigs API
type RANConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RANConfigSpec   `json:"spec,omitempty"`
	Status RANConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RANConfigList contains a list of RANConfig
type RANConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RANConfig `json:"items"`
}
