package logging

import (
	"context"

	aulogging "github.com/StephanHCB/go-autumn-logging"
	auzerolog "github.com/StephanHCB/go-autumn-logging-zerolog"
	auloggingapi "github.com/StephanHCB/go-autumn-logging/api"
	"github.com/rs/zerolog"
)

type Logger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})

	// expected to terminate the process
	Fatal(format string, v ...interface{})
}

type loggingWrapper struct {
	logger auloggingapi.ContextAwareLoggingImplementation
}

func (l *loggingWrapper) Debug(format string, v ...interface{}) {
	l.logger.Debug().Printf(format, v...)
}

func (l *loggingWrapper) Info(format string, v ...interface{}) {
	l.logger.Info().Printf(format, v...)
}

func (l *loggingWrapper) Warn(format string, v ...interface{}) {
	l.logger.Warn().Printf(format, v...)
}

func (l *loggingWrapper) Error(format string, v ...interface{}) {
	l.logger.Error().Printf(format, v...)
}

// expected to terminate the process.
func (l *loggingWrapper) Fatal(format string, v ...interface{}) {
	l.logger.Fatal().Printf(format, v...)
}

// context key with a separate type, so no other package has a chance of accessing it.
type key int

// the value actually doesn't matter, the type alone will guarantee no package gets at this context value.
const RequestIDKey key = 0

const defaultReqID = "00000000"

func GetRequestID(ctx context.Context) string {
	reqID := ctx.Value(RequestIDKey)
	if reqID == nil {
		return defaultReqID
	}
	return reqID.(string)
}

func SetupLogging(applicationName string, useEcsLogging bool) {
	aulogging.RequestIdRetriever = GetRequestID
	if useEcsLogging {
		auzerolog.SetupJsonLogging(applicationName)
	} else {
		aulogging.DefaultRequestIdValue = defaultReqID
		auzerolog.SetupPlaintextLogging()
	}
}

func SetLoglevel(severity string) {
	switch severity {
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "WARN":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
	}
}

func LoggerFromContext(ctx context.Context) Logger {
	logger := aulogging.Logger.Ctx(ctx)

	return &loggingWrapper{
		logger: logger,
	}
}

func NewLogger() Logger {
	logger := aulogging.Logger.NoCtx()

	return &loggingWrapper{
		logger: logger,
	}
}

func ChildCtxWithRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, reqID)
}

func NewNoopLogger() Logger {
	return &noopLogger{}
}

type noopLogger struct{}

func (l *noopLogger) Debug(format string, v ...interface{}) {
}

func (l *noopLogger) Info(format string, v ...interface{}) {
}

func (l *noopLogger) Warn(format string, v ...interface{}) {
}

func (l *noopLogger) Error(format string, v ...interface{}) {
}

func (l *noopLogger) Fatal(format string, v ...interface{}) {
}
