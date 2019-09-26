package switcher6x2

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/byuoitav/common/nerr"
)

//AddHeaders adds all of the security headers needed for making the call
func AddHeaders(req *http.Request) *http.Request {
	req.Header.Add("Content-Type", "application/json")
	//This needs to be replaced with an environmental variable
	req.Header.Add("Authorization", "Basic YWRtaW46QXRsb25h")
	return req
}

// SetInput changes the input on the given output to input
func SetInput(address, output, input string) *nerr.E {
	in, err := strconv.Atoi(input)
	if err != nil {
		return nerr.Translate(err).Addf("error when making call: %s", err)
	}
	url := fmt.Sprintf("http://%s/cgi-bin/config.cgi", address)
	payload := strings.NewReader("")
	if output == "1" {
		payload = strings.NewReader(fmt.Sprintf(`
		{
			"setConfig":{
				"video":{
					"vidOut":{
						"hdmiOut":{
							"hdmiOutA":{
								"videoSrc":%v
							}
						}
					}
				}
			}
		}`, in))
	} else if output == "2" {
		payload = strings.NewReader(fmt.Sprintf(`
		{
			"setConfig":{
				"video":{
					"vidOut":{
						"hdmiOut":{
							"hdmiOutB":{
								"videoSrc":%v
							}
						}
					}
				}
			}
		}`, in))
	} else {
		return nerr.Create("Invalid Output. Valid Output names are 1 and 2", "")
	}
	req, _ := http.NewRequest("POST", url, payload)
	req = AddHeaders(req)
	res, gerr := http.DefaultClient.Do(req)
	if gerr != nil {
		return nerr.Translate(gerr).Addf("error when making call: %s", gerr)
	}
	defer res.Body.Close()
	return nil
}

// SetVolume changes the input on the given output to input
func SetVolume(address, output, level string) *nerr.E {
	//Atlona volume levels are from -90 to 10 and the number we recieve is 0-100
	volumeLevel, err := strconv.Atoi(level)
	if err != nil {
		return nerr.Translate(err).Add("unable to switch change volume")
	}
	//if volume level is supposed to be zero set it to zero (which is -90) on atlona
	if volumeLevel == 0 {
		volumeLevel = -90
	} else {
		convertedVolume := -40 + math.Round(float64(volumeLevel/2))
		volumeLevel = int(convertedVolume)
	}
	err = SetVolumeHelper(address, output, strconv.Itoa(volumeLevel))
	if err != nil {
		return nerr.Translate(err).Add("unable to switch change volume")
	}
	return nil
}

//SetVolumeHelper .
func SetVolumeHelper(address, output, level string) *nerr.E {
	url := fmt.Sprintf("http://%s/cgi-bin/config.cgi", address)
	if output == "1" || output == "2" {
		body := fmt.Sprintf(`
		{
			"setConfig": {
				"audio": {
					"audOut": {
						"zoneOut%s": {
							"audioVol": %s
						}
					}
				}
			}
		}`, output, level)
		payload := strings.NewReader(body)
		req, _ := http.NewRequest("POST", url, payload)
		req = AddHeaders(req)
		res, gerr := http.DefaultClient.Do(req)
		if gerr != nil {
			return nerr.Translate(gerr).Addf("error when making call: %s", gerr)
		}
		defer res.Body.Close()
	} else {
		return nerr.Create("Invalid Output. Valid Audio Output names are Audio1 and Audio2", "")
	}

	return nil
}

// SetMute changes the input on the given output to input
func SetMute(address, output, mute string) *nerr.E {

	var err *nerr.E
	//Now we need to find out which input is being routed to the output
	err = SetMuteHelper(address, output, mute)
	if err != nil {
		return nerr.Translate(err).Addf("error when making call: %s", err)
	}
	return nil
}

//SetMuteHelper .
func SetMuteHelper(address, output, mute string) *nerr.E {
	url := fmt.Sprintf("http://%s/cgi-bin/config.cgi", address)
	if output == "1" || output == "2" {
		body := fmt.Sprintf(`
		{
			"setConfig": {
				"audio": {
					"audOut": {
						"zoneOut%s": {
							"analogOut": {
								"audioMute": %s
							}
						}
					}
				}
			}
		}`, output, mute)
		payload := strings.NewReader(body)
		req, _ := http.NewRequest("POST", url, payload)
		req = AddHeaders(req)

		res, gerr := http.DefaultClient.Do(req)
		if gerr != nil {
			return nerr.Translate(gerr).Addf("error when making call: %s", gerr)
		}
		defer res.Body.Close()
	} else {
		return nerr.Create("Invalid Output. Valid Output names are Audio1 and Audio2", "")
	}
	return nil
}
