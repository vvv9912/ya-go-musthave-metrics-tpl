package metrics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetrics_UpdateMetricsCounter(t *testing.T) {

	tests := []struct {
		name string
		//values map[string]uint64
		n    int
		want map[string]uint64
	}{
		{
			name: "test #1",
			n:    3,

			want: map[string]uint64{
				"PollCount": 3,
			},
		},
		{
			name: "test #2",
			n:    100,

			want: map[string]uint64{
				"PollCount": 100,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cout := NewMetriсs()
			var c uint64
			var err error
			for i := 0; i < tt.n; i++ {
				c, err = cout.UpdateMetricsCounter()
				if err != nil {
					t.Error(err)
				}
			}
			assert.Equal(t, c, tt.want["PollCount"])

		})
	}
}
