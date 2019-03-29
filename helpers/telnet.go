package helpers

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/pooled"

	telnet "github.com/reiver/go-telnet"
)

const (
	// CR is a carriage return
	CR = '\r'
	// LF is a line feed
	LF = '\n'
)

var tlsConfig *tls.Config
var caller telnet.Caller

func init() {
	tlsConfig = &tls.Config{}
}

func readUntil(delimeter byte, conn pooled.Conn, timeoutInSeconds int) ([]byte, error) {
	conn.SetReadDeadline(time.Now().Add(time.Duration(int64(timeoutInSeconds)) * time.Second))

	b, err := conn.ReadWriter().ReadBytes(delimeter)
	if err != nil {
		err = fmt.Errorf("Error reading response: %s", err.Error())
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

	return pconn, err
}
