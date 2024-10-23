package prismalogger_test

import (
	"testing"

	prismalogger "github.com/prisma-devlab/prisma-logger-go"
	"github.com/stretchr/testify/suite"
)

type LoggerTestSuite struct {
	suite.Suite
}

func (k *LoggerTestSuite) TestLogLevels() {
	tests := []struct {
		logLevel string
		message  string
		expected string
	}{
		{"debug", "I am debug", "I am debug"},
		{"info", "I am info", "I am info"},
		{"warn", "I am warn", "I am warn"},
		{"error", "I am error", "I am error"},
	}

	for _, test := range tests {
		prismalogger.Init(test.logLevel)

		msg := prismalogger.Format{
			RequestID: "RequestId",
			Event:     "EventName",
			Endpoint:  "",
			Message:   test.message,
			Data: map[string]string{
				"one": "one",
				"two": "two",
			},
		}

		prismalogger.Debug(msg)
		prismalogger.Info(msg)
		prismalogger.Warn(msg)
		prismalogger.Error(msg)
	}
}

func (k *LoggerTestSuite) TestShouldInitOnce() {
	prismalogger.Init("debug")
	prismalogger.Init("info")
	prismalogger.Init("error")
}

func TestLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}
