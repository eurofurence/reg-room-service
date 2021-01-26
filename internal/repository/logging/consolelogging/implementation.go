package consolelogging

import (
	"fmt"
	"log"
	"os"
)

const severityDEBUG = "DEBUG"
const severityINFO = "INFO"
const severityWARN = "WARN"
const severityERROR = "ERROR"
const severityFATAL = "FATAL"

type ConsoleLoggingImpl struct {
	RequestId string
}

func (l *ConsoleLoggingImpl) isEnabled(severity string) bool {
	// TODO read from config
	return true
}

func (l *ConsoleLoggingImpl) print(severity string, v ...interface{}) {
	if l.isEnabled(severity) {
		args := []interface{}{fmt.Sprintf("%-5s [%s] ", severity, l.RequestId)}
		args = append(args, v...)
		log.Print(args...)
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
