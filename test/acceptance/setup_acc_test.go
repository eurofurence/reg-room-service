package acceptance

import (
	"context"
	"net/http/httptest"

	"github.com/eurofurence/reg-room-service/internal/web/app"

	"github.com/eurofurence/reg-room-service/internal/repository/config"
)

var ts *httptest.Server

const (
	tstDefaultConfigFileBeforeLaunch      = "../resources/testconfig_beforeLaunch.yaml"
	tstDefaultConfigFileAfterStaffLaunch  = "../resources/testconfig_afterStaffLaunch.yaml"
	tstDefaultConfigFileAfterPublicLaunch = "../resources/testconfig_afterPublicLaunch.yaml"
)

func tstSetup(configfile string) {
	tstLoadConfig(configfile)
	tstSetupHttpTestServer()
}

func tstSetupHttpTestServer() {
	router := app.CreateRouter(context.Background())
	ts = httptest.NewServer(router)
}

func tstLoadConfig(configfile string) {
	config.LoadConfiguration(configfile)
}

func tstShutdown() {
	ts.Close()
}
