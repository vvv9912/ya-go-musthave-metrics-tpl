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

func TestMww(t *testing.T) {
	//type fields struct {
	//	GaugeStorage   storage.GaugeStorager
	//	CounterStorage storage.CounterStorager
	//	Log            *log.Logger
	//}
	//type args struct {
	//	next http.Handler
	//}
	//tests := []struct {
	//	name   string
	//	fields fields
	//	args   args
	//	want   http.Handler
	//}{
	//	// TODO: Add test cases.
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//
	//		log := &log.Logger{}
	//		mw := &Mw{
	//			GaugeStorage:   tt.fields.GaugeStorage,
	//			CounterStorage: tt.fields.CounterStorage,
	//			Log:            log,
	//		}
	//		if got := mw.Middlware(tt.args.next); !reflect.DeepEqual(got, tt.want) {
	//			t.Errorf("Middlware() = %v, want %v", got, tt.want)
	//		}
	//	})
	//}
}
