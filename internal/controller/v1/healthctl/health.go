package healthctl

import (
	"github.com/eurofurence/reg-room-service/internal/application/web"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// Handler implements methods, which satisfy the endpoint format.
type Handler struct{}

func InitRoutes(router chi.Router) {
	h := &Handler{}

	router.Route("/", func(sr chi.Router) {
		initGetRoutes(sr, h)
	})
}

func initGetRoutes(router chi.Router, h *Handler) {
	router.Method(
		http.MethodGet,
		"/",
		web.CreateHandler(
			h.GetHealth,
			h.GetHealthRequest,
			h.GetHealthResponse,
		),
	)
}
