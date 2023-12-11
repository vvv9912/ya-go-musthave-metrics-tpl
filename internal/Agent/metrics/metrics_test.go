package metrics

import (
	"reflect"
	"testing"
)

func TestMetrics_UpdateMetricsGauge(t *testing.T) {
	type fields struct {
		MetricsGauge   map[string]string
		MetricsCounter map[string]uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   *map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metrics{
				MetricsGauge:   tt.fields.MetricsGauge,
				MetricsCounter: tt.fields.MetricsCounter,
			}
			if got := m.UpdateMetricsGauge(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateMetricsGauge() = %v, want %v", got, tt.want)
			}
		})
	}
}
