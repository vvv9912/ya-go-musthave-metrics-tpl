package grpcServer

import (
	pb "github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/grpcServer/proto"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/service"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

func NewGrpcServer(service *service.Service, addr string) error {
	//grpc
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Log.Error("listen failed", zap.Error(err))
		return err
	}

	HandlersGrpc := Metrics{
		Service: service,
	}

	grpcNewServer1 := grpc.NewServer(grpc.UnaryInterceptor(UnaryInterceptor))

	pb.RegisterMetricsServer(grpcNewServer1, &HandlersGrpc)

	go func() {

		if err := grpcNewServer1.Serve(listen); err != nil {
			logger.Log.Error("grpc server start failed", zap.Error(err))
		}

	}()

	return err
}
