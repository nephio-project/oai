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
	"reflect"
	"testing"

	configref "github.com/nephio-project/api/references/v1alpha1"
	workloadv1alpha1 "github.com/nephio-project/api/workload/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

/*
Generate a List of Kind:Config, With NF-providers as provided in ConfigProvider
*/
func generateConfigInstanceForTesting(configProviders []string) []*configref.Config {

	configInstance := []*configref.Config{}

	for i := 0; i < len(configProviders); i++ {
		nf := workloadv1alpha1.NFDeployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "nf-" + configProviders[i],
				Namespace: "nf-dummy-du-ns",
			},
			Spec: workloadv1alpha1.NFDeploymentSpec{
				Provider: configProviders[i],
			},
		}

		raw, _ := json.Marshal(nf)
		if configProviders[i] == "generate-json-error" {
			raw = []byte("String Not JSON Marshal in order to actuate JSON error")
		}
		yourConfigInstance := configref.Config{
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
		configInstance = append(configInstance, &yourConfigInstance)

	}
	return configInstance
}

/*
	configInstancesMap is a Map of []

If Dependency = UPFDeployment, which will lead to kind:Config with Spec.config Having the Dependency kind (UPFDeployment)
So, configInstanceMap: {DependencyKind : [All the dependencies of that king injected]}
*/
func TestGetConfigInstanceByProvider(t *testing.T) {
	logger := log.Log

	cases := map[string]struct {
		configProviders []string
		providerToLook  string
		want            string
	}{
		"Normal": {
			configProviders: []string{"cucp.openairinterface.org", "du.openairinterface.org"},
			providerToLook:  "du.openairinterface.org",
			want:            "du.openairinterface.org",
		},
		"Provider Not Found": {
			configProviders: []string{"cucp.openairinterface.org", "cuup.openairinterface.org"},
			providerToLook:  "du.openairinterface.org",
			want:            "error",
		},
		"Json Unmarshal Error": {
			configProviders: []string{"generate-json-error", "cucp.openairinterface.org", "du.openairinterface.org"},
			providerToLook:  "du.openairinterface.org",
			want:            "error",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := getConfigInstanceByProvider(logger, generateConfigInstanceForTesting(tc.configProviders), tc.providerToLook)
			if tc.want != "error" {
				wanted := workloadv1alpha1.NFDeployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "nf-" + tc.want,
						Namespace: "nf-dummy-du-ns",
					},
					Spec: workloadv1alpha1.NFDeploymentSpec{
						Provider: tc.want,
					},
				}
				if !reflect.DeepEqual(*got, wanted) {
					t.Errorf("getConfigInstanceByProvider returns %v | Wanted %v", got, wanted)
				}
			} else {
				if got != nil {
					t.Errorf("getConfigInstanceByProvider returns %v | Wanted nil", got)
				}
			}

		})
	}
}

func TestCheckMandatoryKinds(t *testing.T) {
	cases := map[string]struct {
		configSelfInfo map[string]runtime.RawExtension
		want           bool
	}{
		"Normal": {
			configSelfInfo: map[string]runtime.RawExtension{
				"PLMN":      runtime.RawExtension{},
				"RANConfig": runtime.RawExtension{},
				"OAIConfig": runtime.RawExtension{},
			},
			want: true,
		},
		"PLMN missing": {
			configSelfInfo: map[string]runtime.RawExtension{
				"RANConfig": runtime.RawExtension{},
			},
			want: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := CheckMandatoryKinds(tc.configSelfInfo)
			if got != tc.want {
				t.Errorf("CheckMandatoryKinds Returned %t wanted %t", got, tc.want)
			}

		})
	}

}
