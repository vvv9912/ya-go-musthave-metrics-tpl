package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostRequest_PostReq(t *testing.T) {
	type fields struct {
		PostRequester PostRequester
	}
	type args struct {
		ctx context.Context
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "#1 start server ",
			wantErr: false,
			args: args{
				ctx: context.TODO(),
				url: "http://localhost:8080",
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				//логика сервера
			}))
			defer server.Close()

			p := &PostRequest{
				PostRequester: tt.fields.PostRequester,
			}

			// Update the URL to the test server URL
			tt.args.url = server.URL
			if err := p.PostReq(tt.args.ctx, tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("PostReq() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
