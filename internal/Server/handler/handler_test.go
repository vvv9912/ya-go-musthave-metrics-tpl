package handler

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		HandlerGauge(w, request)

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
