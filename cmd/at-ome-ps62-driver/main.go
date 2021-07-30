package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	atgain60 "github.com/byuoitav/atlona/AT-GAIN-60"
	atomeps62 "github.com/byuoitav/atlona/AT-OME-PS62"
	atuhdsw52ed "github.com/byuoitav/atlona/AT-UHD-SW-52ED"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	var (
		port     int
		username string
		password string
	)

	pflag.IntVarP(&port, "port", "P", 8080, "port to run the server on")
	pflag.StringVarP(&username, "username", "u", "", "username for device")
	pflag.StringVarP(&password, "password", "p", "", "password for device")
	pflag.Parse()

	addr := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("failed to start server: %s\n", err)
		os.Exit(1)
	}

	switchers := &sync.Map{}
	amps := &sync.Map{}

	cfg := zap.NewProductionConfig()
	cfg.Level.SetLevel(zapcore.DebugLevel)
	zapLog, _ := cfg.Build()

	handlers := Handlers{
		CreateVideoSwitcher6x2: func(addr string) *atomeps62.AtlonaVideoSwitcher6x2 {
			if vs, ok := switchers.Load(addr); ok {
				return vs.(*atomeps62.AtlonaVideoSwitcher6x2)
			}

			vs := &atomeps62.AtlonaVideoSwitcher6x2{
				Address:      addr,
				Username:     username,
				Password:     password,
				RequestDelay: 500 * time.Millisecond,
			}

			switchers.Store(addr, vs)
			return vs
		},
		CreateVideoSwitcher5x1: func(addr string) *atuhdsw52ed.AtlonaVideoSwitcher5x1 {
			if vs, ok := switchers.Load(addr); ok {
				return vs.(*atuhdsw52ed.AtlonaVideoSwitcher5x1)
			}

			vs := atuhdsw52ed.NewAtlonaVideoSwitcher5x1(addr)

			switchers.Store(addr, vs)
			return vs
		},
		CreateAmp: func(addr string) *atgain60.Amp {
			if amp, ok := amps.Load(addr); ok {
				return amp.(*atgain60.Amp)
			}

			amp := &atgain60.Amp{
				Address:      addr,
				Username:     username,
				Password:     password,
				Log:          zapLog.Named(addr),
				RequestDelay: 500 * time.Millisecond,
			}

			amps.Store(addr, amp)
			return amp
		},
	}

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())

	api := e.Group("/api/v1")
	handlers.RegisterRoutes(api)

	log.Printf("Server started on %v", lis.Addr())
	if err := e.Server.Serve(lis); err != nil {
		log.Printf("unable to serve: %s", err)
	}
}
