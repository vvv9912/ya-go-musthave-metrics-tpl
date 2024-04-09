package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/service"
	service_mock "github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/service/mock"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/store/repo_mock"
	storage "github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/store/storage"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/typeconst"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

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
func TestHandlerGetCounter(t *testing.T) {
	type want struct {
		code        int
		contentType string
		val         []byte
	}
	tests := []struct {
		name string
		want want
		url  string
		val  string
	}{
		{
			name: "positive test #1",
			want: want{
				code:        200,
				contentType: "text/plain; charset=utf-8",
				val:         []byte("527"),
			},
			url: "/update/counter/someMetric/527",
			val: "527",
		},
	}
	for _, test := range tests {
		t.Log(test.name)
		t.Log(test.url)
		t.Log(test)
		request := httptest.NewRequest(http.MethodGet, test.url, nil)
		w := httptest.NewRecorder()

		ctx := context.WithValue(request.Context(), typeconst.UserIDContextKey, test.val)
		HandlerGetCounter(w, request.WithContext(ctx))

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
		require.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		require.Equal(t, test.want.val, resBody)
	}
}
func TestHandlerGetGauge(t *testing.T) {
	type want struct {
		code        int
		contentType string
		val         []byte
	}
	tests := []struct {
		name string
		want want
		url  string
		val  string
	}{
		{
			name: "positive test #1",
			want: want{
				code:        200,
				contentType: "text/plain; charset=utf-8",
				val:         []byte("5244"),
			},
			url: "/update/gauge/someMetric/5244",
			val: "5244",
		},
	}
	for _, test := range tests {
		t.Log(test.name)
		t.Log(test.url)
		t.Log(test)
		request := httptest.NewRequest(http.MethodGet, test.url, nil)
		w := httptest.NewRecorder()

		ctx := context.WithValue(request.Context(), typeconst.UserIDContextKey, test.val)
		HandlerGetCounter(w, request.WithContext(ctx))

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
		require.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		require.Equal(t, test.want.val, resBody)
	}
}

func TestHandlerGetMetrics(t *testing.T) {
	type mockBehavior func(s *repo_mock.MockStorager, ctx context.Context)

	type args struct {
		mockBehavior mockBehavior
	}
	tests := []struct {
		name                string
		args                args
		want                func(res http.ResponseWriter, req *http.Request)
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "positive test #1",
			args: args{mockBehavior: func(s *repo_mock.MockStorager, ctx context.Context) {
				s.EXPECT().GetAllCounter(gomock.Any()).Return(map[string]int64{
					"someCounter": 1,
				}, nil)
				s.EXPECT().GetAllGauge(gomock.Any()).Return(map[string]float64{
					"someGauge": 1.1,
				}, nil)
			}},
			expectedStatusCode:  200,
			expectedRequestBody: "someGauge: 1.100000\nsomeCounter: 1\n",
		},
		// TODO: Add test cases.
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			stor := repo_mock.NewMockStorager(c)

			testCase.args.mockBehavior(stor, context.Background())

			w := httptest.NewRecorder()
			req := httptest.NewRequest("get", "/", nil)

			h := HandlerGetMetrics(stor)
			h(w, req)
			assert.Equal(t, w.Code, testCase.expectedStatusCode)

			res := w.Result()
			defer res.Body.Close()

			fmt.Println(res.StatusCode)

			resBody, _ := io.ReadAll(res.Body)

			require.Equal(t, string(resBody), string(testCase.expectedRequestBody))
		})
	}
}

