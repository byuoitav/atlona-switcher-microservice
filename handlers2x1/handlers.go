package handlers2x1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/byuoitav/atlona-switcher-microservice/switcher2x1"
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
	if intInput != 1 && intInput != 2 {
		l.Warnf("The input requested must be 1 or 2")
		return ectx.String(http.StatusInternalServerError, "Invalid Input")
	}
	ctx := ectx.Request().Context()
	er := switcher2x1.SetInput(ctx, address, input)
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
	ctx := ectx.Request().Context()
	l := log.L.Named(address)
	l.Infof("Getting input")

	input, err := switcher2x1.GetInput(ctx, address)
	if err != nil {
		l.Warnf("%s", err.Error())
		return ectx.String(http.StatusInternalServerError, err.Error())
	}

	l.Infof("Input is %s", input)

	return ectx.JSON(http.StatusOK, status.Input{
		Input: fmt.Sprintf("%s:1", input),
	})
}
