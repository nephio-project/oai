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
	workloadv1alpha1 "github.com/nephio-project/api/workload/v1alpha1"
	workloadnfconfig "workload.nephio.org/ran_deployment/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

type CuUpResources struct {
}

func (resource CuUpResources) createNetworkAttachmentDefinitionNetworks(templateName string, ranDeploymentSpec *workloadv1alpha1.NFDeploymentSpec) (string, error) {
	return CreateNetworkAttachmentDefinitionNetworks(templateName, map[string][]workloadv1alpha1.InterfaceConfig{
		"e1":  GetInterfaceConfigs(ranDeploymentSpec.Interfaces, "e1"),
		"n3":  GetInterfaceConfigs(ranDeploymentSpec.Interfaces, "n3"),
		"f1u": GetInterfaceConfigs(ranDeploymentSpec.Interfaces, "f1u"),
	})
}
func (resource CuUpResources) GetDeployment(log logr.Logger, ranDeployment *workloadv1alpha1.NFDeployment, configInfo *ConfigInfo) []*appsv1.Deployment {

	spec := ranDeployment.Spec

	networkAttachmentDefinitionNetworks, err := resource.createNetworkAttachmentDefinitionNetworks(ranDeployment.Name, &spec)

	if err != nil {
		return nil
	}

	paramsOAI := &workloadnfconfig.OAIConfig{}
	if err := json.Unmarshal(configInfo.ConfigSelfInfo["OAIConfig"].Raw, paramsOAI); err != nil {
		log.Error(err, "Cannot Unmarshal OAIConfig")
		return nil
	}

	podAnnotations := make(map[string]string)
	podAnnotations[NetworksAnnotation] = networkAttachmentDefinitionNetworks

	deployment1 := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app.kubernetes.io/name": "oai-cu-up",
			},
			Name: "oai-cu-up",
		},
		Spec: appsv1.DeploymentSpec{
			Paused:   false,
			Replicas: ptr.To(int32(1)),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/name": "oai-cu-up",
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.DeploymentStrategyType("Recreate"),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: podAnnotations,
					Labels: map[string]string{
						"app":                    "oai-cu-up",
						"app.kubernetes.io/name": "oai-cu-up",
					},
				},
				Spec: corev1.PodSpec{
					HostIPC:                       false,
					HostNetwork:                   false,
					HostPID:                       false,
					RestartPolicy:                 corev1.RestartPolicy("Always"),
					TerminationGracePeriodSeconds: ptr.To(int64(5)),
					Containers: []corev1.Container{

						corev1.Container{
							VolumeMounts: []corev1.VolumeMount{

								corev1.VolumeMount{
									MountPath: "/opt/oai-gnb/etc/gnb.conf",
									Name:      "configuration",
									ReadOnly:  false,
									SubPath:   "gnb.conf",
								},
							},
							Env: []corev1.EnvVar{

								corev1.EnvVar{
									Name:  "TZ",
									Value: "Europe/Paris",
								},
								corev1.EnvVar{
									Name:  "USE_ADDITIONAL_OPTIONS",
									Value: "--sa --log_config.global_log_options level,nocolor,time",
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
								Privileged: ptr.To(true),
							},
							Stdin:     false,
							StdinOnce: false,
							TTY:       false,
							Image:     paramsOAI.Spec.Image,
							Name:      "cuup",
						},
					},
					DNSPolicy:          corev1.DNSPolicy("ClusterFirst"),
					SchedulerName:      "default-scheduler",
					ServiceAccountName: "oai-cu-up-sa",
					Volumes: []corev1.Volume{

						corev1.Volume{
							Name: "configuration",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "oai-cu-up-configmap",
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
			Name: "oai-cu-up-sa",
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
	}

	return []*corev1.ServiceAccount{serviceAccount1}
}

func (resource CuUpResources) GetConfigMap(log logr.Logger, ranDeployment *workloadv1alpha1.NFDeployment, configInfo *ConfigInfo) []*corev1.ConfigMap {

	n3Ip, err := GetFirstInterfaceConfigIPv4(ranDeployment.Spec.Interfaces, "n3")
	if err != nil {
		log.Error(err, "Interface n3 not found in RANDeployment Spec")
		return nil
	}

	quotedN3Ip := strconv.Quote(n3Ip)

	e1Ip, err := GetFirstInterfaceConfigIPv4(ranDeployment.Spec.Interfaces, "e1")
	if err != nil {
		log.Error(err, "Interface e1 not found in RANDeployment Spec")
		return nil
	}

	quotedE1Ip := strconv.Quote(e1Ip)

	f1uIp, err := GetFirstInterfaceConfigIPv4(ranDeployment.Spec.Interfaces, "f1u")
	if err != nil {
		log.Error(err, "Interface F1 U not found in RANDeployment Spec")
		return nil
	}

	quotedF1UIp := strconv.Quote(f1uIp)

	ranDeploymentConfigRef := getConfigInstanceByProvider(log, configInfo.ConfigRefInfo["NFDeployment"], "cucp.openairinterface.org")

	cuCpIp, err := GetFirstInterfaceConfigIPv4(ranDeploymentConfigRef.Spec.Interfaces, "e1")
	if err != nil {
		log.Error(err, "CU CP IP not found in Config Refs RANDeployment")
		return nil
	}

	quotedCuCpIp := strconv.Quote(cuCpIp)

	paramsPlmn := &workloadnfconfig.PLMN{}
	if err := json.Unmarshal(configInfo.ConfigSelfInfo["PLMN"].Raw, paramsPlmn); err != nil {
		log.Error(err, "Cannot Unmarshal PLMN")
		return nil
	}

	templateValues := configurationTemplateValuesForCuUp{
		E1_IP:           quotedE1Ip,
		F1U_IP:          quotedF1UIp,
		N3_IP:           quotedN3Ip,
		CUCP_E1:         quotedCuCpIp,
		TAC:             paramsPlmn.Spec.PLMNInfo[0].TAC,
		PLMN_MCC:        paramsPlmn.Spec.PLMNInfo[0].PLMNID.MCC,
		PLMN_MNC:        paramsPlmn.Spec.PLMNInfo[0].PLMNID.MNC,
		PLMN_MNC_LENGTH: strconv.Itoa(int(len(paramsPlmn.Spec.PLMNInfo[0].PLMNID.MNC))),
		NSSAI_SST:       paramsPlmn.Spec.PLMNInfo[0].NSSAI[0].SST,
		NSSAI_SD:        *paramsPlmn.Spec.PLMNInfo[0].NSSAI[0].SD,
	}

	configuration, err := renderConfigurationTemplateForCuUp(templateValues)
	if err != nil {
		log.Error(err, "Could not render CU UP configuration template.")
		return nil
	}

	configMap1 := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "oai-cu-up-configmap",
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		Data: map[string]string{
			"gnb.conf": configuration,
		},
	}

	return []*corev1.ConfigMap{configMap1}
}

func (resource CuUpResources) GetService() []*corev1.Service {
	return []*corev1.Service{}
}
