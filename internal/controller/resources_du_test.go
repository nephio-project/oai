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
	"reflect"
	"strconv"
	"strings"
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

/*
The pointer library doesn't has alternate of Uint16
*/
func uint16Ptr(val int) *uint16 {
	a := uint16(val)
	return &a
}

func marshalJsonReturnByteOnly(obj any) []byte {
	marshalObj, err := json.Marshal(obj)
	if err != nil {
		return nil
	}
	return marshalObj
}

func compareStringLineByLineTrimmed(str1 string, str2 string) bool {
	str1List := strings.Split(str1, "\n")
	str2List := strings.Split(str2, "\n")
	if len(str1List) != len(str2List) {
		fmt.Println("length doesn't match in compare string LineByLine | ", len(str1List), " | ", len(str2List))
		return false
	}
	for i := 0; i < len(str1List); i++ {
		str1List[i] = strings.Trim(str1List[i], "\t")
		str2List[i] = strings.Trim(str2List[i], "\t")
		if strings.TrimSpace(str1List[i]) != strings.TrimSpace(str2List[i]) {
			return false
		}
	}
	return true
}

func TestCreateNetworkAttachmentDefinitionNetworksDu(t *testing.T) {
	/*
		Since CreateNetworkAttachmentDefinitionNetworks are rigorously tested in the network_attachment_defination_tests in which error handling is tested
		Therefore, Here testing only the normal cases
	*/
	dummyNfSpec := workloadv1alpha1.NFDeploymentSpec{
		Provider:   "du.openairinterface.org",
		Interfaces: []workloadv1alpha1.InterfaceConfig{},
	}

	cases := map[string]struct {
		inputInterfaceConfig []workloadv1alpha1.InterfaceConfig
		want                 string
	}{
		"Normal": {
			inputInterfaceConfig: []workloadv1alpha1.InterfaceConfig{
				{
					Name: "f1",
					IPv4: &workloadv1alpha1.IPv4{
						Address: "172.5.1.3/24",
						Gateway: pointer.String("172.5.1.1"),
					},
					VLANID: uint16Ptr(2),
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
	}
	duresource := DuResources{}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			dummyNfSpec.Interfaces = tc.inputInterfaceConfig
			got, err := duresource.createNetworkAttachmentDefinitionNetworks("abc", &dummyNfSpec)
			if err != nil {
				t.Errorf("DuResource| createNetworkAttachmentDefinitionNetworks Error %v ", err)
			}
			if !compareStringLineByLineTrimmed(got, tc.want) {
				t.Errorf("DuResource| createNetworkAttachmentDefinitionNetworks returned %s wanted %s", got, tc.want)
			}
		})
	}
}

func generateConfigInstancesMapForTesting(nfF1cSpec workloadv1alpha1.NFDeploymentSpec) map[string][]*configref.Config {
	// Creating the F1C INterface ConfigRef
	nfF1c := workloadv1alpha1.NFDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nf-",
			Namespace: "nf-dummy-du-ns",
		},
		Spec: nfF1cSpec,
	}
	raw, _ := json.Marshal(nfF1c)
	nfF1cConfigInstance := configref.Config{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "ref.nephio.org/v1alpha1",
			Kind:       "Config",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dummy-name",
			Namespace: "dummy-ns",
		},
		Spec: configref.ConfigSpec{
			Config: runtime.RawExtension{
				Raw: raw,
			},
		},
	}

	configInstanceMap := map[string][]*configref.Config{
		"NFDeployment": []*configref.Config{&nfF1cConfigInstance},
	}

	return configInstanceMap
}

