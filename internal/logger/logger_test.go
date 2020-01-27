package logger_test

import (
	"testing"

	. "github.com/coby9241/frontend-service/internal/logger"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestGetLogger(t *testing.T) {
	// default GetInstance should set to INFO
	logInstance := GetInstance()
	hook := test.NewLocal(logInstance)

	// test error level should be logged
	logInstance.Error("Helloerror")
	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, "Helloerror", hook.LastEntry().Message)

	hook.Reset()
	// test info level should be logged
	logInstance.Info("Helloinfo")
	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.InfoLevel, hook.LastEntry().Level)
	assert.Equal(t, "Helloinfo", hook.LastEntry().Message)

	hook.Reset()
	// test debug level should not be logged
	logInstance.Debug("Hellodebug")
	assert.Equal(t, 0, len(hook.Entries))
}
