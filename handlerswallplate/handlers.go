package handlerswallplate

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/atlona-switcher-microservice/switcherwallplate"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/status"
	"github.com/labstack/echo"
)

// SetInput works
func SetInput(ectx echo.Context) error {
	address := ectx.Param("address")
	input := ectx.Param("input")

	l := log.L.Named(address)
	l.Infof("Switching input for output to %s", input)

	err := switcherwallplate.SetInput(address, input)
	if err != nil {
		l.Warnf("%s", err.Error())
		return ectx.String(http.StatusInternalServerError, err.Error())
	}

	l.Infof("Successfully changed input for output to %s", input)
	return ectx.JSON(http.StatusOK, status.Input{
		Input: fmt.Sprintf("%v:1", input),
	})
}

// GetInput .
func GetInput(ectx echo.Context) error {
	address := ectx.Param("address")

	l := log.L.Named(address)
	l.Info("Getting input")

	input, err := switcherwallplate.GetInput(address)
	if err != nil {
		l.Warnf("%s", err.Error())
		return ectx.String(http.StatusInternalServerError, err.Error())
	}

	l.Infof("Input for output is %v", input)

	return ectx.JSON(http.StatusOK, status.Input{
		Input: fmt.Sprintf("%v:1", input),
	})
}
