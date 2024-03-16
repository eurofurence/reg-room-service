package v1

import (
	"github.com/eurofurence/reg-room-service/internal/repository/database"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/eurofurence/reg-room-service/internal/service/groups"
	"github.com/eurofurence/reg-room-service/internal/web/v1/countdown"
	"github.com/eurofurence/reg-room-service/internal/web/v1/groups"
	"github.com/eurofurence/reg-room-service/internal/web/v1/rooms"
)

func Router(db database.Repository) http.Handler {
	router := chi.NewMux()

	groups.InitRoutes(router, groupservice.NewService(db))
	rooms.InitRoutes(router, nil)
	countdown.InitRoutes(router, nil)

	return router
}
