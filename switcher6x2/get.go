package switcher6x2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/byuoitav/atlona-switcher-microservice/structs"
	"github.com/byuoitav/common/nerr"
)

// GetInput returns the current input for the requested output
func GetInput(address, output string) (string, *nerr.E) {
	var resp structs.AtlonaVideo
	url := fmt.Sprintf("http://%s/cgi-bin/config.cgi", address)

	requestBody := fmt.Sprintf(`
	{
		"getConfig": {
			"video": {
				"vidOut": {
					"hdmiOut": {
					}
				}
			}
		}
	}`)
	payload := strings.NewReader(requestBody)
	req, _ := http.NewRequest("POST", url, payload)
	req = AddHeaders(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", nerr.Translate(err).Addf("error when making call: %s", err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		fmt.Printf("%s/n", body)
		return "", nerr.Translate(err).Addf("error when unmarshalling the response: %s", err)
	}
	//Get the inputsrc for the requested output
	input := ""
	if output == "1" {
		input = strconv.Itoa(resp.Video.VidOut.HdmiOut.HdmiOutA.VideoSrc)
	} else if output == "2" {
		input = strconv.Itoa(resp.Video.VidOut.HdmiOut.HdmiOutB.VideoSrc)
	} else {
		return input, nerr.Create("Invalid Output. Valid Output names are 1 and 2 you gave us "+output, "")
	}
	return input, nil
}

// GetHardwareInfo returns a hardware info struct         Change to structs.HardwareInfo
func GetHardwareInfo(address string) (structs.AtlonaNetwork, *nerr.E) {
	var resp structs.AtlonaNetwork
	url := fmt.Sprintf("http://%s/cgi-bin/config.cgi", address)

	payload := strings.NewReader(`
	{
		"getConfig": {
			"network": {
				"eth0":{
				}
			}
		}
	}`)
	req, _ := http.NewRequest("POST", url, payload)

	req = AddHeaders(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return resp, nerr.Translate(err).Addf("error when making call: %s", err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	err = json.Unmarshal([]byte(body), &resp) // here!
	if err != nil {
		return resp, nerr.Translate(err).Addf("error when unmarshalling the response: %s", err)
	}
	fmt.Print(resp)
	return resp, nil
}

// GetMute .
func GetMute(address, output string) (bool, *nerr.E) {
	var resp structs.AtlonaAudio
	if output == "1" || output == "2" {
		url := fmt.Sprintf("http://%s/cgi-bin/config.cgi", address)
		requestBody := fmt.Sprintf(`
		{
			"getConfig": {
				"audio":{
					"audOut":{
						"zoneOut%s":{
							"analogOut": {				
							}
						}
					}
				}	
			}	
		}`, output)
		payload := strings.NewReader(requestBody)

		req, _ := http.NewRequest("POST", url, payload)

		req = AddHeaders(req)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return false, nerr.Translate(err).Addf("error when making call: %s", err)
		}
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		err = json.Unmarshal([]byte(body), &resp)
		if err != nil {
			return false, nerr.Translate(err).Addf("error when unmarshalling the response: %s", err)
		}
	} else {
		return false, nerr.Create("Invalid Output. Valid Output names are 1 and 2 you gave us "+output, "")
	}
	if output == "1" {
		return resp.Audio.AudOut.ZoneOut1.AnalogOut.AudioMute, nil
	} else if output == "2" {
		return resp.Audio.AudOut.ZoneOut2.AnalogOut.AudioMute, nil
	} else {
		return false, nerr.Create("Invalid Output. Valid Output names are 1 and 2 you gave us "+output, "")
	}

}

// GetVolume TODO: convert the value returned back to a 0-100 value
func GetVolume(address, output string) (int, *nerr.E) {
	var resp structs.AtlonaAudio
	url := fmt.Sprintf("http://%s/cgi-bin/config.cgi", address)
	requestBody := fmt.Sprintf(`
	{
		"getConfig": {
			"audio": {
				"audOut": {
					}
				}
			}
	}`)
	payload := strings.NewReader(requestBody)

	req, _ := http.NewRequest("POST", url, payload)

	req = AddHeaders(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nerr.Translate(err).Addf("error when making call: %s", err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	err = json.Unmarshal([]byte(body), &resp) // here!
	if err != nil {
		return 0, nerr.Translate(err).Addf("error when unmarshalling the response: %s", err)
	}
	if output == "1" {
		if resp.Audio.AudOut.ZoneOut1.AudioVol < -40 {
			return 0, nil
		} else {
			volume := ((resp.Audio.AudOut.ZoneOut1.AudioVol + 40) * 2)
			return volume, nil
		}

	} else if output == "2" {
		return resp.Audio.AudOut.ZoneOut2.AudioVol + 90, nil
	} else {
		return 0, nerr.Create("Invalid Output. Valid Output names are 1 and 2 you gave us "+output, "")
	}
}
