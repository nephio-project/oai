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
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/log"
	workloadnfconfig "workload.nephio.org/ran_deployment/api/v1alpha1"
)

func TestCreateNetworkAttachmentDefinitionNetworksCuCP(t *testing.T) {
	/*
		Since CreateNetworkAttachmentDefinitionNetworks are rigorously tested in the network_attachment_defination_tests in which error handling is tested
		Therefore, Here testing only the normal cases
	*/
	dummyNfSpec := workloadv1alpha1.NFDeploymentSpec{
		Provider:   "cucp.openairinterface.org",
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
						Address: "172.5.3.6/24",
						Gateway: pointer.String("172.7.1.1"),
					},
					VLANID: uint16Ptr(2),
				},
				{
					Name: "f1c",
					IPv4: &workloadv1alpha1.IPv4{
						Address: "172.5.1.3/24",
						Gateway: pointer.String("172.5.1.1"),
					},
					VLANID: uint16Ptr(2),
				}, {
					Name: "n2",
					IPv4: &workloadv1alpha1.IPv4{
						Address: "172.6.0.254/24",
						Gateway: pointer.String("172.6.0.1"),
					},
					VLANID: uint16Ptr(6),
				},
			},
			want: `[
				{
				 "name": "abc-e1",
				 "interface": "e1",
				 "ips": ["172.5.3.6/24"],
				 "gateways": ["172.7.1.1"]
				},		  
				{
				 "name": "abc-f1c",
				 "interface": "f1c",
				 "ips": ["172.5.1.3/24"],
				 "gateways": ["172.5.1.1"]
				},
				{
				 "name": "abc-n2",
				 "interface": "n2",
				 "ips": ["172.6.0.254/24"],
				 "gateways": ["172.6.0.1"]
				}
			   ] `,
		},
	}
	cucpResource := CuCpResources{}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			dummyNfSpec.Interfaces = tc.inputInterfaceConfig
			got, err := cucpResource.createNetworkAttachmentDefinitionNetworks("abc", &dummyNfSpec)
			if err != nil {
				t.Errorf("CucpResource| createNetworkAttachmentDefinitionNetworks Error %v ", err)
			}
			if !compareStringLineByLineTrimmed(got, tc.want) {
				t.Errorf("CucpResource| createNetworkAttachmentDefinitionNetworks returned %s wanted %s", got, tc.want)
			}
		})
	}
}

