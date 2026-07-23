package video

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, v *Video) error {
	query := `INSERT INTO videos (id, title, description, original_filename, content_type, size, storage_key, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.pool.Exec(ctx, query, v.ID, v.Title, v.Description, v.OriginalFileName, v.ContentType, v.Size, v.StorageKey, v.Status)
	return err
}

func (r *Repository) MarkUploaded(ctx context.Context, id string, size int64) (bool, error) {
	query := `UPDATE videos 
	SET status = $1, size = $2, updated_at = now()
	WHERE id = $3 AND status = $4
	`
	tag, err := r.pool.Exec(ctx, query, StatusUploaded, size, id, StatusUploading)
	if err != nil {
		return false, err
	}

	return tag.RowsAffected() > 0, nil
}
