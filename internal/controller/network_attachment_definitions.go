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

// TODO: These functions are similar to functions used in free5gc, need to part of common repository so that all NFs can refer.

package controller

import (
	"fmt"
	"sort"
	"strings"

	workloadv1alpha1 "github.com/nephio-project/api/workload/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const NetworksAnnotation = "k8s.v1.cni.cncf.io/networks"

var NetworkAttachmentDefinitionGVK = schema.GroupVersionKind{
	Group:   "k8s.cni.cncf.io",
	Kind:    "NetworkAttachmentDefinition",
	Version: "v1",
}

type networkAttachmentDefinitionNetwork struct {
	Name      string `json:"name"`
	Interface string `json:"interface"`
	IP        string `json:"ip"`
	Gateway   string `json:"gateway"`
}

func CreateNetworkAttachmentDefinitionNetworks(templateName string, interfaceConfigs map[string][]workloadv1alpha1.InterfaceConfig) (string, error) {

	interfaceNames := make([]string, 0, len(interfaceConfigs))
	for interfaceName := range interfaceConfigs {
		interfaceNames = append(interfaceNames, interfaceName)
	}

	sort.Strings(interfaceNames) // ensure consistent return value for unit tests

	var networksJson []string
	for _, interfaceName := range interfaceNames {
		for _, interfaceConfig := range interfaceConfigs[interfaceName] {
			if interfaceConfig.IPv4.Gateway == nil {
				return "", fmt.Errorf("missing `InterfaceConfig.IPv4.Gateway` for %q", interfaceName)
			}

			networksJson = append(networksJson, fmt.Sprintf(` {
  "name": %q,
  "interface": %q,
  "ips": [%q],
  "gateways": [%q]
 }`,
				CreateNetworkAttachmentDefinitionName(templateName, interfaceName),
				interfaceConfig.Name,
				interfaceConfig.IPv4.Address,
				*interfaceConfig.IPv4.Gateway))
		}
	}

	return "[\n" + strings.Join(networksJson, ",\n") + "\n]", nil
}

func CreateNetworkAttachmentDefinitionName(templateName string, suffix string) string {
	return templateName + "-" + suffix
}
