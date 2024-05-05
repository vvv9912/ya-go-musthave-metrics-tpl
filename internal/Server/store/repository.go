package store

import (
	"context"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"
)

//go:generate mockgen -source=repository.go -destination=repo_mock/repo_mock.go -package=repo_mock

// Интерфейс Storager.
type Storager interface {
	// UpdateGauge - обновляет значение метрики типа Gauge по ключу.
	UpdateGauge(ctx context.Context, key string, val float64) error
	// GetGauge - возвращает значение метрики типа Gauge по ключу.
	GetGauge(ctx context.Context, key string) (float64, error)
	// GetAllGauge - возвращает все значения метрик типа Gauge.
	GetAllGauge(ctx context.Context) (map[string]float64, error)
	//counter
	// UpdateCounter - обновляет значение метрики типа Counter по ключу.
	UpdateCounter(ctx context.Context, key string, val int64) error
	// GetCounter - возвращает значение метрики типа Counter по ключу.
	GetCounter(ctx context.Context, key string) (int64, error)
	// GetAllCounter - возвращает все значения метрик типа Counter.
	GetAllCounter(ctx context.Context) (map[string]int64, error)
	// UpdateMetricsBatch - обновляет значения метрик батчем.
	UpdateMetricsBatch(ctx context.Context, metrics []model.Metrics) error
	// Ping - проверяет доступность хранилища.
	Ping(ctx context.Context) error
}

// Repository - структура, реализует методы интерфейса Storager.
type Repository struct {
	Storager
}

// Конструктор репозитория.
func NewRepository(storager Storager) *Repository {
	return &Repository{Storager: storager}
}
