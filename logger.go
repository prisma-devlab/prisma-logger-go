package prismalogger

import (
	"time"

	log "github.com/prisma-devlab/prisma-logger-go/log/zap"
	"go.uber.org/zap"
)

type Format struct {
	RequestID string
	Event     string
	Endpoint  string
	Message   string
	Data      map[string]string
}

type PrismaLogger struct {
	*zap.Logger
}

func addIfNonEmpty(key string, value string, result log.F) {
	if value != "" {
		result[key] = value
	}
}

func getZapFields(logFormat Format) log.F {
	result := log.F{}
	addIfNonEmpty("request_id", logFormat.RequestID, result)
	addIfNonEmpty("event", logFormat.Event, result)
	addIfNonEmpty("endpoint", logFormat.Endpoint, result)
	addIfNonEmpty("timestamp", time.Now().Format(time.RFC3339), result)

	if logFormat.Data != nil {
		result["data"] = logFormat.Data
	}

	return result
}

func Debug(logFormat Format) {
	log.Debug(log.Context(), logFormat.Message, getZapFields(logFormat))
}

func Info(logFormat Format) {
	log.Info(log.Context(), logFormat.Message, getZapFields(logFormat))
}

func Warn(logFormat Format) {
	log.Warn(log.Context(), logFormat.Message, getZapFields(logFormat))
}

func Error(logFormat Format) {
	log.Error(log.Context(), logFormat.Message, getZapFields(logFormat))
}

func Init(logLevel string) {
	log.Init(logLevel)
}
