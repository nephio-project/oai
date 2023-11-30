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

const configurationTemplateSourceForCuCp = `
Active_gNBs = ( "oai-cu-cp");
# Asn1_verbosity, choice in: none, info, annoying
Asn1_verbosity = "none";
Num_Threads_PUSCH = 8;
sa = 1;

gNBs =
(
 {
    ////////// Identification parameters:
    gNB_CU_ID = 0xe00;

#     cell_type =  "CELL_MACRO_GNB";

    gNB_name  =  "oai-cu-cp";

    // Tracking area code, 0x0000 and 0xfffe are reserved values
    tracking_area_code  =  {{ .TAC }};
    plmn_list = ({ mcc = {{ .PLMN_MCC }};
                   mnc = {{ .PLMN_MNC }};
                   mnc_length ={{ .PLMN_MNC_LENGTH }};
                   snssaiList = ({ sst = {{ .NSSAI_SST }}, sd = 0x{{ .NSSAI_SD }} })
                });


    nr_cellid = {{ .CELL_ID }};
    force_256qam_off = 1;

    tr_s_preference = "f1";

    local_s_if_name = "f1c";
    local_s_address = {{ .F1C_IP }};
    remote_s_address = "0.0.0.0";
    local_s_portc   = 501;
    local_s_portd   = 2152;
    remote_s_portc  = 500;
    remote_s_portd  = 2152;

    ssb_SubcarrierOffset                                      = 0;
    min_rxtxtime                                              = 6;

    servingCellConfigCommon = (
    {
 #spCellConfigCommon

      physCellId                                                    = {{ .PHY_CELL_ID }};

#  downlinkConfigCommon
    #frequencyInfoDL
      # this is 3600 MHz + 43 PRBs@30kHz SCS (same as initial BWP)
      absoluteFrequencySSB                                             = 641280;
      dl_frequencyBand                                                 = {{ .DL_FREQ_BAND }};
      # this is 3600 MHz
      dl_absoluteFrequencyPointA                                       = 640008;
      #scs-SpecificCarrierList
        dl_offstToCarrier                                              = 0;
# subcarrierSpacing
# 0=kHz15, 1=kHz30, 2=kHz60, 3=kHz120
        dl_subcarrierSpacing                                           = {{ .DL_SCS }};
        dl_carrierBandwidth                                            = {{ .DL_CARRIER_BW }};
     #initialDownlinkBWP
      #genericParameters
        # this is RBstart=27,L=48 (275*(L-1))+RBstart
        initialDLBWPlocationAndBandwidth                               = 28875; # 6366 12925 12956 28875 12952
# subcarrierSpacing
# 0=kHz15, 1=kHz30, 2=kHz60, 3=kHz120
        initialDLBWPsubcarrierSpacing                                   = 1;
      #pdcch-ConfigCommon
        initialDLBWPcontrolResourceSetZero                              = 11;
        initialDLBWPsearchSpaceZero                                     = 0;

  #uplinkConfigCommon
     #frequencyInfoUL
      ul_frequencyBand                                              = {{ .UL_FREQ_BAND }};
      #scs-SpecificCarrierList
      ul_offstToCarrier                                             = 0;
# subcarrierSpacing
# 0=kHz15, 1=kHz30, 2=kHz60, 3=kHz120
      ul_subcarrierSpacing                                          = {{ .UL_SCS }};
      ul_carrierBandwidth                                           = {{ .UL_CARRIER_BW }};
      pMax                                                          = 20;
     #initialUplinkBWP
      #genericParameters
        initialULBWPlocationAndBandwidth                            = 28875;
# subcarrierSpacing
# 0=kHz15, 1=kHz30, 2=kHz60, 3=kHz120
        initialULBWPsubcarrierSpacing                               = 1;
      #rach-ConfigCommon
        #rach-ConfigGeneric
          prach_ConfigurationIndex                                  = 98;
#prach_msg1_FDM
#0 = one, 1=two, 2=four, 3=eight
          prach_msg1_FDM                                            = 0;
          prach_msg1_FrequencyStart                                 = 0;
          zeroCorrelationZoneConfig                                 = 13;
          preambleReceivedTargetPower                               = -96;
#preamblTransMax (0...10) = (3,4,5,6,7,8,10,20,50,100,200)
          preambleTransMax                                          = 6;
#powerRampingStep
# 0=dB0,1=dB2,2=dB4,3=dB6
        powerRampingStep                                            = 1;
#ra_ReponseWindow
#1,2,4,8,10,20,40,80
        ra_ResponseWindow                                           = 4;
#ssb_perRACH_OccasionAndCB_PreamblesPerSSB_PR
#1=oneeighth,2=onefourth,3=half,4=one,5=two,6=four,7=eight,8=sixteen
        ssb_perRACH_OccasionAndCB_PreamblesPerSSB_PR                = 4;
#oneHalf (0..15) 4,8,12,16,...60,64
        ssb_perRACH_OccasionAndCB_PreamblesPerSSB                   = 14;
#ra_ContentionResolutionTimer
#(0..7) 8,16,24,32,40,48,56,64
        ra_ContentionResolutionTimer                                = 7;
        rsrp_ThresholdSSB                                           = 19;
#prach-RootSequenceIndex_PR
#1 = 839, 2 = 139
        prach_RootSequenceIndex_PR                                  = 2;
        prach_RootSequenceIndex                                     = 1;
        # SCS for msg1, can only be 15 for 30 kHz < 6 GHz, takes precendence over the one derived from prach-ConfigIndex
        #
        msg1_SubcarrierSpacing                                      = 1,
# restrictedSetConfig
# 0=unrestricted, 1=restricted type A, 2=restricted type B
        restrictedSetConfig                                         = 0,

        msg3_DeltaPreamble                                          = 1;
        p0_NominalWithGrant                                         =-90;

# pucch-ConfigCommon setup :
# pucchGroupHopping
# 0 = neither, 1= group hopping, 2=sequence hopping
        pucchGroupHopping                                           = 0;
        hoppingId                                                   = 40;
        p0_nominal                                                  = -90;
# ssb_PositionsInBurs_BitmapPR
# 1=short, 2=medium, 3=long
      ssb_PositionsInBurst_PR                                       = 2;
      ssb_PositionsInBurst_Bitmap                                   = 1;

# ssb_periodicityServingCell
# 0 = ms5, 1=ms10, 2=ms20, 3=ms40, 4=ms80, 5=ms160, 6=spare2, 7=spare1
      ssb_periodicityServingCell                                    = 2;

# dmrs_TypeA_position
# 0 = pos2, 1 = pos3
      dmrs_TypeA_Position                                           = 0;

# subcarrierSpacing
# 0=kHz15, 1=kHz30, 2=kHz60, 3=kHz120
      subcarrierSpacing                                             = 1;


  #tdd-UL-DL-ConfigurationCommon
# subcarrierSpacing
# 0=kHz15, 1=kHz30, 2=kHz60, 3=kHz120
      referenceSubcarrierSpacing                                    = 1;
      # pattern1
      # dl_UL_TransmissionPeriodicity
      # 0=ms0p5, 1=ms0p625, 2=ms1, 3=ms1p25, 4=ms2, 5=ms2p5, 6=ms5, 7=ms10
      dl_UL_TransmissionPeriodicity                                 = 6;
      nrofDownlinkSlots                                             = 7;
      nrofDownlinkSymbols                                           = 6;
      nrofUplinkSlots                                               = 2;
      nrofUplinkSymbols                                             = 4;

      ssPBCH_BlockPower                                             = -25;
  }

  );
    # ------- SCTP definitions
    SCTP :
    {
        # Number of streams to use in input/output
        SCTP_INSTREAMS  = 2;
        SCTP_OUTSTREAMS = 2;
    };


    ////////// AMF parameters:
    amf_ip_address      = ( { ipv4       = {{ .AMF_IP }};
                              ipv6       = "0:0:0::0";
                              active     = "yes";
                              preference = "ipv4";
                            }
                          );

    E1_INTERFACE =
    (
      {
        type = "cp";
        ipv4_cucp = {{ .E1_IP }};
        port_cucp = 38462;
        ipv4_cuup = "0.0.0.0";
        port_cuup = 38462;
      }
    )

    NETWORK_INTERFACES :
    {
        GNB_INTERFACE_NAME_FOR_NG_AMF            = "n2";
        GNB_IPV4_ADDRESS_FOR_NG_AMF              = {{ .N2_IP }};
    };
  }
);

security = {
  # preferred ciphering algorithms
  # the first one of the list that an UE supports in chosen
  # valid values: nea0, nea1, nea2, nea3
  ciphering_algorithms = ( "nea0" );

  # preferred integrity algorithms
  # the first one of the list that an UE supports in chosen
  # valid values: nia0, nia1, nia2, nia3
  integrity_algorithms = ( "nia2", "nia0" );

  # setting 'drb_ciphering' to "no" disables ciphering for DRBs, no matter
  # what 'ciphering_algorithms' configures; same thing for 'drb_integrity'
  drb_ciphering = "yes";
  drb_integrity = "no";
};
     log_config :
     {
       global_log_level                      ="info";
       hw_log_level                          ="info";
       phy_log_level                         ="info";
       mac_log_level                         ="info";
       rlc_log_level                         ="debug";
       pdcp_log_level                        ="info";
       rrc_log_level                         ="info";
       f1ap_log_level                         ="info";
       ngap_log_level                         ="debug";
    };

`
const configurationTemplateSourceForCuUp = `
Active_gNBs = ( "oai-cu-up");
# Asn1_verbosity, choice in: none, info, annoying
Asn1_verbosity = "none";
sa = 1;
gNBs =
(
 {
    ////////// Identification parameters:
    gNB_CU_ID = 0xe00;

#     cell_type =  "CELL_MACRO_GNB";

    gNB_name  =  "oai-cu-up";

    // Tracking area code, 0x0000 and 0xfffe are reserved values
    tracking_area_code  =  {{ .TAC }};
    plmn_list = ({ mcc = {{ .PLMN_MCC }}; 
                   mnc = {{ .PLMN_MNC }}; 
                   mnc_length ={{ .PLMN_MNC_LENGTH }}; 
                   snssaiList = ({ sst = {{ .NSSAI_SST }}, sd = 0x{{ .NSSAI_SD }} }) 
                });

    tr_s_preference = "f1";

    local_s_if_name = "f1u";
    local_s_address = {{ .F1U_IP }};
    remote_s_address = "0.0.0.0";
    local_s_portc   = 501;
    local_s_portd   = 2152;
    remote_s_portc  = 500;
    remote_s_portd  = 2152;

    # ------- SCTP definitions
    SCTP :
    {
        # Number of streams to use in input/output
        SCTP_INSTREAMS  = 2;
        SCTP_OUTSTREAMS = 2;
    };

    E1_INTERFACE =
    (
      {
        type = "up";
        ipv4_cucp = {{ .CUCP_E1 }};
        ipv4_cuup = {{ .E1_IP }};
      }
    )

    NETWORK_INTERFACES :
    {
        GNB_INTERFACE_NAME_FOR_NG_AMF            = "n3";
        GNB_IPV4_ADDRESS_FOR_NG_AMF              = {{ .N3_IP }};
        GNB_INTERFACE_NAME_FOR_NGU               = "n3";
        GNB_IPV4_ADDRESS_FOR_NGU                 = {{ .N3_IP }};
        GNB_PORT_FOR_S1U                         = 2152; # Spec 2152
    };
  }
);

security = {
  # preferred ciphering algorithms
  # the first one of the list that an UE supports in chosen
  # valid values: nea0, nea1, nea2, nea3
  ciphering_algorithms = ( "nea0" );

  # preferred integrity algorithms
  # the first one of the list that an UE supports in chosen
  # valid values: nia0, nia1, nia2, nia3
  integrity_algorithms = ( "nia2", "nia0" );

  # setting 'drb_ciphering' to "no" disables ciphering for DRBs, no matter
  # what 'ciphering_algorithms' configures; same thing for 'drb_integrity'
  drb_ciphering = "yes";
  drb_integrity = "no";
};
     log_config :
     {
       global_log_level                      ="info";
       hw_log_level                          ="info";
       phy_log_level                         ="info";
       mac_log_level                         ="info";
       rlc_log_level                         ="debug";
       pdcp_log_level                        ="info";
       rrc_log_level                         ="info";
       f1ap_log_level                         ="info";
       ngap_log_level                         ="debug";
    };
`
const configurationTemplateSourceForDu = `
Active_gNBs = ( "oai-du-rfsim");
# Asn1_verbosity, choice in: none, info, annoying
Asn1_verbosity = "none";

gNBs =
(
 {
    ////////// Identification parameters:
    gNB_ID = 0xe00;

#     cell_type =  "CELL_MACRO_GNB";

    gNB_name  =  "oai-du-rfsim";

    // Tracking area code, 0x0000 and 0xfffe are reserved values
    tracking_area_code  =  {{ .TAC }};
    plmn_list = ({ mcc = {{ .PLMN_MCC }}; mnc = {{ .PLMN_MNC }}; mnc_length = {{ .PLMN_MNC_LENGTH }}; snssaiList = ({ sst = {{ .NSSAI_SST }}, sd = 0x{{ .NSSAI_SD }} }) });


    nr_cellid = {{ .CELL_ID }};

    ////////// Physical parameters:

    min_rxtxtime                                              = 6;
    force_256qam_off = 1;

    servingCellConfigCommon = (
    {
 #spCellConfigCommon

      physCellId                                                    = {{ .PHY_CELL_ID }};

#  downlinkConfigCommon
    #frequencyInfoDL
      # this is 3600 MHz + 43 PRBs@30kHz SCS (same as initial BWP)
      absoluteFrequencySSB                                          = 641280;
      dl_frequencyBand                                                 = {{ .DL_FREQ_BAND }};
      # this is 3600 MHz
      dl_absoluteFrequencyPointA                                       = 640008;
      #scs-SpecificCarrierList
        dl_offstToCarrier                                              = 0;
# subcarrierSpacing
# 0=kHz15, 1=kHz30, 2=kHz60, 3=kHz120  
        dl_subcarrierSpacing                                           = {{ .DL_SCS }};
        dl_carrierBandwidth                                            = {{ .DL_CARRIER_BW }};
     #initialDownlinkBWP
      #genericParameters
        # this is RBstart=27,L=48 (275*(L-1))+RBstart
        initialDLBWPlocationAndBandwidth                               = 28875; # 6366 12925 12956 28875 12952
# subcarrierSpacing
# 0=kHz15, 1=kHz30, 2=kHz60, 3=kHz120  
        initialDLBWPsubcarrierSpacing                                           = 1;
      #pdcch-ConfigCommon
        initialDLBWPcontrolResourceSetZero                              = 12;
        initialDLBWPsearchSpaceZero                                             = 0;

  #uplinkConfigCommon 
     #frequencyInfoUL
      ul_frequencyBand                                                 = {{ .UL_FREQ_BAND }};
      #scs-SpecificCarrierList
      ul_offstToCarrier                                              = 0;
# subcarrierSpacing
# 0=kHz15, 1=kHz30, 2=kHz60, 3=kHz120  
      ul_subcarrierSpacing                                           = {{ .UL_SCS }};
      ul_carrierBandwidth                                            = {{ .UL_CARRIER_BW }};
      pMax                                                          = 20;
     #initialUplinkBWP
      #genericParameters
        initialULBWPlocationAndBandwidth                            = 28875;
# subcarrierSpacing
# 0=kHz15, 1=kHz30, 2=kHz60, 3=kHz120  
        initialULBWPsubcarrierSpacing                                           = 1;
      #rach-ConfigCommon
        #rach-ConfigGeneric
          prach_ConfigurationIndex                                  = 98;
#prach_msg1_FDM
#0 = one, 1=two, 2=four, 3=eight
          prach_msg1_FDM                                            = 0;
          prach_msg1_FrequencyStart                                 = 0;
          zeroCorrelationZoneConfig                                 = 13;
          preambleReceivedTargetPower                               = -96;
#preamblTransMax (0...10) = (3,4,5,6,7,8,10,20,50,100,200)
          preambleTransMax                                          = 6;
#powerRampingStep
# 0=dB0,1=dB2,2=dB4,3=dB6
        powerRampingStep                                            = 1;
#ra_ReponseWindow
#1,2,4,8,10,20,40,80
        ra_ResponseWindow                                           = 4;
#ssb_perRACH_OccasionAndCB_PreamblesPerSSB_PR
#1=oneeighth,2=onefourth,3=half,4=one,5=two,6=four,7=eight,8=sixteen
        ssb_perRACH_OccasionAndCB_PreamblesPerSSB_PR                = 4;
#oneHalf (0..15) 4,8,12,16,...60,64
        ssb_perRACH_OccasionAndCB_PreamblesPerSSB                   = 14;
#ra_ContentionResolutionTimer
#(0..7) 8,16,24,32,40,48,56,64
        ra_ContentionResolutionTimer                                = 7;
        rsrp_ThresholdSSB                                           = 19;
#prach-RootSequenceIndex_PR
#1 = 839, 2 = 139
        prach_RootSequenceIndex_PR                                  = 2;
        prach_RootSequenceIndex                                     = 1;
        # SCS for msg1, can only be 15 for 30 kHz < 6 GHz, takes precendence over the one derived from prach-ConfigIndex
        #  
        msg1_SubcarrierSpacing                                      = 1,
# restrictedSetConfig
# 0=unrestricted, 1=restricted type A, 2=restricted type B
        restrictedSetConfig                                         = 0,

        msg3_DeltaPreamble                                          = 1;
        p0_NominalWithGrant                                         =-90;

# pucch-ConfigCommon setup :
# pucchGroupHopping
# 0 = neither, 1= group hopping, 2=sequence hopping
        pucchGroupHopping                                           = 0;
        hoppingId                                                   = 40;
        p0_nominal                                                  = -90;
# ssb_PositionsInBurs_BitmapPR
# 1=short, 2=medium, 3=long
      ssb_PositionsInBurst_PR                                       = 2;
      ssb_PositionsInBurst_Bitmap                                   = 1;

# ssb_periodicityServingCell
# 0 = ms5, 1=ms10, 2=ms20, 3=ms40, 4=ms80, 5=ms160, 6=spare2, 7=spare1 
      ssb_periodicityServingCell                                    = 2;

# dmrs_TypeA_position
# 0 = pos2, 1 = pos3
      dmrs_TypeA_Position                                           = 0;

# subcarrierSpacing
# 0=kHz15, 1=kHz30, 2=kHz60, 3=kHz120  
      subcarrierSpacing                                             = 1;


  #tdd-UL-DL-ConfigurationCommon
# subcarrierSpacing
# 0=kHz15, 1=kHz30, 2=kHz60, 3=kHz120  
      referenceSubcarrierSpacing                                    = 1;
      # pattern1 
      # dl_UL_TransmissionPeriodicity
      # 0=ms0p5, 1=ms0p625, 2=ms1, 3=ms1p25, 4=ms2, 5=ms2p5, 6=ms5, 7=ms10
      dl_UL_TransmissionPeriodicity                                 = 6;
      nrofDownlinkSlots                                             = 7;
      nrofDownlinkSymbols                                           = 6;
      nrofUplinkSlots                                               = 2;
      nrofUplinkSymbols                                             = 4;

      ssPBCH_BlockPower                                             = -25;
     }

  );


    # ------- SCTP definitions
    SCTP :
    {
        # Number of streams to use in input/output
        SCTP_INSTREAMS  = 2;
        SCTP_OUTSTREAMS = 2;
    };
  }
);

MACRLCs = (
  {
    num_cc           = 1;
    tr_s_preference  = "local_L1";
    tr_n_preference  = "f1";
    local_n_if_name = "f1";
    local_n_address = {{ .F1C_DU_IP }};
    remote_n_address = {{ .F1C_CU_IP }};
    local_n_portc   = 500;
    local_n_portd   = 2152;
    remote_n_portc  = 501;
    remote_n_portd  = 2152;
    pusch_TargetSNRx10          = 200;
    pucch_TargetSNRx10          = 200;
    ulsch_max_frame_inactivity = 1;
  }
);

L1s = (
{
  num_cc = 1;
  tr_n_preference = "local_mac";
  prach_dtx_threshold = 200;
  pucch0_dtx_threshold = 150;
  ofdm_offset_divisor = 8; #set this to UINT_MAX for offset 0
}
);

RUs = (
    {     
       local_rf       = "yes"
         nb_tx          = 1
         nb_rx          = 1
         att_tx         = 0
         att_rx         = 0;
         bands          = [78];
         max_pdschReferenceSignalPower = -27;
         max_rxgain                    = 114;
         eNB_instances  = [0];
         #beamforming 1x4 matrix:
         bf_weights = [0x00007fff, 0x0000, 0x0000, 0x0000];
         clock_src = "internal";
    }
);  

THREAD_STRUCT = (
  {
    #three config for level of parallelism "PARALLEL_SINGLE_THREAD", "PARALLEL_RU_L1_SPLIT", or "PARALLEL_RU_L1_TRX_SPLIT"
    parallel_config    = "PARALLEL_SINGLE_THREAD";
    #two option for worker "WORKER_DISABLE" or "WORKER_ENABLE"
    worker_config      = "WORKER_ENABLE";
  }
);
rfsimulator: {
serveraddr = "server";
    serverport = "4043";
    options = (); #("saviq"); or/and "chanmod"
    modelname = "AWGN";
    IQfile = "/tmp/rfsimulator.iqs"
}

     log_config :
     {
       global_log_level                      ="info";
       hw_log_level                          ="info";
       phy_log_level                         ="info";
       mac_log_level                         ="info";
       rlc_log_level                         ="info";
       pdcp_log_level                        ="info";
       rrc_log_level                         ="info";
       f1ap_log_level                         ="info";
       ngap_log_level                         ="debug";
    };
`

var (
	configurationTemplateForCuCp = template.Must(template.New("RanCuCpConfiguration").Parse(configurationTemplateSourceForCuCp))
	configurationTemplateForCuUp = template.Must(template.New("RanCuUpConfiguration").Parse(configurationTemplateSourceForCuUp))
	configurationTemplateForDu   = template.Must(template.New("RanDuConfiguration").Parse(configurationTemplateSourceForDu))
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
	if err := configurationTemplateForCuCp.Execute(&buffer, values); err == nil {
		return buffer.String(), nil
	} else {
		return "", err
	}
}

func renderConfigurationTemplateForCuUp(values configurationTemplateValuesForCuUp) (string, error) {
	var buffer bytes.Buffer
	if err := configurationTemplateForCuUp.Execute(&buffer, values); err == nil {
		return buffer.String(), nil
	} else {
		return "", err
	}
}

func renderConfigurationTemplateForDu(values configurationTemplateValuesForDu) (string, error) {
	var buffer bytes.Buffer
	if err := configurationTemplateForDu.Execute(&buffer, values); err == nil {
		return buffer.String(), nil
	} else {
		return "", err
	}
}
