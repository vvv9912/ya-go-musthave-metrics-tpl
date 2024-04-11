// Пакет postgresql реализует взаимодейсвие с БД.
package postgresql

import (
	"context"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"
)

//go:generate mockgen -source=store.go -destination=mock/store_mock.go

// Ping - проверяет доступность хранилища.
func (db *Database) Ping(ctx context.Context) error {
	return db.pgx.PingContext(ctx)
}

// UpdateMetricsBatch - обновляет значения метрик батчем.
func (db *Database) UpdateMetricsBatch(ctx context.Context, metrics []model.Metrics) error {

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

// UpdateGauge - обновляет значение метрики типа Gauge по ключу.
func (db *Database) UpdateGauge(ctx context.Context, key string, val float64) error {

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

// GetGauge - возвращает значение метрики типа Gauge по ключу.
func (db *Database) GetGauge(ctx context.Context, key string) (float64, error) {
	return db.getGauge(ctx, key)
}

// GetAllGauge - возвращает все значения метрик типа Gauge.
func (db *Database) GetAllGauge(ctx context.Context) (map[string]float64, error) {
	return db.getAllGauge(ctx)
}

// UpdateCounter - обновляет значение метрики типа Counter по ключу.
func (db *Database) UpdateCounter(ctx context.Context, key string, val int64) error {
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

// GetCounter - возвращает значение метрики типа Counter по ключу.
func (db *Database) GetCounter(ctx context.Context, key string) (int64, error) {
	return db.getCounter(ctx, key)
}

// GetAllCounter - возвращает все значения метрик типа Counter.
func (db *Database) GetAllCounter(ctx context.Context) (map[string]int64, error) {
	return db.getAllCounter(ctx)
}
