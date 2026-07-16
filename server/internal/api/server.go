package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Config struct {
	Addr string
	UploadDir string
	MaxUploadBytes int64
}

type Server struct {
	cfg Config
	router *chi.Mux
	logger *slog.Logger
}

func NewServer(cfg Config, logger *slog.Logger) *Server {
	s := &Server{
		cfg: cfg,
		router: chi.NewRouter(),
		logger: logger.With("component", "api"),
	}

	s.routes()

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w , r)
}