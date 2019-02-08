package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/byuoitav/atlona-switcher-microservice/helpers"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/status"
	"github.com/byuoitav/common/structs"
	"github.com/labstack/echo"
)

func SwitchInput(context echo.Context) error {
	output := context.Param("output")

	outport, _ := strconv.Atoi(output)
	outport = outport + 1

	input := context.Param("input")

	inport, _ := strconv.Atoi(input)
	inport = inport + 1

	address := context.Param("address")

	resp, err := helpers.SwitchInput(address, fmt.Sprintf("%v", outport), fmt.Sprintf("%v", inport))
	if err != nil {
		log.L.Errorf("Failed to establish connection with %s : %s", address, err.Error())
		return context.JSON(http.StatusInternalServerError, err)
	}

	//decrement response by 1
	response, _ := strconv.Atoi(resp)
	response = response - 1
	//in:out
	return context.JSON(http.StatusOK, status.Input{Input: fmt.Sprintf("%v:%v", input, output)})
}

func ShowOutput(context echo.Context) error {
	address := context.Param("address")

	inport, outport, err := helpers.GetOutput(address)
	if err != nil {
		log.L.Errorf("Failed to establish connection with %s : %s", address, err.Error())
		return context.JSON(http.StatusInternalServerError, err)
	}
	input, _ := strconv.Atoi(inport)
	input = input - 1
	log.L.Infof("input: %d", input)
	output, _ := strconv.Atoi(outport)
	log.L.Infof("ouput: %d", output)

	return context.JSON(http.StatusOK, status.Input{Input: fmt.Sprintf("%v:%v", input, output)})
}

func HardwareInfo(context echo.Context) error {
	address := context.Param("address")
	ipaddr, macaddr, verdata, err := helpers.GetHardware(address)
	if err != nil {
		log.L.Errorf("Failed to establish connection with %s : %s", address, err.Error())
		return context.JSON(http.StatusInternalServerError, err)
	}
	return context.JSON(http.StatusOK, structs.HardwareInfo{
		NetworkInfo: structs.NetworkInfo{
			IPAddress:  ipaddr,
			MACAddress: macaddr,
		},
		FirmwareVersion: verdata,
	})
}
