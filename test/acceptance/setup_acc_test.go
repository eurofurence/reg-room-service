package acceptance

import (
	"context"
	"github.com/eurofurence/reg-room-service/internal/repository/database"
	"github.com/eurofurence/reg-room-service/internal/repository/database/historizeddb"
	"github.com/eurofurence/reg-room-service/internal/repository/database/inmemorydb"
	v1 "github.com/eurofurence/reg-room-service/internal/web/v1"
	"net/http/httptest"
)

var ts *httptest.Server
var db database.Repository

const (
	tstDefaultConfigFileBeforeLaunch      = "../resources/testconfig_beforeLaunch.yaml"
	tstDefaultConfigFileAfterStaffLaunch  = "../resources/testconfig_afterStaffLaunch.yaml"
	tstDefaultConfigFileAfterPublicLaunch = "../resources/testconfig_afterPublicLaunch.yaml"
	tstDefaultConfigFileRoomGroups        = "../resources/testconfig_roomgroups.yaml"
)

func tstSetup(configfile string) {
	tstLoadConfig(configfile)
	db = tstCreateInmemoryDatabase()
	tstSetupHttpTestServer(db)
}

func tstSetupHttpTestServer(db database.Repository) {
	router := v1.Router(db)
	ts = httptest.NewServer(router)
}

func tstCreateInmemoryDatabase() database.Repository {
	db := historizeddb.New(inmemorydb.New())
	if err := db.Open(context.TODO()); err != nil {
		panic("failed to open inmemory database")
	}
	return db
}

func tstLoadConfig(configfile string) {
	// config.LoadConfiguration(configfile)
}

func tstShutdown() {
	ts.Close()
	db.Close(context.TODO())
}
