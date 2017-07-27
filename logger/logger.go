// Package logger provides functions to retrieve and set loggers on a context
package logger

import (
	"io"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"regexp"
)

var (
	// Key is used for contexts to retrieve and store loggers
	Key = "logger"

	// LogFormatter is the logrus formatting function
	// It defaults to logrus.JSONFormatter.
	Formatter logrus.Formatter = new(logrus.JSONFormatter)

	// LogWriter is where the logs go. By default this empty buffer blackholes
	// logs. It must be an io.Writer such as os.Stdout
	Writer io.Writer = ioutil.Discard

	// regex to format grpc call logs
	// "/service.Service/Endpt" -> "Endpt"
	FormatServiceEndpoint = regexp.MustCompile(`\/([A-Za-z])\w+\.([A-Za-z])\w+\/`)
)

// FromContext takes a context and returns a pointer to a logrus.Entry
// If there is an existing entry on ctx.LoggerKey, that is returned
// Otherwise, a new logger is created and returned that blackholes logs
func FromContext(ctx context.Context) *logrus.Entry {
	entry, ok := ctx.Value(Key).(*logrus.Entry)
	if !ok {

		return logrus.NewEntry(&logrus.Logger{
			Formatter: Formatter,
			Out:       Writer,
			Level:     logrus.InfoLevel,
		})
	}

	return entry
}

// SetEntry replaces ctx.LoggerKey with the given log entry
func SetEntry(ctx context.Context, e *logrus.Entry) context.Context {
	return context.WithValue(ctx, Key, e)
}
