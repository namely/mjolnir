package health

import (
	"fmt"
	"net"
)

const (
	// HealthMessage is the string that the TCP listener writes upon connections
	HealthMessage = `mjolnir/health: bow before me`
)

// StartTCP starts a health endpoint using the TCP protocol. The host should
// be in the format "localhost:1234" for the binding to take place correctly.
func StartTCP(host string) error {
	l, err := net.Listen("tcp", host)
	if err != nil {
		return fmt.Errorf("mjolnir/health: could not create TCP listener: %q", err.Error())
	}

	fmt.Println("listening?")

	defer l.Close()
	for {
		conn, err := l.Accept()
		fmt.Println("accepted")
		if err != nil {
			return fmt.Errorf("mjolnir/health: could not accept connection: %q", err.Error())
		}
		defer conn.Close()

		if _, err := conn.Write([]byte(HealthMessage)); err != nil {
			return fmt.Errorf("mjolnir/health: could not write response: %q", err.Error())
		}
		fmt.Println("write")
	}
}

func StartHTTP(host string) error {
	return nil
}
