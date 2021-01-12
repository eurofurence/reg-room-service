package logging

import (
	"context"
	"github.com/eurofurence/reg-room-service/internal/repository/logging/consolelogging"
)

type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})

	// expected to terminate the process
	Fatal(v ...interface{})
}

// context key with a separate type, so no other package has a chance of accessing it
type key int

// the value actually doesn't matter, the type alone will guarantee no package gets at this context value
const loggerKey key = 0


var defaultLogger = createLogger("00000000")

func createLogger(requestId string) Logger {
	return &consolelogging.ConsoleLoggingImpl{RequestId: requestId}
}

func CreateContextWithLoggerForRequestId(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, loggerKey, createLogger(requestId))
}

// you should only use this when your code really does not belong to request processing.
// otherwise be a good citizen and do pass down the context, so log output can be associated with
// the request being processed!
func NoCtx() Logger {
	return defaultLogger
}

// whenever processing a specific request, use this and give it the context.
func Ctx(ctx context.Context) Logger {
	logger, ok := ctx.Value(loggerKey).(Logger)
	if !ok {
		// better than no logger at all
		return defaultLogger
	}
	return logger
}
