package handler

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	service_mock "github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/service/mock"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/service"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"
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
		//{
		//	name: "neg test #2",
		//	want: want{
		//		code: 400,
		//		//response:
		//		contentType: "",
		//	},
		//	url: "/update/cs",
		//},
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

		//	})
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
				s.EXPECT().GetMetrics(gomock.Eq(metric)).Return(model.Metrics{
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
			// создаем mock проекта
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
