package grpcServer

import (
	"context"
	"crypto/rsa"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/grpcServer/proto"
	"log"
)

type GrpcRequest struct {
	keyAuth   string
	publicKey *rsa.PublicKey
	host      string
	Client    proto.MetricsClient
}

func (m *GrpcRequest) UpdateGauge(ctx context.Context, update *proto.Update) error {
	resp, err := m.Client.UpdateGauge(ctx, update)
	if err != nil {
		log.Println("Bad grpcServer req:", err)
		return err
	}
	log.Println("Okay! resp:", resp)
	return nil
}

//func (m *GrpcRequest) UpdateCounter(ctx context.Context, in *proto.Update) (*proto.Response, error) {
//	return nil, nil
//}
//func (m *GrpcRequest) UpdateGaugeJson(ctx context.Context, in *proto.UpdateSlice) (*proto.Response, error) {
//	return nil, nil
//}
//func (m *GrpcRequest) UpdateCounterJson(ctx context.Context, in *proto.UpdateSlice) (*proto.Response, error) {
//	return nil, nil
//}
//func (m *GrpcRequest) UpdatesBatched(ctx context.Context, in *proto.UpdateSlice) (*proto.Response, error) {
//	return nil, nil
//}
