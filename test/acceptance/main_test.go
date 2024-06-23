package acceptance

import (
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// squelch normal log output
	aulogging.SetupNoLoggerForTesting()

	code := m.Run()
	os.Exit(code)
}
