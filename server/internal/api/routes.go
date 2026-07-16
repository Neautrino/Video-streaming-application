package api

import "github.com/go-chi/chi/v5/middleware"

func (s *Server) routes() {
	s.router.Use(middleware.Logger)
	s.router.Post("/videos/upload", s.handleUpload)
	s.router.Get("/videos/{id}", s.getVideoById)
}