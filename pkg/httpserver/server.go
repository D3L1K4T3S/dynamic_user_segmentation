package httpserver

import (
	"context"
	"net/http"
	"time"
)

const (
	defaultAddress         = "localhost:80"
	defaultWriteTimeout    = 10 * time.Second
	defaultReadTimeout     = 10 * time.Second
	defaultShutdownTimeout = 10 * time.Second
)

type Server struct {
	server          *http.Server
	shutdownTimeout time.Duration
	notify          chan error
}

func NewServer(handler http.Handler, options ...Option) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		Addr:         defaultAddress,
	}

	server := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, option := range options {
		option(server)
	}

	return server
}

func (s *Server) Start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}

func (s *Server) Notify() <-chan error {
	return s.notify
}
