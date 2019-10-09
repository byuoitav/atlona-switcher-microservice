package switcher5x1

import (
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strings"

	"github.com/byuoitav/common/nerr"
)

// SetInput changes the input on the given output to input
func SetInput(ctx context.Context, address, input string) *nerr.E {
	loginerr := Login(ctx, address)
	if loginerr != nil {
		return nerr.Translate(loginerr).Addf("error logging in to make a request: %s", loginerr)
	}
	url := fmt.Sprintf("http://%s/ajstatus.html?value=status&uid=Y1&mlf=1&inp=%s", address, input)
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

// SetVolume changes the input on the given output to input
func SetVolume(ctx context.Context, address string, level int) *nerr.E {
	loginerr := Login(ctx, address)
	if loginerr != nil {
		return nerr.Translate(loginerr).Addf("error logging in to make a request: %s", loginerr)
	}
	//Atlona volume levels are from -80 to 15 and the number we recieve is 0-100
	//if volume level is supposed to be zero set it to zero (which is -90) on atlona
	if level == 0 {
		level = -80
	} else {
		convertedVolume := -35 + math.Round(float64(level/2))
		level = int(convertedVolume)
	}
	err := SetVolumeHelper(ctx, address, level)
	if err != nil {
		return nerr.Translate(err).Add("unable to switch change volume")
	}
	return nil
}

//SetVolumeHelper .
func SetVolumeHelper(ctx context.Context, address string, level int) *nerr.E {
	loginerr := Login(ctx, address)
	if loginerr != nil {
		return nerr.Translate(loginerr).Addf("error logging in to make a request: %s", loginerr)
	}
	url := fmt.Sprintf("http://%s/ajstatus.html?value=status&uid=Z3&mlf=1&vol=%v", address, level)
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
	fmt.Println(resp)
	splitRes := strings.Split(resp, ";")
	if splitRes[0] == "ER" {
		return nerr.Create("There is an error being returned from atlona", "")
	}
	defer res.Body.Close()
	return nil
}

// SetMute changes the input on the given output to input
func SetMute(ctx context.Context, address string) *nerr.E {
	loginerr := Login(ctx, address)
	if loginerr != nil {
		return nerr.Translate(loginerr).Addf("error logging in to make a request: %s", loginerr)
	}
	url := fmt.Sprintf("http://%s/ajstatus.html?value=status&uid=Y1&mlf=1&lraud=1", address)
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

// SetUnmute changes the input on the given output to input
func SetUnmute(ctx context.Context, address string) *nerr.E {
	loginerr := Login(ctx, address)
	if loginerr != nil {
		return nerr.Translate(loginerr).Addf("error logging in to make a request: %s", loginerr)
	}
	url := fmt.Sprintf("http://%s/ajstatus.html?value=status&uid=Y1&mlf=1&lraud=0", address)
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
