# oai-ran-operators

This repository contains source code for k8s custom operator for OAI RAN functions (CU-CP, CU-UP, DU) which can be deployed in Nephio. <br />
The operator is currently common for all the RAN functions (CU-CP, CU-UP, DU). <br />

The operator listens to the NFDeployment CRD as shown here: <br />
https://github.com/nephio-project/api/blob/main/workload/v1alpha1/nf_deployment_types.go. <br />
The operator decides which Network function the CR is intented for based on the Provider field in the NFDeplymentSpec. <br />
The operator also retrieves the custom configuration required for the Network function from the NFConfig CR. <br />
The NFConfig CRD is available here: <br />
https://github.com/nephio-project/api/blob/main/workload/v1alpha1/nf_config_types.go <br />

An example for the NfDeployment CR for CU-CP is available here: <br />
https://github.com/nephio-project/catalog/blob/main/workloads/oai/pkg-example-cucp-bp/cucpdeployment.yaml <br />
The CR also contains the reference to custom configuraton NFConfig CR as part of the parametersRefs. <br />
Note that this CR is only the initial CR and it will get specialized by Nephio and more fields will added before the CR is applied in the cluster. <br />


The KPT package used for deploying the RAN operator is located in the nephio/catalog repository. <br />
 https://github.com/nephio-project/catalog/tree/main/workloads/oai/oai-ran-operator.

The KPT packages for deploying the RAN Network functions (CU-CP, CU-UP, DU) are located here: <br />
 https://github.com/nephio-project/catalog/tree/main/workloads/oai

Note that dynamic updates of the NFdeployment CR is currently not supported by the operator.

The directory structure of this repository is as follows: <br />

```bash
.
├── Dockerfile
├── LICENSE
├── Makefile
├── OWNERS
├── README.md
├── api
│   └── v1alpha1
│       ├── oai_ran_nf_types.go
│       ├── plmn_types.go
│       └── ranconfig_types.go
├── cmd
│   └── main.go
├── go.mod
├── go.sum
└── internal
    └── controller
        ├── helper.go
        ├── helper_test.go
        ├── interface_configs.go
        ├── mock_Client_test.go
        ├── mock_NfResource_test.go
        ├── network_attachment_defination_test.go
        ├── network_attachment_definitions.go
        ├── randeployment_controller.go
        ├── randeployment_controller_test.go
        ├── resources_cucp.go
        ├── resources_cucp_test.go
        ├── resources_cuup.go
        ├── resources_cuup_test.go
        ├── resources_du.go
        ├── resources_du_test.go
        └── templates.go

```

More details on how to deploy the controllers and the RAN Network functions using Nephio is described here: https://github.com/nephio-project/catalog/blob/main/workloads/oai/README.md