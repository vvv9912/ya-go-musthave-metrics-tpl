package mw

//
//import (
//	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/storage"
//	"net/http"
//	"reflect"
//	"testing"
//)
//
//func TestMw_Middlware(t *testing.T) {
//	type fields struct {
//		GaugeStorage   storage.GaugeStorager
//		CounterStorage storage.CounterStorager
//	}
//	type args struct {
//		next http.Handler
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		want   http.Handler
//	}{
//		{
//			name: "positive test #1",
//		},
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			m := &Mw{
//				GaugeStorage:   tt.fields.GaugeStorage,
//				CounterStorage: tt.fields.CounterStorage,
//			}
//			if got := m.Middlware(tt.args.next); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Middlware() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestMw_MiddlwareCounter(t *testing.T) {
//	type fields struct {
//		GaugeStorage   storage.GaugeStorager
//		CounterStorage storage.CounterStorager
//	}
//	type args struct {
//		next http.Handler
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		want   http.Handler
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			m := &Mw{
//				GaugeStorage:   tt.fields.GaugeStorage,
//				CounterStorage: tt.fields.CounterStorage,
//			}
//			if got := m.MiddlwareCounter(tt.args.next); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("MiddlwareCounter() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestMw_MiddlwareGauge(t *testing.T) {
//	type fields struct {
//		GaugeStorage   storage.GaugeStorager
//		CounterStorage storage.CounterStorager
//	}
//	type args struct {
//		next http.Handler
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		want   http.Handler
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			m := &Mw{
//				GaugeStorage:   tt.fields.GaugeStorage,
//				CounterStorage: tt.fields.CounterStorage,
//			}
//			if got := m.MiddlwareGauge(tt.args.next); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("MiddlwareGauge() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
