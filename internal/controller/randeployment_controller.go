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
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/strings/slices"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	configref "github.com/nephio-project/api/references/v1alpha1"
	workloadv1alpha1 "github.com/nephio-project/api/workload/v1alpha1"
)

func GetSupportedProviders() []string {
	return []string{"cucp.openairinterface.org", "cuup.openairinterface.org", "du.openairinterface.org"}
}

func GetMandatoryNfKinds() []string {
	return []string{"PLMN", "RANConfig", "OAIConfig"}
}

type ConfigInfo struct {
	ConfigRefInfo  map[string][]*configref.Config
	ConfigSelfInfo map[string]runtime.RawExtension
}

func NewConfigInfo() *ConfigInfo {
	return &ConfigInfo{ConfigRefInfo: make(map[string][]*configref.Config), ConfigSelfInfo: make(map[string]runtime.RawExtension)}
}

// RANDeploymentReconciler reconciles a RANDeployment object
type RANDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Interface definition for NfResource
type NfResource interface {
	GetServiceAccount() []*corev1.ServiceAccount
	GetConfigMap(logr.Logger, *workloadv1alpha1.NFDeployment, *ConfigInfo) []*corev1.ConfigMap
	createNetworkAttachmentDefinitionNetworks(string, *workloadv1alpha1.NFDeploymentSpec) (string, error)
	GetDeployment(logr.Logger, *workloadv1alpha1.NFDeployment, *ConfigInfo) []*appsv1.Deployment
	GetService() []*corev1.Service
}

func (r *RANDeploymentReconciler) CreateAll(ctx context.Context, ranDeployment *workloadv1alpha1.NFDeployment, nfResource NfResource, configInfo *ConfigInfo) {
	namespacedName := types.NamespacedName{Namespace: ranDeployment.Namespace, Name: ranDeployment.Name}
	logger := log.FromContext(ctx).WithValues("RANDeployment", namespacedName)
	var err error
	namespaceProvided := ranDeployment.Namespace

	for _, resource := range nfResource.GetServiceAccount() {
		if resource.ObjectMeta.Namespace == "" {
			resource.ObjectMeta.Namespace = namespaceProvided
		}
		err = r.Create(ctx, resource)
		if err != nil {
			logger.Error(err, "Error During Creating resource of GetServiceAccount()")
		}
	}

	for _, resource := range nfResource.GetConfigMap(logger, ranDeployment, configInfo) {
		if resource.ObjectMeta.Namespace == "" {
			resource.ObjectMeta.Namespace = namespaceProvided
		}
		err = r.Create(ctx, resource)
		if err != nil {
			logger.Error(err, "Error During Creating resource of GetConfigMap()")
		}
	}

	for _, resource := range nfResource.GetDeployment(logger, ranDeployment, configInfo) {
		if resource.ObjectMeta.Namespace == "" {
			resource.ObjectMeta.Namespace = namespaceProvided
		}
		err = r.Create(ctx, resource)
		if err != nil {
			logger.Error(err, "Error During Creating resource of GetDeployment()")
		}
	}
	for _, resource := range nfResource.GetService() {
		if resource.ObjectMeta.Namespace == "" {
			resource.ObjectMeta.Namespace = namespaceProvided
		}
		err = r.Create(ctx, resource)
		if err != nil {
			logger.Error(err, "Error During Creating resource of GetService()")
		}
	}

}

func (r *RANDeploymentReconciler) DeleteAll(ctx context.Context, ranDeployment *workloadv1alpha1.NFDeployment, nfResource NfResource, configInfo *ConfigInfo) {
	namespacedName := types.NamespacedName{Namespace: ranDeployment.Namespace, Name: ranDeployment.Name}
	logger := log.FromContext(ctx).WithValues("RANDeployment", namespacedName)
	var err error
	namespaceProvided := ranDeployment.Namespace

	for _, resource := range nfResource.GetServiceAccount() {
		if resource.ObjectMeta.Namespace == "" {
			resource.ObjectMeta.Namespace = namespaceProvided
		}
		err = r.Delete(ctx, resource)
		if err != nil {
			logger.Error(err, "Error During Deleting resource of GetServiceAccount()")
		}
	}

	for _, resource := range nfResource.GetConfigMap(logger, ranDeployment, configInfo) {
		if resource.ObjectMeta.Namespace == "" {
			resource.ObjectMeta.Namespace = namespaceProvided
		}
		err = r.Delete(ctx, resource)
		if err != nil {
			logger.Error(err, "Error During Deleting resource of GetConfigMap()")
		}
	}

	for _, resource := range nfResource.GetDeployment(logger, ranDeployment, configInfo) {
		if resource.ObjectMeta.Namespace == "" {
			resource.ObjectMeta.Namespace = namespaceProvided
		}
		err = r.Delete(ctx, resource)
		if err != nil {
			logger.Error(err, "Error During Deleting resource of GetDeployment()")
		}

	}

	for _, resource := range nfResource.GetService() {
		if resource.ObjectMeta.Namespace == "" {
			resource.ObjectMeta.Namespace = namespaceProvided
		}
		err = r.Delete(ctx, resource)
		if err != nil {
			logger.Error(err, "Error During Deleting resource of GetService()")
		}

	}

}

