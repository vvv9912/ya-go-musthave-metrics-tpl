package mw

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/handler"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/service"
	service_mock "github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/service/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

//func TestMw_MwLogger(t *testing.T) {
//	type want struct {
//		code        int
//		contentType string
//	}
//	host := "http://localhost:8080"
//	tests := []struct {
//		name string
//		//args args
//		path string
//		want want
//	}{
//		{name: "#1 positive test",
//			//args: ,
//			path: "/update/gauge/testGauge/123",
//			want: want{code: 200, contentType: "text/plain; charset=utf-8"},
//		},
//		{name: "#2 negative test",
//			//args: ,
//			path: "/update/counter/testGauge/453",
//			want: want{code: 200, contentType: "text/plain; charset=utf-8"},
//		},
//		{name: "#3 negative test",
//			//args: args{next: m.Middlware(m.MiddlwareGauge(http.HandlerFunc(handler.HandlerGauge)))},
//			path: "/update/counter/testGauge/none",
//			want: want{code: 400, contentType: "text/plain; charset=utf-8"},
//		},
//		{name: "#4 negative test",
//			//args: args{next: m.Middlware(m.MiddlwareGauge(http.HandlerFunc(handler.HandlerGauge)))},
//			path: "/update/gauge/testGauge/none",
//			want: want{code: 400, contentType: "text/plain; charset=utf-8"},
//		},
//		{name: "#5 negative test",
//			//args: args{next: m.Middlware(m.MiddlwareGauge(http.HandlerFunc(handler.HandlerGauge)))},
//			path: "/update/dsad/dsae/none",
//			want: want{code: 404, contentType: "text/plain; charset=utf-8"},
//		},
//		{name: "#5 negative test",
//			//args: args{next: m.Middlware(m.MiddlwareGauge(http.HandlerFunc(handler.HandlerGauge)))},
//			path: "/update/counter/testGauge/123/123",
//			want: want{code: 400, contentType: "text/plain; charset=utf-8"},
//		},
//	}
//	for _, test := range tests {
//		//t.Run(tt.name, func(t *testing.T) {
//		URL := host + test.path
//		request := httptest.NewRequest(http.MethodPost, URL, nil)
//		w := httptest.NewRecorder()
//
//		//m.Middlware(m.MiddlwareGauge(http.HandlerFunc(handler.HandlerGauge))).ServeHTTP(w, request)
//
//		res := w.Result()
//		assert.Equal(t, test.want.code, res.StatusCode)
//		defer res.Body.Close()
//		resBody, err := io.ReadAll(res.Body)
//		require.NoError(t, err)
//		t.Log("----------///\nContent type:", res.Header.Get("Content-Type"))
//		t.Log("----------///\nres body:", string(resBody))
//		assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
//
//	}
//}
//
//func TestMww(t *testing.T) {
//	type want struct {
//		code        int
//		contentType string
//	}
//	type args struct {
//		next http.Handler
//	}
//
//	counter := storage.NewCounterStorage()
//	gauge := storage.NewGaugeStorage()
//	m := Mw{
//		GaugeStorage:   gauge,
//		CounterStorage: counter,
//	}
//	tests := []struct {
//		name string
//		args args
//		URL  string
//		want want
//	}{
//		{name: "#1 positive test",
//			args: args{next: m.Middlware(m.MiddlwareGauge(http.HandlerFunc(handler.HandlerGauge)))},
//			URL:  "/update/gauge/someMetric/527",
//			want: want{code: 200, contentType: "text/plain; charset=utf-8"},
//		},
//		{name: "#2 negative test",
//			args: args{next: m.Middlware(m.MiddlwareGauge(http.HandlerFunc(handler.HandlerGauge)))},
//			URL:  "/update/gauge",
//			want: want{code: 404, contentType: "text/plain; charset=utf-8"},
//		},
//		{name: "#3 negative test",
//			args: args{next: m.Middlware(m.MiddlwareGauge(http.HandlerFunc(handler.HandlerGauge)))},
//			URL:  "/update/gauge/dasd",
//			want: want{code: 404, contentType: "text/plain; charset=utf-8"},
//		},
//		{name: "#4 negative test",
//			args: args{next: m.Middlware(m.MiddlwareGauge(http.HandlerFunc(handler.HandlerGauge)))},
//			URL:  "/update",
//			want: want{code: 404, contentType: "text/plain; charset=utf-8"},
//		},
//	}
//	for _, test := range tests {
//		//t.Run(tt.name, func(t *testing.T) {
//
//		request := httptest.NewRequest(http.MethodPost, test.URL, nil)
//		w := httptest.NewRecorder()
//
//		m.Middlware(m.MiddlwareGauge(http.HandlerFunc(handler.HandlerGauge))).ServeHTTP(w, request)
//		res := w.Result()
//		assert.Equal(t, test.want.code, res.StatusCode)
//		defer res.Body.Close()
//		resBody, err := io.ReadAll(res.Body)
//		require.NoError(t, err)
//		t.Log("----------///\nContent type:", res.Header.Get("Content-Type"))
//		t.Log("----------///\nres body:", string(resBody))
//		assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
//
//	}
//}

