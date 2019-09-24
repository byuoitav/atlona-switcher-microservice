package handlers6x2

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/atlona-switcher-microservice/switcher6x2"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/status"
	"github.com/labstack/echo"
)

// SetInput .
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
		Input: fmt.Sprintf("%v:%v", input, output),
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
		Input: fmt.Sprintf("%v:%v", input, output),
	})
}

//GetMute .
func GetMute(ectx echo.Context) error {
	address := ectx.Param("address")
	output := ectx.Param("output")

	l := log.L.Named(address)
	l.Infof("Getting input for output %s", output)

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
	output := ectx.Param("output")
	level := ectx.Param("level")

	l := log.L.Named(address)
	l.Infof("Switching output %s to volume %s", output, level)

	resp, err := switcher6x2.SetVolume(address, output, level)
	if err != nil {
		l.Warnf("%s", err.Error())
		return ectx.String(http.StatusInternalServerError, err.Error())
	}
	return ectx.JSON(http.StatusOK, resp)
}

// SetMute .
func SetMute(ectx echo.Context) error {
	address := ectx.Param("address")
	output := ectx.Param("output")
	isMuted := ectx.Param("isMuted")

	l := log.L.Named(address)
	l.Infof("Output %s: Mute = %s", output, isMuted)

	resp, err := switcher6x2.SetMute(address, output, isMuted)
	if err != nil {
		l.Warnf("%s", err.Error())
		return ectx.String(http.StatusInternalServerError, err.Error())
	}
	return ectx.JSON(http.StatusOK, resp)
}
