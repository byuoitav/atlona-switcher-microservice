package main

import (
	"net/http"

	"github.com/byuoitav/atlona-switcher-microservice/handlers"
	"github.com/byuoitav/atlona-switcher-microservice/handlers5x1"
	"github.com/byuoitav/atlona-switcher-microservice/handlers6x2"
	"github.com/byuoitav/common"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/v2/auth"
)

func main() {
	log.SetLevel("info")
	port := ":8026"
	router := common.NewRouter()

	// Functionality Endpoints
	write := router.Group("", auth.AuthorizeRequest("write-state", "room", auth.LookupResourceFromAddress))
	write.GET("/:address/output/:output/input/:input", handlers.SwitchInput)
	// 6x2 functionality Endpoints
	write.GET("/:address/output/:output/input/:input/6x2", handlers6x2.SetInput)
	write.GET("/:address/output/:output/volume/:level/6x2", handlers6x2.SetVolume)
	write.GET("/:address/output/:output/mute/:isMuted/6x2", handlers6x2.SetMute)
	// 5x1 functionality Endpoints
	write.GET("/:address/output/:output/input/:input/5x1", handlers5x1.SetInput)
	write.GET("/:address/volume/:level/5x1/:input", handlers5x1.SetVolume)
	write.GET("/:address/mute/5x1", handlers5x1.SetMute)
	write.GET("/:address/unmute/5x1", handlers5x1.SetUnmute)

	// Status/Hardware Info Endpoints
	read := router.Group("", auth.AuthorizeRequest("read-state", "room", auth.LookupResourceFromAddress))
	read.GET("/:address/output/:output/input", handlers.GetInput)
	read.GET("/:address/hardware", handlers.HardwareInfo)
	read.GET("/:address/hardware/6x2", handlers6x2.GetHardware)
	read.GET("/:address/output/:output/input/6x2", handlers6x2.GetInput)
	read.GET("/:address/output/:output/mute/6x2", handlers6x2.GetMute)
	read.GET("/:address/output/:output/volume/6x2", handlers6x2.GetVolume)
	// 5x1 enpoints
	write.GET("/:address/output/:output/input/5x1", handlers5x1.GetInput)
	write.GET("/:address/volume/5x1/:input", handlers5x1.GetVolume)
	write.GET("/:address/muteStatus/5x1", handlers5x1.GetMute)

	// log level endpoints
	router.PUT("/log-level/:level", log.SetLogLevel)
	router.GET("/log-level", log.GetLogLevel)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
