package grpcServer

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/grpcServer/proto"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/service"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/delaysend"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/logger"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"strconv"
	"syscall"
)

type Metrics struct {
	Service *service.Service
	proto.UnimplementedMetricsServer
}

func (m *Metrics) UpdateGauge(ctx context.Context, in *proto.Update) (*proto.Response, error) {
	//var resp proto.Response
	value, err := strconv.ParseFloat(in.Values, 64)
	if err != nil {
		logger.Log.Error("parse float error", zap.Error(err))
		return &proto.Response{Error: "parse float error"}, err
	}

	err = m.Service.Storage.UpdateGauge(ctx, in.Key, value)
	if err != nil {
		return &proto.Response{Error: "update gauge in database"}, err
	}

	err = m.Service.Metrics.SendMetricstoFile(ctx)
	if err != nil {
		return &proto.Response{Error: "send metrics to file"}, err
	}
	return nil, nil
}
func (m *Metrics) UpdateCounter(ctx context.Context, in *proto.Update) (*proto.Response, error) {
	var resp proto.Response

	value, err := strconv.ParseInt(in.Values, 10, 64)
	if err != nil {
		logger.Log.Error("parse int error", zap.Error(err))
		return &proto.Response{Error: "parse int error"}, err
	}

	err = m.Service.Storage.UpdateCounter(ctx, in.Key, value)
	if err != nil {
		return &proto.Response{Error: "update counter in database"}, err
	}

	err = m.Service.Metrics.SendMetricstoFile(ctx)
	if err != nil {
		return &proto.Response{Error: "send metrics to file"}, err
	}

	return nil, nil
}
func (m *Metrics) UpdateGaugeJson(ctx context.Context, in *proto.UpdateSlice) (*proto.Response, error) {
	var metrics model.Metrics

	err := json.Unmarshal(in.Data, &metrics)
	if err != nil {
		return &proto.Response{Error: "parse json error"}, err
	}

	err = delaysend.NewDelaySend().SetDelay([]int{1, 3, 5}).
		AddExpectedError(syscall.ECONNREFUSED).
		SendDelayed(func() error {
			return m.Service.Metrics.PutMetrics(ctx, metrics)
		})
	if err != nil {
		logger.Log.Info("Failed to put metrics", zap.Error(err))
		return &proto.Response{Error: "Failed to put metrics"}, err
	}
	err = m.Service.Metrics.SendMetricstoFile(ctx)
	if err != nil {
		logger.Log.Error("Failed to send metrics to file", zap.Error(err))
		return &proto.Response{Error: "Failed to send metrics to file"}, err
	}

	return nil, nil
}

//	func (m *Metrics) UpdateCounterJson(ctx context.Context, in *proto.UpdateSlice) (*proto.Response, error) {
//		var resp proto.Response
//
//		log.Println(in.Data)
//		//log.Println(in.Values)
//		_ = resp
//		return nil, nil
//	}
func (m *Metrics) UpdatesBatched(ctx context.Context, in *proto.UpdateSlice) (*proto.Response, error) {
	var metrics []model.Metrics

	err := json.Unmarshal(in.Data, &metrics)
	if err != nil {
		logger.Log.Info("Failed to read request body", zap.Error(err))
		return &proto.Response{Error: "Failed to unmarshal request body"}, err
	}

	err = delaysend.NewDelaySend().SetDelay([]int{1, 3, 5}).
		AddExpectedError(syscall.ECONNREFUSED).
		SendDelayed(func() error {
			return m.Service.Storage.UpdateMetricsBatch(ctx, metrics)
		})

	if err != nil {
		logger.Log.Info("Failed to send metrics to batch", zap.Error(err))
		return &proto.Response{Error: "Failed to send metrics to batch"}, err
	}

	return nil, nil
}

func (m *Metrics) unGzip(in []byte) (out []byte, err error) {
	bb := bytes.NewReader(in)

	r, err := gzip.NewReader(bb)
	if err != nil {
		logger.Log.Error("Error gzip", zap.Error(err))
		return nil, err
	}
	data, err := io.ReadAll(r)
	if err != nil {
		fmt.Println("Error read", err)

	}

	return data, nil
}
func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var token string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("token")
		if len(values) > 0 {
			token = values[0]
		}
	}
	_ = token
	return handler(ctx, req)
}

//// Регистрация RPC методов с разными UnaryInterceptor
//// Метод 1 с UnaryInterceptor myUnaryInterceptor1
//RegisterService1Server(server, &service1{ /* ... */ }, grpc.WithUnaryInterceptor(myUnaryInterceptor1))
//// Метод 2 с UnaryInterceptor myUnaryInterceptor2
//RegisterService2Server(server, &service2{ /* ... */ }, grpc.WithUnaryInterceptor(myUnaryInterceptor2))
//
//// Запуск сервера и обработка RPC методов
//// ...
//
//fmt.Println("Сервер gRPC запущен")
