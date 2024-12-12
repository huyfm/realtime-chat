package pg

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func OpenDB(ctx context.Context, connStr string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, connStr)
}

func CloseDB(db *pgxpool.Pool) {
	db.Close()
}
