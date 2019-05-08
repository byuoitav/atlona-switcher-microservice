package helpers

import (
	"fmt"
	"strings"
	"time"

	"github.com/byuoitav/common/pooled"
)

//SwitchInput takes the IP address, the output and the input from the user and
//switches the input to the one requested
func SwitchInput(address, output, input string) (string, error) {
	// validate that the current input doesn't match
	curin, _, err := GetOutput(address)
	if err != nil {
		return "", fmt.Errorf("unable to check if the current input matches the desired input: %s", err)
	}

	if strings.EqualFold(curin, input) {
		return input, nil // we don't need to change the input
	}

	work := func(conn pooled.Conn) error {
		conn.Log().Infof("Setting output %s to %s", output, input)

		cmd := []byte(fmt.Sprintf("x%vAVx%v\r", input, output))

		// send command
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

		if strings.Contains(string(b), "FAILED") {
			return fmt.Errorf("input or output is out of range")
		}

		conn.Log().Infof("Set input to %s successful", input)
		return nil
	}

	err = pool.Do(address, work)
	if err != nil {
		return "", fmt.Errorf("failed to switch input: %s", err)
	}

	return input, nil
}
