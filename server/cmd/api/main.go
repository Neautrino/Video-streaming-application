package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/neautrino/video-streaming/internal/api"
)

func main() {

	cfg := api.Config{Addr: ":8080", UploadDir: "data/uploads", MaxUploadBytes: 1 << 30}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if err := os.MkdirAll(cfg.UploadDir, 0o755); err != nil {
      logger.Error("creating upload dir", "err", err)
      os.Exit(1)
	}

	srv := api.NewServer(cfg, logger)

	logger.Info("listening", "addr", cfg.Addr)

	if err := http.ListenAndServe(cfg.Addr, srv); err != nil {
		logger.Error("server failed", "err", err)
		os.Exit(1)
	}
}