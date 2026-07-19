package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/neautrino/video-streaming/internal/api"
	"github.com/neautrino/video-streaming/internal/storage"
	"github.com/neautrino/video-streaming/internal/video"
)

func env(key, fallback string) string {
      if v := os.Getenv(key); v != "" {
              return v
      }
      return fallback
}

func main() {
	_ = godotenv.Load()

	cfg := api.Config{Addr: ":8080", MaxVideoBytes: 1 << 30}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	dbURL := env("DATABASE_URL", "postgres://postgres:streaming_dev@localhost:5432/streaming")
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		logger.Error("Database Connection Failed", "err", err)
		os.Exit(1)
	}
	
	defer pool.Close()

	if err = pool.Ping(context.Background()); err != nil {
		logger.Error("Database Ping Failed", "err", err)
		os.Exit(1)
	}

	logger.Info("Database connection successfully")

	repo := video.NewRepository(pool)

	bucket := os.Getenv("S3_BUCKET")
	if bucket == "" {
		logger.Error("S3_BUCKET is required")
		os.Exit(1) 
	}

	store, err := storage.New(context.Background(), storage.Config{Bucket: bucket})
	if err != nil {
		logger.Error("storage init failed", "err", err)
		os.Exit(1)
	}

	srv := api.NewServer(cfg, logger, repo, store)

	logger.Info("listening", "addr", cfg.Addr)

	if err := http.ListenAndServe(cfg.Addr, srv); err != nil {
		logger.Error("server failed", "err", err)
		os.Exit(1)
	}
}