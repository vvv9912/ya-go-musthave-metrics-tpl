package notifier

import (
	"context"
	"testing"
	"time"
)

func TestNotifier_StartNotifyCron(t *testing.T) {
	type fields struct {
		EventsMetric  EventsMetric
		PostRequester PostRequester
		TimerUpdate   time.Duration
		TimerSend     time.Duration
		URL           string
	}
	type args struct {
		ctx context.Context
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
			n := &Notifier{
				EventsMetric:  tt.fields.EventsMetric,
				PostRequester: tt.fields.PostRequester,
				TimerUpdate:   tt.fields.TimerUpdate,
				TimerSend:     tt.fields.TimerSend,
				URL:           tt.fields.URL,
			}
			if err := n.StartNotifyCron(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("StartNotifyCron() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
