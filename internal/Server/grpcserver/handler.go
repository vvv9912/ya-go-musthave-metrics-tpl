package grpcserver

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/grpcserver/proto"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/service"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/delaysend"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/logger"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"
	"go.uber.org/zap"
	"io"
	"strconv"
	"syscall"
)

type Metrics struct {
	Service *service.Service
	pb.UnimplementedMetricsServer
}

// UpdateGauge - обработчик для обновления значения метрики Gauge.
func (m *Metrics) UpdateGauge(ctx context.Context, in *pb.Update) (*pb.Response, error) {
	//var resp pb.Response
	value, err := strconv.ParseFloat(in.Values, 64)
	if err != nil {
		logger.Log.Error("parse float error", zap.Error(err))
		return &pb.Response{Error: "parse float error"}, err
	}

	err = m.Service.Storage.UpdateGauge(ctx, in.Key, value)
	if err != nil {
		return &pb.Response{Error: "update gauge in database"}, err
	}

	err = m.Service.Metrics.SendMetricstoFile(ctx)
	if err != nil {
		return &pb.Response{Error: "send metrics to file"}, err
	}
	return &pb.Response{}, nil
}

// UpdateCounter - обработчик для обновления значения метрики Counter.
func (m *Metrics) UpdateCounter(ctx context.Context, in *pb.Update) (*pb.Response, error) {

	value, err := strconv.ParseInt(in.Values, 10, 64)
	if err != nil {
		logger.Log.Error("parse int error", zap.Error(err))
		return &pb.Response{Error: "parse int error"}, err
	}

	err = m.Service.Storage.UpdateCounter(ctx, in.Key, value)
	if err != nil {
		return &pb.Response{Error: "update counter in database"}, err
	}

	err = m.Service.Metrics.SendMetricstoFile(ctx)
	if err != nil {
		return &pb.Response{Error: "send metrics to file"}, err
	}

	return nil, nil
}

// UpdateJson - обработчик для обновления значения метрик в формате Json.
func (m *Metrics) UpdateJSON(ctx context.Context, in *pb.UpdateSlice) (*pb.Response, error) {
	var metrics model.Metrics

	err := json.Unmarshal(in.Data, &metrics)
	if err != nil {
		return &pb.Response{Error: "parse json error"}, err
	}

	err = delaysend.NewDelaySend().SetDelay([]int{1, 3, 5}).
		AddExpectedError(syscall.ECONNREFUSED).
		SendDelayed(func() error {
			return m.Service.Metrics.PutMetrics(ctx, metrics)
		})
	if err != nil {
		logger.Log.Info("Failed to put metrics", zap.Error(err))
		return &pb.Response{Error: "Failed to put metrics"}, err
	}
	err = m.Service.Metrics.SendMetricstoFile(ctx)
	if err != nil {
		logger.Log.Error("Failed to send metrics to file", zap.Error(err))
		return &pb.Response{Error: "Failed to send metrics to file"}, err
	}

	return nil, nil
}

// UpdatesBatched - обработчик для обновления значения метрик в формате Batch.
func (m *Metrics) UpdatesBatched(ctx context.Context, in *pb.UpdateSlice) (*pb.Response, error) {
	var metrics []model.Metrics

	err := json.Unmarshal(in.Data, &metrics)
	if err != nil {
		logger.Log.Info("Failed to read request body", zap.Error(err))
		return &pb.Response{Error: "Failed to unmarshal request body"}, err
	}

	err = delaysend.NewDelaySend().SetDelay([]int{1, 3, 5}).
		AddExpectedError(syscall.ECONNREFUSED).
		SendDelayed(func() error {
			return m.Service.Storage.UpdateMetricsBatch(ctx, metrics)
		})

	if err != nil {
		logger.Log.Info("Failed to send metrics to batch", zap.Error(err))
		return &pb.Response{Error: "Failed to send metrics to batch"}, err
	}

	return nil, nil
}

// unGzip - распаковка значений.
func unGzip(in []byte) (out []byte, err error) {
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
