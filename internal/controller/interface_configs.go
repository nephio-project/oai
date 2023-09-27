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
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	nephiov1alpha1 "github.com/nephio-project/api/nf_deployments/v1alpha1"
	configref "github.com/nephio-project/api/references/v1alpha1"
	free5gccontrollers "github.com/nephio-project/free5gc/controllers"
	"k8s.io/apimachinery/pkg/runtime/schema"
	workloadnephioorgv1alpha1 "workload.nephio.org/ran_deployment/api/v1alpha1"
)

func GetInterfaceConfigs(interfaceConfigs []workloadnephioorgv1alpha1.InterfaceConfig, interfaceName string) []workloadnephioorgv1alpha1.InterfaceConfig {
	var selectedInterfaceConfigs []workloadnephioorgv1alpha1.InterfaceConfig

	for _, interfaceConfig := range interfaceConfigs {
		if interfaceConfig.Name == interfaceName {
			selectedInterfaceConfigs = append(selectedInterfaceConfigs, interfaceConfig)
		}
	}

	return selectedInterfaceConfigs
}

func GetFirstInterfaceConfig(interfaceConfigs []workloadnephioorgv1alpha1.InterfaceConfig, interfaceName string) (*workloadnephioorgv1alpha1.InterfaceConfig, error) {
	for _, interfaceConfig := range interfaceConfigs {
		if interfaceConfig.Name == interfaceName {
			return &interfaceConfig, nil
		}
	}

	return nil, fmt.Errorf("Interface %q not found", interfaceName)
}

func GetFirstInterfaceConfigIPv4(interfaceConfigs []workloadnephioorgv1alpha1.InterfaceConfig, interfaceName string, quotes bool) (string, error) {
	interfaceConfig, err := GetFirstInterfaceConfig(interfaceConfigs, interfaceName)
	if err != nil {
		return "", err
	}

	ip, _, err := net.ParseCIDR(interfaceConfig.IPv4.Address)
	if err != nil {
		return "", err
	}

	if quotes {
		return strconv.Quote(ip.String()), nil
	} else {
		return ip.String(), nil
	}
}

func GetAmfIpFromConfigRef(configRefInstances []*configref.Config, gvk schema.GroupVersionKind) (string, error) {

	var extractedIp = ""
	for _, ref := range configRefInstances {
		var b []byte
		if ref.Spec.Config.Object == nil {
			b = ref.Spec.Config.Raw
		} else {
			if ref.Spec.Config.Object.GetObjectKind().GroupVersionKind() == gvk {
				var err error
				if b, err = json.Marshal(ref.Spec.Config.Object); err != nil {
					return extractedIp, err
				}
			} else {
				continue
			}
		}
		amfDeployment := &nephiov1alpha1.AMFDeployment{}
		if err := json.Unmarshal(b, amfDeployment); err != nil {
			return extractedIp, err
		} else {
			extractedIp, err = free5gccontrollers.GetFirstInterfaceConfigIPv4(amfDeployment.Spec.Interfaces, "n2")
			if err != nil {
				extractedIp = ""
				return extractedIp, nil
			}
		}
	}

	return extractedIp, nil
}

func GetCuCpIpFromConfigRef(configRefInstances []*configref.Config, gvk schema.GroupVersionKind, interfaceName string) (string, error) {

	var extractedIp = ""
	for _, ref := range configRefInstances {
		var b []byte
		if ref.Spec.Config.Object == nil {
			b = ref.Spec.Config.Raw
		} else {
			if ref.Spec.Config.Object.GetObjectKind().GroupVersionKind() == gvk {
				var err error
				if b, err = json.Marshal(ref.Spec.Config.Object); err != nil {
					return extractedIp, err
				}
			} else {
				continue
			}
		}
		ranDeployment := &workloadnephioorgv1alpha1.RANDeployment{}
		if err := json.Unmarshal(b, ranDeployment); err != nil {
			return extractedIp, err
		} else {

			extractedIp, err = GetFirstInterfaceConfigIPv4(ranDeployment.Spec.Interfaces, interfaceName, false)
			if err != nil {
				extractedIp = ""
				return extractedIp, err
			}
		}
	}

	return extractedIp, nil
}