func TestHandler_HandlerPostJSON(t *testing.T) {
	type mockBehavior func(s *service_mock.MockMetrics, ctx context.Context, metrics model.Metrics)

	//пример json
	reqBody := []byte(`{
  "id": "example_metric",
  "type": "gauge",
  "value": 10.5
}`)

	type args struct {
		mockBehavior mockBehavior
	}
	tests := []struct {
		name                string
		args                args
		want                func(res http.ResponseWriter, req *http.Request)
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "positive test #1",
			args: args{func(s *service_mock.MockMetrics, ctx context.Context, metrics model.Metrics) {
				s.EXPECT().PutMetrics(gomock.Any(), gomock.Any()).Return(nil)
				s.EXPECT().SendMetricstoFile(gomock.Any()).Return(nil)
			}},
			expectedStatusCode:  200,
			expectedRequestBody: string(reqBody),
		},
		// TODO: Add test cases.
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			stor := service_mock.NewMockMetrics(c)

			var metrics model.Metrics

			err := json.Unmarshal(reqBody, &metrics)
			require.NoError(t, err)

			testCase.args.mockBehavior(stor, context.Background(), metrics)

			handler := Handler{Service: &service.Service{Metrics: stor}}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/update/", bytes.NewBuffer(reqBody))

			handler.HandlerPostJSON(w, req)
			assert.Equal(t, w.Code, testCase.expectedStatusCode)

			res := w.Result()
			defer res.Body.Close()

			fmt.Println(res.StatusCode)

			resBody, _ := io.ReadAll(res.Body)

			require.Equal(t, string(resBody), string(testCase.expectedRequestBody))
		})
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

func TestHandler_HandlerPingDatabase(t *testing.T) {
	type mockBehavior func(s *repo_mock.MockStorager, ctx context.Context)

	type args struct {
		mockBehavior mockBehavior
	}
	var tests = []struct {
		name               string
		args               args
		expectedStatusCode int
	}{
		{
			name: "positive test #1",
			args: args{mockBehavior: func(s *repo_mock.MockStorager, ctx context.Context) {
				s.EXPECT().Ping(gomock.Any()).Return(nil)
			},
			},
			expectedStatusCode: 200,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			stor := repo_mock.NewMockStorager(c)

			h := &Handler{
				&service.Service{Storage: stor},
			}

			tt.args.mockBehavior(stor, context.Background())

			w := httptest.NewRecorder()
			req := httptest.NewRequest("get", "/", nil)

			h.HandlerPingDatabase(w, req)

			assert.Equal(t, w.Code, tt.expectedStatusCode)
		})
	}
}

func ExampleHandlerSucess() {

	http.HandleFunc("/success", HandlerSucess)
	http.ListenAndServe(":8080", nil)
}

func ExampleHandlerGetCounter() {
	path := "/update/counter/someMetric/527"
	val := "527"

	request := httptest.NewRequest(http.MethodGet, path, nil)
	w := httptest.NewRecorder()
	ctx := context.WithValue(request.Context(), typeconst.UserIDContextKey, val)

	HandlerGetCounter(w, request.WithContext(ctx))

	res := w.Result()
	defer res.Body.Close()

	fmt.Println(res.StatusCode)

	resBody, _ := io.ReadAll(res.Body)
	fmt.Println(string(resBody))

	// Output:
	// 200
	// 527
}

func ExampleHandlerGetGauge() {
	path := "/update/gauge/someMetric/5244"
	val := "5244"

	request := httptest.NewRequest(http.MethodGet, path, nil)
	w := httptest.NewRecorder()
	ctx := context.WithValue(request.Context(), typeconst.UserIDContextKey, val)

	HandlerGetCounter(w, request.WithContext(ctx))

	res := w.Result()
	defer res.Body.Close()

	fmt.Println(res.StatusCode)

	resBody, _ := io.ReadAll(res.Body)
	fmt.Println(string(resBody))

	// Output:
	// 200
	// 5244
}

func ExampleHandlerGetMetrics() {
	s := storage.NewStorage()

	http.HandleFunc("/GetMetrics", HandlerGetMetrics(s))
	http.ListenAndServe(":8080", nil)

}

type Reporter struct {
}

