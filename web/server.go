package web

import (
	"github.com/eurofurence/reg-room-service/web/controller/healthctl"
	"github.com/go-chi/chi"
	"net/http"
)

func Create() chi.Router {
	server := chi.NewRouter()

	healthctl.Create(server)
	return server
}

func Serve(server chi.Router) {
	address := ":8080"
	err := http.ListenAndServe(address, server)
	if err != nil {
		// TODO
	}
}