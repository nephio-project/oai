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
	"reflect"
	"strconv"
	"testing"

	configref "github.com/nephio-project/api/references/v1alpha1"
	workloadv1alpha1 "github.com/nephio-project/api/workload/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/log"
	workloadnfconfig "workload.nephio.org/ran_deployment/api/v1alpha1"
)

func TestCreateNetworkAttachmentDefinitionNetworksCuUP(t *testing.T) {
	/*
		Since CreateNetworkAttachmentDefinitionNetworks are rigorously tested in the network_attachment_defination_tests in which error handling is tested
		Therefore, Here testing only the normal cases
	*/
	dummyNfSpec := workloadv1alpha1.NFDeploymentSpec{
		Provider:   "cuup.openairinterface.org",
		Interfaces: []workloadv1alpha1.InterfaceConfig{},
	}

	cases := map[string]struct {
		inputInterfaceConfig []workloadv1alpha1.InterfaceConfig
		want                 string
	}{
		"Normal": {
			inputInterfaceConfig: []workloadv1alpha1.InterfaceConfig{
				{
					Name: "e1",
					IPv4: &workloadv1alpha1.IPv4{
						Address: "172.5.5.7/24",
						Gateway: ptr.To("172.5.1.2"),
					},
					VLANID: uint16Ptr(3),
				}, {
					Name: "f1u",
					IPv4: &workloadv1alpha1.IPv4{
						Address: "172.5.1.3/24",
						Gateway: ptr.To("172.5.1.1"),
					},
					VLANID: uint16Ptr(2),
				}, {
					Name: "n3",
					IPv4: &workloadv1alpha1.IPv4{
						Address: "172.6.0.254/24",
						Gateway: ptr.To("172.6.0.1"),
					},
					VLANID: uint16Ptr(6),
				},
			},
			want: `[
				{
				 "name": "abc-e1",
				 "interface": "e1",
				 "ips": ["172.5.5.7/24"],
				 "gateways": ["172.5.1.2"]
				},		  
				{
				 "name": "abc-f1u",
				 "interface": "f1u",
				 "ips": ["172.5.1.3/24"],
				 "gateways": ["172.5.1.1"]
				},
				{
				 "name": "abc-n3",
				 "interface": "n3",
				 "ips": ["172.6.0.254/24"],
				 "gateways": ["172.6.0.1"]
				}
			   ] `,
		},
	}
	cuupResource := CuUpResources{}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			dummyNfSpec.Interfaces = tc.inputInterfaceConfig
			got, err := cuupResource.createNetworkAttachmentDefinitionNetworks("abc", &dummyNfSpec)
			if err != nil {
				t.Errorf("CuupResource| createNetworkAttachmentDefinitionNetworks Error %v ", err)
			}
			if !compareStringLineByLineTrimmed(got, tc.want) {
				t.Errorf("CuupResource| createNetworkAttachmentDefinitionNetworks returned %s wanted %s", got, tc.want)
			}
		})
	}
}

