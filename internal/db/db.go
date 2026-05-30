package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	db *pgxpool.Pool
}

func NewDB(pgx *pgxpool.Pool) *DB {
	return &DB{db: pgx}
}

func (d *DB) Ping(ctx context.Context) error {
	return d.db.Ping(ctx)
}
