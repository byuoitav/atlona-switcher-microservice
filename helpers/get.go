package helpers

import (
	"fmt"
	"net"
	"strings"

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

	conn.Write([]byte(fmt.Sprintf("Status\r\n")))
	b, err := readUntil(CARRIAGE_RETURN, conn, 10)
	if err != nil {
		return "", "", nerr.Translate(err).Add("failed to read from connection")
	}

	response := strings.Split(fmt.Sprintf("%s", b), "AV")

	log.L.Infof("response: '%s'", response[0])
	log.L.Infof("response: '%s'", response[1])
	input := string(response[0])
	input = input[len(input)-1:]
	output := string(response[1])
	output = output[1:]

	log.L.Infof("input: '%s'", input)
	log.L.Infof("output: '%s'", output)

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

	ipaddr, err := getIPAddress(address, conn)
	if err != nil {
		log.L.Errorf("Failed to establish connection with %s : %s", address, err.Error())
		return "", "", "", err.Add("Telnet connection failed")
	}

	verdata, err := getVerData(address, conn)
	if err != nil {
		log.L.Errorf("Failed to establish connection with %s : %s", address, err.Error())
		return "", "", "", err.Add("Telnet connection failed")
	}

	macaddr, err := getMacAddress(address, conn)
	if err != nil {
		log.L.Errorf("Failed to establish connection with %s : %s", address, err.Error())
		return "", "", "", err.Add("Telnet connection failed")
	}

	conn.Close()
	return ipaddr, macaddr, verdata, nil
}

func getIPAddress(address string, conn *net.TCPConn) (string, *nerr.E) {
	conn.Write([]byte("IPCFG\r\n"))
	b, err := readUntil(LINE_FEED, conn, 10)
	if err != nil {
		return "", nerr.Translate(err).Add("failed to read IP address from connection")
	}
	response := strings.Split(string(b), "IP Addr: ")
	ipaddr := strings.Split(response[1], "Netmask")
	ipaddr[0] = strings.TrimSpace(ipaddr[0])
	log.L.Infof("IP address: %s", ipaddr[0])
	return ipaddr[0], nil
}

//gets software and hardware data
func getVerData(address string, conn *net.TCPConn) (string, *nerr.E) {
	conn.Write([]byte("Version\r\n"))
	log.L.Info("Just wrote the command for version")
	b, err := readUntil(LINE_FEED, conn, 10)
	if err != nil {
		return "", nerr.Translate(err).Add("failed to read VerData from connection")
	}
	verdata := fmt.Sprintf("%s", b)
	verdata = strings.TrimSpace(verdata)

	log.L.Infof("version: %s", verdata)
	return verdata, nil
}

//gets macaddress of device
func getMacAddress(address string, conn *net.TCPConn) (string, *nerr.E) {
	conn.Write([]byte("RAtlMac\r\n"))
	b, err := readUntil(CARRIAGE_RETURN, conn, 10)
	if err != nil {
		return "", nerr.Translate(err).Add("failed to read Mac Address from connection")
	}
	macaddr := fmt.Sprintf("%s", b)
	macaddr = strings.TrimSpace(macaddr)

	log.L.Infof("macaddress: %s", macaddr)
	return macaddr, nil
}

//check to make sure that the device is awake (is no longer needed,
//but will keep in case needed for the future)
// func getPowerStatus(conn *net.TCPConn) string {
// 	conn.Write([]byte("PWSTA\r\n"))
// 	b, _ := readUntil(CARRIAGE_RETURN, conn, 10)
// 	status := ""
// 	log.L.Infof("Response: %s", b)
// 	test := strings.Split(fmt.Sprintf("%s", b), "PW")
// 	log.L.Infof("Split: %s", test[1])

// 	if test[1] != "ON" {
// 		conn.Write([]byte("PWON\r\n"))
// 		status = "ON"
// 	} else {
// 		status = "ON"
// 	}
// 	return status
// }
