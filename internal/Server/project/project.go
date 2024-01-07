package project

import "github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"

type Project struct {
	counter model.CounterStorager
	gauge   model.GaugeStorager
}

func NewProject(counter model.CounterStorager, gauge model.GaugeStorager) *Project {
	return &Project{counter: counter, gauge: gauge}
}
func (p *Project) GetMetrics(metrics model.Metrics) (model.Metrics, error) {
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
func (p *Project) PutMetrics(metrics model.Metrics) error {
	switch metrics.MType {
	case "counter":
		return p.PutCounter(metrics.ID, uint64(*metrics.Delta))
	case "gauge":
		return p.PutGauge(metrics.ID, *metrics.Value)
	}
	return nil
}
func (p *Project) GetCounter(key string) (uint64, error) {
	return p.counter.GetCounter(key)
}
func (p *Project) GetGauge(key string) (float64, error) {
	return p.gauge.GetGauge(key)
}
func (p *Project) PutGauge(key string, val float64) error {
	return p.gauge.UpdateGauge(key, val)
}

func (p *Project) PutCounter(key string, val uint64) error {
	return p.counter.UpdateCounter(key, val)
}
