package mw

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/handler"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

//	func TestMw_MiddlwareGauge(t *testing.T) {
//		type fields struct {
//			GaugeStorage   storage.GaugeStorager
//			CounterStorage storage.CounterStorager
//		}
//		type want struct {
//			code int
//			//	response string
//			contentType string
//		}
//		type args struct {
//			next http.Handler
//		}
//		f := fields{}
//		m := Mw{
//			GaugeStorage:   f.GaugeStorage,
//			CounterStorage: f.CounterStorage,
//		}
//		tests := []struct {
//			name   string
//			fields fields
//			args   args
//			URL    string
//			want   want
//		}{
//			{name: "first test",
//				args: args{next: m.Middlware(m.MiddlwareGauge(http.HandlerFunc(handler.HandlerGauge)))},
//				URL:  "/update/counter/someMetric/527",
//				want: want{code: 200, contentType: "text/plain; charset=utf-8"},
//			},
//			// TODO: Add test cases.
//		}
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				m := &Mw{
//					GaugeStorage:   tt.fields.GaugeStorage,
//					CounterStorage: tt.fields.CounterStorage,
//				}
//				if got := m.MiddlwareGauge(tt.args .next); !reflect.DeepEqual(got, tt.want) {
//					t.Errorf("MiddlwareGauge() = %v, want %v", got, tt.want)
//				}
//			})
//		}
//	}
func TestMw_MiddlwareGauge(t *testing.T) {

	type want struct {
		code        int
		contentType string
	}
	type args struct {
		next http.Handler
	}

	counter := storage.NewCounterStorage()
	gauge := storage.NewGaugeStorage()
	m := Mw{
		GaugeStorage:   gauge,
		CounterStorage: counter,
	}
	tests := []struct {
		name string
		args args
		URL  string
		want want
	}{
		{name: "#1 positive test",
			args: args{next: m.Middlware(m.MiddlwareGauge(http.HandlerFunc(handler.HandlerGauge)))},
			URL:  "/update/gauge/someMetric/527",
			want: want{code: 200, contentType: "text/plain; charset=utf-8"},
		},
		{name: "#2 negative test",
			args: args{next: m.Middlware(m.MiddlwareGauge(http.HandlerFunc(handler.HandlerGauge)))},
			URL:  "/update/gauge",
			want: want{code: 404, contentType: "text/plain; charset=utf-8"},
		},
		{name: "#3 negative test",
			args: args{next: m.Middlware(m.MiddlwareGauge(http.HandlerFunc(handler.HandlerGauge)))},
			URL:  "/update/gauge/dasd",
			want: want{code: 404, contentType: "text/plain; charset=utf-8"},
		},
		{name: "#4 negative test",
			args: args{next: m.Middlware(m.MiddlwareGauge(http.HandlerFunc(handler.HandlerGauge)))},
			URL:  "/update",
			want: want{code: 404, contentType: "text/plain; charset=utf-8"},
		},
	}
	for _, test := range tests {
		//t.Run(tt.name, func(t *testing.T) {

		request := httptest.NewRequest(http.MethodPost, test.URL, nil)
		w := httptest.NewRecorder()

		m.Middlware(m.MiddlwareGauge(http.HandlerFunc(handler.HandlerGauge))).ServeHTTP(w, request)
		res := w.Result()
		assert.Equal(t, test.want.code, res.StatusCode)
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		t.Log("----------///\nContent type:", res.Header.Get("Content-Type"))
		t.Log("----------///\nres body:", string(resBody))
		assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

	}
}

func TestMw_MiddlwareCounter(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	type args struct {
		next http.Handler
	}

	counter := storage.NewCounterStorage()
	gauge := storage.NewGaugeStorage()
	m := Mw{
		GaugeStorage:   gauge,
		CounterStorage: counter,
	}
	tests := []struct {
		name string
		args args
		URL  string
		want want
	}{
		{name: "#1 positive test",
			args: args{next: m.Middlware(m.MiddlwareGauge(http.HandlerFunc(handler.HandlerGauge)))},
			URL:  "/update/counter/someMetric/527",
			want: want{code: 200, contentType: "text/plain; charset=utf-8"},
		},
		{name: "#2 negative test",
			args: args{next: m.Middlware(m.MiddlwareGauge(http.HandlerFunc(handler.HandlerGauge)))},
			URL:  "/update/counter",
			want: want{code: 404, contentType: "text/plain; charset=utf-8"},
		},
		{name: "#3 negative test",
			args: args{next: m.Middlware(m.MiddlwareGauge(http.HandlerFunc(handler.HandlerGauge)))},
			URL:  "/update/counter/dasd",
			want: want{code: 404, contentType: "text/plain; charset=utf-8"},
		},
		{name: "#4 negative test",
			args: args{next: m.Middlware(m.MiddlwareGauge(http.HandlerFunc(handler.HandlerGauge)))},
			URL:  "/update",
			want: want{code: 404, contentType: "text/plain; charset=utf-8"},
		},
	}
	for _, test := range tests {
		//t.Run(tt.name, func(t *testing.T) {

		request := httptest.NewRequest(http.MethodPost, test.URL, nil)
		w := httptest.NewRecorder()

		m.Middlware(m.MiddlwareCounter(http.HandlerFunc(handler.HandlerGauge))).ServeHTTP(w, request)
		res := w.Result()
		assert.Equal(t, test.want.code, res.StatusCode)
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		t.Log("----------///\nContent type:", res.Header.Get("Content-Type"))
		t.Log("----------///\nres body:", string(resBody))
		assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

	}
}
