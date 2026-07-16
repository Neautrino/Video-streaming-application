package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, s.cfg.MaxUploadBytes)

	file, header, err := r.FormFile("file")
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			s.logger.Error("File too large")
			http.Error(w, "file too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "Missing or Invalid 'file' field", http.StatusBadRequest)
		return
	}

	buf := make([]byte, 16)
	rand.Read(buf)
	id := hex.EncodeToString(buf)

	defer file.Close()

	path := filepath.Join(s.cfg.UploadDir, id+filepath.Ext(header.Filename))
	dist, err := os.Create(path)

	if err != nil {
		s.logger.Error("creating upload path", "err", err)
		http.Error(w, "Creating upload path", http.StatusInternalServerError)
		return
	}
	
	defer dist.Close()

	if _, err := io.Copy(dist, file); err != nil {
		s.logger.Error("Uploading File", "err", err)
		http.Error(w, "Uploading File", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id" : id})
}