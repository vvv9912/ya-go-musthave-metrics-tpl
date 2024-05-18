package grpcserver

import (
	pb "github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/grpcserver/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func NewGrpcServer(addr string) (*GrpcRequest, *grpc.ClientConn) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	metricsClient := pb.NewMetricsClient(conn)

	return &GrpcRequest{Client: metricsClient}, conn
}
