package controller

import (
	"encoding/json"
	"fmt"
)

/*
If the cmd is successful, the cmdOutput must end with 'OK\n'
*/
func ValidateCommandRunStatus(cmdOutput string) bool {
	n := len(cmdOutput)
	if n < 3 {
		return false
	}

	return cmdOutput[n-3:] == "OK\n"
}

func GetO1Stats(o1Obj *O1TelnetClient) (any, error) {
	cmd := "o1 stats"
	cmdOutput, err := o1Obj.RunCommand(cmd)
	if err != nil {
		return nil, fmt.Errorf("Error in running cmd " + cmd + " :: Err --> " + err.Error())
	}
	if !ValidateCommandRunStatus(cmdOutput) {
		return nil, fmt.Errorf("OK Not found in cmd output :: cmd --> " + cmd + " :: CmdOutput --> " + cmdOutput)
	}

	var data any
	err = json.Unmarshal([]byte(cmdOutput[:len(cmdOutput)-3]), &data)
	if err != nil {
		fmt.Println("Error in Decoding JSON :: " + err.Error())
		return nil, fmt.Errorf("Error in Decoding JSON :: " + err.Error())
	}

	return data, nil
}

func GetCurrentBW(o1Obj *O1TelnetClient) (string, error) {
	data, err := GetO1Stats(o1Obj)
	if err != nil {
		return "", err
	}
	// ToDO: Use a whole struct to decode full o1-stats (if useful)
	ptr, _ := data.(map[string]any)
	for _, fields := range []string{"o1-config", "NRCELLDU"} {
		next := ptr[fields]
		ptr = next.(map[string]any)
	}

	bw := ptr["nrcelldu3gpp:bSChannelBwDL"].(float64)
	return fmt.Sprint(bw), nil
}

/*
The Current Procedure will reconfigure the Bandwidth
ToDo: What if modem is already-Stopped, In that case first cmd will not return 'OK'
*/
func BandWidthReconfigureProcedure(o1Obj *O1TelnetClient, toValue string) error {
	cmds := []string{
		"o1 stop_modem",
		"o1 bwconfig " + toValue,
		"o1 start_modem",
	}

	for _, cmd := range cmds {
		cmdOutput, err := o1Obj.RunCommand(cmd)
		if err != nil {
			return fmt.Errorf("Error in running cmd " + cmd + " :: Err --> " + err.Error())
		}
		if !ValidateCommandRunStatus(cmdOutput) {
			return fmt.Errorf("OK Not found in cmd output :: cmd --> " + cmd + " :: CmdOutput --> " + cmdOutput)
		}
	}
	return nil
}

func ReconfigureBandWidth(o1Obj *O1TelnetClient, toValue string) error {
	if toValue != "20" && toValue != "40" {
		return fmt.Errorf("Allowed values of Bandwidth are either 20 or 40 :: Recieved-Value : " + toValue)
	}

	curBW, err := GetCurrentBW(o1Obj)
	if err != nil {
		return fmt.Errorf("Error Recieved in Quering for BW before reconfiguration| Err --> " + err.Error())
	}
	if curBW == toValue {
		fmt.Println("Current BW is the required BW | Returning")
		return nil
	}

	fmt.Println("Reconfiguring BW from " + curBW + " to " + toValue)
	err = BandWidthReconfigureProcedure(o1Obj, toValue)
	if err != nil {
		return fmt.Errorf("Error Recieved in Reconfiguring the BW| Err --> " + err.Error())
	}
	// BW is reconfigured (Confirm it)
	newBW, err := GetCurrentBW(o1Obj)
	if err != nil {
		return fmt.Errorf("Error Recieved in Quering for BW after reconfiguration| Err --> " + err.Error())
	}
	if newBW != toValue {
		return fmt.Errorf("Unkown Error:: BW should have been reconfigured but it is not :(")
	}
	return nil
}
