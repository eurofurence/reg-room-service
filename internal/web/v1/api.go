package v1

import (
	"net/http"

	"github.com/eurofurence/reg-room-service/internal/controller"
	"github.com/eurofurence/reg-room-service/internal/web/v1/countdown"
	"github.com/eurofurence/reg-room-service/internal/web/v1/groups"
	"github.com/eurofurence/reg-room-service/internal/web/v1/rooms"
	"github.com/go-chi/chi/v5"
)

func Router(ctrl controller.Controller) http.Handler {
	router := chi.NewMux()

	groups.InitRoutes(router, ctrl)
	rooms.InitRoutes(router, ctrl)
	countdown.InitRoutes(router, ctrl)

	return router
}
