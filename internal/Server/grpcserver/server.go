package grpcserver

import (
	"crypto/rsa"
	pb "github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/grpcserver/proto"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/service"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

func NewGrpcServer(service *service.Service, privateKey *rsa.PrivateKey, trustedSubnet string, KeyAuth string, addr string) (*grpc.Server, error) {
	//grpc
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Log.Error("listen failed", zap.Error(err))
		return nil, err
	}

	HandlersGrpc := Metrics{
		Service: service,
	}
	i := NewInterceptor(privateKey, trustedSubnet, KeyAuth)

	unaryInterceptors := []grpc.UnaryServerInterceptor{}

	unaryInterceptors = append(unaryInterceptors, i.GzipInterceptor)
	unaryInterceptors = append(unaryInterceptors, i.HashInterceptor)
	unaryInterceptors = append(unaryInterceptors, i.TrustedSubnetInterceptor)

	grpcNewServer1 := grpc.NewServer(grpc.ChainUnaryInterceptor(unaryInterceptors...))

	pb.RegisterMetricsServer(grpcNewServer1, &HandlersGrpc)

	go func() {

		if err := grpcNewServer1.Serve(listen); err != nil {
			logger.Log.Error("grpc server start failed", zap.Error(err))
		}

	}()
	return grpcNewServer1, err
}
