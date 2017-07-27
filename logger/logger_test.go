package logger

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var ctx0 = context.Background()
var testEntry1 = logrus.WithFields(logrus.Fields{
	"test": 1,
})
var testEntry2 = logrus.WithFields(logrus.Fields{
	"test": 2,
})

func TestFromContext(t *testing.T) {
	t.Run("given a context with no entry, it returns an entry", func(t *testing.T) {
		result := FromContext(ctx0)
		assert.IsType(t, &logrus.Entry{}, result)
	})

	t.Run("given a context with an entry, it returns that entry", func(t *testing.T) {
		ctx1 := context.WithValue(ctx0, Key, testEntry1)
		ctx2 := FromContext(ctx1)
		assert.Equal(t, testEntry1, ctx2)
	})
}

func TestSetEntry(t *testing.T) {
	t.Run("it returns a new context with the given entry...", func(t *testing.T) {
		tcs := []struct {
			scenario string
			oldCtx   context.Context
		}{
			{
				"given a context with no entry",
				ctx0,
			},
			{
				"given a context with the same entry",
				context.WithValue(ctx0, Key, testEntry1),
			},
			{
				"given a context with a different entry",
				context.WithValue(ctx0, Key, testEntry2),
			},
		}
		for _, tc := range tcs {
			t.Run(tc.scenario, func(t *testing.T) {
				ctx1 := SetEntry(tc.oldCtx, testEntry1)
				assert.NotEqual(t, tc.oldCtx, ctx1)
				l := ctx1.Value(Key)
				assert.NotNil(t, l)
				assert.IsType(t, &logrus.Entry{}, l)
				assert.Equal(t, testEntry1, l.(*logrus.Entry))
			})
		}
	})
}
