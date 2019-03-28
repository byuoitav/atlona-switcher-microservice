package helpers

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/fatih/color"

	telnet "github.com/reiver/go-telnet"
)

const (
	CARRIAGE_RETURN           = 0x0D
	LINE_FEED                 = 0x0A
	SPACE                     = 0x20
	DELAY_BETWEEN_CONNECTIONS = time.Second * 10
)

var tlsConfig *tls.Config
var caller telnet.Caller

func init() {
	tlsConfig = &tls.Config{}
}

func readUntil(delimeter byte, conn net.Conn, timeoutInSeconds int) ([]byte, error) {
	conn.SetReadDeadline(time.Now().Add(time.Duration(int64(timeoutInSeconds)) * time.Second))

	reader := bufio.NewReader(conn)
	b, err := reader.ReadBytes(delimeter)

	if err != nil {
		err = errors.New(fmt.Sprintf("Error reading response: %s", err.Error()))
		log.L.Infof("%s", err.Error())
		return []byte{}, err
	}
	return b, nil
}

func getConnection(key interface{}) (net.Conn, error) {
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

	color.Set(color.FgMagenta)
	log.L.Infof("Reading welcome message")
	color.Unset()

	_, err = readUntil(CARRIAGE_RETURN, conn, 3)
	if err != nil {
		return conn, err
	}

	return conn, err
}
