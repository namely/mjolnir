package logger

import (
	"context"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var ctx0 = context.Background()
var L1 = logrus.WithFields(logrus.Fields{
	"test": 1,
})
var L2 = logrus.WithFields(logrus.Fields{
	"test": 2,
})

func TestFromContext(t *testing.T) {
	t.Run("given a context with no entry, it returns a entry", func(t *testing.T) {
		result := FromContext(ctx0)
		assert.IsType(t, &logrus.Entry{}, result)
	})

	t.Run("given a context with a entry, it returns that entry", func(t *testing.T) {
		ctx1 := context.WithValue(ctx0, Key, L1)
		ctx2 := FromContext(ctx1)
		assert.Equal(t, L1, ctx2)
	})
}

func TestSetEntry(t *testing.T) {

	t.Run("given a context with no entry, it returns a new context with the given entry", func(t *testing.T) {
		ctx1 := SetEntry(ctx0, L1)
		assert.NotEqual(t, ctx0, ctx1)
		l := ctx1.Value(Key)
		assert.NotNil(t, l)
		assert.IsType(t, &logrus.Entry{}, l)
		assert.Equal(t, L1, l.(*logrus.Entry))
	})

	t.Run("given a context with a entry, and the same entry, it returns the same context with the same entry", func(t *testing.T) {
		ctx1 := context.WithValue(ctx0, Key, L1)
		ctx2 := SetEntry(ctx0, L1)
		assert.Equal(t, ctx1, ctx2)
		l := ctx2.Value(Key)
		assert.NotNil(t, l)
		assert.IsType(t, &logrus.Entry{}, l)
		assert.Equal(t, L1, l.(*logrus.Entry))
	})

	t.Run("given a context with a entry, and a different entry, it returns a new context with the new entry", func(t *testing.T) {
		ctx1 := context.WithValue(ctx0, Key, L1)
		ctx2 := SetEntry(ctx0, L2)
		assert.NotEqual(t, ctx1, ctx2)
		l := ctx2.Value(Key)
		assert.NotNil(t, l)
		assert.IsType(t, &logrus.Entry{}, l)
		assert.Equal(t, L2, l.(*logrus.Entry))
	})
}
