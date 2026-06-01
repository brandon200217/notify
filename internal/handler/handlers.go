package handler

import (
	"net/http"

	"github.com/brandon200217/NOTIFY/internal/channel"
	"github.com/brandon200217/NOTIFY/internal/config"
	"github.com/brandon200217/NOTIFY/internal/middleware"
)

type Server struct {
	cfg         *config.Config
	router      *http.ServeMux
	registry    *channel.Registry
	rateLimiter *middleware.RateLimiter
}

func NewServer(cfg *config.Config, registry *channel.Registry, rateLimiter *middleware.RateLimiter) *Server {
	s := &Server{
		cfg:         cfg,
		router:      http.NewServeMux(),
		registry:    registry,
		rateLimiter: rateLimiter,
	}
	s.registerRoutes()
	return s
}

func (s *Server) Router() http.Handler {
	return s.router
}

func (s *Server) registerRoutes() {
	handler := http.HandlerFunc(s.handleNotify)
	wrapped := middleware.RequestID(s.rateLimiter.MiddlewareLimit(handler))
	s.router.Handle("POST /notify", wrapped)

}
