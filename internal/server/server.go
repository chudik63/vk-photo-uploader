package server

import (
	"context"
	"net/http"
	"time"
	"vk-photo-uploader/internal/infrastructure"
)

const shutdownTimeout = 5 * time.Second

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *infrastructure.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + cfg.Server.Port,
			Handler: handler,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop() error {
	ctx, shutdown := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdown()

	return s.httpServer.Shutdown(ctx)
}
