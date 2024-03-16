package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/netip"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eurofurence/reg-room-service/internal/repository/database"

	"github.com/pkg/errors"

	"github.com/eurofurence/reg-room-service/internal/controller"
	v1 "github.com/eurofurence/reg-room-service/internal/web/v1"
)

type server struct {
	ctx  context.Context
	host string
	port string

	hasListeners bool

	ctrl controller.Controller

	listener net.Listener
	srv      *http.Server

	interrupt chan os.Signal
	shutdown  chan struct{}
}

var _ Server = (*server)(nil)

type Server interface {
	Serve(context.Context, database.Repository) error
	Shutdown() error
}

func NewServer() Server {
	s := new(server)

	s.interrupt = make(chan os.Signal, 1)
	s.shutdown = make(chan struct{})

	return s
}

func (s *server) Serve(ctx context.Context, db database.Repository) error {
	if err := s.Listen(); err != nil {
		return err
	}

	if err := s.setupTCPServer(db); err != nil {
		return errors.Wrap(err, "couldn't setup http server")
	}

	s.setupSignalHandler()
	go s.handleInterrupt()

	if err := s.srv.Serve(s.listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
	<-s.shutdown

	return nil
}

func (s *server) setupTCPServer(db database.Repository) error {
	// TODO

	srv := new(http.Server)

	srv.BaseContext = func(l net.Listener) context.Context {
		return s.ctx
	}
	srv.Handler = v1.Router(db)
	srv.IdleTimeout = time.Minute
	srv.ReadTimeout = time.Minute
	srv.WriteTimeout = time.Minute

	s.srv = srv

	return nil
}

func (s *server) Listen() error {
	if s.hasListeners {
		// Server was already set up
		return nil
	}

	addr, err := netip.ParseAddrPort(net.JoinHostPort(s.host, s.port))
	if err != nil {
		return err
	}

	l, err := net.ListenTCP("tcp", net.TCPAddrFromAddrPort(addr))
	if err != nil {
		return err
	}

	s.listener = l

	s.hasListeners = true
	return nil
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

	fmt.Println("gracefully shutting down server")

	tCtx, cancel := context.WithTimeout(s.ctx, time.Second*20)
	defer cancel()

	if err := s.srv.Shutdown(tCtx); err != nil {
		return errors.Wrap(err, "couldn't gracefully shut down server")
	}

	return nil
}
