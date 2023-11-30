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

	"github.com/go-logr/logr"
	configref "github.com/nephio-project/api/references/v1alpha1"
	workloadv1alpha1 "github.com/nephio-project/api/workload/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
)

func getConfigInstanceByProvider(log logr.Logger, configInstances []*configref.Config, provider string) *workloadv1alpha1.NFDeployment {
	for _, configRef := range configInstances {
		b := configRef.Spec.Config.Raw
		nfDeployment := &workloadv1alpha1.NFDeployment{}
		if err := json.Unmarshal(b, nfDeployment); err != nil {
			log.Error(err, "Cannot Unmarshal NFDeployment")
			return nil
		}
		if nfDeployment.Spec.Provider == provider {
			return nfDeployment
		}
	}
	log.Error(fmt.Errorf("Provider %s not found", provider), "Cannot find provider in Config NFDeployment")
	return nil
}

func CheckMandatoryKinds(configSelfInfo map[string]runtime.RawExtension) bool {

	for _, kind := range GetMandatoryNfKinds() {
		if _, available := configSelfInfo[kind]; !available {
			return available
		}
	}
	return true
}
