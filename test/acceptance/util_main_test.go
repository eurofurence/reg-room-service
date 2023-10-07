package acceptance

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	auzerolog "github.com/StephanHCB/go-autumn-logging-zerolog"
	"github.com/eurofurence/reg-room-service/internal/repository/database/dbrepo"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	tstStartInMemoryDB()

	code := m.Run()

	tstStopInMemoryDB()

	os.Exit(code)
}

func tstStartInMemoryDB() {
	ctx := auzerolog.AddLoggerToCtx(context.Background())
	err := dbrepo.Open(ctx, "inmemory", "irrelevant")
	if err != nil {
		aulogging.ErrorErrf(ctx, err, "failed to open inmemory db: %s", err.Error())
	}
}

func tstStopInMemoryDB() {
	ctx := auzerolog.AddLoggerToCtx(context.Background())
	dbrepo.Close(ctx)
}
