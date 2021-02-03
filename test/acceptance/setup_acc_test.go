package acceptance

import (
	"github.com/eurofurence/reg-room-service/internal/repository/config"
	"github.com/eurofurence/reg-room-service/web"
	"net/http/httptest"
)

var (
	ts *httptest.Server
)

func tstSetup() {
	tstSetupHttpTestServer()
	tstLoadConfig()
}

func tstSetupHttpTestServer() {
	router := web.Create()
	ts = httptest.NewServer(router)
}

func tstLoadConfig() {
	config.LoadConfiguration("../resources/testconfig.yaml")
}

func tstShutdown() {
	ts.Close()
}
