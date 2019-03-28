package helpers

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/common/pooled"
)

var pool = pooled.NewMap(30*time.Second, getConnection)

var responseRE = `x(\d)AVx(\d)`
var re *regexp.Regexp

func init() {
	re = regexp.MustCompile(responseRE)
}

//GetOutput This function returns the current input that is being shown as the output
func GetOutput(address string) (string, string, *nerr.E) {
	var input, output string

	work := func(conn net.Conn) error {
		conn.Write([]byte(fmt.Sprintf("Status\r\n")))
		b, err := readUntil(CARRIAGE_RETURN, conn, 10)
		if err != nil {
			return nerr.Translate(err).Add("failed to read from connection")
		}

		log.L.Infof("Get status returned %s", b)

		match := re.FindStringSubmatch(string(b))

		if len(match) == 0 {
			log.L.Errorf("Invalid status response returned")
			return fmt.Errorf("Invalid status returned")
		}

		input = match[1]
		output = match[2]

		log.L.Infof("Parsed response of %s, input %s, Output %s", b, input, output)

		return nil
	}

	err := pool.Do(address, work)
	if err != nil {
		return "", "", nerr.Translate(err)
	}

	return input, output, nil
}

//GetHardware This function gets the IP Address (ipaddr), Software and hardware
//version (verdata), and mac address (macaddr) of the device
func GetHardware(address string) (string, string, string, *nerr.E) {
	var ipaddr, verdata, macaddr string

	work := func(conn net.Conn) error {
		var err error
		ipaddr, err = getIPAddress(address, conn)
		if err != nil {
			log.L.Errorf("Failed to establish connection with %s : %s", address, err.Error())
			return err
		}

		verdata, err = getVerData(address, conn)
		if err != nil {
			log.L.Errorf("Failed to establish connection with %s : %s", address, err.Error())
			return err
		}

		macaddr, err = getMacAddress(address, conn)
		if err != nil {
			log.L.Errorf("Failed to establish connection with %s : %s", address, err.Error())
			return err
		}

		return nil
	}

	err := pool.Do(address, work)
	if err != nil {
		return "", "", "", nerr.Translate(err)
	}

	return ipaddr, macaddr, verdata, nil
}

func getIPAddress(address string, conn net.Conn) (string, *nerr.E) {
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
func getVerData(address string, conn net.Conn) (string, *nerr.E) {
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
func getMacAddress(address string, conn net.Conn) (string, *nerr.E) {
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
