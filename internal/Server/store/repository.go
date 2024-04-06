package store

import (
	"context"

	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"
)

//go:generate mockgen -source=repository.go -destination=repo_mock/repo_mock.go -package=repo_mock
type Storager interface {
	UpdateGauge(ctx context.Context, key string, val float64) error
	GetGauge(ctx context.Context, key string) (float64, error)
	GetAllGauge(ctx context.Context) (map[string]float64, error)
	//counter
	UpdateCounter(ctx context.Context, key string, val int64) error
	GetCounter(ctx context.Context, key string) (int64, error)
	GetAllCounter(ctx context.Context) (map[string]int64, error)
	UpdateMetricsBatch(ctx context.Context, metrics []model.Metrics) error
	Ping(ctx context.Context) error
}

type Repository struct {
	Storager
}

func NewRepository(storager Storager) *Repository {
	return &Repository{Storager: storager}
}
