package switcher6x2

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/byuoitav/common/nerr"
)

//AddHeaders adds all of the security headers needed for making the call
func AddHeaders(req *http.Request) *http.Request {
	req.Header.Add("Content-Type", "application/json")
	//This needs to be replaced with an environmental variable
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(os.Getenv("ATLONA_AUTH"))))
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
	if output == "zoneOut1" {
		payload = strings.NewReader(fmt.Sprintf("{ \n   \"setConfig\":{ \n      \"video\":{ \n         \"vidOut\":{ \n            \"hdmiOut\":{ \n               \"hdmiOutA\":{ \n                  \"videoSrc\":%v\n               }\n            }\n         }\n      }\n   }\n}", in))
	} else if output == "zoneOut2" {
		payload = strings.NewReader(fmt.Sprintf("{ \n   \"setConfig\":{ \n      \"video\":{ \n         \"vidOut\":{ \n            \"hdmiOut\":{ \n               \"hdmiOutB\":{ \n                  \"videoSrc\":%v\n               }\n            }\n         }\n      }\n   }\n}", in))
	} else {
		return nerr.Create("Invalid Output. Valid Output names are zoneOut1 and zoneOut2", "")
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
	volumeLevel = volumeLevel - 90
	err = SetVolumeHelper(address, output, strconv.Itoa(volumeLevel))
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

//SetVolumeHelper .
func SetVolumeHelper(address, output, level string) *nerr.E {
	url := fmt.Sprintf("http://%s/cgi-bin/config.cgi", address)
	if output == "zoneOut1" || output == "zoneOut2" {
		body := fmt.Sprintf("{\n    \"setConfig\": {\n        \"audio\": {\n            \"audOut\": {\n                \"%s\": {\n                    \"audioVol\": %s\n                    }\n                }\n            }\n        }\n    \n}", output, level)
		payload := strings.NewReader(body)
		req, _ := http.NewRequest("POST", url, payload)

		req = AddHeaders(req)

		res, gerr := http.DefaultClient.Do(req)
		if gerr != nil {
			return nerr.Translate(gerr).Addf("error when making call: %s", gerr)
		}
		defer res.Body.Close()
	} else {
		return nerr.Create("Invalid Output. Valid Output names are zoneOut1 and zoneOut2", "")
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
	if output == "zoneOut1" || output == "zoneOut2" {
		body := fmt.Sprintf("{\n    \"setConfig\": {\n        \"audio\": {\n            \"audOut\": {\n                \"%s\": {\n                \t\"analogOut\": {\n                \t\t\"audioMute\": %s\n                \t}\n                    }\n                }\n            }\n        }\n    \n}", output, mute)
		payload := strings.NewReader(body)
		req, _ := http.NewRequest("POST", url, payload)
		req = AddHeaders(req)

		res, gerr := http.DefaultClient.Do(req)
		if gerr != nil {
			return nerr.Translate(gerr).Addf("error when making call: %s", gerr)
		}
		defer res.Body.Close()
	} else {
		return nerr.Create("Invalid Output. Valid Output names are zoneOut1 and zoneOut2", "")
	}
	return nil
}
