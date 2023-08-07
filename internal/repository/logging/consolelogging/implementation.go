package consolelogging

import (
	"log"
	"os"

	"github.com/eurofurence/reg-room-service/internal/repository/logging/consolelogging/logformat"
)

const (
	severityDEBUG = "DEBUG"
	severityINFO  = "INFO"
	severityWARN  = "WARN"
	severityERROR = "ERROR"
	severityFATAL = "FATAL"
)

type ConsoleLoggingImpl struct {
	RequestId string
}

func (l *ConsoleLoggingImpl) isEnabled(severity string) bool {
	return true
}

func (l *ConsoleLoggingImpl) print(severity string, v ...interface{}) {
	if l.isEnabled(severity) {
		log.Print(logformat.Logformat(severity, l.RequestId, v...))
	}
}

func (l *ConsoleLoggingImpl) Debug(v ...interface{}) {
	l.print(severityDEBUG, v...)
}

func (l *ConsoleLoggingImpl) Info(v ...interface{}) {
	l.print(severityINFO, v...)
}

func (l *ConsoleLoggingImpl) Warn(v ...interface{}) {
	l.print(severityWARN, v...)
}

func (l *ConsoleLoggingImpl) Error(v ...interface{}) {
	l.print(severityERROR, v...)
}

func (l *ConsoleLoggingImpl) Fatal(v ...interface{}) {
	l.print(severityFATAL, v...)
	os.Exit(1)
}
