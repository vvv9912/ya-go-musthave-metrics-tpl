package handler

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/service"
	service_mock "github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/service/mock"
	storage "github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/store/storage"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"
	"io"
	"net/http"
	"net/http/httptest"

	"testing"
)

func TestHandlerGauge(t *testing.T) {
	type want struct {
		code int
		//	response string
		contentType string
	}
	tests := []struct {
		name string
		want want
		url  string
	}{
		{
			name: "positive test #1",
			want: want{
				code: 200,
				//response:
				contentType: "text/plain; charset=utf-8",
			},
			url: "/update/counter/someMetric/527",
		},
	}
	for _, test := range tests {
		//t.Run(test.name, func(t *testing.T) {
		t.Log(test.name)
		t.Log(test.url)
		t.Log(test)
		request := httptest.NewRequest(http.MethodPost, test.url, nil)
		w := httptest.NewRecorder()
		HandlerSucess(w, request)

		res := w.Result()
		assert.Equal(t, test.want.code, res.StatusCode)

		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		t.Log("----------///\nAll content type:")
		for _, n := range res.Header {
			t.Log(n, "\n")
		}
		t.Log("----------///\nContent type:", res.Header.Get("Content-Type"))
		t.Log("----------///\nres body:", string(resBody))
		assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

	}
}
func TestHandler_HandlerGetJSON(t *testing.T) {
	type mockBehavior func(s *service_mock.MockMetrics, metric model.Metrics)

	val := 0.123

	forTestMetric := model.Metrics{
		ID:    "123",
		MType: "gauge",
		Delta: nil,
		Value: &val,
	}

	jsonforTestMetric, err := json.Marshal(forTestMetric)
	if err != nil {
		t.Fatal(err)
	}
	testTable := []struct {
		name                string
		inputBody           string
		inputMetric         model.Metrics
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "positive test #1",
			inputBody:   string(jsonforTestMetric),
			inputMetric: forTestMetric,
			mockBehavior: func(s *service_mock.MockMetrics, metric model.Metrics) {
				s.EXPECT().GetMetrics(gomock.Any(), gomock.Eq(metric)).Return(model.Metrics{
					ID:    "123",
					MType: "gauge",
					Delta: nil,
					Value: &val,
				}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: string(jsonforTestMetric),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			// создаем repository_mock проекта
			metrics := service_mock.NewMockMetrics(c)

			serviceMock := &service.Service{
				Metrics: metrics,
			}

			handler := Handler{Service: serviceMock}
			testCase.mockBehavior(metrics, testCase.inputMetric)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("get", "/", bytes.NewBufferString(testCase.inputBody))

			handler.HandlerGetJSON(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)

		})
	}
}

func TestHandlerSucess(t *testing.T) {
	type args struct {
		expectedRequest int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive test #1",
			args: args{expectedRequest: http.StatusOK},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("get", "/", bytes.NewBufferString(""))
			HandlerSucess(w, req)
			assert.Equal(t, w.Code, tt.args.expectedRequest)
		})
	}
}
func ExampleHandlerSucess() {
	http.HandleFunc("/success", HandlerSucess)
	http.ListenAndServe(":8080", nil)
}

func ExampleHandlerGetCounter() {
	http.HandleFunc("/GetCounter", HandlerGetCounter)
	http.ListenAndServe(":8080", nil)
}

func ExampleHandlerGetGauge() {
	http.HandleFunc("/GetGauge", HandlerGetGauge)
	http.ListenAndServe(":8080", nil)
}

func ExampleHandlerGetMetrics() {
	s := storage.NewStorage()

	http.HandleFunc("/GetMetrics", HandlerGetMetrics(s))
	http.ListenAndServe(":8080", nil)
}

func ExampleHandlerPostJSON() {
	s := service.Service{}
	h := Handler{Service: &s}

	http.HandleFunc("/PostJson", h.HandlerPostJSON)
	http.ListenAndServe(":8080", nil)
}

func ExampleHandlerGetJSON() {
	s := service.Service{}
	h := Handler{Service: &s}

	http.HandleFunc("/GetJson", h.HandlerGetJSON)
	http.ListenAndServe(":8080", nil)
}

func ExampleHandlerGauge() {
	s := service.Service{}
	h := Handler{Service: &s}

	http.HandleFunc("/Gauge", h.HandlerGauge)
	http.ListenAndServe(":8080", nil)
}

func ExampleHandlerPingDatabase() {
	s := service.Service{}
	h := Handler{Service: &s}

	http.HandleFunc("/ping", h.HandlerPingDatabase)
	http.ListenAndServe(":8080", nil)
}
func ExampleHandlerPostBatched() {
	s := service.Service{}
	h := Handler{Service: &s}

	http.HandleFunc("/PostBatched", h.HandlerPostBatched)
	http.ListenAndServe(":8080", nil)
}
