package logger

import (
	"context"
	"testing"

	"github.com/Sirupsen/logrus"
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
	t.Run("given a context with no entry, it returns a entry", func(t *testing.T) {
		result := FromContext(ctx0)
		assert.IsType(t, &logrus.Entry{}, result)
	})

	t.Run("given a context with a entry, it returns that entry", func(t *testing.T) {
		ctx1 := context.WithValue(ctx0, Key, testEntry1)
		ctx2 := FromContext(ctx1)
		assert.Equal(t, testEntry1, ctx2)
	})
}

func TestSetEntry(t *testing.T) {
	t.Run("it returns a new context with the given entry...", func(t *testing.T) {
		tcs := []struct {
			scenario string
			oldE     *logrus.Entry
		}{
			{
				"given a context with no entry",
				nil,
			},
			{
				"given a context with the same entry",
				testEntry1,
			},
			{
				"given a context with a different entry",
				testEntry2,
			},
		}
		for _, tc := range tcs {
			t.Run(tc.scenario, func(t *testing.T) {
				var ctx1 context.Context
				if tc.oldE != nil {
					ctx1 = context.WithValue(ctx0, Key, tc.oldE)
				} else {
					ctx1 = ctx0
				}
				ctx2 := SetEntry(ctx1, testEntry1)
				assert.NotEqual(t, ctx1, ctx2)
				l := ctx2.Value(Key)
				assert.NotNil(t, l)
				assert.IsType(t, &logrus.Entry{}, l)
				assert.Equal(t, testEntry1, l.(*logrus.Entry))
			})
		}
	})
}
