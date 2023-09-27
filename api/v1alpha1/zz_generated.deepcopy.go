//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	nf_requirementsv1alpha1 "github.com/nephio-project/api/nf_requirements/v1alpha1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BGPConfig) DeepCopyInto(out *BGPConfig) {
	*out = *in
	if in.BGPNeigbors != nil {
		in, out := &in.BGPNeigbors, &out.BGPNeigbors
		*out = make([]BGPNeighbor, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BGPConfig.
func (in *BGPConfig) DeepCopy() *BGPConfig {
	if in == nil {
		return nil
	}
	out := new(BGPConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BGPNeighbor) DeepCopyInto(out *BGPNeighbor) {
	*out = *in
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BGPNeighbor.
func (in *BGPNeighbor) DeepCopy() *BGPNeighbor {
	if in == nil {
		return nil
	}
	out := new(BGPNeighbor)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DataNetwork) DeepCopyInto(out *DataNetwork) {
	*out = *in
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.Pool != nil {
		in, out := &in.Pool, &out.Pool
		*out = make([]Pool, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataNetwork.
func (in *DataNetwork) DeepCopy() *DataNetwork {
	if in == nil {
		return nil
	}
	out := new(DataNetwork)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IPv4) DeepCopyInto(out *IPv4) {
	*out = *in
	if in.Gateway != nil {
		in, out := &in.Gateway, &out.Gateway
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IPv4.
func (in *IPv4) DeepCopy() *IPv4 {
	if in == nil {
		return nil
	}
	out := new(IPv4)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IPv6) DeepCopyInto(out *IPv6) {
	*out = *in
	if in.Gateway != nil {
		in, out := &in.Gateway, &out.Gateway
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IPv6.
func (in *IPv6) DeepCopy() *IPv6 {
	if in == nil {
		return nil
	}
	out := new(IPv6)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InterfaceConfig) DeepCopyInto(out *InterfaceConfig) {
	*out = *in
	if in.IPv4 != nil {
		in, out := &in.IPv4, &out.IPv4
		*out = new(IPv4)
		(*in).DeepCopyInto(*out)
	}
	if in.IPv6 != nil {
		in, out := &in.IPv6, &out.IPv6
		*out = new(IPv6)
		(*in).DeepCopyInto(*out)
	}
	if in.VLANID != nil {
		in, out := &in.VLANID, &out.VLANID
		*out = new(uint16)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InterfaceConfig.
func (in *InterfaceConfig) DeepCopy() *InterfaceConfig {
	if in == nil {
		return nil
	}
	out := new(InterfaceConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NFDeploymentSpec) DeepCopyInto(out *NFDeploymentSpec) {
	*out = *in
	if in.Capacity != nil {
		in, out := &in.Capacity, &out.Capacity
		*out = new(nf_requirementsv1alpha1.CapacitySpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Interfaces != nil {
		in, out := &in.Interfaces, &out.Interfaces
		*out = make([]InterfaceConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.NetworkInstances != nil {
		in, out := &in.NetworkInstances, &out.NetworkInstances
		*out = make([]NetworkInstance, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ConfigRefs != nil {
		in, out := &in.ConfigRefs, &out.ConfigRefs
		*out = make([]v1.ObjectReference, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NFDeploymentSpec.
func (in *NFDeploymentSpec) DeepCopy() *NFDeploymentSpec {
	if in == nil {
		return nil
	}
	out := new(NFDeploymentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NFDeploymentStatus) DeepCopyInto(out *NFDeploymentStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]metav1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NFDeploymentStatus.
func (in *NFDeploymentStatus) DeepCopy() *NFDeploymentStatus {
	if in == nil {
		return nil
	}
	out := new(NFDeploymentStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkInstance) DeepCopyInto(out *NetworkInstance) {
	*out = *in
	if in.Interfaces != nil {
		in, out := &in.Interfaces, &out.Interfaces
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Peers != nil {
		in, out := &in.Peers, &out.Peers
		*out = make([]PeerConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.DataNetworks != nil {
		in, out := &in.DataNetworks, &out.DataNetworks
		*out = make([]DataNetwork, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.BGP != nil {
		in, out := &in.BGP, &out.BGP
		*out = new(BGPConfig)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetworkInstance.
func (in *NetworkInstance) DeepCopy() *NetworkInstance {
	if in == nil {
		return nil
	}
	out := new(NetworkInstance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PeerConfig) DeepCopyInto(out *PeerConfig) {
	*out = *in
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.IPv4 != nil {
		in, out := &in.IPv4, &out.IPv4
		*out = new(IPv4)
		(*in).DeepCopyInto(*out)
	}
	if in.IPv6 != nil {
		in, out := &in.IPv6, &out.IPv6
		*out = new(IPv6)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PeerConfig.
func (in *PeerConfig) DeepCopy() *PeerConfig {
	if in == nil {
		return nil
	}
	out := new(PeerConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Pool) DeepCopyInto(out *Pool) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Pool.
func (in *Pool) DeepCopy() *Pool {
	if in == nil {
		return nil
	}
	out := new(Pool)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RANDeployment) DeepCopyInto(out *RANDeployment) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RANDeployment.
func (in *RANDeployment) DeepCopy() *RANDeployment {
	if in == nil {
		return nil
	}
	out := new(RANDeployment)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RANDeployment) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RANDeploymentList) DeepCopyInto(out *RANDeploymentList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]RANDeployment, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RANDeploymentList.
func (in *RANDeploymentList) DeepCopy() *RANDeploymentList {
	if in == nil {
		return nil
	}
	out := new(RANDeploymentList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RANDeploymentList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RANDeploymentSpec) DeepCopyInto(out *RANDeploymentSpec) {
	*out = *in
	in.NFDeploymentSpec.DeepCopyInto(&out.NFDeploymentSpec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RANDeploymentSpec.
func (in *RANDeploymentSpec) DeepCopy() *RANDeploymentSpec {
	if in == nil {
		return nil
	}
	out := new(RANDeploymentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RANDeploymentStatus) DeepCopyInto(out *RANDeploymentStatus) {
	*out = *in
	in.NFDeploymentStatus.DeepCopyInto(&out.NFDeploymentStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RANDeploymentStatus.
func (in *RANDeploymentStatus) DeepCopy() *RANDeploymentStatus {
	if in == nil {
		return nil
	}
	out := new(RANDeploymentStatus)
	in.DeepCopyInto(out)
	return out
}
