package interceptor

import (
	"golang.org/x/net/context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Interceptor makes it easy to return unary interceptors for grpc servers
type Interceptor struct {
	logger *logrus.Logger
	chain  []grpc.UnaryServerInterceptor
}

// New initializes and returns an Interceptor
func New() *Interceptor {
	return &Interceptor{}
}

// Use adds a middleware to the chain for the server interceptor
func (i *Interceptor) Use(m grpc.UnaryServerInterceptor) {
	i.chain = append(i.chain, m)
}

// Middleware returns a grpc.UnaryServerHandler compatible function
// documented here https://github.com/grpc/grpc-go/tree/master/interceptor.go#L73
// It will handle intercepting errors and unwrapping the original error for logging.
//
// A grpc endpoint should return generic errors such as "company not found" by wrapping
// the original error.
// This chaining logic shamelessly stolen from https://github.com/mwitkow/go-grpc-middleware/blob/master/chain.go
func (i *Interceptor) Middleware() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		buildChain := func(current grpc.UnaryServerInterceptor, next grpc.UnaryHandler) grpc.UnaryHandler {
			return func(currentCtx context.Context, currentReq interface{}) (interface{}, error) {
				return current(currentCtx, currentReq, info, next)
			}
		}

		chain := handler

		for ii := len(i.chain) - 1; ii >= 0; ii-- {
			chain = buildChain(i.chain[ii], chain)
		}

		return chain(ctx, req)
	}
}
