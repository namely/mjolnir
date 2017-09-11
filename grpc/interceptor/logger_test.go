package interceptor

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/namely/mjolnir/logger"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// LogLine is only used to parse log lines from interceptors making it
// easier to test
type logLine struct {
	Msg       string      `json:"msg"`
	RequestID string      `json:"request_id"`
	Duration  json.Number `json:"core.duration,Number"`
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
		e := logger.FromContext(ctx)
		e.Info("from handler")
		return nil, nil
	}

	middleware(ctx, nil, info, final)

	decoder := json.NewDecoder(buf)
	var firstLine logLine
	require.NoError(t, decoder.Decode(&firstLine))
	assert.Equal(t, "processing rpc", firstLine.Msg)
	assert.NotEmpty(t, firstLine.RequestID)

	var secondLine logLine
	require.NoError(t, decoder.Decode(&secondLine))
	assert.Equal(t, "from handler", secondLine.Msg)
	assert.NotEmpty(t, secondLine.RequestID)

	var thirdLine logLine
	require.NoError(t, decoder.Decode(&thirdLine))
	assert.Equal(t, "finished rpc", thirdLine.Msg)
	assert.NotEmpty(t, thirdLine.RequestID)
	assert.NotEmpty(t, thirdLine.Duration)
}
