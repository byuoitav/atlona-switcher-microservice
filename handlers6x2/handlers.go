package handlers6x2

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/atlona-switcher-microservice/switcher6x2"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/status"
	"github.com/labstack/echo"
)

// SetInput works
func SetInput(ectx echo.Context) error {
	address := ectx.Param("address")
	output := ectx.Param("output")
	input := ectx.Param("input")

	l := log.L.Named(address)
	l.Infof("Switching input for output %s to %s", output, input)

	err := switcher6x2.SetInput(address, output, input)
	if err != nil {
		l.Warnf("%s", err.Error())
		return ectx.String(http.StatusInternalServerError, err.Error())
	}

	l.Infof("Successfully changed input for output %s to %s", output, input)
	return ectx.JSON(http.StatusOK, status.Input{
		Input: fmt.Sprintf("%v:%v", input, output[(len(output)-1):]),
	})
}

// GetInput .
func GetInput(ectx echo.Context) error {
	address := ectx.Param("address")
	output := ectx.Param("output")

	l := log.L.Named(address)
	l.Infof("Getting input for output %s", output)

	input, err := switcher6x2.GetInput(address, output)
	if err != nil {
		l.Warnf("%s", err.Error())
		return ectx.String(http.StatusInternalServerError, err.Error())
	}

	l.Infof("Input for output %v is %v", output, input)

	return ectx.JSON(http.StatusOK, status.Input{
		Input: fmt.Sprintf("%v:%v", input, output[(len(output)-1):]),
	})
}

//GetMute .
func GetMute(ectx echo.Context) error {
	address := ectx.Param("address")
	output := ectx.Param("output")

	l := log.L.Named(address)
	l.Infof("Getting mute status for output %s", output)

	resp, err := switcher6x2.GetMute(address, output)
	if err != nil {
		l.Warnf("%s", err.Error())
		return ectx.String(http.StatusInternalServerError, err.Error())
	}

	l.Infof("%s Mute: %s", output, resp)

	return ectx.JSON(http.StatusOK, resp)
}

//GetVolume .
func GetVolume(ectx echo.Context) error {
	address := ectx.Param("address")
	output := ectx.Param("output")

	l := log.L.Named(address)
	l.Infof("Getting volume for output %s", output)

	resp, err := switcher6x2.GetVolume(address, output)
	if err != nil {
		l.Warnf("%s", err.Error())
		return ectx.String(http.StatusInternalServerError, err.Error())
	}

	return ectx.JSON(http.StatusOK, resp)
}

//GetHardware .
func GetHardware(ectx echo.Context) error {
	address := ectx.Param("address")

	l := log.L.Named(address)
	l.Infof("Getting Hardware for device %s", address)

	resp, err := switcher6x2.GetHardwareInfo(address)
	if err != nil {
		l.Warnf("%s", err.Error())
		return ectx.String(http.StatusInternalServerError, err.Error())
	}

	return ectx.JSON(http.StatusOK, resp)
}

// SetVolume .
func SetVolume(ectx echo.Context) error {
	address := ectx.Param("address")
	level := ectx.Param("level")
	output := ectx.Param("output")

	l := log.L.Named(address)
	l.Infof("Changing Volume on Output %s to %s", output, level)

	err := switcher6x2.SetVolume(address, output, level)
	if err != nil {
		l.Warnf("%s", err.Error())
		return ectx.String(http.StatusInternalServerError, err.Error())
	}
	return ectx.JSON(http.StatusOK, "ok")
}

// SetMute .
func SetMute(ectx echo.Context) error {
	address := ectx.Param("address")
	isMuted := ectx.Param("isMuted")
	output := ectx.Param("output")
	l := log.L.Named(address)
	l.Infof("Mute = %s", isMuted)

	err := switcher6x2.SetMute(address, output, isMuted)
	if err != nil {
		l.Warnf("%s", err.Error())
		return ectx.String(http.StatusInternalServerError, err.Error())
	}
	return ectx.JSON(http.StatusOK, "ok")
}
