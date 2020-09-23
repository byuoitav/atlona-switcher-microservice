package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/byuoitav/atlona-driver"
	"github.com/labstack/echo"
	"github.com/spf13/pflag"
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

	handlers := Handlers{
		CreateVideoSwitcher: func(addr string) *atlona.AtOmePs62 {
			// TODO need a map
			vs := &atlona.AtOmePs62{
				Address:      addr,
				Username:     username,
				Password:     password,
				RequestDelay: 500 * time.Millisecond,
			}

			return vs
		},
	}

	e := echo.New()

	api := e.Group("/api/v1")
	handlers.RegisterRoutes(api)

	if err := e.Server.Serve(lis); err != nil {
	}
}
