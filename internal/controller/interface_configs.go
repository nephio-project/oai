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
	"fmt"
	"net"
	"strconv"

	workloadnephioorgv1alpha1 "workload.nephio.org/ran_deployment/api/v1alpha1"
)

// TODO: will be removed and common functions defined for free5gc can be used when using NFdeploy
func GetInterfaceConfigs(interfaceConfigs []workloadnephioorgv1alpha1.InterfaceConfig, interfaceName string) []workloadnephioorgv1alpha1.InterfaceConfig {
	var selectedInterfaceConfigs []workloadnephioorgv1alpha1.InterfaceConfig

	for _, interfaceConfig := range interfaceConfigs {
		if interfaceConfig.Name == interfaceName {
			selectedInterfaceConfigs = append(selectedInterfaceConfigs, interfaceConfig)
		}
	}

	return selectedInterfaceConfigs
}

// TODO: will be removed and common functions defined for free5gc can be used when using NFdeploy
func GetFirstInterfaceConfig(interfaceConfigs []workloadnephioorgv1alpha1.InterfaceConfig, interfaceName string) (*workloadnephioorgv1alpha1.InterfaceConfig, error) {
	for _, interfaceConfig := range interfaceConfigs {
		if interfaceConfig.Name == interfaceName {
			return &interfaceConfig, nil
		}
	}

	return nil, fmt.Errorf("Interface %q not found", interfaceName)
}

// TODO: will be removed and common functions defined for free5gc can be used when using NFdeploy
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
