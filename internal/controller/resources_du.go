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

type DuResources struct {
}

func (resource DuResources) createNetworkAttachmentDefinitionNetworks(templateName string, ranDeploymentSpec *workloadnephioorgv1alpha1.RANDeploymentSpec) (string, error) {
	return CreateNetworkAttachmentDefinitionNetworksConfigs(templateName, map[string][]workloadnephioorgv1alpha1.InterfaceConfig{
		"f1": GetInterfaceConfigs(ranDeploymentSpec.Interfaces, "f1-du"),
	})
}

func (resource DuResources) GetConfigMap(log logr.Logger, ranDeployment *workloadnephioorgv1alpha1.RANDeployment, configInstancesMap map[string]*configref.Config) []*corev1.ConfigMap {

	quotedF1Ip, err := GetFirstInterfaceConfigIPv4(ranDeployment.Spec.Interfaces, "f1-du", false)
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

	quotedCuCpIp, err := GetFirstInterfaceConfigIPv4(ranDeploymentConfigRef.Spec.Interfaces, "f1c", false)
	if err != nil {
		log.Error(err, "AMF IP not found in Config Refs AMFDeployment")
		return nil
	}

	configMap1 := &corev1.ConfigMap{
		Data: map[string]string{
			"nssaiSst":             ranDeployment.Spec.NssaiList[0].Sst,
			"timeZone":             "Europe/Paris",
			"f1IfName":             "f1",
			"f1cuIpAddress":        quotedCuCpIp,
			"gnbNgaIfName":         "eth0",
			"gnbNguIpAddress":      "status.podIP",
			"mcc":                  ranDeployment.Spec.Mcc,
			"mnc":                  ranDeployment.Spec.Mnc,
			"rfSimulator":          "server",
			"tac":                  ranDeployment.Spec.Tac,
			"amfIpAddress":         "127.0.0.1",
			"f1duIpAddress":        quotedF1Ip,
			"gnbNguIfName":         "eth0",
			"useSaTDDdu":           "yes",
			"mountConfig":          "false",
			"f1duPort":             "2152",
			"gnbduName":            "oai-du-rfsim",
			"mncLength":            strconv.Itoa(ranDeployment.Spec.MncLength),
			"useAdditionalOptions": "--sa --rfsim --log_config.global_log_options level,nocolor,time",
			"f1cuPort":             "2152",
			"gnbNgaIpAddress":      "status.podIP",
			"nssaiSd0":             ranDeployment.Spec.NssaiList[0].Sd,
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

func (resource DuResources) GetDeployment(ranDeployment *workloadnephioorgv1alpha1.RANDeployment) []*appsv1.Deployment {

	spec := ranDeployment.Spec

	networkAttachmentDefinitionNetworks, err := resource.createNetworkAttachmentDefinitionNetworks("oai-ran-du", &spec)

	if err != nil {
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
					TerminationGracePeriodSeconds: int64Ptr(5),
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
									Name:  "TZ",
									Value: "Europe/Paris",
								},
								corev1.EnvVar{
									Name:  "RFSIMULATOR",
									Value: "server",
								},
								corev1.EnvVar{
									Name:  "USE_ADDITIONAL_OPTIONS",
									Value: "--sa --rfsim --log_config.global_log_options level,nocolor,time",
								},
								corev1.EnvVar{
									Name: "USE_SA_TDD_DU",
									ValueFrom: &corev1.EnvVarSource{
										ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
											Key: "useSaTDDdu",
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "oai-gnb-du-configmap",
											},
										},
									},
								},
								corev1.EnvVar{
									Name: "GNB_NAME",
									ValueFrom: &corev1.EnvVarSource{
										ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "oai-gnb-du-configmap",
											},
											Key: "gnbduName",
										},
									},
								},
								corev1.EnvVar{
									Name: "MCC",
									ValueFrom: &corev1.EnvVarSource{
										ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
											Key: "mcc",
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "oai-gnb-du-configmap",
											},
										},
									},
								},
								corev1.EnvVar{
									Name: "MNC",
									ValueFrom: &corev1.EnvVarSource{
										ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
											Key: "mnc",
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "oai-gnb-du-configmap",
											},
										},
									},
								},
								corev1.EnvVar{
									Name: "MNC_LENGTH",
									ValueFrom: &corev1.EnvVarSource{
										ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "oai-gnb-du-configmap",
											},
											Key: "mncLength",
										},
									},
								},
								corev1.EnvVar{
									Name: "TAC",
									ValueFrom: &corev1.EnvVarSource{
										ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "oai-gnb-du-configmap",
											},
											Key: "tac",
										},
									},
								},
								corev1.EnvVar{
									ValueFrom: &corev1.EnvVarSource{
										ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
											Key: "nssaiSst",
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "oai-gnb-du-configmap",
											},
										},
									},
									Name: "NSSAI_SST",
								},
								corev1.EnvVar{
									Name: "NSSAI_SD0",
									ValueFrom: &corev1.EnvVarSource{
										ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
											Key: "nssaiSd0",
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "oai-gnb-du-configmap",
											},
										},
									},
								},
								corev1.EnvVar{
									Name: "AMF_IP_ADDRESS",
									ValueFrom: &corev1.EnvVarSource{
										ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
											Key: "amfIpAddress",
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "oai-gnb-du-configmap",
											},
										},
									},
								},
								corev1.EnvVar{
									Name: "GNB_NGA_IF_NAME",
									ValueFrom: &corev1.EnvVarSource{
										ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
											Key: "gnbNgaIfName",
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "oai-gnb-du-configmap",
											},
										},
									},
								},
								corev1.EnvVar{
									Name: "GNB_NGA_IP_ADDRESS",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
								corev1.EnvVar{
									Name: "GNB_NGU_IF_NAME",
									ValueFrom: &corev1.EnvVarSource{
										ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "oai-gnb-du-configmap",
											},
											Key: "gnbNguIfName",
										},
									},
								},
								corev1.EnvVar{
									Name: "GNB_NGU_IP_ADDRESS",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
								corev1.EnvVar{
									Name: "F1_IF_NAME",
									ValueFrom: &corev1.EnvVarSource{
										ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
											Key: "f1IfName",
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "oai-gnb-du-configmap",
											},
										},
									},
								},
								corev1.EnvVar{
									Name: "F1_DU_IP_ADDRESS",
									ValueFrom: &corev1.EnvVarSource{
										ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
											Key: "f1duIpAddress",
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "oai-gnb-du-configmap",
											},
										},
									},
								},
								corev1.EnvVar{
									Name: "F1_CU_IP_ADDRESS",
									ValueFrom: &corev1.EnvVarSource{
										ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
											Key: "f1cuIpAddress",
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "oai-gnb-du-configmap",
											},
										},
									},
								},
								corev1.EnvVar{
									Name: "F1_CU_D_PORT",
									ValueFrom: &corev1.EnvVarSource{
										ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
											Key: "f1cuPort",
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "oai-gnb-du-configmap",
											},
										},
									},
								},
								corev1.EnvVar{
									Name: "F1_DU_D_PORT",
									ValueFrom: &corev1.EnvVarSource{
										ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
											Key: "f1duPort",
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "oai-gnb-du-configmap",
											},
										},
									},
								},
							},
							Image: "docker.io/oaisoftwarealliance/oai-gnb:2023.w19",
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
									SubPath:   "mounted.conf",
									MountPath: "/opt/oai-gnb/etc/mounted.conf",
								},
							},
							Name: "gnbdu",
							SecurityContext: &corev1.SecurityContext{
								Privileged: boolPtr(true),
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
