package logger

import (
	"context"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var ctx0 = context.Background()

func TestFromContext(t *testing.T) {
	t.Run("given a context with no logger it returns a logger", func(t *testing.T) {
		ctx, _ := context.WithTimeout(ctx0, 1*time.Second)
		result := FromContext(ctx)
		assert.IsType(t, &logrus.Entry{}, result)
	})

	t.Run("given a context with a logger it returns that logger", func(t *testing.T) {
		L1 := logrus.WithFields(logrus.Fields{
			"test": 1,
		})
		ctx := context.WithValue(ctx0, Key, L1)
		result := FromContext(ctx)
		assert.Equal(t, L1, result)
	})
}

func TestSetEntry(t *testing.T) {
	L2 := logrus.WithFields(logrus.Fields{
		"test": 2,
	})

	t.Run("given a context with no logger, it returns a new context with the given logger", func(t *testing.T) {
		ctx := SetEntry(ctx0, L2)
		assert.NotEqual(t, ctx0, ctx)
		l := ctx.Value(Key)
		assert.NotNil(t, l)
		assert.IsType(t, &logrus.Entry{}, l)
		assert.Equal(t, L2, l.(*logrus.Entry))
	})

	t.Run("given a context with a logger, it returns the same context", func(t *testing.T) {
		ctx1 := context.WithValue(ctx0, Key, L2)
		ctx2 := SetEntry(ctx0, L2)
		assert.Equal(t, ctx1, ctx2)
		l := ctx2.Value(Key)
		assert.NotNil(t, l)
		assert.IsType(t, &logrus.Entry{}, l)
		assert.Equal(t, L2, l.(*logrus.Entry))
	})
}
