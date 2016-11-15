package interceptor

import (
	"time"

	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
	"github.com/namely/mjolnir/logger"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
)

// Interceptor makes it easy to return unary interceptors for grpc servers
type Interceptor struct {
	logger *logrus.Logger
}

// New initializes and returns an Interceptor
func New(logger *logrus.Logger) *Interceptor {
	return &Interceptor{logger: logger}
}

func (i *Interceptor) addLoggerToContext(ctx context.Context) context.Context {
	entry := i.logger.WithField("request_id", uuid.NewV4().String())

	return logger.SetEntry(ctx, entry)
}

// Middleware returns a grpc.UnaryServerHandler compatible function
// documented here https://github.com/grpc/grpc-go/tree/master/interceptor.go#L73
// It will handle intercepting errors and unwrapping the original error for logging.
//
// A grpc endpoint should return generic errors such as "company not found" by wrapping
// the original error.
func (i *Interceptor) Middleware() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = i.addLoggerToContext(ctx)
		entry := logger.FromContext(ctx)
		name := info.FullMethod

		entry.WithField("endpoint", name).Info("processing rpc")

		start := time.Now()
		out, err := handler(ctx, req)
		if err != nil {
			entry.WithError(err).Error("rpc endpoint failed")
			return nil, err
		}

		entry.WithFields(logrus.Fields{
			"endpoint": name,
			"duration": time.Since(start).String(),
		}).Info("finished rpc")

		return out, err
	}
}
