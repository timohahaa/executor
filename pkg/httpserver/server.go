package httpserver

import (
	"context"
	"net/http"
	"time"
)

const (
	defaultAddres          = ":80"
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultShutdownTimeout = 5 * time.Second
)

type Server struct {
	server          *http.Server
	shutDownTimeout time.Duration
}

func New(handler http.Handler, opts ...Option) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		Addr:         defaultAddres,
	}

	s := &Server{
		server:          httpServer,
		shutDownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(s)
	}

	s.start()
	return s
}

func (s *Server) start() {
	go func() {
		s.server.ListenAndServe()
	}()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutDownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
