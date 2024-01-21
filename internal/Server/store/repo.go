package store

import (
	"context"
)

func (db *Database) updateGauge(ctx context.Context, key string, val float64) error {
	//_, err := db.pgx.ExecContext(ctx, "UPDATE GaugeMetrics SET val=$1 WHERE key=$2", val, key)
	_, err := db.pgx.ExecContext(ctx, "INSERT INTO GaugeMetrics (key, val) VALUES ($1, $2) ON CONFLICT (key) DO UPDATE SET val = $2;", key, val)
	if err != nil {
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

func (db *Database) updateCounter(ctx context.Context, key string, val uint64) error {
	//_, err := db.pgx.ExecContext(ctx, "UPDATE CounterMetrics SET val=$1 WHERE key=$2", val, key)
	_, err := db.pgx.ExecContext(ctx, "INSERT INTO CounterMetrics (key, val) VALUES ($1, $2) ON CONFLICT (key) DO UPDATE SET val = $2;", key, val)
	if err != nil {
		return err
		//todo добавить
	}
	return nil
}
func (db *Database) getCounter(ctx context.Context, key string) (uint64, error) {
	rows, err := db.pgx.QueryContext(ctx, "SELECT val FROM CounterMetrics where key=$1", key)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return 0, rows.Err()
	}

	var val int
	for rows.Next() {
		err = rows.Scan(&val)
		if err != nil {
			return 0, err
		}
	}

	return uint64(val), nil

}
func (db *Database) getAllCounter(ctx context.Context) (map[string]uint64, error) {
	rows, err := db.pgx.QueryContext(ctx, "SELECT * FROM CounterMetrics")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	metrics := make(map[string]uint64)

	for rows.Next() {
		var key string
		var val int
		err = rows.Scan(&key, &val)
		if err != nil {
			return nil, err
		}
		metrics[key] = uint64(val)
	}

	return metrics, nil

}
