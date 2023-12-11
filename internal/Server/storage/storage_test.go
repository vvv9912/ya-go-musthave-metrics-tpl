package storage

import (
	"reflect"
	"sync"
	"testing"
)

func TestMemStorage_GetAllCounter(t *testing.T) {
	type fields struct {
		gaugeStorage   map[string]float64
		counterStorage map[string]uint64
		gaugeMutex     sync.Mutex
		counterMutex   sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			S := &MemStorage{
				gaugeStorage:   tt.fields.gaugeStorage,
				counterStorage: tt.fields.counterStorage,
				gaugeMutex:     tt.fields.gaugeMutex,
				counterMutex:   tt.fields.counterMutex,
			}
			if got := S.GetAllCounter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_GetAllGauge(t *testing.T) {
	type fields struct {
		gaugeStorage   map[string]float64
		counterStorage map[string]uint64
		gaugeMutex     sync.Mutex
		counterMutex   sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			S := &MemStorage{
				gaugeStorage:   tt.fields.gaugeStorage,
				counterStorage: tt.fields.counterStorage,
				gaugeMutex:     tt.fields.gaugeMutex,
				counterMutex:   tt.fields.counterMutex,
			}
			if got := S.GetAllGauge(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllGauge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_GetCounter(t *testing.T) {
	type fields struct {
		gaugeStorage   map[string]float64
		counterStorage map[string]uint64
		gaugeMutex     sync.Mutex
		counterMutex   sync.Mutex
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uint64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			S := &MemStorage{
				gaugeStorage:   tt.fields.gaugeStorage,
				counterStorage: tt.fields.counterStorage,
				gaugeMutex:     tt.fields.gaugeMutex,
				counterMutex:   tt.fields.counterMutex,
			}
			got, err := S.GetCounter(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCounter() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
	type fields struct {
		gaugeStorage   map[string]float64
		counterStorage map[string]uint64
		gaugeMutex     sync.Mutex
		counterMutex   sync.Mutex
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			S := &MemStorage{
				gaugeStorage:   tt.fields.gaugeStorage,
				counterStorage: tt.fields.counterStorage,
				gaugeMutex:     tt.fields.gaugeMutex,
				counterMutex:   tt.fields.counterMutex,
			}
			got, err := S.GetGauge(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGauge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetGauge() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_UpdateCounter(t *testing.T) {
	type fields struct {
		gaugeStorage   map[string]float64
		counterStorage map[string]uint64
		gaugeMutex     sync.Mutex
		counterMutex   sync.Mutex
	}
	type args struct {
		key string
		val uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			S := &MemStorage{
				gaugeStorage:   tt.fields.gaugeStorage,
				counterStorage: tt.fields.counterStorage,
				gaugeMutex:     tt.fields.gaugeMutex,
				counterMutex:   tt.fields.counterMutex,
			}
			if err := S.UpdateCounter(tt.args.key, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("UpdateCounter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemStorage_UpdateGauge(t *testing.T) {
	type fields struct {
		gaugeStorage   map[string]float64
		counterStorage map[string]uint64
		gaugeMutex     sync.Mutex
		counterMutex   sync.Mutex
	}
	type args struct {
		key string
		val float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			S := &MemStorage{
				gaugeStorage:   tt.fields.gaugeStorage,
				counterStorage: tt.fields.counterStorage,
				gaugeMutex:     tt.fields.gaugeMutex,
				counterMutex:   tt.fields.counterMutex,
			}
			if err := S.UpdateGauge(tt.args.key, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("UpdateGauge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewCounterStorage(t *testing.T) {
	tests := []struct {
		name string
		want CounterStorager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCounterStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCounterStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGaugeStorage(t *testing.T) {
	tests := []struct {
		name string
		want GaugeStorager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGaugeStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGaugeStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}
