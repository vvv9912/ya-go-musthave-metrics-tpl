package service

import (
	"context"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"
)

//go:generate mockgen -source=service.go -destination=mock/service_mock.go -package=service_mock

type GaugeStorager interface {
	UpdateGauge(key string, val float64) error
	GetGauge(key string) (float64, error)
	GetAllGauge() map[string]float64
}

type CounterStorager interface {
	UpdateCounter(key string, val uint64) error
	GetCounter(key string) (uint64, error)
	GetAllCounter() map[string]uint64
}

type NotifierSend interface {
	NotifierPending() error
}
type Metrics interface {
	GetMetrics(metrics model.Metrics) (model.Metrics, error)
	PutMetrics(metrics model.Metrics) error
	GetCounter(key string) (uint64, error)
	GetGauge(key string) (float64, error)
	PutGauge(key string, val float64) error
	PutCounter(key string, val uint64) error
	SendMetricstoFile() error
}

type Store interface {
	Ping(ctx context.Context) error
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
