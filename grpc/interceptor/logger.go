package interceptor

import (
	"time"

	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
	"github.com/namely/mjolnir/logger"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
)

// Logger returns an interceptor that will set on the context a *logrus.Entry
// that will automatically be tagged with the request_id UUIDv4.
// It will log the start and end of the request including the duration of
// the call as well.
func Logger(l *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = addLoggerToContext(l, ctx)
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

func addLoggerToContext(l *logrus.Logger, ctx context.Context) context.Context {
	entry := l.WithField("request_id", uuid.NewV4().String())
	return logger.SetEntry(ctx, entry)
}
