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

// GetInput returns the current input
func GetInput(address, output string) (string, *nerr.E) {
	var resp structs.AtlonaVideoWrapper
	url := fmt.Sprintf("http://%s/cgi-bin/config.cgi", address)

	payload := strings.NewReader("{\n\t\"getConfig\": {\n\t\t\"video\": {\n\t\t\t\"vidOut\": {\n\t\t\t\t\"hdmiOut\": {\n\t\t\t\t\t\n\t\t\t\t}\n\t\t\t}\n\t\t}\n\t\t\n\t}\n\t\n}")

	req, _ := http.NewRequest("POST", url, payload)

	req = AddHeaders(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	err = json.Unmarshal([]byte(body), &resp) // here!
	if err != nil {
		panic(err)
	}
	//Get the inputsrc for the requested output
	input := ""
	if output == "out1" || output == "1" {
		input = strconv.Itoa(resp.Video.VidOut.HdmiOut.Output1.Src)
	}
	if output == "out2" || output == "2" {
		input = strconv.Itoa(resp.Video.VidOut.HdmiOut.Output2.Src)
	}
	return input, nil
}

// GetHardwareInfo returns a hardware info struct         Change to structs.HardwareInfo
func GetHardwareInfo(address string) (structs.NetworkInfo, *nerr.E) {
	var resp structs.NetworkWrapper
	url := fmt.Sprintf("http://%s/cgi-bin/config.cgi", address)

	payload := strings.NewReader("{\n\t\"getConfig\": {\n\t\t\"network\": {\n\t\t\t\"eth0\":{}\n\t\t}\n\t}\n\t\n}")
	req, _ := http.NewRequest("POST", url, payload)

	req = AddHeaders(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	err = json.Unmarshal([]byte(body), &resp) // here!
	if err != nil {
		panic(err)
	}
	fmt.Print(resp)
	return resp.Network.Ethernet, nil
}

// GetMute .
func GetMute(address, output string) (bool, *nerr.E) {
	var resp structs.AtlonaAudioWrapper
	var isMuted bool
	url := fmt.Sprintf("http://%s/cgi-bin/config.cgi", address)

	payload := strings.NewReader("{\n\t\"getConfig\": {\n\t\t\"audio\": {\n\t\t\t\"audIn\":{}\n\t\t}\n\t\t\n\t}\n\t\n}")

	req, _ := http.NewRequest("POST", url, payload)

	req = AddHeaders(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	err = json.Unmarshal([]byte(body), &resp) // here!
	if err != nil {
		panic(err)
	}
	fmt.Print(resp)
	//output will either be out1, out2, aux1, aux2, aux3
	if output == "out1" || output == "out2" || output == "1" || output == "2" {
		//get what input is routed to the output
		input, nerr := GetInput(address, output[len(output)-1:])
		if nerr != nil {
			panic(err)
		}
		//check if the input is muted
		fmt.Println(input)
		switch {
		case input == "1":
			isMuted = resp.Audio.AudIn.Input1.Mute
		case input == "2":
			isMuted = resp.Audio.AudIn.Input2.Mute
		case input == "3":
			isMuted = resp.Audio.AudIn.Input3.Mute
		case input == "4":
			isMuted = resp.Audio.AudIn.Input4.Mute
		case input == "5":
			isMuted = resp.Audio.AudIn.Input5.Mute
		case input == "6":
			isMuted = resp.Audio.AudIn.Input6.Mute
		}
	}
	switch {
	case output == "aux1":
		isMuted = resp.Audio.AudIn.Aux1.Mute
	case output == "aux2":
		isMuted = resp.Audio.AudIn.Aux2.Mute
	case output == "aux3":
		isMuted = resp.Audio.AudIn.Mic1.Mute
	}
	return isMuted, nil
}

// GetVolume .
func GetVolume(address, output string) (int, *nerr.E) {
	var resp structs.AtlonaAudioWrapper
	var volume int
	url := fmt.Sprintf("http://%s/cgi-bin/config.cgi", address)

	payload := strings.NewReader("{\n\t\"getConfig\": {\n\t\t\"audio\": {\n\t\t\t\"audIn\":{}\n\t\t}\n\t\t\n\t}\n\t\n}")

	req, _ := http.NewRequest("POST", url, payload)

	req = AddHeaders(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	err = json.Unmarshal([]byte(body), &resp) // here!
	if err != nil {
		panic(err)
	}
	fmt.Print(resp)
	//output will either be out1, out2, aux1, aux2, aux3
	if output == "out1" || output == "out2" {
		//get what input is routed to the output
		input, nerr := GetInput(address, output[len(output)-1:])
		if nerr != nil {
			panic(err)
		}
		//check if the input is muted
		fmt.Println(input)
		switch {
		case input == "1":
			volume = resp.Audio.AudIn.Input1.Volume
		case input == "2":
			volume = resp.Audio.AudIn.Input2.Volume
		case input == "3":
			volume = resp.Audio.AudIn.Input3.Volume
		case input == "4":
			volume = resp.Audio.AudIn.Input4.Volume
		case input == "5":
			volume = resp.Audio.AudIn.Input5.Volume
		case input == "6":
			volume = resp.Audio.AudIn.Input6.Volume
		}
	}
	switch {
	case output == "aux1":
		volume = resp.Audio.AudIn.Aux1.Volume
	case output == "aux2":
		volume = resp.Audio.AudIn.Aux2.Volume
	case output == "aux3":
		volume = resp.Audio.AudIn.Mic1.Volume
	}
	//Add 90 because volume range is from -90 to 10 (convert to 0-100)
	return volume + 90, nil
}
