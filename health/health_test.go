package health

import (
	"bufio"
	"net"
	"testing"
	"time"
)

func TestTCPListener(t *testing.T) {
	host := "localhost:28347"

	var tcpErr error
	go func() {
		tcpErr = StartTCP(host)
		t.Logf("Hello world")
	}()

	// wait 100 ms to wait for TCP listener to get ready
	time.Sleep(100 * time.Millisecond)

	conn, err := net.DialTimeout("tcp", host, 100*time.Millisecond)
	if err != nil {
		t.Errorf("could not dial connection: %v", err)
	}
	time.Sleep(1000 * time.Millisecond)
	defer conn.Close()

	response, err := bufio.NewReader(conn).ReadString('#')
	if err != nil {
		t.Errorf("Error Reading from TCP:  %v", err)
	}
	if HealthMessage != response {
		t.Errorf("unexpected message received, got %q but wanted %q", response, HealthMessage)
	}

	if tcpErr != nil {
		t.Errorf("tcp error: %v", tcpErr)
	}
}
