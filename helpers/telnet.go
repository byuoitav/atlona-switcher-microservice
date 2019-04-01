package helpers

import (
	"fmt"
	"net"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/pooled"
)

const (
	// CR is a carriage return
	CR = '\r'
	// LF is a line feed
	LF = '\n'
)

func readUntil(delimeter byte, conn pooled.Conn, timeoutInSeconds int) ([]byte, error) {
	conn.SetReadDeadline(time.Now().Add(time.Duration(int64(timeoutInSeconds)) * time.Second))

	b, err := conn.ReadWriter().ReadBytes(delimeter)
	if err != nil {
		return []byte{}, err
	}

	return b, nil
}

func getConnection(key interface{}) (pooled.Conn, error) {
	address, ok := key.(string)
	if !ok {
		return nil, fmt.Errorf("key must be a string")
	}

	addr, err := net.ResolveTCPAddr("tcp", address+":23")
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}

	pconn := pooled.Wrap(conn)
	log.L.Infof("Reading welcome message")

	// read first new line
	_, err = readUntil(LF, pconn, 3)
	if err != nil {
		return nil, err
	}

	// read welcome to telnet message
	_, err = readUntil(LF, pconn, 3)
	if err != nil {
		return nil, err
	}

	time.Sleep(750 * time.Millisecond) // time for the switcher to chill out

	/*
		// don't broadcast our messages
		_, err = pconn.Write([]byte("Broadcast off"))
		if err != nil {
			return nil, err
		}

		log.L.Infof("Sent 'Broadcast off' message")
		b, err := readUntil(LF, pconn, 20)
		if err != nil {
			return nil, err
		}
		log.L.Infof("Received '%s' message", string(b))

		resp := strings.TrimSpace(string(b))
		if !strings.EqualFold(resp, "Broadcast off") {
			return nil, fmt.Errorf("unable to disable broadcasting messages")
		}
	*/

	return pconn, nil
}
