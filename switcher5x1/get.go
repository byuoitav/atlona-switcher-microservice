package switcher5x1

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/byuoitav/common/nerr"
)

// Login Logs the user in... this should be called before every command.
func Login(ctx context.Context, address string) *nerr.E {
	url := fmt.Sprintf("http://%s/ajlogin.html?value=login&usn=root&pwd=Atlona", address)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nerr.Translate(err).Addf("error when making request: %s", err)
	}
	req = req.WithContext(ctx)
	res, gerr := http.DefaultClient.Do(req)
	if gerr != nil {
		return nerr.Translate(gerr).Addf("error when making call: %s", gerr)
	}
	body, _ := ioutil.ReadAll(res.Body)
	resp := string(body)
	splitRes := strings.Split(resp, ";")
	if splitRes[0] == "ER" {
		return nerr.Create("There is an error being returned from atlona", "")
	}
	defer res.Body.Close()
	return nil
}

// GetInput returns the current input for the requested output
func GetInput(ctx context.Context, address string) (string, *nerr.E) {
	loginerr := Login(ctx, address)
	if loginerr != nil {
		return "", nerr.Translate(loginerr).Addf("error logging in to make a request: %s", loginerr)
	}
	url := fmt.Sprintf("http://%s/ajstatus.html?value=status&uid=Y1&ro=0", address)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", nerr.Translate(err).Addf("error when making request: %s", err)
	}
	req = req.WithContext(ctx)
	res, gerr := http.DefaultClient.Do(req)
	if gerr != nil {
		return "", nerr.Translate(gerr).Addf("error when making call: %s", gerr)
	}
	body, _ := ioutil.ReadAll(res.Body)
	resp := string(body)
	splitRes := strings.Split(resp, ";")
	if splitRes[0] == "ER" {
		return "", nerr.Create("There is an error being returned from atlona", "")
	}
	defer res.Body.Close()
	return splitRes[1], nil
}

// GetHardwareInfo returns a hardware info struct         Change to structs.HardwareInfo
func GetHardwareInfo(ctx context.Context, address string) (string, *nerr.E) {
	loginerr := Login(ctx, address)
	if loginerr != nil {
		return "", nerr.Translate(loginerr).Addf("error logging in to make a request: %s", loginerr)
	}
	url := fmt.Sprintf("http://%s/ajstatus.html?value=status&uid=Y1&mlf=1&inp=", address)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", nerr.Translate(err).Addf("error when making request: %s", err)
	}
	req = req.WithContext(ctx)
	res, gerr := http.DefaultClient.Do(req)
	if gerr != nil {
		return "", nerr.Translate(gerr).Addf("error when making call: %s", gerr)
	}
	body, _ := ioutil.ReadAll(res.Body)
	resp := string(body)
	splitRes := strings.Split(resp, ";")
	if splitRes[0] == "ER" {
		return "", nerr.Create("There is an error being returned from atlona", "")
	}
	defer res.Body.Close()
	return "", nil
}

// GetMute .
func GetMute(ctx context.Context, address string) (bool, *nerr.E) {
	loginerr := Login(ctx, address)
	if loginerr != nil {
		return false, nerr.Translate(loginerr).Addf("error logging in to make a request: %s", loginerr)
	}
	url := fmt.Sprintf("http://%s/ajstatus.html?value=status&uid=Y1&ro=0", address)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, nerr.Translate(err).Addf("error when making request: %s", err)
	}
	req = req.WithContext(ctx)
	res, gerr := http.DefaultClient.Do(req)
	if gerr != nil {
		return false, nerr.Translate(gerr).Addf("error when making call: %s", gerr)
	}
	body, _ := ioutil.ReadAll(res.Body)
	resp := string(body)
	splitRes := strings.Split(resp, ";")
	if splitRes[0] == "ER" {
		fmt.Println(resp)
		return false, nerr.Create("There is an error being returned from atlona", "")
	}
	defer res.Body.Close()
	if splitRes[(len(splitRes)-2)] == "1" {
		return true, nil
	}
	return false, nil
}

//GetVolume .
func GetVolume(ctx context.Context, address string) (int, *nerr.E) {
	loginerr := Login(ctx, address)
	if loginerr != nil {
		return 0, nerr.Translate(loginerr).Addf("error logging in to make a request: %s", loginerr)
	}
	url := fmt.Sprintf("http://%s/ajstatus.html?value=status&uid=Y1&mlf=1&inp=", address)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, nerr.Translate(err).Addf("error when making request: %s", err)
	}
	req = req.WithContext(ctx)
	res, gerr := http.DefaultClient.Do(req)
	if gerr != nil {
		return 0, nerr.Translate(gerr).Addf("error when making call: %s", gerr)
	}
	body, _ := ioutil.ReadAll(res.Body)
	resp := string(body)
	splitRes := strings.Split(resp, ";")
	if splitRes[0] == "ER" {
		return 0, nerr.Create("There is an error being returned from atlona", "")
	}
	defer res.Body.Close()
	//convert response back to 0-100 value
	volumeLevel, gerr := strconv.Atoi(splitRes[2])
	if gerr != nil {
		return 0, nerr.Translate(gerr).Addf("error when making call: %s", gerr)
	}
	if volumeLevel < -35 {
		return 0, nil
	} else {
		volume := ((volumeLevel + 35) * 2)
		if volume%2 != 0 {
			volume = volume + 1
		}
		return volume, nil
	}
}
