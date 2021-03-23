package web

import (
	"github.com/eurofurence/reg-room-service/internal/repository/config"
	"github.com/eurofurence/reg-room-service/internal/repository/logging"
	"github.com/eurofurence/reg-room-service/web/controller/countdownctl"
	"github.com/eurofurence/reg-room-service/web/controller/healthctl"
	"github.com/eurofurence/reg-room-service/web/filter/corsfilter"
	"github.com/eurofurence/reg-room-service/web/filter/jwt"
	"github.com/eurofurence/reg-room-service/web/filter/logreqid"
	"github.com/eurofurence/reg-room-service/web/filter/reqid"
	"github.com/go-chi/chi"
	"net/http"
)

func Create() chi.Router {
	logging.NoCtx().Info("Building routers...")
	server := chi.NewRouter()

	server.Use(reqid.RequestIdMiddleware())
	server.Use(logreqid.LogRequestIdMiddleware())
	server.Use(jwt.JwtMiddleware(config.JWTPublicKey()))
	server.Use(corsfilter.CorsHeadersMiddleware())

	healthctl.Create(server)
	countdownctl.Create(server)
	return server
}

func Serve(server chi.Router) {
	address := config.ServerAddr()
	logging.NoCtx().Info("Listening on " + address)
	err := http.ListenAndServe(address, server)
	if err != nil {
		logging.NoCtx().Error(err)
	}
}