func TestGetDeploymentCuUp(t *testing.T) {
	//Since we have done rigrous testing of createNetworkAttachmentDefinitionNetworks(), so, we are skipping corner-cases for that function
	logger := log.Log
	cases := map[string]struct {
		ranDeployment workloadv1alpha1.NFDeployment
		configInfo    *ConfigInfo
		want          string
	}{
		"Normal": {
			ranDeployment: workloadv1alpha1.NFDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dummy-du",
					Namespace: "dummy-du-ns",
				},
				Spec: workloadv1alpha1.NFDeploymentSpec{
					Provider: "cuup.openairinterface.org",
					Interfaces: []workloadv1alpha1.InterfaceConfig{
						{
							Name: "e1",
							IPv4: &workloadv1alpha1.IPv4{
								Address: "172.5.1.3/24",
								Gateway: ptr.To("172.5.1.1"),
							},
							VLANID: uint16Ptr(2),
						}, {
							Name: "n3",
							IPv4: &workloadv1alpha1.IPv4{
								Address: "172.6.0.254/24",
								Gateway: ptr.To("172.6.0.1"),
							},
							VLANID: uint16Ptr(6),
						}, {
							Name: "f1u",
							IPv4: &workloadv1alpha1.IPv4{
								Address: "172.6.0.7/24",
								Gateway: ptr.To("172.6.0.1"),
							},
							VLANID: uint16Ptr(7),
						},
					},
				},
			},
			configInfo: &ConfigInfo{
				ConfigSelfInfo: map[string]runtime.RawExtension{
					"OAIConfig": runtime.RawExtension{
						Raw: marshalJsonReturnByteOnly(&workloadnfconfig.OAIConfig{Spec: workloadnfconfig.OAIConfigSpec{Image: "dummy-image"}}),
					},
				},
			},
			want: "Pod-Annotations",
		},
		"NAD Error": {
			ranDeployment: workloadv1alpha1.NFDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "dummy-du",
				},
				Spec: workloadv1alpha1.NFDeploymentSpec{
					Interfaces: []workloadv1alpha1.InterfaceConfig{
						{
							Name: "e1",
							IPv4: &workloadv1alpha1.IPv4{
								Address: "172.5.1.3/24", // Gateway is removed
							},
							VLANID: uint16Ptr(2),
						},
					},
				},
			},
			want: "error",
		},
		"OAI-Config Unmarshal Error": {
			ranDeployment: workloadv1alpha1.NFDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dummy-du",
					Namespace: "dummy-du-ns",
				},
				Spec: workloadv1alpha1.NFDeploymentSpec{
					Provider: "cuup.openairinterface.org",
					Interfaces: []workloadv1alpha1.InterfaceConfig{
						{
							Name: "e1",
							IPv4: &workloadv1alpha1.IPv4{
								Address: "172.5.1.3/24",
								Gateway: ptr.To("172.5.1.1"),
							},
							VLANID: uint16Ptr(2),
						}, {
							Name: "n3",
							IPv4: &workloadv1alpha1.IPv4{
								Address: "172.6.0.254/24",
								Gateway: ptr.To("172.6.0.1"),
							},
							VLANID: uint16Ptr(6),
						}, {
							Name: "f1u",
							IPv4: &workloadv1alpha1.IPv4{
								Address: "172.6.0.7/24",
								Gateway: ptr.To("172.6.0.1"),
							},
							VLANID: uint16Ptr(7),
						},
					},
				},
			},
			configInfo: &ConfigInfo{
				ConfigSelfInfo: map[string]runtime.RawExtension{
					"OAIConfig": runtime.RawExtension{
						Raw: []byte(" "),
					},
				},
			},
			want: "error",
		},
	}

	cuupResource := CuUpResources{}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := cuupResource.GetDeployment(logger, &tc.ranDeployment, tc.configInfo)
			if tc.want == "error" {
				if got != nil {
					t.Errorf("GetDeployment returned %v wanted nil", got)
				}
			} else {
				if got == nil {
					t.Error("GetDeployment returned nil wanted DeploymentObject")
					return
				}
				gotPodAnnotations := got[0].Spec.Template.Annotations
				if len(gotPodAnnotations) == 0 {
					t.Error("PodAnnotations Not Set During GetDeployment")
				}
				gotImage := got[0].Spec.Template.Spec.Containers[0].Image
				if gotImage != "dummy-image" {
					t.Errorf("Image Got %s wanted %s", gotImage, "dummy-image")
				}
			}

		})

	}

}

func TestGetServiceAccountCuUp(t *testing.T) {

	cuUpResource := CuUpResources{}
	actual := cuUpResource.GetServiceAccount()
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: "oai-cu-up-sa",
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
	}
	expected := []*corev1.ServiceAccount{serviceAccount}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("GetServiceAccount returned %v, expected %v", actual, expected)
	}
}

