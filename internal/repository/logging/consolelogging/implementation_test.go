package consolelogging

import (
	"testing"
)

var cut = &ConsoleLoggingImpl{RequestID: "00000000"}

func TestConsoleLoggingImpl_Debug(t *testing.T) {
	cut.Debug("a", "b", "c")
}

func TestConsoleLoggingImpl_Info(t *testing.T) {
	cut.Info("d", "e", "f")
}

func TestConsoleLoggingImpl_Warn(t *testing.T) {
	cut.Warn("x", "y", "z")
}

func TestConsoleLoggingImpl_Error(t *testing.T) {
	cut.Error("some error happened")
}
