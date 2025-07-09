# OpenAirInterface (OAI) RAN Operators

This repository contains source code for k8s custom operator for [OpenAirInterface](https://gitlab.eurecom.fr/oai/openairinterface5g/-/tree/develop?ref_type=heads) RAN functions (CU-CP, CU-UP, DU) which can be deployed in Nephio. <br />
The operator is currently common for all the RAN functions (CU-CP, CU-UP, DU). <br />

The operator listens to the [NFDeployment CRD](https://github.com/nephio-project/api/blob/main/workload/v1alpha1/nf_deployment_types.go). The operator decides which Network function the CR is intended for based on the Provider field in the `NFDeplymentSpec`. The operator retrieves the custom configuration required for the Network function from the [NFConfig CR](https://github.com/nephio-project/api/blob/main/workload/v1alpha1/nf_config_types.go). <br />

Here is an [example](https://github.com/nephio-project/catalog/blob/main/workloads/oai/pkg-example-cucp-bp/cucpdeployment.yaml) for the NfDeployment CR for CU-CP. The CR also contains the reference to custom configuration NFConfig CR as part of the parametersRefs.<br />

The KPT package used for deploying the RAN operator is located in the [nephio/catalog repository](https://github.com/nephio-project/catalog/tree/main/workloads/oai/oai-ran-operator). <br />

The KPT packages for deploying the RAN Network functions (CU-CP, CU-UP, DU) are located in the [catalog repository](https://github.com/nephio-project/catalog/tree/main/workloads/oai). <br />

**Note**: 
1. The CR is only the initial CR and it will get specialized by Nephio and more fields will added before the CR is applied in the cluster. <br />
2. Dynamic updates of the NFdeployment CR is currently not supported by the operator.

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

The [document](https://github.com/nephio-project/catalog/blob/main/workloads/oai/README.md) contains details on how to deploy the controllers and the RAN Network functions using Nephio. 

To know more about the OpenAirInterface project [checkout their website](https://openairinterface.org/).  


