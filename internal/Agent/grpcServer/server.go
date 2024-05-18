package grpcServer

import (
	pb "github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/grpcServer/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func NewGrpcServer(addr string) *GrpcRequest {
	// устанавливаем соединение с сервером grpc
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	metricsClient := pb.NewMetricsClient(conn)

	return &GrpcRequest{Client: metricsClient}
}
