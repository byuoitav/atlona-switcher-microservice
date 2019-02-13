package helpers

import (
	"fmt"
	"strings"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
)

//SwitchInput takes the IP address, the output and the input from the user and
//switches the input to the one requested

//TODO: clear the buffer before everything!!! implement other regex commands
func SwitchInput(address, ouput, input string) (string, *nerr.E) {
	//establish telnet connection to device
	conn, err := getConnection(address, true)
	var b []byte
	if err != nil {
		log.L.Errorf("Failed to establish connection with %s : %s", address, err.Error())
		return "", nerr.Translate(err).Add("Telnet connection failed")
	}

	//execute telnet command to switch input
	command := fmt.Sprintf("SPO0" + ouput + "SI0" + input)
	log.L.Infof("%s", command)
	conn.Write([]byte("\r\n"))
	conn.Write([]byte(command + "\r\n"))
	b, err = readUntil(CARRIAGE_RETURN, conn, 10)
	if err != nil {
		return "", nerr.Translate(err).Add("failed to read from connection")
	}

	if strings.Contains(string(b), "FAILED") {
		return "", nerr.Create("Input or Output is out of range", "Error")
	}
	log.L.Infof("response: %s", b)
	response := strings.Split(fmt.Sprintf("%s", b), "")
	test := strings.Split(fmt.Sprintf("%s", response), "x")
	log.L.Infof("test: %s", test)
	defer conn.Close()
	return fmt.Sprintf("%s", input), nil
}
