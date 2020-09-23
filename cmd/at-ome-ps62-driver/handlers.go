package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/byuoitav/atlona-driver"
	"github.com/byuoitav/common/status"
	"github.com/labstack/echo"
)

type Handlers struct {
	CreateVideoSwitcher func(string) *atlona.AtOmePs62
}

func (h *Handlers) RegisterRoutes(group *echo.Group) {
	ps62 := group.Group("/AT-OME-PS62/:address")

	// TODO singleflight?

	// get state
	ps62.GET("/output/:output/input", func(c echo.Context) error {
		vs := h.CreateVideoSwitcher(c.Param("address"))

		inputs, err := vs.GetAudioVideoInputs(c.Request().Context())
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}

		out := c.Param("output")
		in := inputs[out]

		return c.JSON(http.StatusOK, status.Input{
			Input: fmt.Sprintf("%v:%v", in, out),
		})
	})

	ps62.GET("/block/:block/volume", func(c echo.Context) error {
		vs := h.CreateVideoSwitcher(c.Param("address"))

		vols, err := vs.GetVolumes(c.Request().Context(), []string{})
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, status.Volume{
			Volume: vols[c.Param("block")],
		})
	})

	ps62.GET("/block/:block/muted", func(c echo.Context) error {
		vs := h.CreateVideoSwitcher(c.Param("address"))

		mutes, err := vs.GetMutes(c.Request().Context(), []string{})
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, status.Mute{
			Muted: mutes[c.Param("block")],
		})
	})

	// set state
	ps62.GET("/output/:output/input/:input", func(c echo.Context) error {
		vs := h.CreateVideoSwitcher(c.Param("address"))
		out := c.Param("output")
		in := c.Param("input")

		err := vs.SetAudioVideoInput(c.Request().Context(), out, in)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, status.Input{
			Input: fmt.Sprintf("%v:%v", in, out),
		})
	})

	ps62.GET("/block/:block/volume/:volume", func(c echo.Context) error {
		vs := h.CreateVideoSwitcher(c.Param("address"))

		vol, err := strconv.Atoi(c.Param("volume"))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}

		err = vs.SetVolume(c.Request().Context(), c.Param("block"), vol)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, status.Volume{
			Volume: vol,
		})
	})

	ps62.GET("/block/:block/muted/:mute", func(c echo.Context) error {
		vs := h.CreateVideoSwitcher(c.Param("address"))

		mute, err := strconv.ParseBool(c.Param("mute"))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}

		err = vs.SetMute(c.Request().Context(), c.Param("block"), mute)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, status.Mute{
			Muted: mute,
		})
	})
}
