package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/byuoitav/atlona-driver"
	atgain60 "github.com/byuoitav/atlona/AT-GAIN-60"
	"github.com/byuoitav/common/status"
	"github.com/labstack/echo"
)

type Handlers struct {
	CreateVideoSwitcher func(string) *atlona.AtOmePs62
	CreateAmp           func(string) *atgain60.Amp
}

func (h *Handlers) RegisterRoutes(group *echo.Group) {
	ps62 := group.Group("/AT-OME-PS62/:address")

	// TODO singleflight?

	// get state
	ps62.GET("/output/:output/input", func(c echo.Context) error {
		addr := c.Param("address")
		vs := h.CreateVideoSwitcher(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)

		l.Printf("Getting inputs")

		inputs, err := vs.GetAudioVideoInputs(c.Request().Context())
		if err != nil {
			l.Printf("unable to get inputs: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		out := c.Param("output")
		in, ok := inputs[out]
		if !ok {
			l.Printf("invalid output %q requested", out)
			return c.String(http.StatusBadRequest, "invalid output")
		}

		l.Printf("Got inputs: %+v", inputs)
		return c.JSON(http.StatusOK, status.Input{
			Input: fmt.Sprintf("%v:%v", in, out),
		})
	})

	ps62.GET("/block/:block/volume", func(c echo.Context) error {
		addr := c.Param("address")
		vs := h.CreateVideoSwitcher(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)

		l.Printf("Getting volumes")

		vols, err := vs.GetVolumes(c.Request().Context(), []string{})
		if err != nil {
			l.Printf("unable to get volumes: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		block := c.Param("block")
		vol, ok := vols[block]
		if !ok {
			l.Printf("invalid block %q requested", block)
			return c.String(http.StatusBadRequest, "invalid block")
		}

		l.Printf("Got volumes: %+v", vols)
		return c.JSON(http.StatusOK, status.Volume{
			Volume: vol,
		})
	})

	ps62.GET("/block/:block/muted", func(c echo.Context) error {
		addr := c.Param("address")
		vs := h.CreateVideoSwitcher(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)

		l.Printf("Getting mutes")

		mutes, err := vs.GetMutes(c.Request().Context(), []string{})
		if err != nil {
			l.Printf("unable to get mutes: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		block := c.Param("block")
		mute, ok := mutes[block]
		if !ok {
			l.Printf("invalid block %q requested", block)
			return c.String(http.StatusBadRequest, "invalid block")
		}

		l.Printf("Got mutes: %+v", mutes)
		return c.JSON(http.StatusOK, status.Mute{
			Muted: mute,
		})
	})

	// set state
	ps62.GET("/output/:output/input/:input", func(c echo.Context) error {
		addr := c.Param("address")
		vs := h.CreateVideoSwitcher(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)
		out := c.Param("output")
		in := c.Param("input")

		l.Printf("Setting AV input on %q to %q", out, in)

		err := vs.SetAudioVideoInput(c.Request().Context(), out, in)
		if err != nil {
			l.Printf("unable to set AV input: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		l.Printf("Set AV input")
		return c.JSON(http.StatusOK, status.Input{
			Input: fmt.Sprintf("%v:%v", in, out),
		})
	})

	ps62.GET("/block/:block/volume/:volume", func(c echo.Context) error {
		addr := c.Param("address")
		vs := h.CreateVideoSwitcher(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)
		block := c.Param("block")

		vol, err := strconv.Atoi(c.Param("volume"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		l.Printf("Setting volume on %q to %d", block, vol)

		err = vs.SetVolume(c.Request().Context(), block, vol)
		if err != nil {
			l.Printf("unable to set volume: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		l.Printf("Set volume")
		return c.JSON(http.StatusOK, status.Volume{
			Volume: vol,
		})
	})

	ps62.GET("/block/:block/muted/:mute", func(c echo.Context) error {
		addr := c.Param("address")
		vs := h.CreateVideoSwitcher(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)
		block := c.Param("block")

		mute, err := strconv.ParseBool(c.Param("mute"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		l.Printf("Setting mute on %q to %t", block, mute)

		err = vs.SetMute(c.Request().Context(), block, mute)
		if err != nil {
			l.Printf("unable to set mute: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		l.Printf("Set mute")
		return c.JSON(http.StatusOK, status.Mute{
			Muted: mute,
		})
	})

	gain60 := group.Group("/AT-GAIN-60/:address")

	// get state
	gain60.GET("/block/:block/volume", func(c echo.Context) error {
		addr := c.Param("address")
		amp := h.CreateAmp(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)

		l.Printf("Getting volumes")

		ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
		defer cancel()

		vols, err := amp.Volumes(ctx, []string{})
		if err != nil {
			l.Printf("unable to get volumes: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		l.Printf("Got volumes: %+v", vols)

		block := c.Param("block")
		vol, ok := vols[block]
		if !ok {
			l.Printf("invalid block %q requested", block)
			return c.String(http.StatusBadRequest, "invalid block")
		}

		return c.JSON(http.StatusOK, status.Volume{
			Volume: vol,
		})
	})

	gain60.GET("/block/:block/muted", func(c echo.Context) error {
		addr := c.Param("address")
		amp := h.CreateAmp(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)

		l.Printf("Getting mutes")

		ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
		defer cancel()

		mutes, err := amp.Mutes(ctx, []string{})
		if err != nil {
			l.Printf("unable to get mutes: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		l.Printf("Got mutes: %+v", mutes)

		block := c.Param("block")
		mute, ok := mutes[block]
		if !ok {
			l.Printf("invalid block %q requested", block)
			return c.String(http.StatusBadRequest, "invalid block")
		}

		return c.JSON(http.StatusOK, status.Mute{
			Muted: mute,
		})
	})

	// set state
	gain60.GET("/block/:block/volume/:volume", func(c echo.Context) error {
		addr := c.Param("address")
		amp := h.CreateAmp(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)
		block := c.Param("block")

		vol, err := strconv.Atoi(c.Param("volume"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		l.Printf("Setting volume on %q to %d", block, vol)

		ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
		defer cancel()

		err = amp.SetVolume(ctx, block, vol)
		if err != nil {
			l.Printf("unable to set volume: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		l.Printf("Set volume")
		return c.JSON(http.StatusOK, status.Volume{
			Volume: vol,
		})
	})

	gain60.GET("/block/:block/muted/:mute", func(c echo.Context) error {
		addr := c.Param("address")
		amp := h.CreateAmp(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)
		block := c.Param("block")

		mute, err := strconv.ParseBool(c.Param("mute"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		l.Printf("Setting mute on %q to %t", block, mute)

		ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
		defer cancel()

		err = amp.SetMute(ctx, block, mute)
		if err != nil {
			l.Printf("unable to set mute: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		l.Printf("Set mute")
		return c.JSON(http.StatusOK, status.Mute{
			Muted: mute,
		})
	})
}
