package handler

import (
	"net/http"

	"github.com/brandon200217/NOTIFY/internal/channel"
	"github.com/brandon200217/NOTIFY/internal/config"
)

type Server struct {
	cfg      *config.Config
	router   *http.ServeMux
	registry *channel.Registry
}

func NewServer(cfg *config.Config, registry *channel.Registry) *Server {
	s := &Server{
		cfg:      cfg,
		router:   http.NewServeMux(),
		registry: registry,
	}
	s.registerRoutes()
	return s
}

func (s *Server) Router() http.Handler {
	return s.router
}

func (s *Server) registerRoutes() {
	s.router.HandleFunc("POST /notify", s.handleNotify)
}
