package interceptor

import (
	"errors"
	"regexp"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

// regex to pull out the rpc name from the FullMethod tag on the call info
var methodRegex = regexp.MustCompile(`^\/[\w\.]+\/([\w]+)`)

// StatsdClient is an interface to allow passing in any statsd-esque client.
// Since we rely on the datadog one usually, this is where this is stolen from.
// See: https://godoc.org/github.com/DataDog/datadog-go/statsd
type StatsdClient interface {
	Count(name string, value int64, tags []string, rate float64) error
	Decr(name string, tags []string, rate float64) error
	Gauge(name string, value float64, tags []string, rate float64) error
	Histogram(name string, value float64, tags []string, rate float64) error
	Incr(name string, tags []string, rate float64) error
	Set(name string, value string, tags []string, rate float64) error
	SimpleEvent(title, text string) error
	TimeInMilliseconds(name string, value float64, tags []string, rate float64) error
	Timing(name string, value time.Duration, tags []string, rate float64) error
}

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
func Timer(prefix string, c StatsdClient) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		groups := methodRegex.FindStringSubmatch(info.FullMethod)
		if groups[1] == "" {
			return nil, errors.New("interceptor/timing: could not find method in grpc info")
		}

		key := prefix + "." + groups[1]
		c.Incr(key, nil, 1)
		start := time.Now()
		defer c.Timing(key+".duration", time.Since(start), nil, 0)
		return handler(ctx, req)
	}
}
