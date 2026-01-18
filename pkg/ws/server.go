package ws

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type Server struct {
	server *http.Server
}

func NewServer(port int) *Server {
	if port < 1 || port > 65535 {
		slog.Error("invalid port number, must be between 1 and 65535", "port", port)
		os.Exit(1)
	}

	return &Server{
		server: &http.Server{
			Addr: ":" + fmt.Sprintf("%d", port),
		},
	}
}

func (s *Server) RegisterHandler(handler http.Handler) {
	s.server.Handler = handler
}

func (s *Server) Start() error {
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error("Server error", "error", err)
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
