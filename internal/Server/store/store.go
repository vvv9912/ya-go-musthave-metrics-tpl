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

type GaugeStorager interface {
	UpdateGauge(ctx context.Context, key string, val float64) error
	GetGauge(ctx context.Context, key string) (float64, error)
	GetAllGauge(ctx context.Context) (map[string]float64, error)
}

type CounterStorager interface {
	UpdateCounter(ctx context.Context, key string, val uint64) error
	GetCounter(ctx context.Context, key string) (uint64, error)
	GetAllCounter(ctx context.Context) (map[string]uint64, error)
}

func (db *Database) Ping(ctx context.Context) error {
	return db.pgx.PingContext(ctx)
}

// todo  tx sql.Tx,
func (db *Database) UpdateGauge(ctx context.Context, key string, val float64) error {
	err := db.updateGauge(ctx, key, val)
	return err
}

func (db *Database) GetGauge(ctx context.Context, key string) (float64, error) {
	return db.getGauge(ctx, key)
}
func (db *Database) GetAllGauge(ctx context.Context) (map[string]float64, error) {
	return db.getAllGauge(ctx)
}

func (db *Database) UpdateCounter(ctx context.Context, key string, val uint64) error {
	err := db.updateCounter(ctx, key, val)
	return err

}

func (db *Database) GetCounter(ctx context.Context, key string) (uint64, error) {
	return db.getCounter(ctx, key)
}

func (db *Database) GetAllCounter(ctx context.Context) (map[string]uint64, error) {
	return db.getAllCounter(ctx)
}
