# oai-ran-operators

This repository contains source code for k8s custom operator for OAI RAN functions (CU-CP, CU-UP, DU) which can be deployed in Nephio.
The operator is currently common for all the RAN functions (CU-CP, CU-UP, DU).

The operator listens to the NFDeployment CRD as shown here: https://github.com/nephio-project/api/blob/main/workload/v1alpha1/nf_deployment_types.go
The operator decides which Network function the CR is intented for based on the Provider field in the NFDeplymentSpec.
The operator also retrieves the custom configuration required for the Network function from the NFConfig CR.
The NFConfig CRD is available here: https://github.com/nephio-project/api/blob/main/workload/v1alpha1/nf_config_types.go

An example for the NfDeployment CR for CU-CP is available here: https://github.com/nephio-project/catalog/blob/main/workloads/oai/pkg-example-cucp-bp/cucpdeployment.yaml
The CR also contains the reference to custom configuraton NFConfig CR as part of the parametersRefs.
Note that this CR is only the initial CR and it will get specialized by Nephio and more fields will added before the CR is applied in the cluster.


The KPT package used for deploying the RAN operator is located in the nephio/catalog repository. https://github.com/nephio-project/catalog/tree/main/workloads/oai/oai-ran-operator.

The KPT packages for deploying the RAN Network functions (CU-CP, CU-UP, DU) are located here : https://github.com/nephio-project/catalog/tree/main/workloads/oai

Note that dynamic updates of the NFdeployment CR is currently not supported by the operator.