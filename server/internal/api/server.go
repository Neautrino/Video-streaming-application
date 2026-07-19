package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/neautrino/video-streaming/internal/video"
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
	repo *video.Repository
}

func NewServer(cfg Config, logger *slog.Logger, repo *video.Repository) *Server {
	s := &Server{
		cfg: cfg,
		router: chi.NewRouter(),
		logger: logger.With("component", "api"),
		repo: repo,
	}

	s.routes()

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w , r)
}