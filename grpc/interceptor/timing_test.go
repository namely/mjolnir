package interceptor

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/armon/go-metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

// MockSink is stolen from https://github.com/armon/go-metrics/blob/master/sink_test.go#L9
type MockSink struct {
	keys   [][]string
	vals   []float32
	labels [][]metrics.Label
}

func (m *MockSink) SetGauge(key []string, val float32) {
	m.SetGaugeWithLabels(key, val, nil)
}
func (m *MockSink) SetGaugeWithLabels(key []string, val float32, labels []metrics.Label) {
	m.keys = append(m.keys, key)
	m.vals = append(m.vals, val)
	m.labels = append(m.labels, labels)
}
func (m *MockSink) EmitKey(key []string, val float32) {
	m.keys = append(m.keys, key)
	m.vals = append(m.vals, val)
	m.labels = append(m.labels, nil)
}
func (m *MockSink) IncrCounter(key []string, val float32) {
	m.IncrCounterWithLabels(key, val, nil)
}
func (m *MockSink) IncrCounterWithLabels(key []string, val float32, labels []metrics.Label) {
	m.keys = append(m.keys, key)
	m.vals = append(m.vals, val)
	m.labels = append(m.labels, labels)
}
func (m *MockSink) AddSample(key []string, val float32) {
	m.AddSampleWithLabels(key, val, nil)
}
func (m *MockSink) AddSampleWithLabels(key []string, val float32, labels []metrics.Label) {
	m.keys = append(m.keys, key)
	m.vals = append(m.vals, val)
	m.labels = append(m.labels, labels)
}

func TestTimingInterceptor(t *testing.T) {
	i := New()
	sink := &MockSink{}
	stats, err := metrics.New(metrics.DefaultConfig("timing-test"), sink)
	require.NoError(t, err)

	i.Use(Timer(stats))
	middleware := i.Middleware()
	ctx := context.TODO()
	info := &grpc.UnaryServerInfo{FullMethod: "/package.Service/SuperDopeRPCMethod"}
	final := func(ctx context.Context, in interface{}) (interface{}, error) {
		return nil, nil
	}

	middleware(ctx, nil, info, final)

	assert.Equal(t, "timing-test", sink.keys[0][0])
	assert.Equal(t, "SuperDopeRPCMethod", sink.keys[0][1])

	assert.Equal(t, "timing-test", sink.keys[1][0])
	assert.Equal(t, "SuperDopeRPCMethod", sink.keys[1][1])
}
