package app

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/eurofurence/reg-room-service/internal/repository/config"
	"github.com/eurofurence/reg-room-service/internal/web/controller/countdownctl"
	"github.com/eurofurence/reg-room-service/internal/web/controller/healthctl"
	"github.com/eurofurence/reg-room-service/internal/web/middleware"
)

func CreateRouter(ctx context.Context) chi.Router {
	server := chi.NewRouter()

	// server.Use(middleware.AddRequestIdToContextAndResponse)
	// server.Use(middleware.RequestLogger)
	// server.Use(middleware.PanicRecoverer)
	server.Use(middleware.CorsHandling)
	// server.Use(middleware.TokenValidator)

	countdownctl.Create(server)
	healthctl.Create(server)

	return server
}

func newServer(ctx context.Context, router chi.Router) *http.Server {
	return &http.Server{
		Addr:    config.ServerAddr(),
		Handler: router,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}
}

func runServerWithGracefulShutdown() error {
	ctx, cancel := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	handler := CreateRouter(ctx)
	srv := newServer(ctx, handler)

	go func() {
		<-sig
		defer cancel()

		tCtx, tcancel := context.WithTimeout(ctx, time.Second*5)
		defer tcancel()

		if err := srv.Shutdown(tCtx); err != nil {
			os.Exit(3)
		}
	}()

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
