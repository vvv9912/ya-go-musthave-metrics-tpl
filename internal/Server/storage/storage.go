package storage

import (
	"errors"
	"sync"
)

type GaugeStorager interface {
	//AddGauge(key string, val float64) error
	//GetGauge(key string, val float64) (float64, error)
	//DeleteGauge(key string) error
	//UpdateGauge(key string, val float64) error
	UpdateGauge(key string, val float64) error
	GetGauge(key string) (float64, error)
	GetAllGauge() map[string]float64
}

type CounterStorager interface {
	//AddCounter(key string, val int64) error
	//GetCounter(key string, val int64) (int64, error)
	//DeleteCounter(key string, val int64) error
	//UpdateCounter(key string, val int64) error
	UpdateCounter(key string, val uint64) error
	GetCounter(key string) (uint64, error)
	GetAllCounter() map[string]uint64
}
type MemStorage struct {
	gaugeStorage   map[string]float64 //todo переделать под map[string]string, в соотв. с agent
	counterStorage map[string]uint64
	gaugeMutex     sync.Mutex
	counterMutex   sync.Mutex
}

func NewGaugeStorage() GaugeStorager {
	gaugeStorage := make(map[string]float64)
	return &MemStorage{gaugeStorage: gaugeStorage}
}
func NewCounterStorage() CounterStorager {
	counterStorage := make(map[string]uint64)
	return &MemStorage{counterStorage: counterStorage}
}

/*
Тип gauge, float64 — новое значение должно замещать предыдущее.
Тип counter, int64 — новое значение должно добавляться к предыдущему,
если какое-то значение уже было известно серверу.
*/
func (S *MemStorage) UpdateGauge(key string, val float64) error {
	S.gaugeMutex.Lock()
	defer S.gaugeMutex.Unlock()
	S.gaugeStorage[key] = val
	return nil
}
func (S *MemStorage) GetGauge(key string) (float64, error) {
	S.gaugeMutex.Lock()
	defer S.gaugeMutex.Unlock()
	val, found := S.gaugeStorage[key]
	if !found {
		return 0, errors.New("gauge not found")
	}
	return val, nil
}
func (S *MemStorage) GetAllGauge() map[string]float64 {
	return S.gaugeStorage
}
func (S *MemStorage) UpdateCounter(key string, val uint64) error {
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
func (S *MemStorage) GetCounter(key string) (uint64, error) {
	S.counterMutex.Lock()
	defer S.counterMutex.Unlock()
	val, found := S.counterStorage[key]
	if !found {
		return 0, errors.New("counter not found")
	}
	return val, nil
}
func (S *MemStorage) GetAllCounter() map[string]uint64 {
	return S.counterStorage
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
