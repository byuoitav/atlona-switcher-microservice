package helpers

import (
	"fmt"
	"net"
	"time"

	"github.com/byuoitav/common/pooled"
)

const (
	// CR is a carriage return
	CR = '\r'
	// LF is a line feed
	LF = '\n'
)

func getConnection(key interface{}) (pooled.Conn, error) {
	address, ok := key.(string)
	if !ok {
		return nil, fmt.Errorf("key must be a string")
	}

	conn, err := net.DialTimeout("tcp", address+":23", 10*time.Second)
	if err != nil {
		return nil, err
	}

	pconn := pooled.Wrap(conn)

	// read first new line
	_, err = pconn.ReadUntil(LF, 3*time.Second)
	if err != nil {
		return nil, err
	}

	time.Sleep(1 * time.Second) // time for the switcher to chill out
	return pconn, nil
}
