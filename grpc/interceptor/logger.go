package interceptor

import (
	"time"

	"golang.org/x/net/context"

	"../../data"
	"github.com/Sirupsen/logrus"
	"github.com/namely/mjolnir/logger"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"regexp"
)

// Logger returns an interceptor that will set on the context a *logrus.Entry
// that will automatically be tagged with the request_id UUIDv4.
// It will log the start and end of the request including the duration of
// the call as well.
func Logger(l *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = addLoggerToContext(l, ctx)
		entry := logger.FromContext(ctx)

		// regex to change /service.Service/Endpt -> Endpt
		r := regexp.MustCompile(`\/([A-Za-z])\w+\.([A-Za-z])\w+\/`)
		name := r.ReplaceAllString(info.FullMethod, "")

		entry.WithField("endpoint", name).Info("processing rpc")

		start := time.Now()
		out, err := handler(ctx, req)
		if err != nil {
			if ferr, ok := err.(data.ErrorFielder); ok {
				fields := ferr.Fields()
				entry.WithError(ferr).WithFields(*fields).WithField(
					"duration", time.Since(start).String(),
				).Error("rpc endpoint " + name + " failed")
			} else {
				entry.WithError(err).WithField(
					"duration", time.Since(start).String(),
				).Error("rpc endpoint " + name + " failed")
			}
			return nil, data.ErrGrpcInternalError
		}

		entry.WithFields(logrus.Fields{
			"endpoint":    name,
			"pb_response": out,
			"duration":    time.Since(start).String(),
		}).Info("finished rpc")

		return out, err
	}
}

func addLoggerToContext(l *logrus.Logger, ctx context.Context) context.Context {
	entry := l.WithField("request_id", uuid.NewV4().String())
	return logger.SetEntry(ctx, entry)
}
