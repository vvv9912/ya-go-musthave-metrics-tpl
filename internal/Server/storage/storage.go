package storage

import (
	"context"
	"errors"
	"sync"
)

// TODO переделать под общий интерфейс get, update
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
type MemStorage struct {
	gaugeStorage   map[string]float64 //todo переделать под map[string]string, в соотв. с agent
	counterStorage map[string]int64
	gaugeMutex     sync.Mutex
	counterMutex   sync.Mutex
}

func NewGaugeStorage() GaugeStorager {
	gaugeStorage := make(map[string]float64)
	return &MemStorage{gaugeStorage: gaugeStorage}
}
func NewCounterStorage() CounterStorager {
	counterStorage := make(map[string]int64)
	return &MemStorage{counterStorage: counterStorage}
}

/*
Тип gauge, float64 — новое значение должно замещать предыдущее.
Тип counter, int64 — новое значение должно добавляться к предыдущему,
если какое-то значение уже было известно серверу.
*/
func (S *MemStorage) UpdateGauge(ctx context.Context, key string, val float64) error {
	S.gaugeMutex.Lock()
	defer S.gaugeMutex.Unlock()
	S.gaugeStorage[key] = val
	return nil
}
func (S *MemStorage) GetGauge(ctx context.Context, key string) (float64, error) {
	S.gaugeMutex.Lock()
	defer S.gaugeMutex.Unlock()
	val, found := S.gaugeStorage[key]
	if !found {
		return 0, errors.New("gauge not found")
	}
	return val, nil
}
func (S *MemStorage) GetAllGauge(ctx context.Context) (map[string]float64, error) {
	return S.gaugeStorage, nil
}
func (S *MemStorage) UpdateCounter(ctx context.Context, key string, val int64) error {
	S.counterMutex.Lock()
	defer S.counterMutex.Unlock()
	//fmt.Println("value counter get to map:", val)
	_, found := S.counterStorage[key]
	if found {
		S.counterStorage[key] += val
	} else {
		S.counterStorage[key] = val
	}
	return nil
}
func (S *MemStorage) GetCounter(ctx context.Context, key string) (int64, error) {
	S.counterMutex.Lock()
	defer S.counterMutex.Unlock()
	val, found := S.counterStorage[key]
	if !found {
		return 0, errors.New("counter not found")
	}
	return val, nil
}
func (S *MemStorage) GetAllCounter(ctx context.Context) (map[string]int64, error) {
	return S.counterStorage, nil
}

//func (S *storage) AddGauge(key string, val float64) error {
//	S.gaugeMutex.Lock()
//	defer S.gaugeMutex.Unlock()
//	_, found := S.gaugeStorage[key]
//	if found {
//		return errors.New("gauge exists")
//	}
//	S.gaugeStorage[key] = val
//	return nil
//}
//func (S *storage) AddCounter(key string, val int64) error {
//	S.counterMutex.Lock()
//	defer S.counterMutex.Unlock()
//	_, found := S.counterStorage[key]
//	if found {
//		return errors.New("counter exists")
//	}
//	S.counterStorage[key] = val
//	return nil
//}
//func (S *storage) GetGauge(key string, val float64) (float64, error) {
//	S.gaugeMutex.Lock()
//	defer S.gaugeMutex.Unlock()
//	val, found := S.gaugeStorage[key]
//	if !found {
//		return 0, errors.New("gauge not found")
//	}
//	return val, nil
//}
//func (S *storage) GetCounter(key string, val int64) (int64, error) {
//	S.counterMutex.Lock()
//	defer S.counterMutex.Unlock()
//	val, found := S.counterStorage[key]
//	if !found {
//		return 0, errors.New("counter not found")
//	}
//	return val, nil
//}
//func (S *storage) DeleteGauge(key string) error {
//	S.gaugeMutex.Lock()
//	defer S.gaugeMutex.Unlock()
//	delete(S.gaugeStorage, key)
//	return nil
//}
//func (S *storage) DeleteCounter(key string, val int64) error {
//	S.counterMutex.Lock()
//	defer S.counterMutex.Unlock()
//	delete(S.counterStorage, key)
//	return nil
//}
//func (S *storage) UpdateGauge(key string, val float64) error {
//	S.gaugeMutex.Lock()
//	defer S.gaugeMutex.Unlock()
//	_, found := S.gaugeStorage[key]
//	if !found {
//		return errors.New("gauge not found")
//	}
//	S.gaugeStorage[key] = val
//	return nil
//}
//func (S *storage) UpdateCounter(key string, val int64) error {
//	S.counterMutex.Lock()
//	defer S.counterMutex.Unlock()
//	_, found := S.counterStorage[key]
//	if found {
//		return errors.New("counter not found")
//	}
//	S.counterStorage[key] = val
//	return nil
//}
