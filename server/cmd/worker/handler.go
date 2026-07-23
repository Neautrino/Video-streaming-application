package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/url"
	"path"
	"strings"

	"github.com/neautrino/video-streaming/internal/queue"
	"github.com/neautrino/video-streaming/internal/storage"
	"github.com/neautrino/video-streaming/internal/transcode"
	"github.com/neautrino/video-streaming/internal/video"
)

type s3Event struct {
	Records []struct {
		EventName string `json:"eventName"`
		S3        struct {
			Object struct {
				Key  string `json:"key"`
				Size int64  `json:"size"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}

type handler struct {
	logger   *slog.Logger
	repo     *video.Repository
	consumer *queue.Consumer
	producer *queue.Producer
	storage  *storage.Client
}

func videoIDFromKey(key string) (string, error) {
	decoded, err := url.QueryUnescape(key)
	if err != nil {
		return "", err
	}
	base := path.Base(decoded)
	return strings.TrimSuffix(base, path.Ext(base)), nil
}

func (h *handler) handleMessage(
	ctx context.Context,
	msg queue.Message,
) {
	var event s3Event
	if err := json.Unmarshal([]byte(msg.Body), &event); err != nil {
		h.logger.Error("unparseble message body", "err", err)
		return
	}

	for _, record := range event.Records {
		id, err := videoIDFromKey(record.S3.Object.Key)
		if err != nil {
			h.logger.Error("decoding object key", "key", record.S3.Object.Key, "err", err)
			return
		}

		updated, err := h.repo.MarkUploaded(ctx, id, record.S3.Object.Size)
		if err != nil {
			h.logger.Error("marking uploded", "video_id", id, "err", err)
			return
		}

		if updated {
			h.logger.Info("video uploded", "video_id", id, "size", record.S3.Object.Size)

			v, err := h.repo.GetById(ctx, id)
			if err != nil {
				h.logger.Error("feching video for prepare", "video_id", id, "err", err )
				return
			}

			url, err := h.storage.PresignGet(ctx, v.StorageKey)
			if err != nil {
				h.logger.Error("presigning source url", "video_id", id, "err", err)
				return
			}

			meta, err := transcode.Probe(ctx, url)
			if err != nil {
				h.logger.Error("probing video", "video_id", id, "err", err)
				return
			}

			h.logger.Info("probed", "video_id", id, "duration", meta.Duration,"width", meta.Width, "height", meta.Height)
		} else {
			h.logger.Info("already processed, skipping", "video_id", id)
		}
	}

	if err := h.consumer.Delete(ctx, msg.ReceiptHandle); err != nil {
		h.logger.Error("deleting message", "err", err)
	}
}