func TestGetConfigMapCuUp(t *testing.T) {
	defaultParamsPlmn := workloadnfconfig.PLMN{
		Spec: workloadnfconfig.PLMNSpec{
			PLMNInfo: []workloadnfconfig.PLMNInfo{
				{
					PLMNID: workloadnfconfig.PLMNID{
						MCC: "001",
						MNC: "01",
					},
					TAC: 1,
					NSSAI: []workloadnfconfig.NSSAI{
						{
							SST: 1,
							SD:  ptr.To("ffffff"),
						},
					},
				},
			},
		},
	}
	defaultConfiguration, _ := renderConfigurationTemplateForCuUp(configurationTemplateValuesForCuUp{
		E1_IP:           "\"172.5.1.3\"",
		F1U_IP:          "\"172.6.0.7\"",
		N3_IP:           "\"172.6.0.254\"",
		CUCP_E1:         "\"172.5.1.3\"",
		TAC:             defaultParamsPlmn.Spec.PLMNInfo[0].TAC,
		PLMN_MCC:        defaultParamsPlmn.Spec.PLMNInfo[0].PLMNID.MCC,
		PLMN_MNC:        defaultParamsPlmn.Spec.PLMNInfo[0].PLMNID.MNC,
		PLMN_MNC_LENGTH: strconv.Itoa(int(len(defaultParamsPlmn.Spec.PLMNInfo[0].PLMNID.MNC))),
		NSSAI_SST:       defaultParamsPlmn.Spec.PLMNInfo[0].NSSAI[0].SST,
		NSSAI_SD:        *defaultParamsPlmn.Spec.PLMNInfo[0].NSSAI[0].SD,
	})

	cases := map[string]struct {
		ranDeploymentSpec      workloadv1alpha1.NFDeploymentSpec
		configNfDeploymentSpec workloadv1alpha1.NFDeploymentSpec
		configSelfInfo         map[string]runtime.RawExtension
		wantedConfiguration    string
	}{
		"Normal": {
			ranDeploymentSpec: workloadv1alpha1.NFDeploymentSpec{
				Interfaces: []workloadv1alpha1.InterfaceConfig{
					{
						Name: "e1",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.5.1.3/24",
							Gateway: ptr.To("172.5.1.1"),
						},
						VLANID: uint16Ptr(2),
					}, {
						Name: "n3",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.6.0.254/24",
							Gateway: ptr.To("172.6.0.1"),
						},
						VLANID: uint16Ptr(6),
					}, {
						Name: "f1u",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.6.0.7/24",
							Gateway: ptr.To("172.6.0.1"),
						},
						VLANID: uint16Ptr(7),
					},
				},
			},
			configNfDeploymentSpec: workloadv1alpha1.NFDeploymentSpec{
				Provider: "cucp.openairinterface.org",
				Interfaces: []workloadv1alpha1.InterfaceConfig{
					{
						Name: "e1",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.5.1.3/24",
							Gateway: ptr.To("172.5.1.1"),
						},
						VLANID: uint16Ptr(2),
					},
				},
			},
			configSelfInfo: map[string]runtime.RawExtension{
				"PLMN": runtime.RawExtension{Raw: marshalJsonReturnByteOnly(defaultParamsPlmn)},
			},
			wantedConfiguration: defaultConfiguration,
		},
		"N3 Not Provided": {
			ranDeploymentSpec:      workloadv1alpha1.NFDeploymentSpec{},
			configNfDeploymentSpec: workloadv1alpha1.NFDeploymentSpec{},
			configSelfInfo:         nil,
			wantedConfiguration:    "nil",
		},
		"E1 Not Provided": {
			ranDeploymentSpec: workloadv1alpha1.NFDeploymentSpec{
				Interfaces: []workloadv1alpha1.InterfaceConfig{
					{
						Name: "n3",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.6.0.254/24",
							Gateway: ptr.To("172.6.0.1"),
						},
						VLANID: uint16Ptr(6),
					},
				},
			},
			configSelfInfo:         nil,
			configNfDeploymentSpec: workloadv1alpha1.NFDeploymentSpec{},
			wantedConfiguration:    "nil",
		},
		"F1u NOt Provided": {
			ranDeploymentSpec: workloadv1alpha1.NFDeploymentSpec{
				Interfaces: []workloadv1alpha1.InterfaceConfig{
					{
						Name: "e1",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.5.1.3/24",
							Gateway: ptr.To("172.5.1.1"),
						},
						VLANID: uint16Ptr(2),
					}, {
						Name: "n3",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.6.0.254/24",
							Gateway: ptr.To("172.6.0.1"),
						},
						VLANID: uint16Ptr(6),
					},
				},
			},
			configSelfInfo:         nil,
			configNfDeploymentSpec: workloadv1alpha1.NFDeploymentSpec{},
			wantedConfiguration:    "nil",
		},
		"configNfDeployment not set": {
			ranDeploymentSpec: workloadv1alpha1.NFDeploymentSpec{
				Interfaces: []workloadv1alpha1.InterfaceConfig{
					{
						Name: "e1",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.5.1.3/24",
							Gateway: ptr.To("172.5.1.1"),
						},
						VLANID: uint16Ptr(2),
					}, {
						Name: "n3",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.6.0.254/24",
							Gateway: ptr.To("172.6.0.1"),
						},
						VLANID: uint16Ptr(6),
					}, {
						Name: "f1u",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.6.0.7/24",
							Gateway: ptr.To("172.6.0.1"),
						},
						VLANID: uint16Ptr(7),
					},
				},
			},
			configNfDeploymentSpec: workloadv1alpha1.NFDeploymentSpec{Provider: "cucp.openairinterface.org"},
			configSelfInfo:         nil,
			wantedConfiguration:    "nil",
		},
		"configSelfInfo Unmarshal Error": {
			ranDeploymentSpec: workloadv1alpha1.NFDeploymentSpec{
				Interfaces: []workloadv1alpha1.InterfaceConfig{
					{
						Name: "e1",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.5.1.3/24",
							Gateway: ptr.To("172.5.1.1"),
						},
						VLANID: uint16Ptr(2),
					}, {
						Name: "n3",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.6.0.254/24",
							Gateway: ptr.To("172.6.0.1"),
						},
						VLANID: uint16Ptr(6),
					}, {
						Name: "f1u",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.6.0.7/24",
							Gateway: ptr.To("172.6.0.1"),
						},
						VLANID: uint16Ptr(7),
					},
				},
			},
			configNfDeploymentSpec: workloadv1alpha1.NFDeploymentSpec{
				Provider: "cucp.openairinterface.org",
				Interfaces: []workloadv1alpha1.InterfaceConfig{
					{
						Name: "e1",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.5.1.3/24",
							Gateway: ptr.To("172.5.1.1"),
						},
						VLANID: uint16Ptr(2),
					},
				},
			},
			configSelfInfo: map[string]runtime.RawExtension{
				"PLMN": runtime.RawExtension{Raw: []byte(" ")},
			},
			wantedConfiguration: "nil",
		},
	}

	logger := log.Log
	cuUpResource := CuUpResources{}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			ranDeploymentDummy := workloadv1alpha1.NFDeployment{
				Spec: tc.ranDeploymentSpec,
			}
			configInstanceMap := map[string][]*configref.Config{
				"NFDeployment": []*configref.Config{
					{
						Spec: configref.ConfigSpec{
							Config: runtime.RawExtension{
								Raw: marshalJsonReturnByteOnly(workloadv1alpha1.NFDeployment{
									Spec: tc.configNfDeploymentSpec,
								}),
							},
						},
					},
				},
			}

			configInfo := ConfigInfo{
				ConfigRefInfo:  configInstanceMap,
				ConfigSelfInfo: tc.configSelfInfo,
			}

			got := cuUpResource.GetConfigMap(logger, &ranDeploymentDummy, &configInfo)
			if tc.wantedConfiguration == "nil" {
				if got != nil {
					t.Errorf("GetConfigMap CuUp returned %v  Wanted nil", got)
				}
			} else {
				if got[0].Data["gnb.conf"] != tc.wantedConfiguration {
					t.Errorf("GetConfigMap CuUp returned %v  Wanted %v", got[0].Data["gnb.conf"], tc.wantedConfiguration)
				}
			}

		})
	}
}

func TestGetServiceCuUp(t *testing.T) {
	cuupResource := CuUpResources{}
	got := cuupResource.GetService()
	/*
		More cases will be added when more code will be added to GetService
	*/
	if got == nil {
		t.Errorf("GetService returned nil ")
	}

}
