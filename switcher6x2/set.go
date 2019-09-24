package switcher6x2

import (
	"fmt"
	"io/ioutil"
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
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Postman-Token", "5f099c23-e332-4d44-aeff-6546935ca6b2")
	return req
}

// SetInput changes the input on the given output to input
func SetInput(address, output, input string) *nerr.E {
	out, gerr := strconv.Atoi(output)
	if gerr != nil {
		return nerr.Translate(gerr).Addf("unable to switch input on %s", address)
	}

	in, gerr := strconv.Atoi(input)
	if gerr != nil {
		return nerr.Translate(gerr).Addf("unable to switch input on %s", address)
	}
	//first we need to see what the volume of the current input is
	currentVolume, err := GetVolume(address, fmt.Sprintf("out%s", output))
	if err != nil {
		return nerr.Translate(gerr).Addf("unable to switch input on %s- cannot get current volume", address)
	}
	//Then we need to check to see if other output has the same source input
	otherInput, gerr := GetInput(address, output)
	if otherInput == input {
		SetMute(address, "out2", "true")
	} else {
		SetMute(address, "out2", "false")
	}
	//send command
	url := fmt.Sprintf("http://%s/cgi-bin/config.cgi", address)
	payload := strings.NewReader("")
	if out == 1 {
		payload = strings.NewReader(fmt.Sprintf("{ \n   \"setConfig\":{ \n      \"video\":{ \n         \"vidOut\":{ \n            \"hdmiOut\":{ \n               \"hdmiOutA\":{ \n                  \"videoSrc\":%v\n               }\n            }\n         }\n      }\n   }\n}", in))
	}
	if out == 2 {
		payload = strings.NewReader(fmt.Sprintf("{ \n   \"setConfig\":{ \n      \"video\":{ \n         \"vidOut\":{ \n            \"hdmiOut\":{ \n               \"hdmiOutB\":{ \n                  \"videoSrc\":%v\n               }\n            }\n         }\n      }\n   }\n}", in))
	}

	req, _ := http.NewRequest("POST", url, payload)

	req = AddHeaders(req)

	res, gerr := http.DefaultClient.Do(req)
	if gerr != nil {
		return nerr.Translate(gerr).Addf("error when making call: %s", gerr)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
	//Now we needto set the current volume on the output
	SetVolume(address, fmt.Sprintf("out%s", output), strconv.Itoa(currentVolume))

	return nil
}

// SetVolume changes the input on the given output to input
func SetVolume(address, output, level string) (string, *nerr.E) {
	//Atlona volume levels are from -90 to 10 and the number we recieve is 0-100
	resp := ""
	volumeLevel, err := strconv.Atoi(level)
	if err != nil {
		return "", nerr.Translate(err).Add("unable to switch change volume")
	}
	volumeLevel = volumeLevel - 90

	//Now we need to find out which input is being routed to the output
	if output == "out1" || output == "out2" {
		//get what input is routed to the output
		input, nerr := GetInput(address, output[len(output)-1:])
		if nerr != nil {
			fmt.Println(nerr)
		}
		input = fmt.Sprintf("digitalIn%s", input)
		fmt.Println(input)
		resp, err = SetVolumeHelper(address, input, strconv.Itoa(volumeLevel))
		if nerr != nil {
			fmt.Println(nerr)
		}
	}
	switch {
	case output == "aux1":
		resp, err = SetVolumeHelper(address, "analogIn1", strconv.Itoa(volumeLevel))
	case output == "aux2":
		resp, err = SetVolumeHelper(address, "analogIn2", strconv.Itoa(volumeLevel))
	case output == "aux3":
		resp, err = SetVolumeHelper(address, "mic1", strconv.Itoa(volumeLevel))
	}
	if err != nil {
		fmt.Println(err)
	}
	return resp, nil
}

//SetVolumeHelper .
func SetVolumeHelper(address, input, level string) (string, *nerr.E) {
	url := fmt.Sprintf("http://%s/cgi-bin/config.cgi", address)
	fmt.Println("Input: " + input + "   level:" + level)
	payload := strings.NewReader("{\n    \"setConfig\": {\n        \"audio\": {\n            \"audIn\": {\n                \"" + input + "\": {\n                \t\"audioVol\": " + level + "\n                }\n            }\n        }\n    }\n}")
	req, _ := http.NewRequest("POST", url, payload)

	req = AddHeaders(req)

	res, gerr := http.DefaultClient.Do(req)
	if gerr != nil {
		return "", nerr.Translate(gerr).Addf("error when making call: %s", gerr)
	}
	defer res.Body.Close()
	return "", nil
}

// SetMute changes the input on the given output to input
func SetMute(address, output, mute string) (string, *nerr.E) {
	resp := ""
	var err *nerr.E
	//Now we need to find out which input is being routed to the output
	if output == "out1" || output == "out2" {
		//get what input is routed to the output
		input, nerr := GetInput(address, output[len(output)-1:])
		if nerr != nil {
			fmt.Println(nerr)
		}
		input = fmt.Sprintf("digitalIn%s", input)
		fmt.Println(input)
		resp, err := SetMuteHelper(address, input, mute)
		if err != nil {
			fmt.Println(err)
		}
		return resp, nil
	}
	switch {
	case output == "aux1":
		resp, err := SetMuteHelper(address, "analogIn1", mute)
		if err != nil {
			fmt.Println(err)
		}
		return resp, nil
	case output == "aux2":
		resp, err := SetMuteHelper(address, "analogIn2", mute)
		if err != nil {
			fmt.Println(err)
		}
		return resp, nil
	case output == "aux3":
		resp, err := SetMuteHelper(address, "mic1", mute)
		if err != nil {
			fmt.Println(err)
		}
		return resp, nil
	}
	if err != nil {
		fmt.Println(err)
	}
	return resp, nil
}

//SetMuteHelper .
func SetMuteHelper(address, input, mute string) (string, *nerr.E) {
	url := fmt.Sprintf("http://%s/cgi-bin/config.cgi", address)
	payload := strings.NewReader("{\n    \"setConfig\": {\n        \"audio\": {\n            \"audIn\": {\n                \"" + input + "\": {\n                    \"audioMute\": " + mute + "\n                }\n            }\n        }\n    }\n}")
	req, _ := http.NewRequest("POST", url, payload)

	req = AddHeaders(req)

	res, gerr := http.DefaultClient.Do(req)
	if gerr != nil {
		return "", nerr.Translate(gerr).Addf("error when making call: %s", gerr)
	}
	defer res.Body.Close()
	return "", nil
}
