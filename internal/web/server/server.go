package server

import (
	"context"
	"fmt"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-room-service/internal/config"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eurofurence/reg-room-service/internal/repository/database"

	"github.com/pkg/errors"

	v1 "github.com/eurofurence/reg-room-service/internal/web/v1"
)

type server struct {
	ctx          context.Context
	host         string
	port         int
	idleTimeout  time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration

	srv *http.Server

	interrupt chan os.Signal
	shutdown  chan struct{}
}

var _ Server = (*server)(nil)

type Server interface {
	Serve(database.Repository) error
	Shutdown() error
}

func NewServer(conf *config.Config, baseCtx context.Context) Server {
	s := new(server)

	s.interrupt = make(chan os.Signal, 1)
	s.shutdown = make(chan struct{})

	s.ctx = baseCtx

	// TODO should be in config so it is obvious what they are set to
	s.idleTimeout = time.Minute
	s.readTimeout = time.Minute
	s.writeTimeout = time.Minute

	s.host = conf.Server.BaseAddress
	s.port = conf.Server.Port

	return s
}

func (s *server) Serve(db database.Repository) error {
	handler := v1.Router(db)
	s.srv = s.newServer(handler)

	s.setupSignalHandler()
	go s.handleInterrupt()

	aulogging.Logger.NoCtx().Info().Printf("serving requests on %s...", s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	<-s.shutdown

	return nil
}

func (s *server) newServer(handler http.Handler) *http.Server {
	return &http.Server{
		BaseContext: func(l net.Listener) context.Context {
			return s.ctx
		},
		Handler:      handler,
		IdleTimeout:  s.idleTimeout,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		Addr:         fmt.Sprintf("%s:%d", s.host, s.port),
	}
}

func (s *server) setupSignalHandler() {
	s.interrupt = make(chan os.Signal)
	signal.Notify(s.interrupt, syscall.SIGINT, syscall.SIGTERM)
}

func (s *server) handleInterrupt() {
	<-s.interrupt
	if err := s.Shutdown(); err != nil {
		log.Fatal(err)
	}
}

func (s *server) Shutdown() error {
	defer close(s.shutdown)

	aulogging.Logger.NoCtx().Info().Print("gracefully shutting down server")

	tCtx, cancel := context.WithTimeout(s.ctx, time.Second*20)
	defer cancel()

	if err := s.srv.Shutdown(tCtx); err != nil {
		return errors.Wrap(err, "couldn't gracefully shut down server")
	}

	return nil
}
