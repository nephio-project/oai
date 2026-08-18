/*
Copyright 2026 The Nephio Authors.

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

// Package v1alpha1 holds the OAI RAN configuration payloads (RANConfig,
// OAIConfig, PLMN) that this operator decodes out of NFConfig.spec.configRefs.
//
// They are not independently served Kubernetes resources. Nothing registers
// them with a scheme,
// nothing fetches them through a client, and the API server serves the
// aggregate NFConfig CRD from nephio-project/api instead of one CRD per
// configuration entity. kubebuilder:skip tells controller-gen the same thing:
// generate deepcopy for these types, but do not treat the package as an API
// version and do not emit CRDs for it.
//
// +kubebuilder:object:generate=true
// +kubebuilder:skip
package v1alpha1
