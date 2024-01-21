package service

import (
	"context"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"
)

//go:generate mockgen -source=service.go -destination=mock/service_mock.go -package=service_mock

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
	//SendBatchedMetrcs(ctx context.Context, data []model.Metrics) error
}

type Store interface {
	Ping(ctx context.Context) error
	UpdateMetricsBatch(ctx context.Context, metrics []model.Metrics) error
}

type Service struct {
	Metrics         Metrics
	CounterStorager CounterStorager
	GaugeStorager   GaugeStorager
	Notifier        NotifierSend
	Store           Store
}

func NewService(counter CounterStorager, gauge GaugeStorager, notify NotifierSend, store Store) *Service {
	return &Service{
		CounterStorager: counter,
		GaugeStorager:   gauge,
		Notifier:        notify,
		Metrics:         NewMeticsService(counter, gauge, notify),
		Store:           store}

}