func TestMw_MiddlwareGauge(t *testing.T) {
	type mockGaugeStorage func(s *service_mock.MockGaugeStorager, name string, value float64)
	type mockCounterStorage func(s *service_mock.MockCounterStorager, name string, value uint64)
	type mockMetrics func(s *service_mock.MockMetrics)

	type want struct {
		code        int
		contentType string
	}

	tests := []struct {
		name string
		URL  string
		mockGaugeStorage
		mockCounterStorage
		mockMetrics
		want want
	}{
		{name: "#1 positive test",
			URL: "/update/gauge/someMetric/527",
			mockCounterStorage: mockCounterStorage(func(s *service_mock.MockCounterStorager, name string, value uint64) {
				s.EXPECT().UpdateCounter(gomock.Any(), gomock.Any()).Return(nil)
			}),
			mockGaugeStorage: mockGaugeStorage(func(s *service_mock.MockGaugeStorager, name string, value float64) {
				s.EXPECT().UpdateGauge(gomock.Any(), gomock.Any()).Return(nil)
			}),
			mockMetrics: mockMetrics(func(s *service_mock.MockMetrics) {
				s.EXPECT().SendMetricstoFile().Return(nil)
			}),
			want: want{code: 400, contentType: "text/plain; charset=utf-8"},
		},
		{name: "#2 negative test, URL err",
			URL: "/update/gauge/",
			mockCounterStorage: mockCounterStorage(func(s *service_mock.MockCounterStorager, name string, value uint64) {
				s.EXPECT().UpdateCounter(gomock.Any(), gomock.Any()).Return(nil)
			}),
			mockGaugeStorage: mockGaugeStorage(func(s *service_mock.MockGaugeStorager, name string, value float64) {
				s.EXPECT().UpdateGauge(gomock.Any(), gomock.Any()).Return(nil)
			}),
			mockMetrics: mockMetrics(func(s *service_mock.MockMetrics) {
				s.EXPECT().SendMetricstoFile().Return(nil)
			}),
			want: want{code: 400, contentType: "text/plain; charset=utf-8"},
		},
		{name: "#3 negative test, URL err",
			mockCounterStorage: mockCounterStorage(func(s *service_mock.MockCounterStorager, name string, value uint64) {
				s.EXPECT().UpdateCounter(gomock.Any(), gomock.Any()).Return(nil)
			}),
			mockGaugeStorage: mockGaugeStorage(func(s *service_mock.MockGaugeStorager, name string, value float64) {
				s.EXPECT().UpdateGauge(gomock.Any(), gomock.Any()).Return(nil)
			}),
			mockMetrics: mockMetrics(func(s *service_mock.MockMetrics) {
				s.EXPECT().SendMetricstoFile().Return(nil)
			}),
			URL:  "/update/gauge/dasd",
			want: want{code: 400, contentType: "text/plain; charset=utf-8"},
		},
		{name: "#4 negative test, URL err",
			mockCounterStorage: mockCounterStorage(func(s *service_mock.MockCounterStorager, name string, value uint64) {
				s.EXPECT().UpdateCounter(gomock.Any(), gomock.Any()).Return(nil)
			}),
			mockGaugeStorage: mockGaugeStorage(func(s *service_mock.MockGaugeStorager, name string, value float64) {
				s.EXPECT().UpdateGauge(gomock.Any(), gomock.Any()).Return(nil)
			}),
			mockMetrics: mockMetrics(func(s *service_mock.MockMetrics) {
				s.EXPECT().SendMetricstoFile().Return(nil)
			}),
			URL:  "/update",
			want: want{code: 400, contentType: "text/plain; charset=utf-8"},
		},
		{name: "#5 negative test, counter err",
			mockCounterStorage: mockCounterStorage(func(s *service_mock.MockCounterStorager, name string, value uint64) {
				s.EXPECT().UpdateCounter(gomock.Any(), gomock.Any()).Return(fmt.Errorf("some error"))
			}),
			mockGaugeStorage: mockGaugeStorage(func(s *service_mock.MockGaugeStorager, name string, value float64) {
				s.EXPECT().UpdateGauge(gomock.Any(), gomock.Any()).Return(nil)
			}),
			mockMetrics: mockMetrics(func(s *service_mock.MockMetrics) {
				s.EXPECT().SendMetricstoFile().Return(nil)
			}),
			URL:  "/update/gauge/someMetric/527",
			want: want{code: 400, contentType: "text/plain; charset=utf-8"},
		},
		{name: "#6 negative test, gauge err",
			mockCounterStorage: mockCounterStorage(func(s *service_mock.MockCounterStorager, name string, value uint64) {
				s.EXPECT().UpdateCounter(gomock.Any(), gomock.Any()).Return(nil)
			}),
			mockGaugeStorage: mockGaugeStorage(func(s *service_mock.MockGaugeStorager, name string, value float64) {
				s.EXPECT().UpdateGauge(gomock.Any(), gomock.Any()).Return(fmt.Errorf("some error"))
			}),
			mockMetrics: mockMetrics(func(s *service_mock.MockMetrics) {
				s.EXPECT().SendMetricstoFile().Return(nil)
			}),
			URL:  "/update/gauge/someMetric/527",
			want: want{code: 400, contentType: "text/plain; charset=utf-8"},
		},
		{name: "#7 negative test, sendMetrics err",
			mockCounterStorage: mockCounterStorage(func(s *service_mock.MockCounterStorager, name string, value uint64) {
				s.EXPECT().UpdateCounter(gomock.Any(), gomock.Any()).Return(nil)
			}),
			mockGaugeStorage: mockGaugeStorage(func(s *service_mock.MockGaugeStorager, name string, value float64) {
				s.EXPECT().UpdateGauge(gomock.Any(), gomock.Any()).Return(nil)
			}),
			mockMetrics: mockMetrics(func(s *service_mock.MockMetrics) {
				s.EXPECT().SendMetricstoFile().Return(fmt.Errorf("some error"))
			}),
			URL:  "/update/gauge/someMetric/527",
			want: want{code: 400, contentType: "text/plain; charset=utf-8"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Log(test)

			c := gomock.NewController(t)
			defer c.Finish()

			CounterStorageMock := service_mock.NewMockCounterStorager(c)
			GaugeStorageMock := service_mock.NewMockGaugeStorager(c)
			MetricsMock := service_mock.NewMockMetrics(c)

			serviceMock := &service.Service{
				Metrics:         MetricsMock,
				CounterStorager: CounterStorageMock,
				GaugeStorager:   GaugeStorageMock,
				Notifier:        nil,
			}

			mw := Mw{Service: serviceMock}

			request := httptest.NewRequest(http.MethodPost, test.URL, nil)
			w := httptest.NewRecorder()
			mw.MiddlewareType(http.HandlerFunc(handler.HandlerSucess)).ServeHTTP(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			t.Log("----------///\nContent type:", res.Header.Get("Content-Type"))
			t.Log("----------///\nres body:", string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

		})

	}
}

func TestMw_MiddlewareCounter(t *testing.T) {
	type mockGaugeStorage func(s *service_mock.MockGaugeStorager, name string, value float64)
	type mockCounterStorage func(s *service_mock.MockCounterStorager, name string, value uint64)
	type mockMetrics func(s *service_mock.MockMetrics)

	type want struct {
		code        int
		contentType string
	}

	tests := []struct {
		name string
		URL  string
		mockGaugeStorage
		mockCounterStorage
		mockMetrics
		want want
	}{
		{name: "#1 positive test",
			URL: "/update/counter/someMetric/527",
			mockCounterStorage: mockCounterStorage(func(s *service_mock.MockCounterStorager, name string, value uint64) {
				s.EXPECT().UpdateCounter(gomock.Any(), gomock.Any()).Return(nil)
			}),
			mockGaugeStorage: mockGaugeStorage(func(s *service_mock.MockGaugeStorager, name string, value float64) {
				s.EXPECT().UpdateGauge(gomock.Any(), gomock.Any()).Return(nil)
			}),
			mockMetrics: mockMetrics(func(s *service_mock.MockMetrics) {
				s.EXPECT().SendMetricstoFile().Return(nil)
			}),
			want: want{code: 400, contentType: "text/plain; charset=utf-8"},
		},
		{name: "#2 negative test, URL err",
			URL: "/update/counter/",
			mockCounterStorage: mockCounterStorage(func(s *service_mock.MockCounterStorager, name string, value uint64) {
				s.EXPECT().UpdateCounter(gomock.Any(), gomock.Any()).Return(nil)
			}),
			mockGaugeStorage: mockGaugeStorage(func(s *service_mock.MockGaugeStorager, name string, value float64) {
				s.EXPECT().UpdateGauge(gomock.Any(), gomock.Any()).Return(nil)
			}),
			mockMetrics: mockMetrics(func(s *service_mock.MockMetrics) {
				s.EXPECT().SendMetricstoFile().Return(nil)
			}),
			want: want{code: 400, contentType: "text/plain; charset=utf-8"},
		},
		{name: "#3 negative test, URL err",
			mockCounterStorage: mockCounterStorage(func(s *service_mock.MockCounterStorager, name string, value uint64) {
				s.EXPECT().UpdateCounter(gomock.Any(), gomock.Any()).Return(nil)
			}),
			mockGaugeStorage: mockGaugeStorage(func(s *service_mock.MockGaugeStorager, name string, value float64) {
				s.EXPECT().UpdateGauge(gomock.Any(), gomock.Any()).Return(nil)
			}),
			mockMetrics: mockMetrics(func(s *service_mock.MockMetrics) {
				s.EXPECT().SendMetricstoFile().Return(nil)
			}),
			URL:  "/update/counter/dasd",
			want: want{code: 400, contentType: "text/plain; charset=utf-8"},
		},
		{name: "#4 negative test, URL err",
			mockCounterStorage: mockCounterStorage(func(s *service_mock.MockCounterStorager, name string, value uint64) {
				s.EXPECT().UpdateCounter(gomock.Any(), gomock.Any()).Return(nil)
			}),
			mockGaugeStorage: mockGaugeStorage(func(s *service_mock.MockGaugeStorager, name string, value float64) {
				s.EXPECT().UpdateGauge(gomock.Any(), gomock.Any()).Return(nil)
			}),
			mockMetrics: mockMetrics(func(s *service_mock.MockMetrics) {
				s.EXPECT().SendMetricstoFile().Return(nil)
			}),
			URL:  "/update",
			want: want{code: 400, contentType: "text/plain; charset=utf-8"},
		},
		{name: "#5 negative test, counter err",
			mockCounterStorage: mockCounterStorage(func(s *service_mock.MockCounterStorager, name string, value uint64) {
				s.EXPECT().UpdateCounter(gomock.Any(), gomock.Any()).Return(fmt.Errorf("some error"))
			}),
			mockGaugeStorage: mockGaugeStorage(func(s *service_mock.MockGaugeStorager, name string, value float64) {
				s.EXPECT().UpdateGauge(gomock.Any(), gomock.Any()).Return(nil)
			}),
			mockMetrics: mockMetrics(func(s *service_mock.MockMetrics) {
				s.EXPECT().SendMetricstoFile().Return(nil)
			}),
			URL:  "/update/counter/someMetric/527",
			want: want{code: 400, contentType: "text/plain; charset=utf-8"},
		},
		{name: "#6 negative test, gauge err",
			mockCounterStorage: mockCounterStorage(func(s *service_mock.MockCounterStorager, name string, value uint64) {
				s.EXPECT().UpdateCounter(gomock.Any(), gomock.Any()).Return(nil)
			}),
			mockGaugeStorage: mockGaugeStorage(func(s *service_mock.MockGaugeStorager, name string, value float64) {
				s.EXPECT().UpdateGauge(gomock.Any(), gomock.Any()).Return(fmt.Errorf("some error"))
			}),
			mockMetrics: mockMetrics(func(s *service_mock.MockMetrics) {
				s.EXPECT().SendMetricstoFile().Return(nil)
			}),
			URL:  "/update/counter/someMetric/527",
			want: want{code: 400, contentType: "text/plain; charset=utf-8"},
		},
		{name: "#7 negative test, sendMetrics err",
			mockCounterStorage: mockCounterStorage(func(s *service_mock.MockCounterStorager, name string, value uint64) {
				s.EXPECT().UpdateCounter(gomock.Any(), gomock.Any()).Return(nil)
			}),
			mockGaugeStorage: mockGaugeStorage(func(s *service_mock.MockGaugeStorager, name string, value float64) {
				s.EXPECT().UpdateGauge(gomock.Any(), gomock.Any()).Return(nil)
			}),
			mockMetrics: mockMetrics(func(s *service_mock.MockMetrics) {
				s.EXPECT().SendMetricstoFile().Return(fmt.Errorf("some error"))
			}),
			URL:  "/update/counter/someMetric/527",
			want: want{code: 400, contentType: "text/plain; charset=utf-8"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Log(test)

			c := gomock.NewController(t)
			defer c.Finish()

			CounterStorageMock := service_mock.NewMockCounterStorager(c)
			GaugeStorageMock := service_mock.NewMockGaugeStorager(c)
			MetricsMock := service_mock.NewMockMetrics(c)

			serviceMock := &service.Service{
				Metrics:         MetricsMock,
				CounterStorager: CounterStorageMock,
				GaugeStorager:   GaugeStorageMock,
				Notifier:        nil,
			}

			mw := Mw{Service: serviceMock}

			request := httptest.NewRequest(http.MethodPost, test.URL, nil)
			w := httptest.NewRecorder()
			mw.MiddlewareType(http.HandlerFunc(handler.HandlerSucess)).ServeHTTP(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			t.Log("----------///\nContent type:", res.Header.Get("Content-Type"))
			t.Log("----------///\nres body:", string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

		})

	}
}

func Test_supportEncodingType(t *testing.T) {
	type args struct {
		accpetEncoding    map[string]struct{}
		acceptEncodingReq string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "#1 positive test",
			args: args{
				accpetEncoding:    map[string]struct{}{"gzip": {}},
				acceptEncodingReq: "gzip",
			},
			want: true,
		},

		{
			name: "#2 negative test",
			args: args{
				accpetEncoding:    map[string]struct{}{"gzip": {}},
				acceptEncodingReq: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.name)
			got := supportEnodingType(tt.args.accpetEncoding, tt.args.acceptEncodingReq)
			assert.Equal(t, tt.want, got, "supportEnodingType, want: %v, got: %v", tt.want, got)
		})
	}
}
func Test_supportContentType(t *testing.T) {
	type args struct {
		acceptType    map[string]struct{}
		acceptTypeReq string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "#1 positive test",
			args: args{
				acceptType:    map[string]struct{}{"application/json": {}},
				acceptTypeReq: "application/json",
			},
			want: true,
		},
		{
			name: "#2 positive test",
			args: args{
				acceptType: map[string]struct{}{
					"application/json": struct{}{},
					"text/html":        struct{}{},
					"html/text":        struct{}{},
				},
				acceptTypeReq: "text/html",
			},
			want: true,
		},
		{
			name: "#3 positive test",
			args: args{
				acceptType: map[string]struct{}{
					"application/json": struct{}{},
					"text/html":        struct{}{},
					"html/text":        struct{}{},
				},
				acceptTypeReq: "*/*",
			},
			want: true,
		},
		{
			name: "#4 positive test",
			args: args{
				acceptType: map[string]struct{}{
					"application/json": struct{}{},
					"text/html":        struct{}{},
					"html/text":        struct{}{},
				},
				acceptTypeReq: "application/json; charset=utf-8",
			},
			want: true,
		},
		{
			name: "#5 negative test",
			args: args{
				acceptType: map[string]struct{}{
					"application/json": struct{}{},
					"text/html":        struct{}{},
					"html/text":        struct{}{},
				},
				acceptTypeReq: "text/plain",
			},
			want: false,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.name)
			got := supportAcceptType(tt.args.acceptType, tt.args.acceptTypeReq)
			assert.Equal(t, tt.want, got, "supportType: want: %v, got: %v", tt.want, got)
		})
	}
}
