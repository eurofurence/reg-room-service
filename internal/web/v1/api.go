package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/eurofurence/reg-room-service/internal/web/v1/groups"
)

func buildRouter() http.Handler {
	router := chi.NewMux()

	groups.InitRoutes(router)

	return router
}

func Router() http.Handler {
	return buildRouter()
}
