package handlers

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/atlona-switcher-microservice/switcher"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/status"
	"github.com/labstack/echo"
)

// SwitchInput .
func SwitchInput(ectx echo.Context) error {
	address := ectx.Param("address")
	output := ectx.Param("output")
	input := ectx.Param("input")

	l := log.L.Named(address)
	l.Infof("Switching input for output %s to %s", output, input)

	err := switcher.SwitchInput(ectx.Request().Context(), address, output, input)
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

	input, err := switcher.GetInput(ectx.Request().Context(), address)
	if err != nil {
		l.Warnf("%s", err.Error())
		return ectx.String(http.StatusInternalServerError, err.Error())
	}

	l.Infof("Input for output %v is %v", output, input)

	return ectx.JSON(http.StatusOK, status.Input{
		Input: fmt.Sprintf("%v:%v", input, output),
	})
}

// HardwareInfo .
func HardwareInfo(ectx echo.Context) error {
	address := ectx.Param("address")

	l := log.L.Named(address)
	l.Infof("Getting hardware info")

	info, err := switcher.GetHardwareInfo(ectx.Request().Context(), address)
	if err != nil {
		l.Warnf("%s", err.Error())
		return ectx.String(http.StatusInternalServerError, err.Error())
	}

	l.Infof("Successfully got hardware info")

	return ectx.JSON(http.StatusOK, info)
}
