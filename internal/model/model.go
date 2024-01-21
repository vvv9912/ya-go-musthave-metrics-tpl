package model

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

//type GaugeStorager interface {
//	UpdateGauge(key string, val float64) error
//	GetGauge(key string) (float64, error)
//	GetAllGauge() map[string]float64
//}
//
//type CounterStorager interface {
//	UpdateCounter(key string, val uint64) error
//	GetCounter(key string) (uint64, error)
//	GetAllCounter() map[string]uint64
//}
