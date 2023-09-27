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

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	configref "github.com/nephio-project/api/references/v1alpha1"
	workloadnephioorgv1alpha1 "workload.nephio.org/ran_deployment/api/v1alpha1"
)

// RANDeploymentReconciler reconciles a RANDeployment object
type RANDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Interface definition for NfResource
type NfResource interface {
	GetServiceAccount() []*corev1.ServiceAccount
	GetConfigMap(logr.Logger, *workloadnephioorgv1alpha1.RANDeployment, map[string]*configref.Config) []*corev1.ConfigMap
	createNetworkAttachmentDefinitionNetworks(string, *workloadnephioorgv1alpha1.RANDeploymentSpec) (string, error)
	GetDeployment(*workloadnephioorgv1alpha1.RANDeployment) []*appsv1.Deployment
}

func (r *RANDeploymentReconciler) CreateAll(log logr.Logger, ctx context.Context, ranDeployment *workloadnephioorgv1alpha1.RANDeployment, nfResource NfResource, configInstancesMap map[string]*configref.Config) {
	var err error
	namespaceProvided := ranDeployment.Namespace

	for _, resource := range nfResource.GetServiceAccount() {
		if resource.ObjectMeta.Namespace == "" {
			resource.ObjectMeta.Namespace = namespaceProvided
		}
		err = r.Create(ctx, resource)
		if err != nil {
			log.Error(err, "Error During Creating resource of GetServiceAccount()")
		}
	}

	for _, resource := range nfResource.GetConfigMap(log, ranDeployment, configInstancesMap) {
		if resource.ObjectMeta.Namespace == "" {
			resource.ObjectMeta.Namespace = namespaceProvided
		}
		err = r.Create(ctx, resource)
		if err != nil {
			log.Error(err, "Error During Creating resource of GetConfigMap()")
		}
	}

	for _, resource := range nfResource.GetDeployment(ranDeployment) {
		if resource.ObjectMeta.Namespace == "" {
			resource.ObjectMeta.Namespace = namespaceProvided
		}
		err = r.Create(ctx, resource)
		if err != nil {
			log.Error(err, "Error During Creating resource of GetDeployment()")
		}
	}

}

func (r *RANDeploymentReconciler) DeleteAll(log logr.Logger, ctx context.Context, ranDeployment *workloadnephioorgv1alpha1.RANDeployment, nfResource NfResource, configInstancesMap map[string]*configref.Config) {
	var err error
	namespaceProvided := ranDeployment.Namespace

	for _, resource := range nfResource.GetServiceAccount() {
		if resource.ObjectMeta.Namespace == "" {
			resource.ObjectMeta.Namespace = namespaceProvided
		}
		err = r.Delete(ctx, resource)
		if err != nil {
			log.Error(err, "Error During Deleting resource of GetServiceAccount()")
		}
	}

	for _, resource := range nfResource.GetConfigMap(log, ranDeployment, configInstancesMap) {
		if resource.ObjectMeta.Namespace == "" {
			resource.ObjectMeta.Namespace = namespaceProvided
		}
		err = r.Delete(ctx, resource)
		if err != nil {
			log.Error(err, "Error During Deleting resource of GetConfigMap()")
		}
	}

	for _, resource := range nfResource.GetDeployment(ranDeployment) {
		if resource.ObjectMeta.Namespace == "" {
			resource.ObjectMeta.Namespace = namespaceProvided
		}
		err = r.Delete(ctx, resource)
		if err != nil {
			log.Error(err, "Error During Deleting resource of GetDeployment()")
		}

	}

}

func (r *RANDeploymentReconciler) GetConfigRefs(log logr.Logger, ctx context.Context, configRefList []corev1.ObjectReference) map[string]*configref.Config {

	configInstances := []*configref.Config{}
	configInstancesMap := make(map[string]*configref.Config)
	for _, configRef := range configRefList {
		log.Info("ConfigRefs: ", "configRef.Name", configRef.Name)
		configInstance := &configref.Config{}
		if err := r.Get(ctx, types.NamespacedName{Name: configRef.Name, Namespace: configRef.Namespace}, configInstance); err != nil {
			log.Error(err, "Config ref get error")
		}
		log.Info("Config ref:", "configInstance.Name", configInstance.Name)
		configInstances = append(configInstances, configInstance)
		var result map[string]any
		json.Unmarshal(configInstance.Spec.Config.Raw, &result)
		log.Info("Config ref:", "configInstance.Kind", result["kind"].(string))
		configInstancesMap[result["kind"].(string)] = configInstance
	}
	return configInstancesMap
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
	logger.Info("Reconcile for RANDeployment")
	instance := &workloadnephioorgv1alpha1.RANDeployment{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("RANDeployment resource not found, ignoring because object must be deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get RANDeployment")
		return ctrl.Result{}, err
	}

	logger.Info("RANDeployment", "RANDeployment CR", instance.Spec)

	configInstancesMap := r.GetConfigRefs(logger, ctx, instance.Spec.ConfigRefs)

	// name of our custom finalizer
	myFinalizerName := "batch.tutorial.kubebuilder.io/finalizer"
	// examine DeletionTimestamp to determine if object is under deletion
	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		// Adding a Finaliser also adds the DeletionTimestamp while deleting
		if !controllerutil.ContainsFinalizer(instance, myFinalizerName) {
			// Assumed to be called only during CR-Creation

			switch resourceType := instance.Spec.RanNfType; resourceType {
			case "CU-CP":
				logger.Info("--- Creation for CUCP")
				cucpResource := CuCpResources{}
				r.CreateAll(logger, ctx, instance, cucpResource, configInstancesMap)
				logger.Info("--- CUCP Created")
			case "CU-UP":
				logger.Info("--- Creation for CUUP")
				cuupResource := CuUpResources{}
				r.CreateAll(logger, ctx, instance, cuupResource, configInstancesMap)
				logger.Info("--- CUUP Created")
			case "DU":
				logger.Info("--- Creation for DU")
				duResource := DuResources{}
				r.CreateAll(logger, ctx, instance, duResource, configInstancesMap)
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

			switch resourceType := instance.Spec.RanNfType; resourceType {
			case "CU-CP":
				logger.Info("--- Deletion for CUCP")
				cucpResource := CuCpResources{}
				r.DeleteAll(logger, ctx, instance, cucpResource, configInstancesMap)
				logger.Info("--- CUCP Deleted")
			case "CU-UP":
				logger.Info("--- Deletion for CUUP")
				cuupResource := CuUpResources{}
				r.DeleteAll(logger, ctx, instance, cuupResource, configInstancesMap)
				logger.Info("--- CUUP Deleted")
			case "DU":
				logger.Info("--- Deletion for DU")
				duResource := DuResources{}
				r.DeleteAll(logger, ctx, instance, duResource, configInstancesMap)
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
		For(&workloadnephioorgv1alpha1.RANDeployment{}).
		Complete(r)
}
