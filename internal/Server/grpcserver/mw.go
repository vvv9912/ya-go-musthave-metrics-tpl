package grpcserver

import (
	"context"
	pb "github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/grpcserver/proto"
	"google.golang.org/grpc"
)

// UnaryInterceptor - mw для распаковки значений для нужных обработчиков.
func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	if info.FullMethod == "/grpc.Metrics/UpdatesBatched" || info.FullMethod == "/grpc.Metrics/UpdateJson" {
		a := req.(*pb.UpdateSlice)
		data, err := unGzip(a.Data)
		if err != nil {
			return nil, err
		}
		a.Data = data
	}

	return handler(ctx, req)
}
