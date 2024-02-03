package service

import (
	"context"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/store"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"
)

//go:generate mockgen -source=service.go -destination=mock/service_mock.go -package=service_mock

type NotifierSend interface {
	NotifierPending(ctx context.Context) error
}
type Metrics interface {
	GetMetrics(ctx context.Context, metrics model.Metrics) (model.Metrics, error)
	PutMetrics(ctx context.Context, metrics model.Metrics) error
	GetCounter(ctx context.Context, key string) (int64, error)
	GetGauge(ctx context.Context, key string) (float64, error)
	PutGauge(ctx context.Context, key string, val float64) error
	PutCounter(ctx context.Context, key string, val int64) error
	SendMetricstoFile(ctx context.Context) error
}

type Service struct {
	Metrics  Metrics
	Storage  store.Storager
	Notifier NotifierSend
	DB       store.DB
}

func NewService(storage store.Storager, notify NotifierSend, db store.DB) *Service {
	return &Service{
		Storage:  storage,
		Notifier: notify,
		Metrics:  NewMeticsService(storage, notify),
		DB:       db}

}
