package healthctl

import (
	"github.com/eurofurence/reg-room-service/internal/repository/logging"
	"github.com/go-chi/chi"
	"net/http"
)

func Create(server chi.Router) {
	server.Get("/", health)
}

func health(w http.ResponseWriter, r *http.Request) {
	logging.Ctx(r.Context()).Info("health")
	w.WriteHeader(http.StatusOK)
}
