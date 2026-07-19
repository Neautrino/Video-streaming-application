package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/neautrino/video-streaming/internal/storage"
	"github.com/neautrino/video-streaming/internal/video"
)

type Config struct {
	Addr string
	MaxVideoBytes int64
}

type Server struct {
	cfg Config
	router *chi.Mux
	logger *slog.Logger
	repo *video.Repository
	storage *storage.Client
}

func NewServer(cfg Config, logger *slog.Logger, repo *video.Repository, storage *storage.Client) *Server {
	s := &Server{
		cfg: cfg,
		router: chi.NewRouter(),
		logger: logger.With("component", "api"),
		repo: repo,
		storage: storage,
	}

	s.routes()

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w , r)
}