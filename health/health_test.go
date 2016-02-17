package health

import (
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

	var buf []byte
	if _, err := conn.Read(buf); err != nil {
		t.Errorf("could not read from connection: %v", err)
	}

	if exp, got := HealthMessage, string(buf); exp != got {
		t.Errorf("unexpected message received, got %q but wanted %q", got, exp)
	}

	if tcpErr != nil {
		t.Errorf("tcp error: %v", tcpErr)
	}
}
