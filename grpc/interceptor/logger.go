package interceptor

import (
	"time"

	"golang.org/x/net/context"

	"github.com/namely/mjolnir/logger"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Logger returns an interceptor that will set on the context a *logrus.Entry
// that will automatically be tagged with the request_id UUIDv4.
// It will log the start and end of the request including the duration of
// the call as well.
func Logger(l *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = addLoggerToContext(ctx, l)
		entry := logger.FromContext(ctx)

		// regex to change /service.Service/Endpt -> Endpt
		name := logger.FormatServiceEndpoint.ReplaceAllString(info.FullMethod, "")

		start := time.Now()
		out, err := handler(ctx, req)
		if err != nil {
			if ferr, ok := err.(ErrorFielder); ok {
				fields := ferr.Fields()
				entry.WithError(ferr).WithFields(*fields).WithField(
					"core.duration", float64(time.Since(start))/float64(time.Millisecond),
				).Error("rpc endpoint " + name + " failed")
				return nil, ErrGrpcInternalError
			}

			entry.WithError(err).WithField(
				"core.duration", float64(time.Since(start))/float64(time.Millisecond),
			).Error("rpc endpoint " + name + " failed")
			return nil, err
		}

		entry.WithFields(logrus.Fields{
			"endpoint":      name,
			"core.duration": float64(time.Since(start)) / float64(time.Millisecond),
		}).Info("finished rpc")

		return out, err
	}
}

func addLoggerToContext(ctx context.Context, l *logrus.Logger) context.Context {
	id, err := uuid.NewV4()
	if err != nil {
		logrus.WithError(err).Fatal("could not generate request id")
	}

	entry := l.WithField("request_id", id.String())
	return logger.SetEntry(ctx, entry)
}
