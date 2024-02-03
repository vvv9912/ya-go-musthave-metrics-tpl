package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"
	"sync"
)

type MemStorage struct {
	gaugeStorage   map[string]float64 //todo переделать под map[string]string, в соотв. с agent
	counterStorage map[string]int64
	gaugeMutex     sync.Mutex
	counterMutex   sync.Mutex
}

func NewStorage() *MemStorage {
	return &MemStorage{
		gaugeStorage:   make(map[string]float64),
		counterStorage: make(map[string]int64),
	}
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

func (S *MemStorage) UpdateMetricsBatch(ctx context.Context, metrics []model.Metrics) error {
	return fmt.Errorf("method updateMetricsBatch for metrics are not implemented")
}
func (S *MemStorage) Ping(ctx context.Context) error {
	return fmt.Errorf("method ping for metrics are not implemented")
}
