package grpcServer

import (
	"context"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/grpcServer/proto"
	"log"
)

type Metrics struct {
	proto.UnimplementedMetricsServer
}

func (m *Metrics) UpdateGauge(ctx context.Context, in *proto.Update) (*proto.Response, error) {
	var resp proto.Response

	log.Println(in.Key)
	log.Println(in.Values)
	_ = resp
	return nil, nil
}
func (m *Metrics) UpdateCounter(ctx context.Context, in *proto.Update) (*proto.Response, error) {
	return nil, nil
}
func (m *Metrics) UpdateGaugeJson(ctx context.Context, in *proto.UpdateSlice) (*proto.Response, error) {
	return nil, nil
}
func (m *Metrics) UpdateCounterJson(ctx context.Context, in *proto.UpdateSlice) (*proto.Response, error) {
	return nil, nil
}
func (m *Metrics) UpdatesBatched(ctx context.Context, in *proto.UpdateSlice) (*proto.Response, error) {
	return nil, nil
}
