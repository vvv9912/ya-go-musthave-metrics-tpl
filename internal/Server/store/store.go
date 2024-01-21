package store

import (
	"context"
	"database/sql"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"
	"sync"
)

//go:generate mockgen -source=store.go -destination=mock/store_mock.go
type Database struct {
	pgx *sql.DB
	mu  sync.Mutex
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
	UpdateCounter(ctx context.Context, key string, val int64) error
	GetCounter(ctx context.Context, key string) (int64, error)
	GetAllCounter(ctx context.Context) (map[string]int64, error)
}

func (db *Database) Ping(ctx context.Context) error {
	return db.pgx.PingContext(ctx)
}

func (db *Database) UpdateMetricsBatch(ctx context.Context, metrics []model.Metrics) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	tx, err := db.pgx.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	return db.updateMetricsBatch(ctx, tx, metrics)
}

func (db *Database) UpdateGauge(ctx context.Context, key string, val float64) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	tx, err := db.pgx.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	return db.updateGauge(ctx, tx, key, val)
}

func (db *Database) GetGauge(ctx context.Context, key string) (float64, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.getGauge(ctx, key)
}
func (db *Database) GetAllGauge(ctx context.Context) (map[string]float64, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.getAllGauge(ctx)
}

func (db *Database) UpdateCounter(ctx context.Context, key string, val int64) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	tx, err := db.pgx.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	return db.updateCounter(ctx, tx, key, val)
}

func (db *Database) GetCounter(ctx context.Context, key string) (int64, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.getCounter(ctx, key)
}

func (db *Database) GetAllCounter(ctx context.Context) (map[string]int64, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.getAllCounter(ctx)
}
