package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/neautrino/video-streaming/internal/api"
)

func main() {

	cfg := api.Config{Addr: ":8080", UploadDir: "data/uploads"}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	srv := api.NewServer(cfg)

	logger.Info("listening", "addr", cfg.Addr)

	if err := http.ListenAndServe(cfg.Addr, srv); err != nil {
		logger.Error("server failed", "err", err)
		os.Exit(1)
	}
}