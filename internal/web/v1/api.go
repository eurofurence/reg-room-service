package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/eurofurence/reg-room-service/internal/controller"
	"github.com/eurofurence/reg-room-service/internal/web/v1/groups"
)

func Router(ctrl controller.Controller) http.Handler {
	router := chi.NewMux()

	groups.InitRoutes(router, ctrl)
	// TODO(noroth) create handler and add
	// remaining routes for rooms and countdown

	return router
}
