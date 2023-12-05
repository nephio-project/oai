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
	"testing"

	workloadv1alpha1 "github.com/nephio-project/api/workload/v1alpha1"
	"k8s.io/utils/pointer"
)

func TestCreateNetworkAttachmentDefinitionNetworks(t *testing.T) {
	cases := map[string]struct {
		templateName     string
		interfaceConfigs map[string][]workloadv1alpha1.InterfaceConfig
		want             string
	}{
		"Normal": {
			templateName: "abc",
			interfaceConfigs: map[string][]workloadv1alpha1.InterfaceConfig{
				"f1": []workloadv1alpha1.InterfaceConfig{
					{
						Name: "f1",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.5.1.3/24",
							Gateway: pointer.String("172.5.1.1"),
						},
						VLANID: uint16Ptr(2),
					},
				},
			},
			want: `[
				{
				 "name": "abc-f1",
				 "interface": "f1",
				 "ips": ["172.5.1.3/24"],
				 "gateways": ["172.5.1.1"]
				}
			   ] `,
		},
		"Gateway Not Provided": {
			templateName: "abc",
			interfaceConfigs: map[string][]workloadv1alpha1.InterfaceConfig{
				"f1": []workloadv1alpha1.InterfaceConfig{
					{
						Name: "f1",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.5.1.3/24",
						},
						VLANID: uint16Ptr(2),
					},
				},
			},
			want: "error",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := CreateNetworkAttachmentDefinitionNetworks(tc.templateName, tc.interfaceConfigs)
			if tc.want == "error" {
				if err == nil {
					t.Errorf("createNetworkAttachmentDefinitionNetworks Returned Nil expecting error")
				}
			} else {
				if err != nil {
					t.Errorf("createNetworkAttachmentDefinitionNetworks Error %v ", err)
				}
				if !compareStringLineByLineTrimmed(got, tc.want) {
					t.Errorf("createNetworkAttachmentDefinitionNetworks returned %s wanted %s", got, tc.want)
				}
			}

		})
	}

}
