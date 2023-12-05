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
	context "context"
	"errors"
	"fmt"
	"testing"
	"time"

	configref "github.com/nephio-project/api/references/v1alpha1"
	workloadv1alpha1 "github.com/nephio-project/api/workload/v1alpha1"
	"github.com/stretchr/testify/mock"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func TestGetConfigs(t *testing.T) {
	cases := map[string]struct {
		ranDeploymentParameterRef []workloadv1alpha1.ObjectReference // ParameterRefs Required by Your NF
		mock3rdArgsType           string                             // Represents the type of 3rd argument in r.Get
		mockReturnRefVal          interface{}                        // It will copy this value to the 3rd argument in r.Get
		mockReturnError           error                              // Tells if r.Get needs to return error or not
		wantError                 string                             // Tells the overall testFxn (GetConfig) will return an error or not
	}{
		"Normal | Api-version: ref.nephio.org/v1alpha1 ": {
			ranDeploymentParameterRef: []workloadv1alpha1.ObjectReference{{
				APIVersion: "ref.nephio.org/v1alpha1",
				Name:       pointer.String("ABC"),
			}},
			mock3rdArgsType: "*v1alpha1.Config",
			mockReturnRefVal: &configref.Config{Spec: configref.ConfigSpec{
				Config: runtime.RawExtension{Raw: marshalJsonReturnByteOnly(map[string]any{"kind": "Dummy-Kind"})},
			}},
			mockReturnError: nil,
			wantError:       "nil",
		},
		"Mock-Error (k8s not able to get the object as requested by ranDeploymentParameterRef)| Api-version: ref.nephio.org/v1alpha1 ": {
			ranDeploymentParameterRef: []workloadv1alpha1.ObjectReference{{
				APIVersion: "ref.nephio.org/v1alpha1",
				Name:       pointer.String("ABC"),
			}},
			mock3rdArgsType:  "*v1alpha1.Config",
			mockReturnRefVal: &configref.Config{},
			mockReturnError:  errors.New("Requested Object Not Found"),
			wantError:        "Mock-Error (k8s not able to get the object as requested by ranDeploymentParameterRef)",
		},
		"Mock-Returned-Reference-UnMarshal-Error | Api-version: ref.nephio.org/v1alpha1 ": {
			ranDeploymentParameterRef: []workloadv1alpha1.ObjectReference{{
				APIVersion: "ref.nephio.org/v1alpha1",
				Name:       pointer.String("ABC"),
			}},
			mock3rdArgsType: "*v1alpha1.Config",
			mockReturnRefVal: &configref.Config{Spec: configref.ConfigSpec{
				Config: runtime.RawExtension{Raw: []byte(" ")},
			}},
			mockReturnError: nil,
			wantError:       "Mock-Returned-Reference-UnMarshal-Error",
		},
		"Normal | Api-version: workload.nephio.org/v1alpha1 ": {
			ranDeploymentParameterRef: []workloadv1alpha1.ObjectReference{{
				APIVersion: "workload.nephio.org/v1alpha1",
				Name:       pointer.String("ABC"),
			}},
			mock3rdArgsType: "*v1alpha1.NFConfig",
			mockReturnRefVal: &workloadv1alpha1.NFConfig{Spec: workloadv1alpha1.NFConfigSpec{
				ConfigRefs: []runtime.RawExtension{
					{Raw: marshalJsonReturnByteOnly(map[string]any{"kind": "PLMN"})},
					{Raw: marshalJsonReturnByteOnly(map[string]any{"kind": "RANConfig"})},
					{Raw: marshalJsonReturnByteOnly(map[string]any{"kind": "OAIConfig"})},
				},
			}},
			mockReturnError: nil,
			wantError:       "nil",
		},
		"Mock-Error (k8s not able to get the object as requested by ranDeploymentParameterRef) | Api-version: workload.nephio.org/v1alpha1 ": {
			ranDeploymentParameterRef: []workloadv1alpha1.ObjectReference{{
				APIVersion: "workload.nephio.org/v1alpha1",
				Name:       pointer.String("ABC"),
			}},
			mock3rdArgsType:  "*v1alpha1.NFConfig",
			mockReturnRefVal: &workloadv1alpha1.NFConfig{},
			mockReturnError:  errors.New("Requested Object Not Found"),
			wantError:        "Mock-Error (k8s not able to get the object as requested by ranDeploymentParameterRef)",
		},
		"Mock-Returned-Reference-UnMarshal-Error | Api-version: workload.nephio.org/v1alpha1 ": {
			ranDeploymentParameterRef: []workloadv1alpha1.ObjectReference{{
				APIVersion: "workload.nephio.org/v1alpha1",
				Name:       pointer.String("ABC"),
			}},
			mock3rdArgsType: "*v1alpha1.NFConfig",
			mockReturnRefVal: &workloadv1alpha1.NFConfig{Spec: workloadv1alpha1.NFConfigSpec{
				ConfigRefs: []runtime.RawExtension{
					{Raw: []byte(" ")},
					{Raw: marshalJsonReturnByteOnly(map[string]any{"kind": "RanConfig"})},
				},
			}},
			mockReturnError: nil,
			wantError:       "Mock-Returned-Reference-UnMarshal-Error",
		},
		"Mandatory Kinds Not Found in Mock-Returned-Reference| Api-version: workload.nephio.org/v1alpha1 ": {
			ranDeploymentParameterRef: []workloadv1alpha1.ObjectReference{{
				APIVersion: "workload.nephio.org/v1alpha1",
				Name:       pointer.String("ABC"),
			}},
			mock3rdArgsType: "*v1alpha1.NFConfig",
			mockReturnRefVal: &workloadv1alpha1.NFConfig{Spec: workloadv1alpha1.NFConfigSpec{
				ConfigRefs: []runtime.RawExtension{
					{Raw: marshalJsonReturnByteOnly(map[string]any{"kind": "RANConfig"})}, // PLMN & OAIConfig removed
				},
			}},
			mockReturnError: nil,
			wantError:       "Mandatory Kinds Not Found in Mock-Returned-Reference",
		},
		"Not-Supported API-Version Error": {
			ranDeploymentParameterRef: []workloadv1alpha1.ObjectReference{{
				APIVersion: "dummy-api",
				Name:       pointer.String("ABC"),
			}},
			mock3rdArgsType:  "",
			mockReturnRefVal: nil,
			mockReturnError:  nil,
			wantError:        "Not-Supported API-Version Error",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			clientMock := new(MockClient)
			clientMock.On("Get", context.TODO(), mock.AnythingOfType("types.NamespacedName"), mock.AnythingOfType(tc.mock3rdArgsType)).Return(tc.mockReturnError).Run(func(args mock.Arguments) {
				if tc.mock3rdArgsType == "*v1alpha1.Config" {
					configObj := args.Get(2).(*configref.Config)
					mockReturnVal, ok := tc.mockReturnRefVal.(*configref.Config)
					if !ok {
						t.Errorf("Test-Case not properly written | mockReturnRefVal should be of type *configref.Config")
					}
					*configObj = *mockReturnVal // mockReturnVal is what r.Get will store in 3rd Argument
				} else if tc.mock3rdArgsType == "*v1alpha1.NFConfig" {
					configObj := args.Get(2).(*workloadv1alpha1.NFConfig)
					mockReturnVal, ok := tc.mockReturnRefVal.(*workloadv1alpha1.NFConfig)
					if !ok {
						t.Errorf("Test-Case not properly written | mockReturnRefVal should be of type *workloadv1alpha1.NFConfig")
					}
					*configObj = *mockReturnVal // mockReturnVal is what r.Get will store in 3rd Argument
				}

			})

			ranReconcilerObj := RANDeploymentReconciler{
				clientMock,
				runtime.NewScheme(),
			}

			ranDeploymentObj := workloadv1alpha1.NFDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ranNf",
					Namespace: "myns",
				},
				Spec: workloadv1alpha1.NFDeploymentSpec{
					ParametersRefs: tc.ranDeploymentParameterRef,
				},
			}

			got, err := ranReconcilerObj.GetConfigs(context.TODO(), &ranDeploymentObj)
			if tc.wantError == "nil" {
				if err != nil {
					t.Errorf("GetConfigs Returned Error %v", err)
				}
				if tc.mock3rdArgsType == "*v1alpha1.Config" {
					// Your got Object should have ConfigRefInfo Set, with one kind set to "Dummy-Kind" (which was set while setting the test-case in mockReturnRefVal)
					testPassed := false
					for kind := range got.ConfigRefInfo {
						if kind == "Dummy-Kind" {
							testPassed = true
							break
						}
					}
					if !testPassed {
						t.Errorf("Dummy-Kind is not Present in Kind at Attribute ConfigRefInfo of the returned %v ", got)
					}

				} else if tc.mock3rdArgsType == "*v1alpha1.NFConfig" {
					// Mandatory kinds should be present in Got-ConfigSelfInfo
					if !CheckMandatoryKinds(got.ConfigSelfInfo) {
						t.Errorf("GetConfigs Returned %v which don't have mandatory Kinds ", got)
					}
				}

			} else {
				if err == nil {
					t.Errorf("GetConfigs Returned no Error when expected error %s", tc.wantError)
				}
			}

		})
	}

}

