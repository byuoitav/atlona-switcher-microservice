package switcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/common/structs"
)

const (
	avSettingsPage = "avs"
	infoPage       = "info"
)

// AVSettings is the response from the switcher for the av settings page
type AVSettings struct {
	HDMIInputAudioBreakout int   `json:"ARC"`
	HDCPSettings           []int `json:"HDCPSet"`
	AudioOutput            int   `json:"HDMIAud"`
	Toslink                int   `json:"Toslink"`
	AutoSwitch             int   `json:"asw"`
	Input                  int   `json:"inp"`
	LoggedIn               int   `json:"login_ur"`
}

// Info is the response from the switcher for the info page
type Info struct {
	// SystemInfo should be length 3 with:
	//	0: Model Name
	//  1: Software Version
	//  2: On-Time(d-h:m:s)
	SystemInfo []string `json:"info_val1"`

	// VideoInfo should be length 6 with:
	//	0: Active Input: int
	//  1: Signal Type: string
	//  2: Video Format: string
	//  3: Aspect: string
	//  4: Color Space: string
	//  5: Color Depth: string
	VideoInfo []interface{} `json:"info_val2"`
	LoggedIn  int           `json:"login_ur"`
}

// SystemSettings .
type SystemSettings struct {
}

func getPage(ctx context.Context, address, page string, structToFill interface{}) *nerr.E {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/aj.html?a=%s", address, page), nil)
	if err != nil {
		return nerr.Translate(err).Addf("unable to get page %s on %s", page, address)
	}

	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nerr.Translate(err).Addf("unable to get page %s on %s", page, address)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nerr.Translate(err).Addf("unable to get page %s on %s", page, address)
	}

	if resp.StatusCode/100 != 2 {
		return nerr.Createf("error", "unable to get page %s on %s - %v response recevied. body: %s", page, address, resp.StatusCode, b)
	}

	err = json.Unmarshal(b, structToFill)
	if err != nil {
		return nerr.Translate(err).Addf("unable to get page %s on %s", page, address)
	}

	return nil
}

func sendCommand(ctx context.Context, address, command string) *nerr.E {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%v/aj.html?a=command&cmd=%s", address, command), nil)
	if err != nil {
		return nerr.Translate(err).Addf("unable to send command '%s' to %s", command, address)
	}

	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nerr.Translate(err).Addf("unable to send command '%s' to %s", command, address)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nerr.Translate(err).Addf("unable to send command '%s' to %s", command, address)
	}

	if resp.StatusCode/100 != 2 {
		return nerr.Createf("error", "unable to send command '%s' to %s - %v response received. body: %s", command, address, resp.StatusCode, b)
	}

	return nil
}

// TODO finish this :)
func getNetworkSettings(ctx context.Context, address string) (structs.NetworkInfo, *nerr.E) {
	var info structs.NetworkInfo

	// get the ip info (bleh, gross. it's in the html)
	req, gerr := http.NewRequest("GET", fmt.Sprintf("http://%v", address), nil)
	if gerr != nil {
		return info, nerr.Translate(gerr).Addf("unable to get network settings from %s", address)
	}

	req = req.WithContext(ctx)
	resp, gerr := http.DefaultClient.Do(req)
	if gerr != nil {
		return info, nerr.Translate(gerr).Addf("unable to get network settings from %s", address)
	}
	defer resp.Body.Close()

	b, gerr := ioutil.ReadAll(resp.Body)
	if gerr != nil {
		return info, nerr.Translate(gerr).Addf("unable to get network settings from %s", address)
	}

	if resp.StatusCode/100 != 2 {
		return info, nerr.Createf("error", "unable to get network settings from %s - %v response received. body: %s", address, resp.StatusCode, b)
	}

	return info, nil
}
