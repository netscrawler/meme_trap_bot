package repository

import (
	"database/sql"
	"time"

	"github.com/netscrawler/trap_bot/domain"
)

type Image struct {
	ID int64 `db:"id"`

	FileID       string        `db:"file_id"`
	FileUniqueID string        `db:"file_unique_id"`
	FileSize     sql.NullInt64 `db:"file_size"`

	Width  sql.NullInt64 `db:"width"`
	Height sql.NullInt64 `db:"height"`

	CreatedAt time.Time `db:"created_at"`

	AddedBy   string `db:"added_by"`
	AddedByID int64  `db:"added_by_id"`
}

func NewImageFromDomain(img domain.Image) Image {
	return Image{
		ID:           img.ID,
		FileID:       img.FileID,
		FileUniqueID: img.FileUniqueID,
		FileSize:     sql.NullInt64{Int64: int64(img.FileSize), Valid: img.FileSize != 0},
		Width:        sql.NullInt64{Int64: int64(img.Width), Valid: img.Width != 0},
		Height:       sql.NullInt64{Int64: int64(img.Height), Valid: img.Height != 0},
		CreatedAt:    img.CreatedAt,
		AddedBy:      img.AddedBy,
		AddedByID:    img.AddedByID,
	}
}

func NewImageFromDTO(img Image) domain.Image {
	return domain.Image{
		ID:           img.ID,
		FileID:       img.FileID,
		FileUniqueID: img.FileUniqueID,
		FileSize:     int(img.FileSize.Int64),
		Width:        int(img.Width.Int64),
		Height:       int(img.Height.Int64),
		CreatedAt:    img.CreatedAt,
		AddedBy:      img.AddedBy,
		AddedByID:    img.AddedByID,
	}
}
