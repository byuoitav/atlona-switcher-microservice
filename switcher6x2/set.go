package switcher6x2

import (
	"context"
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
func SetInput(ctx context.Context, address, output, input string) *nerr.E {
	in, err := strconv.Atoi(input)
	if err != nil {
		return nerr.Translate(err).Addf("error when making call: %s", err)
	}

	url := fmt.Sprintf("http://%s/cgi-bin/config.cgi", address)

	var payload *strings.Reader
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
		payload = strings.NewReader(fmt.Sprintf(`{"setConfig":{"video":{"vidOut":{"hdmiOut":{"mirror":{"videoSrc":%v}}}}}}`, in))
	}

	req, _ := http.NewRequest("POST", url, payload)
	req = AddHeaders(req)
	req = req.WithContext(ctx)
	res, gerr := http.DefaultClient.Do(req)
	if gerr != nil {
		return nerr.Translate(gerr).Addf("error when making call: %s", gerr)
	}
	defer res.Body.Close()

	return nil
}

// SetVolume changes the input on the given output to input
func SetVolume(ctx context.Context, address, output string, level int) *nerr.E {
	//Atlona volume levels are from -90 to 10 and the number we recieve is 0-100
	//if volume level is supposed to be zero set it to zero (which is -90) on atlona
	if level == 0 {
		level = -90
	} else {
		convertedVolume := -40 + math.Round(float64(level/2))
		level = int(convertedVolume)
	}
	err := SetVolumeHelper(ctx, address, output, level)
	if err != nil {
		return nerr.Translate(err).Add("unable to switch change volume")
	}
	return nil
}

//SetVolumeHelper .
func SetVolumeHelper(ctx context.Context, address, output string, level int) *nerr.E {
	url := fmt.Sprintf("http://%s/cgi-bin/config.cgi", address)
	if output == "1" || output == "2" {
		body := fmt.Sprintf(`
		{
			"setConfig": {
				"audio": {
					"audOut": {
						"zoneOut%s": {
							"audioVol": %d
						}
					}
				}
			}
		}`, output, level)
		payload := strings.NewReader(body)
		req, _ := http.NewRequest("POST", url, payload)
		req = AddHeaders(req)
		req = req.WithContext(ctx)
		res, gerr := http.DefaultClient.Do(req)
		if gerr != nil {
			return nerr.Translate(gerr).Addf("error when making call: %s", gerr)
		}
		defer res.Body.Close()
	} else {
		return nerr.Create("Invalid Output. Valid Audio Output names are Audio1 and Audio2: you gave us "+output, "")
	}

	return nil
}

// SetMute changes the input on the given output to input
func SetMute(ctx context.Context, address, output, mute string) *nerr.E {

	var err *nerr.E
	//Now we need to find out which input is being routed to the output
	err = SetMuteHelper(ctx, address, output, mute)
	if err != nil {
		return nerr.Translate(err).Addf("error when making call: %s", err)
	}
	return nil
}

//SetMuteHelper .
func SetMuteHelper(ctx context.Context, address, output, mute string) *nerr.E {
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
		req = req.WithContext(ctx)
		res, gerr := http.DefaultClient.Do(req)
		if gerr != nil {
			return nerr.Translate(gerr).Addf("error when making call: %s", gerr)
		}
		defer res.Body.Close()
	} else {
		return nerr.Create("Invalid Output. Valid Output names are Audio1 and Audio2 you gave us "+output, "")
	}
	return nil
}
