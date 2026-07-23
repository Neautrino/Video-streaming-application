package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/neautrino/video-streaming/internal/queue"
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

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With("component", "worker")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbURL := env("DATABASE_URL", "postgres://postgres:streaming_dev@localhost:5432/streaming")
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		logger.Error("Database Connection Failed", "err", err)
		os.Exit(1)
	}

	defer pool.Close()

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

	consumer, err := queue.NewConsumer(ctx, queue.Config{
		QueueURL: os.Getenv("SQS_UPLOADS_QUEUE_URL"),
	})
	if err != nil {
		logger.Error("SQS Uploads queue Connection Failed", "err", err)
		os.Exit(1)
	}

	transcodeURL := os.Getenv("SQS_TRANSCODE_QUEUE_URL")
	if transcodeURL == "" {
		logger.Error("SQS_TRANSCODE_QUEUE_URL is required")
		os.Exit(1)
	}

	producer, err := queue.NewProducer(ctx, queue.Config{QueueURL: transcodeURL})
	if err != nil {
		logger.Error("queue producer init failed", "err", err)
		os.Exit(1)
	}

	h := &handler{
		logger:   logger,
		repo:     repo,
		consumer: consumer,
		producer: producer,
		storage:  store,
	}

	for {
		msgs, err := consumer.Receive(ctx)
		if err != nil {
			if ctx.Err() != nil {
				logger.Info("shutting down")
				return
			}
			logger.Error("receiving messages", "err", err)
			continue
		}

		for _, msg := range msgs {
			h.handleMessage(ctx, msg)
		}
	}

}
