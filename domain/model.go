package domain

import (
	"time"
)

type Image struct {
	ID           int64
	FileID       string
	FileUniqueID string
	FileSize     int

	Width  int
	Height int

	CreatedAt time.Time

	AddedBy   string
	AddedByID int64
}
