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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	workloadnfconfig "workload.nephio.org/ran_deployment/api/v1alpha1"
)

type DuResources struct {
}

func (resource DuResources) createNetworkAttachmentDefinitionNetworks(templateName string, ranDeploymentSpec *workloadv1alpha1.NFDeploymentSpec) (string, error) {
	return CreateNetworkAttachmentDefinitionNetworks(templateName, map[string][]workloadv1alpha1.InterfaceConfig{
		"f1": GetInterfaceConfigs(ranDeploymentSpec.Interfaces, "f1"),
	})
}

func (resource DuResources) GetConfigMap(log logr.Logger, ranDeployment *workloadv1alpha1.NFDeployment, configInfo *ConfigInfo) []*corev1.ConfigMap {

	f1cIp, err := GetFirstInterfaceConfigIPv4(ranDeployment.Spec.Interfaces, "f1")
	if err != nil {
		log.Error(err, "Interface f1-du not found in RANDeployment Spec")
		return nil
	}

	quotedF1Ip := strconv.Quote(f1cIp)

	ranDeploymentConfigRef := getConfigInstanceByProvider(log, configInfo.ConfigRefInfo["NFDeployment"], "cucp.openairinterface.org")

	cuCpIp, err := GetFirstInterfaceConfigIPv4(ranDeploymentConfigRef.Spec.Interfaces, "f1c")
	if err != nil {
		log.Error(err, "f1c not found in Config Refs RANDeployment")
		return nil
	}

	quotedCuCpIp := strconv.Quote(cuCpIp)

	paramsRanNf := &workloadnfconfig.RANConfig{}
	if err := json.Unmarshal(configInfo.ConfigSelfInfo["RANConfig"].Raw, paramsRanNf); err != nil {
		log.Error(err, "Cannot Unmarshal RANConfig")
		return nil
	}

	paramsPlmn := &workloadnfconfig.PLMN{}
	if err := json.Unmarshal(configInfo.ConfigSelfInfo["PLMN"].Raw, paramsPlmn); err != nil {
		log.Error(err, "Cannot Unmarshal PLMN")
		return nil
	}

	templateValues := configurationTemplateValuesForDu{
		F1C_DU_IP:       quotedF1Ip,
		F1C_CU_IP:       quotedCuCpIp,
		TAC:             paramsPlmn.Spec.PLMNInfo[0].TAC,
		CELL_ID:         paramsRanNf.Spec.CellIdentity,
		PHY_CELL_ID:     paramsRanNf.Spec.PhysicalCellID,
		DL_FREQ_BAND:    paramsRanNf.Spec.DownlinkFrequencyBand,
		DL_SCS:          paramsRanNf.Spec.DownlinkSubCarrierSpacing,
		DL_CARRIER_BW:   paramsRanNf.Spec.DownlinkCarrierBandwidth,
		UL_FREQ_BAND:    paramsRanNf.Spec.UplinkFrequencyBand,
		UL_SCS:          paramsRanNf.Spec.UplinkSubCarrierSpacing,
		UL_CARRIER_BW:   paramsRanNf.Spec.UplinkCarrierBandwidth,
		PLMN_MCC:        paramsPlmn.Spec.PLMNInfo[0].PLMNID.MCC,
		PLMN_MNC:        paramsPlmn.Spec.PLMNInfo[0].PLMNID.MNC,
		PLMN_MNC_LENGTH: strconv.Itoa(int(len(paramsPlmn.Spec.PLMNInfo[0].PLMNID.MNC))),
		NSSAI_SST:       paramsPlmn.Spec.PLMNInfo[0].NSSAI[0].SST,
		NSSAI_SD:        *paramsPlmn.Spec.PLMNInfo[0].NSSAI[0].SD,
	}

	configuration, err := renderConfigurationTemplateForDu(templateValues)
	if err != nil {
		log.Error(err, "Could not render CU CP configuration template.")
		return nil
	}

	configMap1 := &corev1.ConfigMap{
		Data: map[string]string{
			"gnb.conf": configuration,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "oai-gnb-du-configmap",
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
	}

	return []*corev1.ConfigMap{configMap1}
}

func (resource DuResources) GetDeployment(log logr.Logger, ranDeployment *workloadv1alpha1.NFDeployment, configInfo *ConfigInfo) []*appsv1.Deployment {
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
				"app.kubernetes.io/name": "oai-gnb-du",
			},
			Name: "oai-gnb-du",
		},
		Spec: appsv1.DeploymentSpec{
			Paused: false,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/name": "oai-gnb-du",
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.DeploymentStrategyType("Recreate"),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: podAnnotations,
					Labels: map[string]string{
						"app":                    "oai-gnb-du",
						"app.kubernetes.io/name": "oai-gnb-du",
					},
				},
				Spec: corev1.PodSpec{
					HostIPC:                       false,
					HostNetwork:                   false,
					ServiceAccountName:            "oai-gnb-du-sa",
					TerminationGracePeriodSeconds: pointer.Int64(5),
					Volumes: []corev1.Volume{

						corev1.Volume{
							Name: "configuration",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "oai-gnb-du-configmap",
									},
								},
							},
						},
					},
					Containers: []corev1.Container{

						corev1.Container{
							Env: []corev1.EnvVar{
								corev1.EnvVar{
									Name: "USE_ADDITIONAL_OPTIONS",
									Value: "--sa --rfsim --log_config.global_log_options level,nocolor,time" +
										// " --MACRLCs.[0].local_n_address 192.168.73.3" +
										// " --MACRLCs.[0].remote_n_address 192.168.72.2" +
										" --telnetsrv --telnetsrv.shrmod o1 --telnetsrv.listenaddr 192.168.74.2",
								},
								corev1.EnvVar{
									Name:  "ASAN_OPTIONS",
									Value: "detect_leaks=0",
								},
							},
							Image: paramsOAI.Spec.Image,
							// Image: "arorasagar/testing-images:oai-gnb-telnet",
							// Image:   "nginx:latest",
							// Command: []string{"tail", "-f", "dev/null"},
							Ports: []corev1.ContainerPort{

								corev1.ContainerPort{
									ContainerPort: 38472,
									Name:          "f1c",
									Protocol:      corev1.Protocol("SCTP"),
								},
								corev1.ContainerPort{
									ContainerPort: 2152,
									Name:          "f1u",
									Protocol:      corev1.Protocol("UDP"),
								},
							},
							Stdin: false,
							TTY:   false,
							VolumeMounts: []corev1.VolumeMount{

								corev1.VolumeMount{
									Name:      "configuration",
									ReadOnly:  false,
									SubPath:   "gnb.conf",
									MountPath: "/opt/oai-gnb/etc/gnb.conf",
								},
							},
							Name: "gnbdu",
							SecurityContext: &corev1.SecurityContext{
								Privileged: pointer.Bool(true),
							},
							StdinOnce: false,
						},
					},
					DNSPolicy:     corev1.DNSPolicy("ClusterFirst"),
					HostPID:       false,
					RestartPolicy: corev1.RestartPolicy("Always"),
					SchedulerName: "default-scheduler",
				},
			},
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
	}

	return []*appsv1.Deployment{deployment1}
}

