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
	"reflect"

	v1alpha1 "github.com/nephio-project/api/nf_deployments/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +kubebuilder:validation:Enum=GNB;CU-CP;CU-UP;DU
type RanNfType string

type Plmn struct {
	Mcc       string `json:"mcc,omitempty"`
	Mnc       string `json:"mnc,omitempty"`
	MncLength int    `json:"mncLength,omitempty"`
}

type Nssai struct {
	Sst string `json:"sst,omitempty"`
	Sd  string `json:"sd,omitempty"`
}

type Params3gpp struct {
	//physicalCellId defines the physical cell identity of a cell
	PhysicalCellId int `json:"physicalCellId,omitempty"`
	//cellIdentity defines the cell identity of a cell
	CellIdentity string `json:"cellIdentity,omitempty"`
	//plmn defines the plmn of a cell
	Plmn `json:"plmn,omitempty"`
	//tac defines the tracking area code to be used by the cell
	Tac string `json:"tac,omitempty"`
	//nssaiList defines the Nssai list to be configured for the cell
	NssaiList []Nssai `json:"nssaiList,omitempty"`
}

// RANDeployment is the Schema for the randeployments API
type RANDeploymentSpec struct {
	//ranNfType defines the ranNfType of RAN network function,
	// GNB/CU-CP/CU-UP/DU.
	RanNfType  `json:"ranNfType"`
	Params3gpp `json:"params3gpp"`
	// latency maximum latency tolerated by this RAN NF. This variable will be responsible for latency of this NF
	NfLatency                 int `json:"nfLatency,omitempty"`
	v1alpha1.NFDeploymentSpec `json:",inline" yaml:",inline"`
}

// RANDeploymentStatus defines the observed state of RANDeployment
type RANDeploymentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	v1alpha1.NFDeploymentStatus `json:",inline" yaml:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RANDeployment is the Schema for the randeployments API
type RANDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RANDeploymentSpec   `json:"spec,omitempty"`
	Status RANDeploymentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RANDeploymentList contains a list of RANDeployment
type RANDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RANDeployment `json:"items"`
}

// Implement NFDeployment interface

func (d *RANDeployment) GetNFDeploymentSpec() *v1alpha1.NFDeploymentSpec {
	return d.Spec.NFDeploymentSpec.DeepCopy()
}
func (d *RANDeployment) GetNFDeploymentStatus() *v1alpha1.NFDeploymentStatus {
	return d.Status.NFDeploymentStatus.DeepCopy()
}
func (d *RANDeployment) SetNFDeploymentSpec(s *v1alpha1.NFDeploymentSpec) {
	s.DeepCopyInto(&d.Spec.NFDeploymentSpec)
}
func (d *RANDeployment) SetNFDeploymentStatus(s *v1alpha1.NFDeploymentStatus) {
	s.DeepCopyInto(&d.Status.NFDeploymentStatus)
}

func init() {
	SchemeBuilder.Register(&RANDeployment{}, &RANDeploymentList{})
}

// Interface type metadata.
var (
	RANDeploymentKind             = reflect.TypeOf(RANDeployment{}).Name()
	RANDeploymentGroupKind        = schema.GroupKind{Group: Group, Kind: RANDeploymentKind}.String()
	RANDeploymentKindAPIVersion   = RANDeploymentKind + "." + GroupVersion.String()
	RANDeploymentGroupVersionKind = GroupVersion.WithKind(RANDeploymentKind)
)
