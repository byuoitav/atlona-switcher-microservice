package helpers

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"regexp"

	"github.com/byuoitav/common/structs"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
)

//This function returns the current input that is being shown as the output
func GetOutput(address string) (string, string, *nerr.E) {
	conn, err := getConnection(address, true)
	if err != nil {
		log.L.Errorf("Failed to establish connection with %s : %s", address, err.Error())
		return "", "", nerr.Translate(err).Add("Telnet connection failed")
	}

	//close connection
	defer conn.Close()

	log.L.Info("This is coming soon")
	conn.Write([]byte(fmt.Sprintf("STA\r\n")))

	//capture response
	var buf bytes.Buffer
	io.Copy(&buf, conn)

	//regex black magic
	reg, err := regexp.Compile("Video Output : Input = ([0-9]{2}),")
	if err != nil {
		log.L.Errorf("Failed to read from %s : %s", address, err.Error())
		return "", "", nerr.Translate(err).Add("failed to create regex")
	}
	ReturnInput := reg.FindAllStringSubmatch(fmt.Sprintf("%s", buf.Bytes()), -1)

	input := ReturnInput[0][1]
	input = input[1:]
	output := "0"

	return fmt.Sprintf("%s", input), fmt.Sprintf("%s", output), nil
}

//This function gets the IP Address (ipaddr), Software and hardware
//version (verdata), and mac address (macaddr) of the device
func GetHardware(address string) (string, string, string, *nerr.E) {
	conn, gerr := getConnection(address, true)
	if gerr != nil {
		log.L.Errorf("Failed to get connection with %s: %s", address, gerr.Error())
		return "", "", "", nerr.Translate(gerr).Add("Telnet connection failed")
	}
	//close connection
	defer conn.Close()

	ipaddr, verdata, macaddr, err := geteverything(address, conn)
	if err != nil {
		log.L.Errorf("Failed to establish connection with %s : %s", address, err.Error())
		return "", "", "", err.Add("Telnet connection failed")
	}

	conn.Close()
	return ipaddr, macaddr, verdata, nil
}

//goes and gets all of the info (I ran into issues when I tried to seperate all of them into seperate functions)
func geteverything(address string, conn *net.TCPConn) (string, string, string, *nerr.E) {
	conn.Write([]byte(fmt.Sprintf("STA\r\n")))

	//capture response
	var buf bytes.Buffer
	io.Copy(&buf, conn)

	//regex black magic
	reg, err := regexp.Compile("Host IP Address = ([0-9]{3}.[0-9]{3}.[0-9]{3}.[0-9]{3})")
	if err != nil {
		log.L.Errorf("Failed to read from %s : %s", address, err.Error())
		return "", "", "", nerr.Translate(err).Add("failed to create regex")
	}
	ReturnInput := reg.FindAllStringSubmatch(fmt.Sprintf("%s", buf.Bytes()), -1)
	ipaddr := ReturnInput[0][1]
	log.L.Infof("IP Addr: %s", ipaddr)

	//Version
	reg, err = regexp.Compile("Version : ([0-9]+.[0-9]+)")
	if err != nil {
		log.L.Errorf("Failed to read from %s : %s", address, err.Error())
		return "", "", "", nerr.Translate(err).Add("failed to create regex")
	}
	ReturnInput = reg.FindAllStringSubmatch(fmt.Sprintf("%s", buf.Bytes()), -1)
	verdata := ReturnInput[0][1]
	log.L.Infof("Version: %s", verdata)

	//MacAddress
	reg, err = regexp.Compile("MAC Address = ([A-Z,0-9]{2}:[A-Z,0-9]{2}:[A-Z,0-9]{2}:[A-Z,0-9]{2}:[A-Z,0-9]{2}:[A-Z,0-9]{2})")
	if err != nil {
		log.L.Errorf("Failed to read from %s : %s", address, err.Error())
		return "", "", "", nerr.Translate(err).Add("failed to create regex")
	}
	ReturnInput = reg.FindAllStringSubmatch(fmt.Sprintf("%s", buf.Bytes()), -1)
	macaddr := ReturnInput[0][1]
	log.L.Infof("MAC: %s", macaddr)

	return ipaddr, verdata, macaddr, nil
}

//Get Link status of all inputs:
func GetActiveSignal(address string, port string) (structs.ActiveSignal, *nerr.E) {
	var toReturn structs.ActiveSignal
	conn, gerr := getConnection(address, true)
	if gerr != nil {
		log.L.Errorf("Failed to get connection with %s: %s", address, gerr.Error())
		return toReturn, nerr.Translate(gerr).Add("Telnet connection failed")
	}
	//close connection
	defer conn.Close()

	conn.Write([]byte(fmt.Sprintf("STA\r\n")))

	//capture response
	var buf bytes.Buffer
	io.Copy(&buf, conn)

	//regex black magic
	reg, err := regexp.Compile("Video Input ([0-9]{2})  : EDID = DEFAULT [0-9]{2},  LINK = ([A-Z]+)")
	if err != nil {
		log.L.Errorf("Failed to read from %s : %s", address, err.Error())
		return toReturn, nerr.Translate(err).Add("failed to create regex")
	}
	ReturnInput := reg.FindAllStringSubmatch(fmt.Sprintf("%s", buf.Bytes()), -1)

	//Loop through all inputs and find requested port status
	for i := 0; i < 4; i++ {
		log.L.Infof("testing this=%s with requested port %s", ReturnInput[i][1][1:], port)
		if port == ReturnInput[i][1][1:] {
			log.L.Info("They are equal! lets check status")
			if ReturnInput[i][2] == "ON" {
				toReturn = structs.ActiveSignal{
					Active: true,
				}
				return toReturn, nil
			} else {
				toReturn = structs.ActiveSignal{
					Active: false,
				}
				return toReturn, nil
			}
		}
	}
	return toReturn, nil
}
