package service

import (
	"context"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"
)

type MeticsService struct {
	counter CounterStorager
	gauge   GaugeStorager
	notify  NotifierSend //н
}

func NewMeticsService(counter CounterStorager, gauge GaugeStorager, notify NotifierSend) *MeticsService {
	return &MeticsService{counter: counter, gauge: gauge, notify: notify}
}
func (p *MeticsService) GetMetrics(ctx context.Context, metrics model.Metrics) (model.Metrics, error) {
	switch metrics.MType {
	case "counter":

		val, err := p.GetCounter(ctx, metrics.ID)
		if err != nil {
			return model.Metrics{}, err
		}
		delta := int64(val)

		metrics.Delta = &delta

		return metrics, nil

	case "gauge":
		val, err := p.GetGauge(ctx, metrics.ID)
		if err != nil {
			return model.Metrics{}, err
		}

		metrics.Value = &val

		return metrics, nil
	}
	return model.Metrics{}, nil
}
func (p *MeticsService) PutMetrics(ctx context.Context, metrics model.Metrics) error {
	switch metrics.MType {
	case "counter":
		return p.PutCounter(ctx, metrics.ID, uint64(*metrics.Delta))
	case "gauge":
		return p.PutGauge(ctx, metrics.ID, *metrics.Value)
	}
	return nil
}
func (p *MeticsService) GetCounter(ctx context.Context, key string) (uint64, error) {
	return p.counter.GetCounter(ctx, key)
}
func (p *MeticsService) GetGauge(ctx context.Context, key string) (float64, error) {
	return p.gauge.GetGauge(ctx, key)
}
func (p *MeticsService) PutGauge(ctx context.Context, key string, val float64) error {
	return p.gauge.UpdateGauge(ctx, key, val)
}

func (p *MeticsService) PutCounter(ctx context.Context, key string, val uint64) error {
	return p.counter.UpdateCounter(ctx, key, val)
}

func (p *MeticsService) SendMetricstoFile(ctx context.Context) error {
	return p.notify.NotifierPending(ctx)
}
