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
	"bytes"
	"text/template"
)

var (
	// O1-Telnet
	configurationTemplateForCuCpO1Telnet = template.Must(template.New("RanCuCpConfigurationO1").Parse(configurationTemplateSourceForCuCpO1Telnet))
	configurationTemplateForCuUpO1Telnet = template.Must(template.New("RanCuUpConfigurationO1").Parse(configurationTemplateSourceForCuUpO1Telnet))
	configurationTemplateForDuO1Telnet   = template.Must(template.New("RanDuConfigurationO1").Parse(configurationTemplateSourceForDuO1Telnet))
)

type configurationTemplateValuesForCuCp struct {
	E1_IP           string
	F1C_IP          string
	N2_IP           string
	AMF_IP          string
	TAC             uint32
	CELL_ID         string
	PHY_CELL_ID     uint32
	PLMN_MCC        string
	PLMN_MNC        string
	PLMN_MNC_LENGTH string
	NSSAI_SST       int
	NSSAI_SD        string
	DL_FREQ_BAND    uint32
	DL_SCS          uint16
	DL_CARRIER_BW   uint32
	UL_FREQ_BAND    uint32
	UL_SCS          uint16
	UL_CARRIER_BW   uint32
}

type configurationTemplateValuesForCuUp struct {
	E1_IP           string
	F1U_IP          string
	N3_IP           string
	CUCP_E1         string
	TAC             uint32
	PLMN_MCC        string
	PLMN_MNC        string
	PLMN_MNC_LENGTH string
	NSSAI_SST       int
	NSSAI_SD        string
}

type configurationTemplateValuesForDu struct {
	F1C_DU_IP       string
	F1C_CU_IP       string
	TAC             uint32
	CELL_ID         string
	PHY_CELL_ID     uint32
	PLMN_MCC        string
	PLMN_MNC        string
	PLMN_MNC_LENGTH string
	NSSAI_SST       int
	NSSAI_SD        string
	DL_FREQ_BAND    uint32
	DL_SCS          uint16
	DL_CARRIER_BW   uint32
	UL_FREQ_BAND    uint32
	UL_SCS          uint16
	UL_CARRIER_BW   uint32
}

func renderConfigurationTemplateForCuCp(values configurationTemplateValuesForCuCp) (string, error) {
	var buffer bytes.Buffer

	if err := configurationTemplateForCuCpO1Telnet.Execute(&buffer, values); err == nil {
		return buffer.String(), nil
	} else {
		return "", err
	}
}

func renderConfigurationTemplateForCuUp(values configurationTemplateValuesForCuUp) (string, error) {
	var buffer bytes.Buffer

	if err := configurationTemplateForCuUpO1Telnet.Execute(&buffer, values); err == nil {
		return buffer.String(), nil
	} else {
		return "", err
	}
}

func renderConfigurationTemplateForDu(values configurationTemplateValuesForDu) (string, error) {
	var buffer bytes.Buffer
	if err := configurationTemplateForDuO1Telnet.Execute(&buffer, values); err == nil {
		return buffer.String(), nil
	} else {
		return "", err
	}
}
