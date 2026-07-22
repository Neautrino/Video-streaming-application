package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/neautrino/video-streaming/internal/video"
)

type createVideoRequest struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Filename string `json:"filename"`
	Size int64 `json:"size"`
	ContentType string `json:"content_type"`
}

const uploadPrefix = "uploads/"

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req createVideoRequest
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	if req.Filename == "" {
		http.Error(w, "filename is required", http.StatusBadRequest)
		return
	}

	if req.Size <= 0 {
		http.Error(w, "size must be greater than 0", http.StatusBadRequest)
		return
	}

	if req.Size > s.cfg.MaxVideoBytes {
		http.Error(w, fmt.Sprintf("file size must be less than %d bytes", s.cfg.MaxVideoBytes), http.StatusRequestHeaderFieldsTooLarge)
		return
	}

	if req.ContentType == "" {
		http.Error(w, "content_type is required", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(req.ContentType, "video/") {
		http.Error(w, "content_type must be video", http.StatusBadRequest)
		return
	}

	buf := make([]byte, 16)
	rand.Read(buf)
	id := hex.EncodeToString(buf)

	v := &video.Video{
		ID: id,
		Title: req.Title,
		Description: req.Description,
		OriginalFileName: req.Filename,
		ContentType: req.ContentType,
		Size: req.Size,
		StorageKey: uploadPrefix + id + filepath.Ext(req.Filename),
		Status: video.StatusUploading,
	}

	if err := s.repo.Create(r.Context(), v); err != nil {
		s.logger.Error("creating video record", "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	url, err := s.storage.PresignPut(r.Context(), v.StorageKey)
	if err != nil {
		s.logger.Error("presigning upload url", "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": v.ID, "upload_url": url})
}