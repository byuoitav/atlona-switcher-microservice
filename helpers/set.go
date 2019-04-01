package helpers

import (
	"fmt"
	"strings"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/pooled"
)

//SwitchInput takes the IP address, the output and the input from the user and
//switches the input to the one requested
func SwitchInput(address, ouput, input string) (string, error) {
	work := func(conn pooled.Conn) error {
		//execute telnet command to switch input
		conn.Write([]byte("x" + input + "AVx" + ouput + "\r\n"))

		b, err := readUntil(LF, conn, 10)
		if err != nil {
			return err
		}

		if strings.Contains(string(b), "FAILED") {
			return fmt.Errorf("input or output is out of range")
		}

		log.L.Infof("Set input to %s returned %s", input, b)
		return nil
	}

	err := pool.Do(address, work)
	if err != nil {
		return "", fmt.Errorf("failed to switch input: %s", err)
	}

	return input, nil
}
