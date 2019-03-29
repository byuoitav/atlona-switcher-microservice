package helpers

import (
	"fmt"
	"net"
	"strings"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
)

//SwitchInput takes the IP address, the output and the input from the user and
//switches the input to the one requested
func SwitchInput(address, ouput, input string) (string, *nerr.E) {
	work := func(conn net.Conn) error {
		//execute telnet command to switch input
		conn.Write([]byte("x" + input + "AVx" + ouput + "\r\n"))
		b, err := readUntil(LF, conn, 10)
		if err != nil {
			return err
		}

		if strings.Contains(string(b), "FAILED") {
			return fmt.Errorf("Input or Output is out of range")
		}

		log.L.Infof("Set input to %s returned %s", input, b)

		return nil
	}

	err := pool.Do(address, work)
	if err != nil {
		return "", nerr.Translate(err)
	}

	return input, nil
}
