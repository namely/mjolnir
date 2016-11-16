package interceptor

import (
	"bytes"
	"encoding/json"
	"testing"

	"google.golang.org/grpc"

	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// LogLine is only used to parse log lines from interceptors making it
// easier to test
type logLine struct {
	Msg       string `json:"msg"`
	RequestID string `json:"request_id"`
	Duration  string `json:"duration"`
}

func TestLoggerInterceptor(t *testing.T) {
	i := New()
	l := logrus.New()
	l.Formatter = &logrus.JSONFormatter{}
	buf := new(bytes.Buffer)
	l.Out = buf

	i.Use(Logger(l))
	middleware := i.Middleware()
	ctx := context.TODO()
	info := &grpc.UnaryServerInfo{FullMethod: "/package.service/method"}
	final := func(ctx context.Context, in interface{}) (interface{}, error) {
		l.Info("from handler")
		return nil, nil
	}

	middleware(ctx, nil, info, final)

	decoder := json.NewDecoder(buf)
	var firstLine logLine
	require.NoError(t, decoder.Decode(&firstLine))
	assert.Equal(t, "processing rpc", firstLine.Msg)

	var secondLine logLine
	require.NoError(t, decoder.Decode(&secondLine))
	assert.Equal(t, "from handler", secondLine.Msg)

	var thirdLine logLine
	require.NoError(t, decoder.Decode(&thirdLine))
	assert.Equal(t, "finished rpc", thirdLine.Msg)
}
