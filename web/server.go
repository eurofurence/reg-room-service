package web

import (
	"github.com/go-chi/chi"
	"net/http"
)

func Create() chi.Router {
	server := chi.NewRouter()
	return server
}

func Serve(server chi.Router) {
	address := ":8080"
	err := http.ListenAndServe(address, server)
	if err != nil {
		// TODO
	}
}