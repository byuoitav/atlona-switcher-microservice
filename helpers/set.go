package helpers

import (
	"fmt"
	"strings"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
)

//SwitchInput takes the IP address, the output and the input from the user and
//switches the input to the one requested

func SwitchInput(address, ouput, input string) (string, *nerr.E) {
	//establish telnet connection to device
	conn, err := getConnection(address, true)

	if err != nil {
		log.L.Errorf("Failed to establish connection with %s : %s", address, err.Error())
		return "", nerr.Translate(err).Add("Telnet connection failed")
	}

	//execute telnet command to switch input
	conn.Write([]byte("x" + input + "AVx" + ouput + "\r\n"))
	b, err := readUntil(CARRIAGE_RETURN, conn, 10)
	if err != nil {
		return "", nerr.Translate(err).Add("failed to read from connection")
	}

	if strings.Contains(string(b), "FAILED") {
		return "", nerr.Create("Input or Output is out of range", "Error")
	}

	response := strings.Split(fmt.Sprintf("%s", b), "AV")
	test := strings.Split(fmt.Sprintf("%s", response), "x")
	log.L.Infof("test: %s", test)
	defer conn.Close()
	return fmt.Sprintf("%s", input), nil
}