func (resource DuResources) GetServiceAccount() []*corev1.ServiceAccount {

	serviceAccount1 := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: "oai-gnb-du-sa",
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
	}

	return []*corev1.ServiceAccount{serviceAccount1}
}

func (resource DuResources) GetService() []*corev1.Service {

	service1 := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app.kubernetes.io/name": "oai-gnb-du",
			},
			Name: "oai-gnb-du",
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app.kubernetes.io/name": "oai-gnb-du",
			},
			Type:      corev1.ServiceType("ClusterIP"),
			ClusterIP: "None",
			Ports: []corev1.ServicePort{

				corev1.ServicePort{
					Name:     "f1c",
					Port:     38472,
					Protocol: corev1.Protocol("SCTP"),
					TargetPort: intstr.IntOrString{
						IntVal: 38472,
					},
				},
				corev1.ServicePort{
					Name:     "f1u",
					Port:     2152,
					Protocol: corev1.Protocol("UDP"),
					TargetPort: intstr.IntOrString{
						IntVal: 2152,
					},
				},
				corev1.ServicePort{
					Name:     "rfsim",
					Port:     4043,
					Protocol: corev1.Protocol("UDP"),
					TargetPort: intstr.IntOrString{
						IntVal: 4043,
					},
				},
			},
			PublishNotReadyAddresses: false,
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
	}

	// O1-Telent Service
	service2 := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app.kubernetes.io/name": "oai-gnb-du-o1-telnet",
			},
			Name: "oai-gnb-du-o1-telnet",
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app.kubernetes.io/name": "oai-gnb-du",
			},
			Type:      corev1.ServiceType("ClusterIP"),
			ClusterIP: "None",
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Port:     9090,
					Protocol: corev1.Protocol("TCP"),
					TargetPort: intstr.IntOrString{
						IntVal: 9090,
					},
				},
			},
		},
	}

	return []*corev1.Service{service1, service2}
}
