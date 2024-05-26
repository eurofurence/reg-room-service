package v1

import (
	"github.com/eurofurence/reg-room-service/internal/web/v1/health"
	"net/http"

	"github.com/StephanHCB/go-autumn-logging-zerolog/loggermiddleware"
	"github.com/eurofurence/reg-room-service/internal/config"
	"github.com/eurofurence/reg-room-service/internal/repository/database"
	"github.com/eurofurence/reg-room-service/internal/web/middleware"

	"github.com/go-chi/chi/v5"

	groupservice "github.com/eurofurence/reg-room-service/internal/service/groups"
	"github.com/eurofurence/reg-room-service/internal/web/v1/countdown"
	"github.com/eurofurence/reg-room-service/internal/web/v1/groups"
	"github.com/eurofurence/reg-room-service/internal/web/v1/rooms"
)

func Router(db database.Repository) http.Handler {
	router := chi.NewMux()

	conf, err := config.GetApplicationConfig()
	if err != nil || conf == nil {
		// TODO
		panic("no config loaded or nil")
	}

	router.Use(middleware.PanicRecoverer)
	router.Use(middleware.RequestIdMiddleware)
	router.Use(loggermiddleware.AddZerologLoggerToContext)
	router.Use(middleware.RequestLoggerMiddleware)
	router.Use(middleware.CorsHeadersMiddleware(&conf.Security))
	router.Use(middleware.CheckRequestAuthorization(&conf.Security))

	groups.InitRoutes(router, groupservice.NewService(db))
	rooms.InitRoutes(router, nil)
	countdown.InitRoutes(router)
	health.InitRoutes(router)

	return router
}
