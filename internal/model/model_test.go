package model

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMetrics(t *testing.T) {
	val := float64(10.5)

	metrics := Metrics{
		ID:    "metric1",
		MType: "gauge",
		Value: &val,
	}

	// Тестирование значения
	require.Equal(t, metrics.ID, "metric1")
	require.Equal(t, metrics.MType, "gauge")
	require.Equal(t, metrics.Value, &val)
	require.Equal(t, metrics.Delta, (*int64)(nil))

}

func Example_metrics() {
	val := float64(10.5)

	metrics := Metrics{
		ID:    "metric1",
		MType: "gauge",
		Value: &val,
	}

	fmt.Println(metrics.ID)
	fmt.Println(metrics.MType)
	fmt.Println(metrics.Delta)
	fmt.Println(metrics.Value)

	// OUT:
	// metric1
	// gauge
	// (*int64)(nil)
	// 10.5
}
