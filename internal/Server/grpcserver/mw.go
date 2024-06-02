package grpcserver

import (
	"context"
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha256"
	pb "github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/grpcserver/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Interceptor struct {
	privateKey    *rsa.PrivateKey
	trustedSubnet string
	KeyAuth       string
}

func NewInterceptor(privateKey *rsa.PrivateKey, trustedSubnet string, KeyAuth string) *Interceptor {
	return &Interceptor{privateKey: privateKey, trustedSubnet: trustedSubnet, KeyAuth: KeyAuth}
}

// UnaryInterceptor - mw для распаковки значений для нужных обработчиков.
func (i *Interceptor) GzipInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod == "/grpc.Metrics/UpdatesBatched" || info.FullMethod == "/grpc.Metrics/UpdateJson" {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}
		values := md.Get("Encoding")
		if len(values) == 0 {
			return handler(ctx, req)
		}
		for _, value := range values {
			if value == "gzip" {
				a := req.(*pb.UpdateSlice)
				data, err := unGzip(a.Data)
				if err != nil {
					return nil, err
				}
				a.Data = data
				return handler(ctx, req)
			}
		}

	}

	return handler(ctx, req)
}
func (i *Interceptor) TrustedSubnetInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return handler(ctx, req)
	}
	values := md.Get("X-Real-IP")
	if len(values) != 1 {
		return handler(ctx, req)
	}
	if values[0] == "" {
		return handler(ctx, req)
	}

	a := req.(*pb.UpdateSlice)
	data, err := unGzip(a.Data)
	if err != nil {
		return nil, status.Error(codes.Aborted, "Invalid unzip")
	}
	a.Data = data
	return handler(ctx, req)

}
func (i *Interceptor) HashInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return handler(ctx, req)
	}
	values := md.Get("HashSHA256")
	if len(values) != 1 {
		return handler(ctx, req)
	}
	if values[0] == "" {
		return handler(ctx, req)
	}
	//считаем хеш
	h := hmac.New(sha256.New, []byte(i.KeyAuth))
	a := req.(*pb.UpdateSlice)
	h.Write(a.Data)
	dst := h.Sum(nil)

	ok = hmac.Equal(dst, []byte(values[0]))
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "invalid signature")
	}

	return handler(ctx, req)

}
