package store

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type Database struct {
	pgx *pgx.Conn
}

func NewDatabase(db *pgx.Conn) *Database {
	return &Database{pgx: db}
}

func (db *Database) Ping(ctx context.Context) error {
	return db.pgx.Ping(ctx)
}
