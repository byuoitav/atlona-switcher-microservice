package helpers

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/pooled"
)

var pool = pooled.NewMap(30*time.Second, getConnection)

var responseRE = `x(\d)AVx(\d)`
var re *regexp.Regexp

func init() {
	re = regexp.MustCompile(responseRE)
}

//GetOutput This function returns the current input that is being shown as the output
func GetOutput(address string) (string, string, error) {
	var input, output string

	work := func(conn pooled.Conn) error {
		conn.Write([]byte(fmt.Sprintf("Status\r")))

		b, err := readUntil(LF, conn, 10)
		if err != nil {
			return err
		}

		responseStr := strings.TrimSpace(string(b))
		log.L.Infof("Get status returned %s", responseStr)

		match := re.FindStringSubmatch(responseStr)
		if len(match) == 0 {
			return fmt.Errorf("Invalid status returned (got: 0x%x)", b)
		}

		input = match[1]
		output = match[2]

		log.L.Infof("Parsed response of %s, input %s, Output %s", responseStr, input, output)
		return nil
	}

	err := pool.Do(address, work)
	if err != nil {
		return "", "", fmt.Errorf("failed to get output: %s", err)
	}

	return input, output, nil
}

//GetHardware This function gets the IP Address (ipaddr), Software and hardware
//version (verdata), and mac address (macaddr) of the device
func GetHardware(address string) (string, string, string, error) {
	var ipaddr, verdata, macaddr string

	work := func(conn pooled.Conn) error {
		var err error

		ipaddr, err = getIPAddress(address, conn)
		if err != nil {
			return err
		}

		verdata, err = getVerData(address, conn)
		if err != nil {
			return err
		}

		macaddr, err = getMacAddress(address, conn)
		if err != nil {
			return err
		}

		return nil
	}

	err := pool.Do(address, work)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get hardware info: %s", err)
	}

	return ipaddr, macaddr, verdata, nil
}

func getIPAddress(address string, conn pooled.Conn) (string, error) {
	conn.Write([]byte("IPCFG\r\n"))

	b, err := readUntil(LF, conn, 10)
	if err != nil {
		return "", fmt.Errorf("failed to get ip address: %s", err)
	}

	response := strings.Split(string(b), "IP Addr: ")
	ipaddr := strings.Split(response[1], "Netmask")
	ipaddr[0] = strings.TrimSpace(ipaddr[0])
	log.L.Infof("IP address: %s", ipaddr[0])

	return ipaddr[0], nil
}

//gets software and hardware data
func getVerData(address string, conn pooled.Conn) (string, error) {
	conn.Write([]byte("Version\r\n"))

	b, err := readUntil(LF, conn, 10)
	if err != nil {
		return "", fmt.Errorf("failed to get ver data: %s", err)
	}

	verdata := fmt.Sprintf("%s", b)
	verdata = strings.TrimSpace(verdata)

	log.L.Infof("version: %s", verdata)
	return verdata, nil
}

//gets macaddress of device
func getMacAddress(address string, conn pooled.Conn) (string, error) {
	conn.Write([]byte("RAtlMac\r\n"))

	b, err := readUntil(LF, conn, 10)
	if err != nil {
		return "", fmt.Errorf("failed to get mac address: %s", err)
	}

	macaddr := fmt.Sprintf("%s", b)
	macaddr = strings.TrimSpace(macaddr)

	log.L.Infof("macaddress: %s", macaddr)
	return macaddr, nil
}
