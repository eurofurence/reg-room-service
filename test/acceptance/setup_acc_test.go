package acceptance

import (
	"github.com/eurofurence/reg-room-service/web"
	"net/http/httptest"
)

var (
	ts *httptest.Server
)

func tstSetup() {
	tstSetupHttpTestServer()
}

func tstSetupHttpTestServer() {
	router := web.Create()
	ts = httptest.NewServer(router)
}

func tstShutdown() {
	ts.Close()
}
