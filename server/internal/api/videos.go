package api

import "net/http"

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}