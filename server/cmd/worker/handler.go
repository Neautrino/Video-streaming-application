package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/url"
	"path"
	"strings"

	"github.com/neautrino/video-streaming/internal/queue"
	"github.com/neautrino/video-streaming/internal/video"
)

type s3Event struct {
	Records []struct {
		EventName string `json:"eventName"`
		S3 struct {
			Object struct {
				Key string `jsong:"key"`
				Size int64 `json:"size"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}

func videoIDFromKey(key string) (string, error) {
	decoded, err := url.QueryUnescape(key)
	if err != nil {
		return "", err
	}
	base := path.Base(decoded)
	return strings.TrimSuffix(base, path.Ext(base)), nil
}

func handleMessage(
	ctx context.Context,
	logger *slog.Logger,
	repo *video.Repository,
	consumer *queue.Consumer,
	msg queue.Message,
) {
	var event s3Event
	if err := json.Unmarshal([]byte(msg.Body), &event); err != nil {
		logger.Error("unparseble message body", "err", err)
		return
	}

	for _, record := range event.Records {
		id, err := videoIDFromKey(record.S3.Object.Key)
		if err != nil {
			logger.Error("decoding object key", "key", record.S3.Object.Key, "err", err)
			return
		}

		updated, err := repo.MarkUploaded(ctx, id, record.S3.Object.Size)
		if err != nil {
			logger.Error("marking uploded", "video_id", id, "err", err)
			return
		}

		if updated {
			logger.Info("video uploded", "video_id", id, "size", record.S3.Object.Size)
		} else {
			logger.Info("already processed, skipping", "video_id", id)
		}
	}

	if err := consumer.Delete(ctx, msg.ReceiptHandle); err != nil {
		logger.Error("deleting message", "err", err)
	}
}