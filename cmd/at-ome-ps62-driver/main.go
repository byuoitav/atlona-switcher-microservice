package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/byuoitav/atlona-driver"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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

	switchers := &sync.Map{}

	handlers := Handlers{
		CreateVideoSwitcher: func(addr string) *atlona.AtOmePs62 {
			if vs, ok := switchers.Load(addr); ok {
				return vs.(*atlona.AtOmePs62)
			}

			vs := &atlona.AtOmePs62{
				Address:      addr,
				Username:     username,
				Password:     password,
				RequestDelay: 500 * time.Millisecond,
			}

			switchers.Store(addr, vs)
			return vs
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
