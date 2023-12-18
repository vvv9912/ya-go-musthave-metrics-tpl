package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCounterStorage(t *testing.T) {
	tests := []struct {
		name   string
		values map[string]uint64
		want   map[string]uint64
	}{
		{
			name: "test #1",
			values: map[string]uint64{
				"1": 0,
				"2": 3,
				"4": 12441324,
			},
			want: map[string]uint64{
				"1": 0,
				"2": 3,
				"4": 12441324,
			},
		},
		{
			name: "test #2 ",
			values: map[string]uint64{
				"1": 0,
				"4": 45645687879,
				"2": 12441324,
			},
			want: map[string]uint64{
				"1": 0,
				"2": 12441324,
				"4": 45645687879,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cout := NewCounterStorage()
			for key, value := range tt.values {
				err := cout.UpdateCounter(key, value)
				if err != nil {
					t.Error(err)
				}
			}
			c := cout.GetAllCounter()
			assert.Equal(t, c, tt.want)
			//if got := NewCounterStorage(); !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("NewCounterStorage() = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestNewGaugeStorage(t *testing.T) {
	tests := []struct {
		name   string
		values map[string]float64
		want   map[string]float64
	}{
		{
			name: "test #1",
			values: map[string]float64{
				"1": -0,
				"2": 1797.969,
				"4": -889.4,
			},
			want: map[string]float64{
				"1": -0,
				"2": 1797.969,
				"4": -889.4,
			},
		},
		{
			name: "test #2 ",
			values: map[string]float64{
				"1": -0,
				"4": -889.4,
				"2": 1797.969,
			},
			want: map[string]float64{
				"1": -0,
				"2": 1797.969,
				"4": -889.4,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cout := NewGaugeStorage()
			for key, value := range tt.values {
				err := cout.UpdateGauge(key, value)
				if err != nil {
					t.Error(err)
				}
			}
			c := cout.GetAllGauge()
			assert.Equal(t, c, tt.want)
			//if got := NewCounterStorage(); !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("NewCounterStorage() = %v, want %v", got, tt.want)
			//}
		})
	}
}
