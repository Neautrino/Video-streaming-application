package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Config struct {
	Addr string
	UploadDir string
}

type Server struct {
	cfg Config
	router *chi.Mux
}

func NewServer(cfg Config) *Server {
	s := &Server{
		cfg: cfg,
		router: chi.NewRouter(),
	}

	s.routes()

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w , r)
}