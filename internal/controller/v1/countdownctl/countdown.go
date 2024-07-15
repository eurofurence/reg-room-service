package countdownctl

import (
	"github.com/eurofurence/reg-room-service/internal/application/web"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// Handler implements methods, which satisfy the endpoint format.
type Handler struct{}

func InitRoutes(router chi.Router) {
	h := &Handler{}

	router.Route("/api/rest/v1/countdown", func(sr chi.Router) {
		initGetRoutes(sr, h)
	})
}

func initGetRoutes(router chi.Router, h *Handler) {
	router.Method(
		http.MethodGet,
		"/",
		web.CreateHandler(
			h.GetCountdown,
			h.GetCountdownRequest,
			h.GetCountdownResponse,
		),
	)
}
