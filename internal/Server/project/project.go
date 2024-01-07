package project

import "github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"

type Project struct {
	model.CounterStorager
	model.GaugeStorager
}

func NewProject(counter model.CounterStorager, gauge model.GaugeStorager) *Project {
	return &Project{CounterStorager: counter, GaugeStorager: gauge}
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

func (p *Project) PutGauge(key string, val float64) error {
	return p.UpdateGauge(key, val)
}

func (p *Project) PutCounter(key string, val uint64) error {
	return p.UpdateCounter(key, val)
}
