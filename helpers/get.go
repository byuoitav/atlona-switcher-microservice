package helpers

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/byuoitav/common/pooled"
)

var pool = pooled.NewMap(100*time.Second, getConnection)

var responseRE = `x(\d)AVx(\d)`
var re *regexp.Regexp

func init() {
	re = regexp.MustCompile(responseRE)
}

//GetOutput This function returns the current input that is being shown as the output
func GetOutput(address string) (string, string, error) {
	var input, output string

	work := func(conn pooled.Conn) error {
		conn.Log().Infof("Getting the current output")

		cmd := []byte(fmt.Sprintf("Status\r"))
		n, err := conn.Write(cmd)
		switch {
		case err != nil:
			return err
		case n != len(cmd):
			return fmt.Errorf("wrote %v/%v bytes of command 0x%x", n, len(cmd), cmd)
		}

		b, err := conn.ReadUntil(LF, 10*time.Second)
		if err != nil {
			return err
		}

		conn.Log().Debugf("Response from command: 0x%x", b)

		responseStr := strings.TrimSpace(string(b))
		match := re.FindStringSubmatch(responseStr)
		if len(match) == 0 {
			return fmt.Errorf("Invalid status returned (got: 0x%x)", b)
		}

		input = match[1]
		output = match[2]

		conn.Log().Infof("Current output (port %s) is %s", output, input)
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
		conn.Log().Infof("Getting hardware info")
		var err error

		verdata, err = getVerData(address, conn)
		if err != nil {
			return err
		}

		macaddr, err = getMacAddress(address, conn)
		if err != nil {
			return err
		}

		ipaddr, err = getIPAddress(address, conn)
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
	conn.Write([]byte("IPCFG\r"))

	b, err := conn.ReadUntil(LF, 10*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to get ip address: %s", err)
	}

	resp := strings.TrimSpace(string(b))
	split := strings.Split(resp, ":")

	if len(split) != 2 {
		return "", fmt.Errorf("invalid response getting ip address. response: 0x%x", b)
	}

	conn.Log().Infof("IP Address: %s", split[1])
	return strings.TrimSpace(split[1]), nil
}

//gets software and hardware data
func getVerData(address string, conn pooled.Conn) (string, error) {
	conn.Write([]byte("Version\r"))

	b, err := conn.ReadUntil(LF, 10*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to get ver data: %s", err)
	}

	verdata := fmt.Sprintf("%s", b)
	verdata = strings.TrimSpace(verdata)

	conn.Log().Infof("Version: %s", verdata)
	return verdata, nil
}

//gets macaddress of device
func getMacAddress(address string, conn pooled.Conn) (string, error) {
	conn.Write([]byte("RAtlMac\r"))

	b, err := conn.ReadUntil(LF, 10*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to get mac address: %s", err)
	}

	macaddr := fmt.Sprintf("%s", b)
	macaddr = strings.TrimSpace(macaddr)

	conn.Log().Infof("Macaddress: %s", macaddr)
	return macaddr, nil
}
