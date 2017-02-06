package interceptor

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

func TestTimingInterceptor(t *testing.T) {
	i := New()
	stats := &MockStatsdClient{}

	stats.On("Incr", "test.SuperDopeRPCMethod", []string(nil), float64(1)).Return(nil)
	stats.On("Timing", "test.SuperDopeRPCMethod.duration", mock.AnythingOfType("time.Duration"), []string(nil), float64(0)).Return(nil)

	i.Use(Timer("test", stats))
	middleware := i.Middleware()
	ctx := context.TODO()
	info := &grpc.UnaryServerInfo{FullMethod: "/package.Service/SuperDopeRPCMethod"}
	final := func(ctx context.Context, in interface{}) (interface{}, error) {
		return nil, nil
	}

	middleware(ctx, nil, info, final)

	stats.AssertExpectations(t)
}
