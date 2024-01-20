package store

import (
	"context"
	"database/sql"
)

//go:generate mockgen -source=store.go -destination=mock/store_mock.go
type Database struct {
	pgx *sql.DB
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{pgx: db}
}

func (db *Database) Ping(ctx context.Context) error {
	return db.pgx.PingContext(ctx)
}
