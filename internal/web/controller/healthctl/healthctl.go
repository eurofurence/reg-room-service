package healthctl

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/eurofurence/reg-room-service/internal/repository/logging"
)

func Create(server chi.Router) {
	server.Get("/", health)
}

func health(w http.ResponseWriter, r *http.Request) {
	logging.Ctx(r.Context()).Info("health")
	w.WriteHeader(http.StatusOK)
}
