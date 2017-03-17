package interceptor

import (
	"errors"
	"regexp"
	"time"

	"golang.org/x/net/context"

	"github.com/armon/go-metrics"
	"google.golang.org/grpc"
)

// regex to pull out the rpc name from the FullMethod tag on the call info
var methodRegex = regexp.MustCompile(`^\/[\w\.]+\/([\w]+)`)

// Timer is a grpc middleware for handling timing how long endpoints take to
// handle requests and send them to the client given.
// Emits 2 metrics to the given statsd client:
//
// * A counter for the endpoint so you can see how many times it has been called
// * A timer to see how long the endpoint takes.
//
// So given a prefix of "production.myapp", and an RPC endpoint called "CreateThing", you'd have 2 keys emitted to statsd:
//
// * INCR production.myapp.CreateThing
// * TIMING production.myapp.CreateThing.duration
func Timer(m *metrics.Metrics) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		groups := methodRegex.FindStringSubmatch(info.FullMethod)
		if groups[1] == "" {
			return nil, errors.New("interceptor/timing: could not find method in grpc info")
		}

		key := []string{groups[1]}
		m.IncrCounter(key, 1)
		start := time.Now()
		defer m.MeasureSince(key, start)
    
		return handler(ctx, req)
	}
}
