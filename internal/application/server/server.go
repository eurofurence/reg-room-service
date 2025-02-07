package server

import (
	"context"
	"fmt"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-room-service/internal/repository/config"
	groupservice "github.com/eurofurence/reg-room-service/internal/service/groups"
	roomservice "github.com/eurofurence/reg-room-service/internal/service/rooms"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
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

	groupsvc groupservice.Service
	roomsvc  roomservice.Service
}

var _ Server = (*server)(nil)

type Server interface {
	Serve() error
	Shutdown() error
}

func New(conf *config.Config, baseCtx context.Context, groupsvc groupservice.Service, roomsvc roomservice.Service) Server {
	s := new(server)

	s.interrupt = make(chan os.Signal, 1)
	s.shutdown = make(chan struct{})

	s.ctx = baseCtx

	s.idleTimeout = time.Duration(conf.Server.IdleTimeout) * time.Second
	s.readTimeout = time.Duration(conf.Server.ReadTimeout) * time.Second
	s.writeTimeout = time.Duration(conf.Server.WriteTimeout) * time.Second

	s.host = conf.Server.BaseAddress
	s.port = conf.Server.Port

	s.groupsvc = groupsvc
	s.roomsvc = roomsvc

	return s
}

func (s *server) Serve() error {
	handler := Router(s.groupsvc, s.roomsvc)
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