func TestCreateAll(t *testing.T) {
	/*
		Now Since the NfResource Methods Are unitested, so in this It only make sense to test the r.Create method and its error scenarios
	*/
	cases := map[string]struct {
		errorGivingMethodIndex int // It represents the index of nfResourceMethods that will give an error
	}{
		"Normal":                           {errorGivingMethodIndex: -1},
		"Service Account Failed to Create": {errorGivingMethodIndex: 0},
		"ConfigMap Failed to Create":       {errorGivingMethodIndex: 1},
		"Deployment Failed to Create":      {errorGivingMethodIndex: 2},
		"Service Failed to Create":         {errorGivingMethodIndex: 3},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {

			nfResourceMethods := []string{"GetServiceAccount", "GetConfigMap", "GetDeployment", "GetService"}
			methodArguments := [][]string{{}, {"Logger", "*v1alpha1.NFDeployment", "*controller.ConfigInfo"}, {"Logger", "*v1alpha1.NFDeployment", "*controller.ConfigInfo"}, {}}
			returnTypes := []string{"*v1.ServiceAccount", "*v1.ConfigMap", "*v1.Deployment", "*v1.Service"}

			clientMock := new(MockClient)
			for i := 0; i < len(nfResourceMethods); i++ {
				if tc.errorGivingMethodIndex == i {
					clientMock.On("Create", context.TODO(), mock.AnythingOfType(returnTypes[i])).Return(errors.New("Unable to create the resource"))
				} else {
					clientMock.On("Create", context.TODO(), mock.AnythingOfType(returnTypes[i])).Return(nil)
				}

			}

			ranReconcilerObj := RANDeploymentReconciler{
				clientMock,
				runtime.NewScheme(),
			}

			nfResourceMock := new(MockNfResource)
			for methodIndex, methodName := range nfResourceMethods {
				call := nfResourceMock.Mock.On(methodName)
				for _, arg := range methodArguments[methodIndex] {
					call.Arguments = append(call.Arguments, mock.AnythingOfType(arg))
				}
				switch methodIndex {
				case 0:
					call.ReturnArguments = append(call.ReturnArguments, []*corev1.ServiceAccount{{}})
				case 1:
					call.ReturnArguments = append(call.ReturnArguments, []*corev1.ConfigMap{{}})
				case 2:
					call.ReturnArguments = append(call.ReturnArguments, []*appsv1.Deployment{{}})
				case 3:
					call.ReturnArguments = append(call.ReturnArguments, []*corev1.Service{{}})
				}
			}

			ranReconcilerObj.CreateAll(context.TODO(), &workloadv1alpha1.NFDeployment{}, nfResourceMock, &ConfigInfo{})

		})
	}

}

