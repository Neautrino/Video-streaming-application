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

func (r *Repository) Init(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS videos (
	id text primary key,
	original_filename text not null,
	storage_key text not null,
	status text not null,
	created_at timestamptz not null default now()
	)`)
	return  err
}