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
	"net"

	workloadv1alpha1 "github.com/nephio-project/api/workload/v1alpha1"
)

func GetInterfaceConfigs(interfaceConfigs []workloadv1alpha1.InterfaceConfig, interfaceName string) []workloadv1alpha1.InterfaceConfig {
	var selectedInterfaceConfigs []workloadv1alpha1.InterfaceConfig

	for _, interfaceConfig := range interfaceConfigs {
		if interfaceConfig.Name == interfaceName {
			selectedInterfaceConfigs = append(selectedInterfaceConfigs, interfaceConfig)
		}
	}

	return selectedInterfaceConfigs
}

func GetFirstInterfaceConfig(interfaceConfigs []workloadv1alpha1.InterfaceConfig, interfaceName string) (*workloadv1alpha1.InterfaceConfig, error) {
	for _, interfaceConfig := range interfaceConfigs {
		if interfaceConfig.Name == interfaceName {
			return &interfaceConfig, nil
		}
	}

	return nil, fmt.Errorf("Interface %q not found", interfaceName)
}

func GetFirstInterfaceConfigIPv4(interfaceConfigs []workloadv1alpha1.InterfaceConfig, interfaceName string) (string, error) {
	interfaceConfig, err := GetFirstInterfaceConfig(interfaceConfigs, interfaceName)
	if err != nil {
		return "", err
	}

	ip, _, err := net.ParseCIDR(interfaceConfig.IPv4.Address)
	if err != nil {
		return "", err
	}

	return ip.String(), nil
}