func TestDeleteAll(t *testing.T) {
	/*
		Now Since the NfResource Methods Are unitested, so in this It only make sense to test the r.Delete method and its error scenarios
	*/
	cases := map[string]struct {
		errorGivingMethodIndex int // It represents the index of nfResourceMethods that will give an error
	}{
		"Normal":                           {errorGivingMethodIndex: -1},
		"Service Account Failed to Delete": {errorGivingMethodIndex: 0},
		"ConfigMap Failed to Delete":       {errorGivingMethodIndex: 1},
		"GetDeployment Failed to Delete":   {errorGivingMethodIndex: 2},
		"Service Failed to Delete":         {errorGivingMethodIndex: 3},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {

			nfResourceMethods := []string{"GetServiceAccount", "GetConfigMap", "GetDeployment", "GetService"}
			methodArguments := [][]string{{}, {"Logger", "*v1alpha1.NFDeployment", "*controller.ConfigInfo"}, {"Logger", "*v1alpha1.NFDeployment", "*controller.ConfigInfo"}, {}}
			returnTypes := []string{"*v1.ServiceAccount", "*v1.ConfigMap", "*v1.Deployment", "*v1.Service"}

			clientMock := new(MockClient)
			for i := 0; i < len(nfResourceMethods); i++ {
				if tc.errorGivingMethodIndex == i {
					clientMock.On("Delete", context.TODO(), mock.AnythingOfType(returnTypes[i])).Return(errors.New("Unable to create the resource"))
				} else {
					clientMock.On("Delete", context.TODO(), mock.AnythingOfType(returnTypes[i])).Return(nil)
				}

			}

			ranReconcilerObj := RANDeploymentReconciler{
				clientMock,
				runtime.NewScheme(),
			}

			nfResourceMock := new(MockNfResource)
			for methodIndex, methodName := range nfResourceMethods {
				call := nfResourceMock.Mock.On(methodName)
				for _, arg := range methodArguments[methodIndex] {
					call.Arguments = append(call.Arguments, mock.AnythingOfType(arg))
				}
				switch methodIndex {
				case 0:
					call.ReturnArguments = append(call.ReturnArguments, []*corev1.ServiceAccount{{}})
				case 1:
					call.ReturnArguments = append(call.ReturnArguments, []*corev1.ConfigMap{{}})
				case 2:
					call.ReturnArguments = append(call.ReturnArguments, []*appsv1.Deployment{{}})
				case 3:
					call.ReturnArguments = append(call.ReturnArguments, []*corev1.Service{{}})
				}
			}
			ranReconcilerObj.DeleteAll(context.TODO(), &workloadv1alpha1.NFDeployment{}, nfResourceMock, &ConfigInfo{})

		})
	}

}

