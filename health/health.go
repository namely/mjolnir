package health

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

const (
	// HealthMessage is the string that the TCP listener writes upon connections
	HealthMessage = `mjolnir/health: bow before me#`
)

// StartTCP starts a health endpoint using the TCP protocol. The host should
// be in the format "localhost:1234" for the binding to take place correctly.
func StartTCP(host string) error {
	l, err := net.Listen("tcp", host)
	if err != nil {
		return fmt.Errorf("mjolnir/health: could not create TCP listener: %q", err.Error())
	}

	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("mjolnir/health: could not accept connection: %q", err.Error())
		}
		go func(conn net.Conn) {
			defer conn.Close()
			if _, err := io.Copy(conn, bytes.NewBufferString(HealthMessage)); err != nil {
				fmt.Errorf("mjolnir/health: Could not write message: %q", err)
			}
		}(conn)
	}
}

func StartHTTP(host string) error {
	return nil
}
