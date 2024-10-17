package acceptance

import (
	"context"
	"github.com/eurofurence/reg-room-service/internal/application/server"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/attendeeservice"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/mailservice"
	groupservice "github.com/eurofurence/reg-room-service/internal/service/groups"
	roomservice "github.com/eurofurence/reg-room-service/internal/service/rooms"
	"net/http/httptest"

	"github.com/eurofurence/reg-room-service/internal/repository/config"
	"github.com/eurofurence/reg-room-service/internal/repository/database"
	"github.com/eurofurence/reg-room-service/internal/repository/database/historizeddb"
	"github.com/eurofurence/reg-room-service/internal/repository/database/inmemorydb"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/authservice"
)

var ts *httptest.Server
var db database.Repository
var authMock authservice.Mock
var attMock attendeeservice.Mock
var mailMock mailservice.Mock

const (
	tstDefaultConfigFileBeforeLaunch      = "../resources/testconfig_beforeLaunch.yaml"
	tstDefaultConfigFileAfterStaffLaunch  = "../resources/testconfig_afterStaffLaunch.yaml"
	tstDefaultConfigFileAfterPublicLaunch = "../resources/testconfig_afterPublicLaunch.yaml"
	tstDefaultConfigFileRoomGroups        = "../resources/testconfig_roomgroups.yaml"
)

func tstSetup(configfile string) {
	tstLoadConfig(configfile)
	db = tstCreateInmemoryDatabase()
	authMock = authservice.CreateMock()
	authMock.Enable()
	attMock = attendeeservice.NewMock()
	mailMock = mailservice.NewMock()

	grpsvc := groupservice.New(db, attMock, mailMock)
	roomsvc := roomservice.New(db, attMock, mailMock)

	tstSetupAuthMockResponses()
	tstSetupHttpTestServer(grpsvc, roomsvc)
}

func tstSetupHttpTestServer(grpsrv groupservice.Service, roomsvc roomservice.Service) {
	router := server.Router(grpsrv, roomsvc)
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
	if _, err := config.UnmarshalFromYamlConfiguration(configfile); err != nil {
		panic("failed to load config")
	}
}

func tstShutdown() {
	ts.Close()
	db.Close(context.TODO())
}