func (r *Reporter) Errorf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
func (r *Reporter) Fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// Example_HandlerPostJSON ExampleHandler_HandlerPostJSON
func ExampleHandler_HandlerPostJSON() {

	//todo пример с моком.
	/* panic: test executed panic(nil) or runtime.Goexit

	goroutine 1 [running]:
	testing.(*InternalExample).processRunResult(0xc0000318b0, {0x0, 0x0}, 0x2ecdcc2f2, 0x0, {0x0, 0x0})
	*/
	//	пример json
	reqBody := []byte(`{
	 "id": "example_metric",
	 "type": "gauge",
	 "value": 10.5
	}`)

	mockFunc := func(s *service_mock.MockMetrics, ctx context.Context, metrics model.Metrics) {
		s.EXPECT().PutMetrics(gomock.Any(), gomock.Any()).Return(nil)
		s.EXPECT().SendMetricstoFile(gomock.Any()).Return(nil)
	}

	r := &Reporter{}
	c := gomock.NewController(r)
	defer c.Finish()

	stor := service_mock.NewMockMetrics(c)

	var metrics model.Metrics

	err := json.Unmarshal(reqBody, &metrics)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()

	mockFunc(stor, ctx, metrics)
	// Создаем запрос

	// Создаем фейк хендлер
	handler := Handler{Service: &service.Service{Metrics: stor}}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/update/", bytes.NewBuffer(reqBody))

	handler.HandlerPostJSON(w, req)

	res := w.Result()
	defer res.Body.Close()

	fmt.Println(res.Header.Get("Content-Type"))
	fmt.Println(res.StatusCode)
	resBody, err := io.ReadAll(res.Body)
	//require.NoError(t, err)
	fmt.Println(string(resBody))

	// Output:
	// application/json
	// 200
	// {
	//	 "id": "example_metric",
	//	 "type": "gauge",
	//	 "value": 10.5
	//	}

}

func ExampleHandler_HandlerGetJSON() {
	val := 0.123

	forTestMetricReq := model.Metrics{
		ID:    "123",
		MType: "gauge",
		Delta: nil,
		Value: nil,
	}
	forTestMetricRes := model.Metrics{
		ID:    "123",
		MType: "gauge",
		Delta: nil,
		Value: &val,
	}

	//jsonforTestMetric, err := json.Marshal(forTestMetric)
	//if err != nil {
	//	fmt.Println(err)
	//}

	mock := func(s *service_mock.MockMetrics, metric model.Metrics) {
		s.EXPECT().GetMetrics(gomock.Any(), gomock.Any()).Return(model.Metrics{
			ID:    "123",
			MType: "gauge",
			Delta: nil,
			Value: &val,
		}, nil)
	}
	r := &Reporter{}
	c := gomock.NewController(r)
	defer c.Finish()

	metrics := service_mock.NewMockMetrics(c)

	serviceMock := &service.Service{
		Metrics: metrics,
	}
	handler := Handler{Service: serviceMock}
	mock(metrics, forTestMetricRes)

	w := httptest.NewRecorder()

	metricsreq, err := json.Marshal(forTestMetricReq)
	if err != nil {
		fmt.Println(err)
	}
	req := httptest.NewRequest("get", "/", bytes.NewBuffer(metricsreq))

	handler.HandlerGetJSON(w, req)

	res := w.Result()
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res.StatusCode)
	fmt.Println(string(resBody))
	// Output:
	// 200
	// {"id":"123","type":"gauge","value":0.123}

}

func ExampleHandler_HandlerGauge() {
	s := service.Service{}
	h := Handler{Service: &s}

	http.HandleFunc("/Gauge", h.HandlerGauge)
	http.ListenAndServe(":8080", nil)
}

func ExampleHandler_HandlerPingDatabase() {
	s := service.Service{}
	h := Handler{Service: &s}

	http.HandleFunc("/ping", h.HandlerPingDatabase)
	http.ListenAndServe(":8080", nil)
}
func ExampleHandler_HandlerPostBatched() {
	s := service.Service{}
	h := Handler{Service: &s}

	http.HandleFunc("/PostBatched", h.HandlerPostBatched)
	http.ListenAndServe(":8080", nil)
}
