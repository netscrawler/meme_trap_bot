package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/netscrawler/trap_bot/domain"
)

const maxParams = 999

const fieldsPerRow = 7

const batchSize = maxParams / fieldsPerRow

type MediaFilesRepository struct {
	db *sql.DB
}

func NewMediaFilesRepository(db *sql.DB) *MediaFilesRepository {
	return &MediaFilesRepository{db: db}
}

func (r *MediaFilesRepository) Save(ctx context.Context, imgs []domain.Image) error {
	if len(imgs) == 0 {
		return fmt.Errorf("no images to save")
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}
	dtos := make([]Image, 0, len(imgs))
	for _, i := range imgs {
		dtos = append(dtos, NewImageFromDomain(i))
	}

	return BulkInsertImages(r.db, dtos)
}

func (r *MediaFilesRepository) GetRandom(ctx context.Context) (domain.Image, error) {
	queryFast := `
		SELECT
			id,
			file_id,
			file_unique_id,
			file_size,
			width,
			height,
			created_at,
			added_by,
			added_by_id
		FROM images
		WHERE id >= (
			ABS(RANDOM()) % (SELECT MAX(id) FROM images)
		)
		LIMIT 1
	`

	var img Image

	err := r.db.QueryRowContext(ctx, queryFast).Scan(
		&img.ID,
		&img.FileID,
		&img.FileUniqueID,
		&img.FileSize,
		&img.Width,
		&img.Height,
		&img.CreatedAt,
		&img.AddedBy,
		&img.AddedByID,
	)

	if err == sql.ErrNoRows {
		err = r.db.QueryRowContext(ctx, `
			SELECT
				id,
				file_id,
				file_unique_id,
				file_size,
				width,
				height,
				created_at,
				added_by,
				added_by_id
			FROM images
			ORDER BY RANDOM()
			LIMIT 1
		`).Scan(
			&img.ID,
			&img.FileID,
			&img.FileUniqueID,
			&img.FileSize,
			&img.Width,
			&img.Height,
			&img.CreatedAt,
			&img.AddedBy,
			&img.AddedByID,
		)
	}

	if err != nil {
		return domain.Image{}, err
	}

	return NewImageFromDTO(img), nil
}

func BulkInsertImages(db *sql.DB, images []Image) error {
	if len(images) == 0 {
		return nil
	}

	for i := 0; i < len(images); i += batchSize {
		end := min(i+batchSize, len(images))

		if err := insertBatch(db, images[i:end]); err != nil {
			return err
		}
	}

	return nil
}

func insertBatch(db *sql.DB, images []Image) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var (
		values []string
		args   []any
	)

	for _, img := range images {
		values = append(values, "(?, ?, ?, ?, ?, ?, ?)")

		args = append(args,
			img.FileID,
			img.FileUniqueID,
			nullInt64ToInterface(img.FileSize),
			nullInt64ToInterface(img.Width),
			nullInt64ToInterface(img.Height),
			img.AddedBy,
			img.AddedByID,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO images (
			file_id,
			file_unique_id,
			file_size,
			width,
			height,
			added_by,
			added_by_id
		) VALUES %s
	`, strings.Join(values, ","))

	_, err = tx.Exec(query, args...)
	if err != nil {
		return err
	}

	return tx.Commit()
}
func nullInt64ToInterface(v sql.NullInt64) any {
	if v.Valid {
		return v.Int64
	}
	return nil
}