func TestReconcileErrorScenarios(t *testing.T) {
	cases := map[string]struct {
		ranDeployment *workloadv1alpha1.NFDeployment
		mockReturnErr error
		expectedError error
	}{
		"Ran-deployment Object not Found in k8s Cluster": {
			ranDeployment: &workloadv1alpha1.NFDeployment{
				Spec: workloadv1alpha1.NFDeploymentSpec{},
			},
			mockReturnErr: errors.New("Not Found"),
			expectedError: errors.New("Not Found"),
		},
		"Ran-deployment Nf Doesn't have a provider field": {
			ranDeployment: &workloadv1alpha1.NFDeployment{
				Spec: workloadv1alpha1.NFDeploymentSpec{},
			},
			mockReturnErr: nil,
			expectedError: nil,
		},
		"Ran-deployment Nf Config-Refs Api-Version is not Supported": {
			ranDeployment: &workloadv1alpha1.NFDeployment{
				Spec: workloadv1alpha1.NFDeploymentSpec{
					Provider: "cucp.openairinterface.org",
					ParametersRefs: []workloadv1alpha1.ObjectReference{
						{
							APIVersion: "dummy-apiversion",
							Name:       pointer.String("ABC"),
						},
					},
				},
			},
			mockReturnErr: nil,
			expectedError: fmt.Errorf("Not supported API version \"dummy-apiversion\""),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			clientMock := new(MockClient)
			clientMock.On("Get", context.TODO(), mock.AnythingOfType("types.NamespacedName"), mock.AnythingOfType("*v1alpha1.NFDeployment")).Return(tc.mockReturnErr).Run(func(args mock.Arguments) {
				configObj := args.Get(2).(*workloadv1alpha1.NFDeployment)
				*configObj = *tc.ranDeployment // tc.ranDeployment is what r.Get will store in 3rd Argument
			})

			ranReconcilerObj := RANDeploymentReconciler{
				clientMock,
				runtime.NewScheme(),
			}

			_, err := ranReconcilerObj.Reconcile(context.TODO(), controllerruntime.Request{NamespacedName: types.NamespacedName{Namespace: "myns", Name: "mynf"}})
			if fmt.Sprint(err) != fmt.Sprint(tc.expectedError) {
				t.Errorf("Reconciled Error Scenario| Returned error %v | expected error %v", err, tc.expectedError)
			}
		})
	}

}

