package transcode

import "math"

type Task struct {
	VideoID string `json:"video_id"`
	SourceKey string `json:"source_key"`
	SegmentIndex int `json:"segment_index"`
	Start float64 `json:"start"`
	Duration float64 `json:"duration"`
	Rendition string `json:"rendition"`
}

const segmentDuration = 6.0

func Plan(videoID, sourceKey string, meta *Metadata) []Task {
	segments := int(math.Ceil(meta.Duration / segmentDuration))
	tasks := make([]Task, 0, segments)

	for i:= 0;i<segments;i++ {
		start := float64(i) * segmentDuration
		tasks = append(tasks, Task{
			VideoID: videoID,
			SourceKey: sourceKey,
			SegmentIndex: i,
			Start: start,
			Duration: math.Min(segmentDuration, meta.Duration-start),
			Rendition: "360p",
		})
	}
	return  tasks
}