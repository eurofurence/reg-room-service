package server

import (
	"github.com/StephanHCB/go-autumn-logging-zerolog/loggermiddleware"
	"github.com/eurofurence/reg-room-service/internal/application/middleware"
	"github.com/eurofurence/reg-room-service/internal/controller/v1/countdownctl"
	"github.com/eurofurence/reg-room-service/internal/controller/v1/groupsctl"
	"github.com/eurofurence/reg-room-service/internal/controller/v1/healthctl"
	"github.com/eurofurence/reg-room-service/internal/controller/v1/roomsctl"
	"github.com/eurofurence/reg-room-service/internal/repository/config"
	groupservice "github.com/eurofurence/reg-room-service/internal/service/groups"
	roomservice "github.com/eurofurence/reg-room-service/internal/service/rooms"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Router(groupsvc groupservice.Service, roomsvc roomservice.Service) http.Handler {
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

	groupsctl.InitRoutes(router, groupsvc)
	roomsctl.InitRoutes(router, roomsvc)
	countdownctl.InitRoutes(router)
	healthctl.InitRoutes(router)

	return router
}