func TestReconcileCreateScenarios(t *testing.T) {
	cases := map[string]struct {
		ranDeployment *workloadv1alpha1.NFDeployment
		expectedError error
	}{
		"Create DU": {
			ranDeployment: &workloadv1alpha1.NFDeployment{
				Spec: workloadv1alpha1.NFDeploymentSpec{Provider: "du.openairinterface.org"},
			},
			expectedError: nil,
		},
		"Create CUCP": {
			ranDeployment: &workloadv1alpha1.NFDeployment{
				Spec: workloadv1alpha1.NFDeploymentSpec{Provider: "cucp.openairinterface.org"},
			},
			expectedError: nil,
		},
		"Create CUUP": {
			ranDeployment: &workloadv1alpha1.NFDeployment{
				Spec: workloadv1alpha1.NFDeploymentSpec{Provider: "cuup.openairinterface.org"},
			},
			expectedError: nil,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			clientMock := new(MockClient)
			clientMock.On("Get", context.TODO(), mock.AnythingOfType("types.NamespacedName"), mock.AnythingOfType("*v1alpha1.NFDeployment")).Return(nil).Run(func(args mock.Arguments) {
				configObj := args.Get(2).(*workloadv1alpha1.NFDeployment)
				*configObj = *tc.ranDeployment // tc.ranDeployment is what r.Get will store in 3rd Argument
			})
			/*
				Rationale Behind the following mock-methods
				GetDeployment, GetConfigMap of NfResource (du, cucp, cuup) will give nil Since we are not providing correct ranDeployment Spec-Values
				This is done because
				1) GetDeployment, GetConfigMap are separatly unit-tested for (corner-scenarios)
				2) Much Significant test would be the Integration test
			*/
			clientMock.On("Create", context.TODO(), mock.AnythingOfType("*v1.ServiceAccount")).Return(nil)     // For GetServiceAccount
			clientMock.On("Create", context.TODO(), mock.AnythingOfType("*v1.Service")).Return(nil)            // For GetService
			clientMock.On("Update", context.TODO(), mock.AnythingOfType("*v1alpha1.NFDeployment")).Return(nil) // For r.Update (whose significance is adding a finalizer)
			ranReconcilerObj := RANDeploymentReconciler{
				clientMock,
				runtime.NewScheme(),
			}
			_, err := ranReconcilerObj.Reconcile(context.TODO(), controllerruntime.Request{NamespacedName: types.NamespacedName{Namespace: "myns", Name: "mynf"}})
			if tc.expectedError == nil {
				if err != nil {
					t.Errorf("Reconcile During Creation gives Error %v while NO Error was expected", err)
				}
			}

		})
	}
}

func TestReconcileDeleteScenarios(t *testing.T) {
	cases := map[string]struct {
		ranDeployment *workloadv1alpha1.NFDeployment
		expectedError error
	}{
		"Delete DU": {
			ranDeployment: &workloadv1alpha1.NFDeployment{
				ObjectMeta: metav1.ObjectMeta{
					DeletionTimestamp: &metav1.Time{Time: time.Now()},
				},
				Spec: workloadv1alpha1.NFDeploymentSpec{Provider: "du.openairinterface.org"},
			},
			expectedError: nil,
		},
		"Delete CUCP": {
			ranDeployment: &workloadv1alpha1.NFDeployment{
				ObjectMeta: metav1.ObjectMeta{
					DeletionTimestamp: &metav1.Time{Time: time.Now()},
				},
				Spec: workloadv1alpha1.NFDeploymentSpec{Provider: "cucp.openairinterface.org"},
			},
			expectedError: nil,
		},
		"Delete CUUP": {
			ranDeployment: &workloadv1alpha1.NFDeployment{
				ObjectMeta: metav1.ObjectMeta{
					DeletionTimestamp: &metav1.Time{Time: time.Now()},
				},
				Spec: workloadv1alpha1.NFDeploymentSpec{Provider: "cuup.openairinterface.org"},
			},
			expectedError: nil,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			clientMock := new(MockClient)
			clientMock.On("Get", context.TODO(), mock.AnythingOfType("types.NamespacedName"), mock.AnythingOfType("*v1alpha1.NFDeployment")).Return(nil).Run(func(args mock.Arguments) {
				configObj := args.Get(2).(*workloadv1alpha1.NFDeployment)
				myFinalizerName := "batch.tutorial.kubebuilder.io/finalizer"
				controllerutil.AddFinalizer(tc.ranDeployment, myFinalizerName) // Simulates Deletion
				*configObj = *tc.ranDeployment                                 // tc.ranDeployment is what r.Get will store in 3rd Argument
			})
			/*
				Rationale Behind the following mock-methods
				GetDeployment, GetConfigMap of NfResource (du, cucp, cuup) will give nil Since we are not providing correct ranDeployment Spec-Values
				This is done because
				1) GetDeployment, GetConfigMap are separatly unit-tested for (corner-scenarios)
				2) Much Significant test would be the Integration test
			*/
			clientMock.On("Delete", context.TODO(), mock.AnythingOfType("*v1.ServiceAccount")).Return(nil)     // For GetServiceAccount
			clientMock.On("Delete", context.TODO(), mock.AnythingOfType("*v1.Service")).Return(nil)            // For GetService
			clientMock.On("Update", context.TODO(), mock.AnythingOfType("*v1alpha1.NFDeployment")).Return(nil) // For r.Update (whose significance is deleting the finalizer)
			ranReconcilerObj := RANDeploymentReconciler{
				clientMock,
				runtime.NewScheme(),
			}
			_, err := ranReconcilerObj.Reconcile(context.TODO(), controllerruntime.Request{NamespacedName: types.NamespacedName{Namespace: "myns", Name: "mynf"}})
			if tc.expectedError == nil {
				if err != nil {
					t.Errorf("Reconcile During Creation gives Error %v while NO Error was expected", err)
				}
			}

		})
	}
}
