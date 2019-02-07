package helpers

import (
	"fmt"
	"net"
	"strings"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
)

//This function returns the current input that is being shown as the output
func GetOutput(address, output string) (string, *nerr.E) {
	conn, err := getConnection(address, true)
	if err != nil {
		log.L.Errorf("Failed to establish connection with %s : %s", address, err.Error())
		return "", nerr.Translate(err).Add("Telnet connection failed")
	}
	//close connection
	defer conn.Close()

	conn.Write([]byte(fmt.Sprintf("n %v\r\n", output)))
	b, err := readUntil(CARRIAGE_RETURN, conn, 3)
	if err != nil {
		return "", nerr.Translate(err).Add("failed to read from connection")
	}

	response := strings.Split(fmt.Sprintf("%s", b), "\r\n")

	log.L.Infof("response: '%s'", response)
	input := string(response[1])
	input = input[1:]

	log.L.Infof("input: '%s'", input)

	return fmt.Sprintf("%s", input), nil
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
	conn.Write([]byte("#show_ip\r\n"))
	b, err := readUntil(CARRIAGE_RETURN, conn, 3)
	if err != nil {
		return "", nerr.Translate(err).Add("failed to read from connection")
	}
	ipaddr := ""
	response := strings.Split(string(b), " : ")
	if len(response) >= 1 {
		ipaddr = strings.Replace(response[1], "telnet->", "", -1)
		ipaddr = strings.TrimSpace(ipaddr)
	}
	log.L.Info(ipaddr)
	return ipaddr, nil
}

//gets software and hardware data
func getVerData(address string, conn *net.TCPConn) (string, *nerr.E) {
	conn.Write([]byte("#show_ver_data\r\n"))
	b, err := readUntil(CARRIAGE_RETURN, conn, 3)
	if err != nil {
		return "", nerr.Translate(err).Add("failed to read from connection")
	}
	verdata := ""
	response := strings.Split(string(b), " : ")
	if len(response) >= 1 {
		verdata = strings.Replace(response[1], "telnet->", "", -1)
		verdata = strings.TrimSpace(verdata)
	}
	log.L.Info(verdata)
	return verdata, nil
}

//gets macaddress of device
func getMacAddress(address string, conn *net.TCPConn) (string, *nerr.E) {
	conn.Write([]byte("#show_mac_addr\r\n"))
	b, err := readUntil(CARRIAGE_RETURN, conn, 3)
	if err != nil {
		return "", nerr.Translate(err).Add("failed to read from connection")
	}
	macaddr := ""
	response := strings.Split(string(b), " : ")
	if len(response) >= 1 {
		macaddr = strings.Replace(response[1], "telnet->", "", -1)
		macaddr = strings.TrimSpace(macaddr)
	}
	log.L.Info(macaddr)
	return macaddr, nil
}
