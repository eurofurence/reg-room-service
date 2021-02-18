package acceptance

import (
	"github.com/eurofurence/reg-room-service/internal/repository/config"
	"github.com/eurofurence/reg-room-service/web"
	"net/http/httptest"
)

var (
	ts *httptest.Server
)

const (
	tstDefaultConfigFileBeforeLaunch = "../resources/testconfig_beforeLaunch.yaml"
	tstDefaultConfigFileAfterLaunch = "../resources/testconfig_afterLaunch.yaml"
)

func tstSetup(configfile string) {
	tstLoadConfig(configfile)
	tstSetupHttpTestServer()
}

func tstSetupHttpTestServer() {
	router := web.Create()
	ts = httptest.NewServer(router)
}

func tstLoadConfig(configfile string) {
	config.LoadConfiguration(configfile)
}

func tstShutdown() {
	ts.Close()
}
