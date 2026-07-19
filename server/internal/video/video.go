package video

import "time"

type Status string

const (
	StatusUploaded Status = "uploaded"
	StatusUploading Status = "uploading"
	StatusFailed Status = "failed"
)

type Video struct {
	ID string
	Title string
	Description string
	OriginalFileName string
	ContentType string
	Size int64
	StorageKey string
	Status Status
	CreatedAt time.Time
	UpdatedAt time.Time
}