func TestGetConfigMapDu(t *testing.T) {
	logger := log.Log

	cases := map[string]struct {
		nfF1Spec    workloadv1alpha1.NFDeploymentSpec
		nfF1cSpec   workloadv1alpha1.NFDeploymentSpec
		paramsRanNf workloadnfconfig.RANConfig
		paramsPlmn  workloadnfconfig.PLMN
		wantedError string
	}{
		"Normal": {
			nfF1Spec: workloadv1alpha1.NFDeploymentSpec{
				Interfaces: []workloadv1alpha1.InterfaceConfig{
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
			nfF1cSpec: workloadv1alpha1.NFDeploymentSpec{
				Provider: "cucp.openairinterface.org",
				Interfaces: []workloadv1alpha1.InterfaceConfig{
					{
						Name: "f1c",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.5.1.254/24",
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
		"F1-Du Not Provided": {
			nfF1Spec:    workloadv1alpha1.NFDeploymentSpec{},
			nfF1cSpec:   workloadv1alpha1.NFDeploymentSpec{},
			paramsRanNf: workloadnfconfig.RANConfig{},
			paramsPlmn:  workloadnfconfig.PLMN{},
			wantedError: "F1-Du not provided",
		},
		"F1C Not Provided": {
			nfF1Spec: workloadv1alpha1.NFDeploymentSpec{
				Interfaces: []workloadv1alpha1.InterfaceConfig{
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
			nfF1cSpec: workloadv1alpha1.NFDeploymentSpec{
				Provider: "cucp.openairinterface.org",
				Interfaces: []workloadv1alpha1.InterfaceConfig{
					{
						Name: "fqrt",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.5.1.254/24",
							Gateway: pointer.String("172.5.1.1"),
						},
						VLANID: uint16Ptr(2),
					},
				},
			},
			paramsRanNf: workloadnfconfig.RANConfig{},
			paramsPlmn:  workloadnfconfig.PLMN{},
			wantedError: "F1c-Du not provided",
		},
		"RANConfig Marshal Error": {
			nfF1Spec: workloadv1alpha1.NFDeploymentSpec{
				Interfaces: []workloadv1alpha1.InterfaceConfig{
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
			nfF1cSpec: workloadv1alpha1.NFDeploymentSpec{
				Provider: "cucp.openairinterface.org",
				Interfaces: []workloadv1alpha1.InterfaceConfig{
					{
						Name: "f1c",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.5.1.254/24",
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
			nfF1Spec: workloadv1alpha1.NFDeploymentSpec{
				Interfaces: []workloadv1alpha1.InterfaceConfig{
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
			nfF1cSpec: workloadv1alpha1.NFDeploymentSpec{
				Provider: "cucp.openairinterface.org",
				Interfaces: []workloadv1alpha1.InterfaceConfig{
					{
						Name: "f1c",
						IPv4: &workloadv1alpha1.IPv4{
							Address: "172.5.1.254/24",
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

	duresource := DuResources{}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			configInstanceMap := generateConfigInstancesMapForTesting(tc.nfF1cSpec)
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

			got := duresource.GetConfigMap(logger, &workloadv1alpha1.NFDeployment{Spec: tc.nfF1Spec}, &configInfo)
			if tc.wantedError == "nil" {
				defaultWantConfigurations, _ := renderConfigurationTemplateForDu(configurationTemplateValuesForDu{
					F1C_DU_IP:       "\"172.5.1.3\"",
					F1C_CU_IP:       "\"172.5.1.254\"",
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

				if !reflect.DeepEqual(got[0].Data["mounted.conf"], defaultWantConfigurations) {
					t.Errorf("GetConfigMap returned %s Wanted %s", got[0].Data["mounted.conf"], defaultWantConfigurations)
				}
			} else {
				if got != nil {
					t.Errorf("GetConfigMap returned %v wanted nil (Error Scenario)", got)
				}
			}

		})
	}
}

func TestGetDeploymentDu(t *testing.T) {
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
					Provider: "du.openairinterface.org",
					Interfaces: []workloadv1alpha1.InterfaceConfig{
						{
							Name: "f1",
							IPv4: &workloadv1alpha1.IPv4{
								Address: "172.5.1.3/24",
								Gateway: pointer.String("172.5.1.1"),
							},
							VLANID: uint16Ptr(2),
						}, {
							Name: "ue-rfsim",
							IPv4: &workloadv1alpha1.IPv4{
								Address: "172.6.0.254/24",
								Gateway: pointer.String("172.6.0.1"),
							},
							VLANID: uint16Ptr(6),
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
							Name: "f1",
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
					Provider: "du.openairinterface.org",
					Interfaces: []workloadv1alpha1.InterfaceConfig{
						{
							Name: "f1",
							IPv4: &workloadv1alpha1.IPv4{
								Address: "172.5.1.3/24",
								Gateway: pointer.String("172.5.1.1"),
							},
							VLANID: uint16Ptr(2),
						}, {
							Name: "ue-rfsim",
							IPv4: &workloadv1alpha1.IPv4{
								Address: "172.6.0.254/24",
								Gateway: pointer.String("172.6.0.1"),
							},
							VLANID: uint16Ptr(6),
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

	duresource := DuResources{}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := duresource.GetDeployment(logger, &tc.ranDeployment, tc.configInfo)
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

func TestGetServiceAccountDu(t *testing.T) {
	cases := map[string]struct {
		want []*corev1.ServiceAccount
	}{
		"Normal": {
			want: []*corev1.ServiceAccount{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "oai-gnb-du-sa",
					},
					TypeMeta: metav1.TypeMeta{
						APIVersion: "v1",
						Kind:       "ServiceAccount",
					},
				},
			},
		},
	}

	duResource := DuResources{}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := duResource.GetServiceAccount()
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("DuResource| GetServiceAccount returned %v Wanted %v", got, tc.want)
			}
		})
	}
}

func TestGetService(t *testing.T) {
	duResource := DuResources{}
	got := duResource.GetService()
	if len(got) == 0 {
		t.Errorf("GetService returned Empty Service ")
	}
}
