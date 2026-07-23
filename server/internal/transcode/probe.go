package transcode

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
)

type Metadata struct {
	Duration float64
	Width int
	Height int
}

func Probe(ctx context.Context, url string) (*Metadata , error) {
	cmd := exec.CommandContext(ctx, "ffprobe",
		"-v", "error",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		url,
	)
	out, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return nil, fmt.Errorf("ffprobe failed: %w: %s", err, exitErr.Stderr)
		}
		return nil, fmt.Errorf("ffprobe: %w", err)
	}

	var probed struct {
		Format struct {
			Duration string `json:"duration"`
		} `json:"format"`
		Streams []struct {
			CodecType string `json:"codec_type"`
			Width int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}

	if err := json.Unmarshal(out, &probed); err != nil {
		return nil, fmt.Errorf("parsing ffprobe output: %w", err)
	}

	duration, err := strconv.ParseFloat(probed.Format.Duration, 64)
	if err != nil {
		return  nil, fmt.Errorf("parsing duration %q: %w", probed.Format.Duration, err)
	}

	for _, s := range probed.Streams {
		if s.CodecType == "video" {
			return &Metadata{
				Duration: duration,
				Width: s.Width,
				Height: s.Height,
			}, nil
		}
	}

	return nil, fmt.Errorf("no video stream found")
}