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

	conn, err := net.DialTimeout("tcp", host, 100*time.Millisecond)
	if err != nil {
		t.Errorf("could not dial connection: %v", err)
	}
	time.Sleep(1000 * time.Millisecond)

	//response, err := bufio.NewReader(conn).ReadString('\n')
	var b []byte
	response, err := bufio.NewReader(conn).Read(b)

	if HealthMessage != string(b) {
		t.Errorf("unexpected message received, got %q but wanted %q", response, HealthMessage)
	}

	if tcpErr != nil {
		t.Errorf("tcp error: %v", tcpErr)
	}
}
