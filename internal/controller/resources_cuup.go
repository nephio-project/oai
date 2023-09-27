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
	"strconv"

	"github.com/go-logr/logr"
	configref "github.com/nephio-project/api/references/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	workloadnephioorgv1alpha1 "workload.nephio.org/ran_deployment/api/v1alpha1"
)

type CuUpResources struct {
}

func (resource CuUpResources) createNetworkAttachmentDefinitionNetworks(templateName string, ranDeploymentSpec *workloadnephioorgv1alpha1.RANDeploymentSpec) (string, error) {
	return CreateNetworkAttachmentDefinitionNetworksConfigs(templateName, map[string][]workloadnephioorgv1alpha1.InterfaceConfig{
		"e1":  GetInterfaceConfigs(ranDeploymentSpec.Interfaces, "e1-cu-up"),
		"n3":  GetInterfaceConfigs(ranDeploymentSpec.Interfaces, "n3"),
		"f1u": GetInterfaceConfigs(ranDeploymentSpec.Interfaces, "f1u"),
	})
}
func (resource CuUpResources) GetDeployment(ranDeployment *workloadnephioorgv1alpha1.RANDeployment) []*appsv1.Deployment {

	spec := ranDeployment.Spec

	networkAttachmentDefinitionNetworks, err := resource.createNetworkAttachmentDefinitionNetworks("oai-ran-cu-up", &spec)

	if err != nil {
		return nil
	}

	podAnnotations := make(map[string]string)
	podAnnotations[NetworksAnnotation] = networkAttachmentDefinitionNetworks

	deployment1 := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app.kubernetes.io/name": "oai-gnb-cu-up",
			},
			Name: "oai-gnb-cu-up",
		},
		Spec: appsv1.DeploymentSpec{
			Paused:   false,
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/name": "oai-gnb-cu-up",
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.DeploymentStrategyType("Recreate"),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: podAnnotations,
					Labels: map[string]string{
						"app":                    "oai-gnb-cu-up",
						"app.kubernetes.io/name": "oai-gnb-cu-up",
					},
				},
				Spec: corev1.PodSpec{
					HostIPC:                       false,
					HostNetwork:                   false,
					HostPID:                       false,
					RestartPolicy:                 corev1.RestartPolicy("Always"),
					TerminationGracePeriodSeconds: int64Ptr(5),
					Containers: []corev1.Container{

						corev1.Container{
							VolumeMounts: []corev1.VolumeMount{

								corev1.VolumeMount{
									MountPath: "/opt/oai-gnb/etc/mounted.conf",
									Name:      "configuration",
									ReadOnly:  false,
									SubPath:   "mounted.conf",
								},
							},
							Env: []corev1.EnvVar{

								corev1.EnvVar{
									Name:  "TZ",
									Value: "Europe/Paris",
								},
								corev1.EnvVar{
									Name:  "USE_ADDITIONAL_OPTIONS",
									Value: "--sa",
								},
								corev1.EnvVar{
									Name:  "USE_VOLUMED_CONF",
									Value: "yes",
								},
							},
							Ports: []corev1.ContainerPort{

								corev1.ContainerPort{
									ContainerPort: 2152,
									Name:          "n3",
									Protocol:      corev1.Protocol("UDP"),
								},
								corev1.ContainerPort{
									ContainerPort: 38462,
									Name:          "e1",
									Protocol:      corev1.Protocol("SCTP"),
								},
								corev1.ContainerPort{
									ContainerPort: 2152,
									Name:          "f1u",
									Protocol:      corev1.Protocol("UDP"),
								},
							},
							SecurityContext: &corev1.SecurityContext{
								Privileged: boolPtr(true),
							},
							Stdin:     false,
							StdinOnce: false,
							TTY:       false,
							Image:     "docker.io/oaisoftwarealliance/oai-nr-cuup:2023.w19",
							Name:      "gnbcuup",
						},
					},
					DNSPolicy:          corev1.DNSPolicy("ClusterFirst"),
					SchedulerName:      "default-scheduler",
					ServiceAccountName: "oai-gnb-cu-up-sa",
					Volumes: []corev1.Volume{

						corev1.Volume{
							Name: "configuration",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "oai-gnb-cu-up-configmap",
									},
								},
							},
						},
					},
				},
			},
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
	}

	return []*appsv1.Deployment{deployment1}
}

func (resource CuUpResources) GetServiceAccount() []*corev1.ServiceAccount {

	serviceAccount1 := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: "oai-gnb-cu-up-sa",
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
	}

	return []*corev1.ServiceAccount{serviceAccount1}
}

func (resource CuUpResources) GetConfigMap(log logr.Logger, ranDeployment *workloadnephioorgv1alpha1.RANDeployment, configInstancesMap map[string]*configref.Config) []*corev1.ConfigMap {

	quotedN3Ip, err := GetFirstInterfaceConfigIPv4(ranDeployment.Spec.Interfaces, "n3", true)
	if err != nil {
		log.Error(err, "Interface n2 not found in RANDeployment Spec")
		return nil
	}

	quotedE1Ip, err := GetFirstInterfaceConfigIPv4(ranDeployment.Spec.Interfaces, "e1-cu-up", true)
	if err != nil {
		log.Error(err, "Interface e1 not found in RANDeployment Spec")
		return nil
	}

	quotedF1UIp, err := GetFirstInterfaceConfigIPv4(ranDeployment.Spec.Interfaces, "f1u", true)
	if err != nil {
		log.Error(err, "Interface f1c not found in RANDeployment Spec")
		return nil
	}

	var b []byte
	ranDeploymentConfigRef := &workloadnephioorgv1alpha1.RANDeployment{}
	b = configInstancesMap["RANDeployment"].Spec.Config.Raw
	if err := json.Unmarshal(b, ranDeploymentConfigRef); err != nil {
		log.Error(err, "Cannot Unmarshal RANDeployment")
		return nil
	}

	cuCpIp, err := GetFirstInterfaceConfigIPv4(ranDeploymentConfigRef.Spec.Interfaces, "e1", false)
	if err != nil {
		log.Error(err, "AMF IP not found in Config Refs AMFDeployment")
		return nil
	}

	quotedCuCpIp := strconv.Quote(cuCpIp)

	templateValues := configurationTemplateValuesForCuUp{
		E1_IP:   quotedE1Ip,
		F1U_IP:  quotedF1UIp,
		N3_IP:   quotedN3Ip,
		CUCP_E1: quotedCuCpIp,
	}

	configuration, err := renderConfigurationTemplateForCuUp(templateValues)
	if err != nil {
		log.Error(err, "Could not render CU UP configuration template.")
		return nil
	}

	configMap1 := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "oai-gnb-cu-up-configmap",
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		Data: map[string]string{
			"mounted.conf": configuration,
		},
	}

	return []*corev1.ConfigMap{configMap1}
}