func (r *RANDeploymentReconciler) GetConfigs(ctx context.Context, ranDeployment *workloadv1alpha1.NFDeployment) (*ConfigInfo, error) {
	namespacedName := types.NamespacedName{Namespace: ranDeployment.Namespace, Name: ranDeployment.Name}
	logger := log.FromContext(ctx).WithValues("RANDeployment", namespacedName)

	configsList := ranDeployment.Spec.ParametersRefs
	configInfo := NewConfigInfo()
	for _, configItem := range configsList {

		logger.Info("Config: ", "config.Name", configItem.Name)
		if configItem.APIVersion == "ref.nephio.org/v1alpha1" {
			configInstance := &configref.Config{}
			if err := r.Get(ctx, types.NamespacedName{Name: *configItem.Name, Namespace: ranDeployment.Namespace}, configInstance); err != nil {
				logger.Error(err, "Config ref get error")
				return configInfo, err
			}
			logger.Info("Config ref:", "configInstance.Name", configInstance.Name)
			var result map[string]any
			if err := json.Unmarshal(configInstance.Spec.Config.Raw, &result); err != nil {
				logger.Error(err, "Unmarshal error")
				return configInfo, err
			}
			logger.Info("Config ref:", "configInstance.Kind", result["kind"].(string))
			kindInfo := result["kind"].(string)
			configInfo.ConfigRefInfo[kindInfo] = append(configInfo.ConfigRefInfo[kindInfo], configInstance)
		} else if configItem.APIVersion == "workload.nephio.org/v1alpha1" {
			configInstance := &workloadv1alpha1.NFConfig{}
			if err := r.Get(ctx, types.NamespacedName{Name: *configItem.Name, Namespace: ranDeployment.Namespace}, configInstance); err != nil {
				logger.Error(err, "Config for Self get error")
				return configInfo, err
			}
			logger.Info("Config for Self:", "configInstance.Name", configInstance.Name)
			for _, configNf := range configInstance.Spec.ConfigRefs {
				var result map[string]any
				if err := json.Unmarshal(configNf.Raw, &result); err != nil {
					logger.Error(err, "Unmarshal error")
					return configInfo, err
				}
				logger.Info("Config for Self:", "configInstance.Kind", result["kind"].(string))
				kindInfo := result["kind"].(string)
				configInfo.ConfigSelfInfo[kindInfo] = configNf
			}

			if !CheckMandatoryKinds(configInfo.ConfigSelfInfo) {
				err := fmt.Errorf("Not all mandatory Kinds available")
				logger.Error(err, "Config for Self get error")
				return configInfo, err
			}
		} else {
			err := fmt.Errorf("Not supported API version %q", configItem.APIVersion)
			logger.Error(err, "Config for Self get error")
			return configInfo, err
		}
	}
	return configInfo, nil
}

//+kubebuilder:rbac:groups=workload.nephio.org,resources=randeployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=workload.nephio.org,resources=randeployments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=workload.nephio.org,resources=randeployments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the RANDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *RANDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	logger := log.FromContext(ctx).WithValues("RANDeployment", req.NamespacedName)
	logger.Info("Overridden Values Done:: Reconcile for RANDeployment")
	instance := &workloadv1alpha1.NFDeployment{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("RANDeployment resource not found, ignoring because object must be deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get RANDeployment")
		return ctrl.Result{}, err
	}

	if !slices.Contains(GetSupportedProviders(), instance.Spec.Provider) {
		logger.Info("Reconcile called for not supported provider", "Provider", instance.Spec.Provider)
		return ctrl.Result{}, nil
	}
	logger.Info("RANDeployment", "RANDeployment CR", instance.Spec)

	configInfo, err := r.GetConfigs(ctx, instance)
	if err != nil || configInfo == nil {
		logger.Error(err, "Failed to get required ConfigInfo")
		return ctrl.Result{}, err
	}
	// name of our custom finalizer
	myFinalizerName := "batch.tutorial.kubebuilder.io/finalizer"
	// examine DeletionTimestamp to determine if object is under deletion
	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		// Adding a Finaliser also adds the DeletionTimestamp while deleting
		if !controllerutil.ContainsFinalizer(instance, myFinalizerName) {
			// Assumed to be called only during CR-Creation

			switch resourceType := instance.Spec.Provider; resourceType {
			case "cucp.openairinterface.org":
				logger.Info("--- Creation for CUCP")
				cucpResource := CuCpResources{}
				r.CreateAll(ctx, instance, cucpResource, configInfo)
				logger.Info("--- CUCP Created")
			case "cuup.openairinterface.org":
				logger.Info("--- Creation for CUUP")
				cuupResource := CuUpResources{}
				r.CreateAll(ctx, instance, cuupResource, configInfo)
				logger.Info("--- CUUP Created")
			case "du.openairinterface.org":
				logger.Info("--- Creation for DU")
				duResource := DuResources{}
				r.CreateAll(ctx, instance, duResource, configInfo)
				logger.Info("--- DU Created")

			}
			controllerutil.AddFinalizer(instance, myFinalizerName)
			if err := r.Update(ctx, instance); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is assumed to be deleted
		if controllerutil.ContainsFinalizer(instance, myFinalizerName) {

			switch resourceType := instance.Spec.Provider; resourceType {
			case "cucp.openairinterface.org":
				logger.Info("--- Deletion for CUCP")
				cucpResource := CuCpResources{}
				r.DeleteAll(ctx, instance, cucpResource, configInfo)
				logger.Info("--- CUCP Deleted")
			case "cuup.openairinterface.org":
				logger.Info("--- Deletion for CUUP")
				cuupResource := CuUpResources{}
				r.DeleteAll(ctx, instance, cuupResource, configInfo)
				logger.Info("--- CUUP Deleted")
			case "du.openairinterface.org":
				logger.Info("--- Deletion for DU")
				duResource := DuResources{}
				r.DeleteAll(ctx, instance, duResource, configInfo)
				logger.Info("--- DU Deleted")

			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(instance, myFinalizerName)
			if err := r.Update(ctx, instance); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RANDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&workloadv1alpha1.NFDeployment{}).
		Complete(r)
}
