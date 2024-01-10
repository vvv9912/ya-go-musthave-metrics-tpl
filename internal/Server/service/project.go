package service

import (
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
func (p *MeticsService) GetMetrics(metrics model.Metrics) (model.Metrics, error) {
	switch metrics.MType {
	case "counter":

		val, err := p.GetCounter(metrics.ID)
		if err != nil {
			return model.Metrics{}, err
		}
		delta := int64(val)

		metrics.Delta = &delta

		return metrics, nil

	case "gauge":
		val, err := p.GetGauge(metrics.ID)
		if err != nil {
			return model.Metrics{}, err
		}

		metrics.Value = &val

		return metrics, nil
	}
	return model.Metrics{}, nil
}
func (p *MeticsService) PutMetrics(metrics model.Metrics) error {
	switch metrics.MType {
	case "counter":
		return p.PutCounter(metrics.ID, uint64(*metrics.Delta))
	case "gauge":
		return p.PutGauge(metrics.ID, *metrics.Value)
	}
	return nil
}
func (p *MeticsService) GetCounter(key string) (uint64, error) {
	return p.counter.GetCounter(key)
}
func (p *MeticsService) GetGauge(key string) (float64, error) {
	return p.gauge.GetGauge(key)
}
func (p *MeticsService) PutGauge(key string, val float64) error {
	return p.gauge.UpdateGauge(key, val)
}

func (p *MeticsService) PutCounter(key string, val uint64) error {
	return p.counter.UpdateCounter(key, val)
}

func (p *MeticsService) SendMetricstoFile() error {
	return p.notify.NotifierPending()
}
