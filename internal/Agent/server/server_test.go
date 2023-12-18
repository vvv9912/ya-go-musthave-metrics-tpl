package server

import (
	"context"
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
			p := &PostRequest{
				PostRequester: tt.fields.PostRequester,
			}
			if err := p.PostReq(tt.args.ctx, tt.args.url); (err != nil) != tt.wantErr {

				t.Errorf("PostReq() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
