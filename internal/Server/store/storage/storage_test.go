package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCounterStorage(t *testing.T) {
	tests := []struct {
		name   string
		values map[string]int64
		want   map[string]int64
	}{
		{
			name: "test #1",
			values: map[string]int64{
				"1": 0,
				"2": 3,
				"4": 12441324,
			},
			want: map[string]int64{
				"1": 0,
				"2": 3,
				"4": 12441324,
			},
		},
		{
			name: "test #2 ",
			values: map[string]int64{
				"1": 0,
				"4": 45645687879,
				"2": 12441324,
			},
			want: map[string]int64{
				"1": 0,
				"2": 12441324,
				"4": 45645687879,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			cout := NewStorage()
			for key, value := range tt.values {
				err := cout.UpdateCounter(ctx, key, value)
				if err != nil {
					t.Error(err)
				}
			}
			c, err := cout.GetAllCounter(ctx)
			require.NoError(t, err)
			require.Equal(t, c, tt.want)

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
			cout := NewStorage()
			ctx := context.TODO()
			for key, value := range tt.values {
				err := cout.UpdateGauge(ctx, key, value)
				if err != nil {
					t.Error(err)
				}
			}
			c, err := cout.GetAllGauge(ctx)
			require.NoError(t, err)
			require.Equal(t, c, tt.want)

		})
	}
}
