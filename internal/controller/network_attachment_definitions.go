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

package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	workloadnephioorgv1alpha1 "workload.nephio.org/ran_deployment/api/v1alpha1"
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

func CreateNetworkAttachmentDefinitionNetworksConfigs(templateName string, interfaceConfigs map[string][]workloadnephioorgv1alpha1.InterfaceConfig) (string, error) {
	//TODO: will be removed and common functions defined for free5gc can be used when using NFdeploy
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

// Gets a Deployment resource and checks that the NetworkAttachmentDefinitions specified in its
// `k8s.v1.cni.cncf.io/networks` annotation exist in the same namespace.
// TODO: will be removed and common functions defined for free5gc can be used when using NFdeploy
func ValidateNetworkAttachmentDefinitions(ctx context.Context, c client.Client, log logr.Logger, kind string, deployment *appsv1.Deployment) bool {
	networksJson, ok := deployment.Spec.Template.Annotations[NetworksAnnotation]
	if !ok {
		log.Info(fmt.Sprintf("Annotation %q not found", NetworksAnnotation), kind+".namespace", deployment.Namespace)
		return false
	}

	var networks []networkAttachmentDefinitionNetwork
	if err := json.Unmarshal([]byte(networksJson), &networks); err != nil {
		log.Error(err, fmt.Sprintf("Failed to parse %q annotation", kind), kind+".namespace", deployment.Namespace)
		return false
	}

	for _, network := range networks {
		var u unstructured.Unstructured
		u.SetGroupVersionKind(NetworkAttachmentDefinitionGVK)
		key := client.ObjectKey{Namespace: deployment.Namespace, Name: network.Name}
		if err := c.Get(ctx, key, &u); err != nil {
			log.Error(err, fmt.Sprintf("Failed to get NetworkAttachmentDefinition %q", network.Name), kind+".namespace", deployment.Namespace)
			return false
		}
	}

	return true
}
