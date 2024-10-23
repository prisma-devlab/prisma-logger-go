package log

import (
	"context"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var enabledLevel zapcore.Level

func init() {
	Init("info")
}

func Init(level string) {
	enabledLevel = getZapLogLevel(level)

	cfg := zap.NewProductionConfig()
	cfg.DisableCaller = true
	cfg.Level = zap.NewAtomicLevelAt(enabledLevel)

	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendInt64(t.Unix())
	}

	logger, _ = cfg.Build(zap.AddCallerSkip(1))
}

func getZapLogLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}

type ctxKey struct{}
type Fields map[string]interface{}
type F = Fields
type contextFields map[string]zap.Field

func disabled(level zapcore.Level) bool {
	return !enabledLevel.Enabled(level)
}

func Context() context.Context {
	return context.WithValue(context.Background(), ctxKey{}, contextFields{})
}

func getContextWithFields(ctx context.Context) (contextFields, context.Context) {
	if fields, ok := ctx.Value(ctxKey{}).(contextFields); ok {
		return fields, ctx
	}
	fields := contextFields{}
	return fields, context.WithValue(ctx, ctxKey{}, fields)
}

func zapIt(ctx context.Context, logFields Fields) []zap.Field {
	ctxFields, _ := getContextWithFields(ctx)

	zapFields := make([]zap.Field, 0, len(ctxFields)+len(logFields))

	for key, val := range logFields {
		zapFields = append(zapFields, zap.Any(key, val))
	}

	for key, v := range ctxFields {
		if logFields[key] == nil {
			zapFields = append(zapFields, v)
		}
	}
	return zapFields
}

func Debug(ctx context.Context, message string, fields Fields) {
	if disabled(zapcore.DebugLevel) {
		return
	}
	logger.Debug(message, zapIt(ctx, fields)...)
}

func Info(ctx context.Context, message string, fields Fields) {
	if disabled(zapcore.InfoLevel) {
		return
	}
	logger.Info(message, zapIt(ctx, fields)...)
}

func Warn(ctx context.Context, message string, fields Fields) {
	if disabled(zapcore.WarnLevel) {
		return
	}
	logger.Warn(message, zapIt(ctx, fields)...)
}

func Error(ctx context.Context, message string, fields Fields) {
	if disabled(zapcore.ErrorLevel) {
		return
	}
	logger.Error(message, zapIt(ctx, fields)...)
}
