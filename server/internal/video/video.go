package video

import "time"

type Status string

const (
	StatusUploaded Status = "uploaded"
)

type Video struct {
	ID string
	OriginalFileName string
	StorageKey string
	Status Status
	CreatedAt time.Time
}