package v1

import (
	"net/http"

	"github.com/eurofurence/reg-room-service/internal/web/v1/groups"
	"github.com/go-chi/chi/v5"
)

func buildRouter() http.Handler {
	router := chi.NewMux()

	groups.InitRoutes(router)

	return router
}

func Router() http.Handler {
	return buildRouter()
}
