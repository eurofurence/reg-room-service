package countdown

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/eurofurence/reg-room-service/internal/web/common"
)

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
		common.CreateHandler(
			h.GetCountdown,
			h.GetCountdownRequest,
			h.GetCountdownResponse,
		),
	)
}
