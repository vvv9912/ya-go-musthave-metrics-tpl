package postgresql

import (
	"context"
	"database/sql"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/logger"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"
	"go.uber.org/zap"
)

type Database struct {
	pgx *sql.DB
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{pgx: db}
}
func (db *Database) updateMetricsBatch(ctx context.Context, tx *sql.Tx, metrics []model.Metrics) error {

	for _, v := range metrics {
		if v.MType == "gauge" {
			err := db.updateGauge(ctx, tx, v.ID, *v.Value)
			if err != nil {
				return err
			}
		} else if v.MType == "counter" {
			err := db.updateCounter(ctx, tx, v.ID, *v.Delta)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
func (db *Database) updateGauge(ctx context.Context, tx *sql.Tx, key string, val float64) error {
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO GaugeMetrics (key, val) VALUES ($1, $2) ON CONFLICT (key) DO UPDATE SET val = $2;")
	if err != nil {
		logger.Log.Info("Failed to create statment for update gauge", zap.Error(err))
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, key, val)
	if err != nil {
		logger.Log.Info("Failed to update gauge", zap.Error(err))
		return err
	}

	return nil
}
func (db *Database) getGauge(ctx context.Context, key string) (float64, error) {
	rows, err := db.pgx.QueryContext(ctx, "SELECT val FROM GaugeMetrics where key=$1", key)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Err() != nil {
		return 0, rows.Err()
	}

	var val float64
	for rows.Next() {
		err = rows.Scan(&val)
		if err != nil {
			return 0, err
		}
	}
	return val, nil
}
func (db *Database) getAllGauge(ctx context.Context) (map[string]float64, error) {
	rows, err := db.pgx.QueryContext(ctx, "SELECT * FROM GaugeMetrics")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	metrics := make(map[string]float64)
	for rows.Next() {
		var key string
		var val float64
		err = rows.Scan(&key, &val)
		if err != nil {
			return nil, err
		}
		metrics[key] = val
	}

	return metrics, nil
}

func (db *Database) updateCounter(ctx context.Context, tx *sql.Tx, key string, val int64) error {
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO CounterMetrics (key, val) VALUES ($1, $2) ON CONFLICT (key) DO UPDATE SET val = CounterMetrics.val + $2;")
	if err != nil {
		logger.Log.Info("Failed to create statment for update counter", zap.Error(err))
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, key, val)
	if err != nil {
		logger.Log.Info("Failed to update counter", zap.Error(err))
		return err
	}

	return nil
}
func (db *Database) getCounter(ctx context.Context, key string) (int64, error) {
	rows, err := db.pgx.QueryContext(ctx, "SELECT val FROM CounterMetrics where key=$1", key)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return 0, rows.Err()
	}

	var val int64
	for rows.Next() {
		err = rows.Scan(&val)
		if err != nil {
			return 0, err
		}
	}

	return val, nil

}
func (db *Database) getAllCounter(ctx context.Context) (map[string]int64, error) {
	rows, err := db.pgx.QueryContext(ctx, "SELECT * FROM CounterMetrics")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	metrics := make(map[string]int64)

	for rows.Next() {
		var key string
		var val int64
		err = rows.Scan(&key, &val)
		if err != nil {
			return nil, err
		}
		metrics[key] = val
	}

	return metrics, nil

}
