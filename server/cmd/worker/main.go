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

	repo := video.NewRepository(pool)

	consumer, err := queue.NewConsumer(ctx, queue.Config{
		QueueURL: os.Getenv("SQS_UPLOADS_QUEUE_URL"),
	})
	if err != nil {
		logger.Error("SQS Uploads queue Connection Failed", "err", err)
		os.Exit(1)
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
			handleMessage(ctx, logger, repo, consumer, msg)
		}
	}

}