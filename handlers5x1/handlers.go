package handlers5x1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/byuoitav/atlona-switcher-microservice/switcher5x1"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/status"
	"github.com/labstack/echo"
)

// SetInput works
func SetInput(ectx echo.Context) error {
	address := ectx.Param("address")
	input := ectx.Param("input")
	l := log.L.Named(address)
	l.Infof("Switching input to %s", input)
	intInput, nerr := strconv.Atoi(input)
	if nerr != nil {
		return ectx.String(http.StatusInternalServerError, nerr.Error())
	}
	if intInput > 4 {
		l.Warnf("The input requested must be between 0-4")
		return ectx.String(http.StatusInternalServerError, "Invalid Input")
	}
	ctx := ectx.Request().Context()
	er := switcher5x1.SetInput(ctx, address, input)
	if er != nil {
		l.Warnf("%s", er.Error())
		return ectx.String(http.StatusInternalServerError, er.Error())
	}

	l.Infof("Successfully changed input to %s", input)
	return ectx.JSON(http.StatusOK, status.Input{
		Input: fmt.Sprintf("%v:1", input),
	})
}

// GetInput works
func GetInput(ectx echo.Context) error {
	address := ectx.Param("address")

	l := log.L.Named(address)
	l.Infof("Getting input")
	ctx := ectx.Request().Context()
	input, err := switcher5x1.GetInput(ctx, address)
	if err != nil {
		l.Warnf("%s", err.Error())
		return ectx.String(http.StatusInternalServerError, err.Error())
	}

	l.Infof("Input is %s", input)

	return ectx.JSON(http.StatusOK, status.Input{
		Input: fmt.Sprintf("%s:1", input),
	})
}

//GetMute works
func GetMute(ectx echo.Context) error {
	address := ectx.Param("address")

	l := log.L.Named(address)
	l.Infof("Getting mute status")
	ctx := ectx.Request().Context()
	resp, err := switcher5x1.GetMute(ctx, address)
	if err != nil {
		l.Warnf("%s", err.Error())
		return ectx.String(http.StatusInternalServerError, err.Error())
	}

	l.Infof("Output Mute Status: %s", resp)

	return ectx.JSON(http.StatusOK, status.Mute{Muted: resp})
}

//GetVolume .
func GetVolume(ectx echo.Context) error {
	address := ectx.Param("address")

	l := log.L.Named(address)
	l.Infof("Getting volume")
	ctx := ectx.Request().Context()
	resp, err := switcher5x1.GetVolume(ctx, address)
	if err != nil {
		l.Warnf("%s", err.Error())
		return ectx.String(http.StatusInternalServerError, err.Error())
	}
	return ectx.JSON(http.StatusOK, status.Volume{Volume: resp})
}

//GetHardware TODO
func GetHardware(ectx echo.Context) error {
	address := ectx.Param("address")

	l := log.L.Named(address)
	l.Infof("Getting Hardware for device %s", address)
	ctx := ectx.Request().Context()
	resp, err := switcher5x1.GetHardwareInfo(ctx, address)
	if err != nil {
		l.Warnf("%s", err.Error())
		return ectx.String(http.StatusInternalServerError, err.Error())
	}

	// TODO return a status.Hardware struct
	return ectx.JSON(http.StatusOK, resp)
}

// SetVolume works
func SetVolume(ectx echo.Context) error {
	address := ectx.Param("address")
	level := ectx.Param("level")

	lev, err := strconv.Atoi(level)
	if err != nil {
		return ectx.String(http.StatusBadRequest, "bad number")
	}

	l := log.L.Named(address)
	l.Infof("Changing Volume to %s", level)
	ctx := ectx.Request().Context()
	er := switcher5x1.SetVolume(ctx, address, lev)
	if er != nil {
		l.Warnf("%s", er.Error())
		return ectx.String(http.StatusInternalServerError, er.Error())
	}

	return ectx.JSON(http.StatusOK, status.Volume{Volume: lev})
}

// SetMute Works
func SetMute(ectx echo.Context) error {
	address := ectx.Param("address")
	ctx := ectx.Request().Context()
	er := switcher5x1.SetMute(ctx, address)

	l := log.L.Named(address)
	l.Infof("Changing Mute to true")
	if er != nil {
		l.Warnf("%s", er.Error())
		return ectx.String(http.StatusInternalServerError, er.Error())
	}

	return ectx.JSON(http.StatusOK, status.Mute{Muted: true})
}

// SetUnmute Works
func SetUnmute(ectx echo.Context) error {
	address := ectx.Param("address")
	ctx := ectx.Request().Context()
	er := switcher5x1.SetUnmute(ctx, address)
	l := log.L.Named(address)
	l.Infof("Changing Mute to false")
	if er != nil {
		l.Warnf("%s", er.Error())
		return ectx.String(http.StatusInternalServerError, er.Error())
	}

	return ectx.JSON(http.StatusOK, status.Mute{Muted: false})
}
