package mw

import (
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

func TestMw_MiddlwareCounter(t *testing.T) {
	//type want struct {
	//	code        int
	//	contentType string
	//}
	//type args struct {
	//	next http.Handler
	//}
	//
	//counter := storage.NewCounterStorage()
	//gauge := storage.NewGaugeStorage()
	//
	//m := Mw{
	//	GaugeStorage:   gauge,
	//	CounterStorage: counter,
	//}
	//tests := []struct {
	//	name string
	//	args args
	//	URL  string
	//	want want
	//}{
	//	{name: "#1 positive test",
	//		args: args{next: m.MiddlewareType(http.HandlerFunc(handler.HandlerSucess))},
	//		URL:  "/update/gauge/someMetric/527",
	//		want: want{code: 400, contentType: "text/plain; charset=utf-8"},
	//	},
	//{name: "#2 negative test",
	//	args: args{next: m.Middlware(m.Middlware(http.HandlerFunc(handler.HandlerSucess)))},
	//	URL:  "/update/counter/",
	//	want: want{code: 404, contentType: "text/plain; charset=utf-8"},
	//},
	//{name: "#3 negative test",
	//	args: args{next: m.Middlware(m.Middlware(http.HandlerFunc(handler.HandlerSucess)))},
	//	URL:  "/update/counter/dasd",
	//	want: want{code: 404, contentType: "text/plain; charset=utf-8"},
	//},
	//{name: "#4 negative test",
	//	args: args{next: m.Middlware(m.Middlware(http.HandlerFunc(handler.HandlerSucess)))},
	//	URL:  "/update",
	//	want: want{code: 404, contentType: "text/plain; charset=utf-8"},
	//},
	//}
	//for _, test := range tests {
	//	t.Log(test)
	//	request := httptest.NewRequest(http.MethodPost, test.URL, nil) //не рааботает тест
	//	w := httptest.NewRecorder()
	//	test.args.next.ServeHTTP(w, request)
	//	res := w.Result()
	//	assert.Equal(t, test.want.code, res.StatusCode)
	//	defer res.Body.Close()
	//	resBody, err := io.ReadAll(res.Body)
	//	require.NoError(t, err)
	//	t.Log("----------///\nContent type:", res.Header.Get("Content-Type"))
	//	t.Log("----------///\nres body:", string(resBody))
	//	assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
	//
	//}
}
