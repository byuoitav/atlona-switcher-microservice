package helpers

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/byuoitav/common/log"

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

func readUntil(delimeter byte, conn net.Conn, timeoutInSeconds int) ([]byte, error) {
	conn.SetReadDeadline(time.Now().Add(time.Duration(int64(timeoutInSeconds)) * time.Second))

	reader := bufio.NewReader(conn)
	b, err := reader.ReadBytes(delimeter)
	if err != nil {
		err = fmt.Errorf("Error reading response: %s", err.Error())
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

	log.L.Infof("Reading welcome message")

	resp, err := readUntil(LF, conn, 3)
	if err != nil {
		return conn, err
	}
	fmt.Printf("resp: '0x%x'", resp)

	conn.Write([]byte("hi\n"))

	resp, err = readUntil(LF, conn, 10)
	if err != nil {
		return conn, err
	}

	fmt.Printf("resp: '0x%x'", resp)

	return conn, err
}