func TestGetDeploymentCuCp(t *testing.T) {
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
					Provider: "cucp.openairinterface.org",
					Interfaces: []workloadv1alpha1.InterfaceConfig{
						{
							Name: "e1",
							IPv4: &workloadv1alpha1.IPv4{
								Address: "172.5.1.3/24",
								Gateway: pointer.String("172.5.1.1"),
							},
							VLANID: uint16Ptr(2),
						}, {
							Name: "n2",
							IPv4: &workloadv1alpha1.IPv4{
								Address: "172.6.0.254/24",
								Gateway: pointer.String("172.6.0.1"),
							},
							VLANID: uint16Ptr(6),
						}, {
							Name: "f1c",
							IPv4: &workloadv1alpha1.IPv4{
								Address: "172.6.0.7/24",
								Gateway: pointer.String("172.6.0.1"),
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
			configInfo: &ConfigInfo{},
			want:       "error",
		},
		"OAI-Config Unmarshal Error": {
			ranDeployment: workloadv1alpha1.NFDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dummy-du",
					Namespace: "dummy-du-ns",
				},
				Spec: workloadv1alpha1.NFDeploymentSpec{
					Provider: "cucp.openairinterface.org",
					Interfaces: []workloadv1alpha1.InterfaceConfig{
						{
							Name: "e1",
							IPv4: &workloadv1alpha1.IPv4{
								Address: "172.5.1.3/24",
								Gateway: pointer.String("172.5.1.1"),
							},
							VLANID: uint16Ptr(2),
						}, {
							Name: "n2",
							IPv4: &workloadv1alpha1.IPv4{
								Address: "172.6.0.254/24",
								Gateway: pointer.String("172.6.0.1"),
							},
							VLANID: uint16Ptr(6),
						}, {
							Name: "f1c",
							IPv4: &workloadv1alpha1.IPv4{
								Address: "172.6.0.7/24",
								Gateway: pointer.String("172.6.0.1"),
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

	cucpResource := CuCpResources{}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := cucpResource.GetDeployment(logger, &tc.ranDeployment, tc.configInfo)
			if tc.want == "error" {
				if got != nil {
					t.Errorf("GetDeployment returned %v wanted nil", got)
				}
			} else {
				if got == nil {
					t.Error("GetDeployment returned nil wanted DeploymentObject")
					return
				}
				gotPodAnnotations := got[0].Spec.Template.ObjectMeta.Annotations
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

func TestGetServiceAccountCuCp(t *testing.T) {

	cucpResource := CuCpResources{}
	actual := cucpResource.GetServiceAccount()
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: "oai-gnb-cu-cp-sa",
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

func TestGetServiceCuCp(t *testing.T) {
	cucpResource := CuCpResources{}
	got := cucpResource.GetService()
	/*
		More cases will be added when more code will be added to GetService
	*/
	if got == nil {
		t.Errorf("GetService returned nil ")
	}

}

func TestGetConfigMapCuCp(t *testing.T) {
	cases := map[string]struct {
		ranDeploymentSpec      workloadv1alpha1.NFDeploymentSpec
		configNfDeploymentSpec workloadv1alpha1.NFDeploymentSpec
		paramsRanNf            workloadnfconfig.RANConfig
		paramsPlmn             workloadnfconfig.PLMN
		wantedError            string
	}{
		"Normal": {
			ranDeploymentSpec: workloadv1alpha1.NFDeploymentSpec{
				Interfaces: []workloadv1alpha1.InterfaceConfig{
					{
						Name: "e1",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.5.1.3/24",
							Gateway: pointer.String("172.5.1.1"),
						},
						VLANID: uint16Ptr(2),
					}, {
						Name: "n2",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.6.0.254/24",
							Gateway: pointer.String("172.6.0.1"),
						},
						VLANID: uint16Ptr(6),
					}, {
						Name: "f1c",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.6.0.7/24",
							Gateway: pointer.String("172.6.0.1"),
						},
						VLANID: uint16Ptr(7),
					},
				},
			},
			configNfDeploymentSpec: workloadv1alpha1.NFDeploymentSpec{
				Provider: "amf.openairinterface.org",
				Interfaces: []workloadv1alpha1.InterfaceConfig{
					{
						Name: "n2",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.5.1.3/24",
							Gateway: pointer.String("172.5.1.1"),
						},
						VLANID: uint16Ptr(2),
					},
				},
			},
			paramsRanNf: workloadnfconfig.RANConfig{
				Spec: workloadnfconfig.RANConfigSpec{
					CellIdentity:   "12345678L",
					PhysicalCellID: uint32(0),
				},
			},
			paramsPlmn: workloadnfconfig.PLMN{
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
									SD:  pointer.String("ffffff"),
								},
							},
						},
					},
				},
			},
			wantedError: "nil",
		},
		"N2 Not Provided": {
			ranDeploymentSpec:      workloadv1alpha1.NFDeploymentSpec{},
			configNfDeploymentSpec: workloadv1alpha1.NFDeploymentSpec{},
			paramsRanNf:            workloadnfconfig.RANConfig{},
			paramsPlmn:             workloadnfconfig.PLMN{},
			wantedError:            "N2 Not Provided",
		},
		"E1 Not Provided": {
			ranDeploymentSpec: workloadv1alpha1.NFDeploymentSpec{
				Interfaces: []workloadv1alpha1.InterfaceConfig{
					{
						Name: "n2",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.6.0.254/24",
							Gateway: pointer.String("172.6.0.1"),
						},
						VLANID: uint16Ptr(6),
					},
				},
			},
			configNfDeploymentSpec: workloadv1alpha1.NFDeploymentSpec{},
			paramsRanNf:            workloadnfconfig.RANConfig{},
			paramsPlmn:             workloadnfconfig.PLMN{},
			wantedError:            "E1 Not Provided",
		},
		"F1c NOt Provided": {
			ranDeploymentSpec: workloadv1alpha1.NFDeploymentSpec{
				Interfaces: []workloadv1alpha1.InterfaceConfig{
					{
						Name: "e1",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.5.1.3/24",
							Gateway: pointer.String("172.5.1.1"),
						},
						VLANID: uint16Ptr(2),
					}, {
						Name: "n2",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.6.0.254/24",
							Gateway: pointer.String("172.6.0.1"),
						},
						VLANID: uint16Ptr(6),
					},
				},
			},
			configNfDeploymentSpec: workloadv1alpha1.NFDeploymentSpec{},
			paramsRanNf:            workloadnfconfig.RANConfig{},
			paramsPlmn:             workloadnfconfig.PLMN{},
			wantedError:            "F1c NOt Provided",
		},
		"configNfDeployment not set": {
			ranDeploymentSpec: workloadv1alpha1.NFDeploymentSpec{
				Interfaces: []workloadv1alpha1.InterfaceConfig{
					{
						Name: "e1",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.5.1.3/24",
							Gateway: pointer.String("172.5.1.1"),
						},
						VLANID: uint16Ptr(2),
					}, {
						Name: "n2",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.6.0.254/24",
							Gateway: pointer.String("172.6.0.1"),
						},
						VLANID: uint16Ptr(6),
					}, {
						Name: "f1c",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.6.0.7/24",
							Gateway: pointer.String("172.6.0.1"),
						},
						VLANID: uint16Ptr(7),
					},
				},
			},
			configNfDeploymentSpec: workloadv1alpha1.NFDeploymentSpec{Provider: "amf.openairinterface.org"},
			paramsRanNf:            workloadnfconfig.RANConfig{},
			paramsPlmn:             workloadnfconfig.PLMN{},
			wantedError:            "configNfDeployment not set",
		},
		"RANConfig Marshal Error": {
			ranDeploymentSpec: workloadv1alpha1.NFDeploymentSpec{
				Interfaces: []workloadv1alpha1.InterfaceConfig{
					{
						Name: "e1",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.5.1.3/24",
							Gateway: pointer.String("172.5.1.1"),
						},
						VLANID: uint16Ptr(2),
					}, {
						Name: "n2",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.6.0.254/24",
							Gateway: pointer.String("172.6.0.1"),
						},
						VLANID: uint16Ptr(6),
					}, {
						Name: "f1c",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.6.0.7/24",
							Gateway: pointer.String("172.6.0.1"),
						},
						VLANID: uint16Ptr(7),
					},
				},
			},
			configNfDeploymentSpec: workloadv1alpha1.NFDeploymentSpec{
				Provider: "amf.openairinterface.org",
				Interfaces: []workloadv1alpha1.InterfaceConfig{
					{
						Name: "n2",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.5.1.3/24",
							Gateway: pointer.String("172.5.1.1"),
						},
						VLANID: uint16Ptr(2),
					},
				},
			},
			paramsRanNf: workloadnfconfig.RANConfig{},
			paramsPlmn:  workloadnfconfig.PLMN{},
			wantedError: "RANConfig Marshal Error",
		},
		"PLMN Marshal Error": {
			ranDeploymentSpec: workloadv1alpha1.NFDeploymentSpec{
				Interfaces: []workloadv1alpha1.InterfaceConfig{
					{
						Name: "e1",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.5.1.3/24",
							Gateway: pointer.String("172.5.1.1"),
						},
						VLANID: uint16Ptr(2),
					}, {
						Name: "n2",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.6.0.254/24",
							Gateway: pointer.String("172.6.0.1"),
						},
						VLANID: uint16Ptr(6),
					}, {
						Name: "f1c",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.6.0.7/24",
							Gateway: pointer.String("172.6.0.1"),
						},
						VLANID: uint16Ptr(7),
					},
				},
			},
			configNfDeploymentSpec: workloadv1alpha1.NFDeploymentSpec{
				Provider: "amf.openairinterface.org",
				Interfaces: []workloadv1alpha1.InterfaceConfig{
					{
						Name: "n2",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.5.1.3/24",
							Gateway: pointer.String("172.5.1.1"),
						},
						VLANID: uint16Ptr(2),
					},
				},
			},
			paramsRanNf: workloadnfconfig.RANConfig{},
			paramsPlmn:  workloadnfconfig.PLMN{},
			wantedError: "PLMN Marshal Error",
		},
	}

	logger := log.Log
	cuCpResource := CuCpResources{}
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
				ConfigRefInfo: configInstanceMap,
				ConfigSelfInfo: map[string]runtime.RawExtension{
					"RANConfig": runtime.RawExtension{Raw: marshalJsonReturnByteOnly(tc.paramsRanNf)},
					"PLMN":      runtime.RawExtension{Raw: marshalJsonReturnByteOnly(tc.paramsPlmn)},
				},
			}
			// Simulating JSON UnMarshal Error:
			if tc.wantedError == "RANConfig Marshal Error" {
				configInfo.ConfigSelfInfo["RANConfig"] = runtime.RawExtension{Raw: []byte("")}
			} else if tc.wantedError == "PLMN Marshal Error" {
				configInfo.ConfigSelfInfo["PLMN"] = runtime.RawExtension{Raw: []byte("")}
			}

			got := cuCpResource.GetConfigMap(logger, &ranDeploymentDummy, &configInfo)
			if tc.wantedError == "nil" {
				defaultWantConfigurations, _ := renderConfigurationTemplateForCuCp(configurationTemplateValuesForCuCp{
					E1_IP:           "\"172.5.1.3\"",
					F1C_IP:          "\"172.6.0.7\"",
					N2_IP:           "\"172.6.0.254\"",
					AMF_IP:          "\"172.5.1.3\"",
					TAC:             tc.paramsPlmn.Spec.PLMNInfo[0].TAC,
					CELL_ID:         tc.paramsRanNf.Spec.CellIdentity,
					PHY_CELL_ID:     tc.paramsRanNf.Spec.PhysicalCellID,
					DL_FREQ_BAND:    tc.paramsRanNf.Spec.DownlinkFrequencyBand,
					DL_SCS:          tc.paramsRanNf.Spec.DownlinkSubCarrierSpacing,
					DL_CARRIER_BW:   tc.paramsRanNf.Spec.DownlinkCarrierBandwidth,
					UL_FREQ_BAND:    tc.paramsRanNf.Spec.UplinkFrequencyBand,
					UL_SCS:          tc.paramsRanNf.Spec.UplinkSubCarrierSpacing,
					UL_CARRIER_BW:   tc.paramsRanNf.Spec.UplinkCarrierBandwidth,
					PLMN_MCC:        tc.paramsPlmn.Spec.PLMNInfo[0].PLMNID.MCC,
					PLMN_MNC:        tc.paramsPlmn.Spec.PLMNInfo[0].PLMNID.MNC,
					PLMN_MNC_LENGTH: strconv.Itoa(int(len(tc.paramsPlmn.Spec.PLMNInfo[0].PLMNID.MNC))),
					NSSAI_SST:       tc.paramsPlmn.Spec.PLMNInfo[0].NSSAI[0].SST,
					NSSAI_SD:        *tc.paramsPlmn.Spec.PLMNInfo[0].NSSAI[0].SD,
				})

				if !reflect.DeepEqual(got[0].Data["gnb.conf"], defaultWantConfigurations) {
					t.Errorf("GetConfigMap returned %s Wanted %s", got[0].Data["gnb.conf"], defaultWantConfigurations)
				}
			} else {
				if got != nil {
					t.Errorf("GetConfigMap returned %v wanted nil (Error Scenario)", got)
				}
			}

		})
	}
}
