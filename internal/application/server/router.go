package server

import (
	"github.com/StephanHCB/go-autumn-logging-zerolog/loggermiddleware"
	"github.com/eurofurence/reg-room-service/internal/application/middleware"
	"github.com/eurofurence/reg-room-service/internal/controller/v1/countdownctl"
	"github.com/eurofurence/reg-room-service/internal/controller/v1/groupsctl"
	"github.com/eurofurence/reg-room-service/internal/controller/v1/healthctl"
	"github.com/eurofurence/reg-room-service/internal/controller/v1/roomsctl"
	"github.com/eurofurence/reg-room-service/internal/repository/config"
	"github.com/eurofurence/reg-room-service/internal/repository/database"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams/attendeeservice"
	groupservice "github.com/eurofurence/reg-room-service/internal/service/groups"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Router(db database.Repository, attsrv attendeeservice.AttendeeService) http.Handler {
	router := chi.NewMux()

	conf, err := config.GetApplicationConfig()
	if err != nil {
		panic("no config loaded - this is a bug")
	}

	router.Use(middleware.PanicRecoverer)
	router.Use(middleware.RequestIdMiddleware)
	router.Use(loggermiddleware.AddZerologLoggerToContext)
	router.Use(middleware.RequestLoggerMiddleware)
	router.Use(middleware.CorsHeadersMiddleware(&conf.Security))
	router.Use(middleware.CheckRequestAuthorization(&conf.Security))

	groupsctl.InitRoutes(router, groupservice.NewService(db, attsrv))
	roomsctl.InitRoutes(router)
	countdownctl.InitRoutes(router)
	healthctl.InitRoutes(router)

	return router
}
