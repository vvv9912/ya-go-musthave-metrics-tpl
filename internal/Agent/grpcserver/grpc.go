package grpcserver

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	pb "github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/grpcserver/proto"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/logger"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type GrpcRequest struct {
	keyAuth   string
	publicKey *rsa.PublicKey
	host      string
	Client    pb.MetricsClient
}

func (m *GrpcRequest) UpdateGauge(ctx context.Context, update *pb.Update) error {

	headres := make(map[string]string)
	if m.host != "" {
		headres["X-Real-IP"] = m.host
	}

	md := metadata.New(headres)
	ctxMd := metadata.NewOutgoingContext(ctx, md)
	_, err := m.Client.UpdateGauge(ctxMd, update)
	if err != nil {
		logger.Log.Error("grpcserver.UpdateGauge failed", zap.Error(err))
		return err
	}

	return nil
}

func (m *GrpcRequest) UpdateCounter(ctx context.Context, update *pb.Update) error {
	headres := make(map[string]string)
	if m.host != "" {
		headres["X-Real-IP"] = m.host
	}

	md := metadata.New(headres)
	ctxMd := metadata.NewOutgoingContext(ctx, md)
	_, err := m.Client.UpdateCounter(ctxMd, update)
	if err != nil {

		logger.Log.Error("grpcserver.UpdateCounter failed", zap.Error(err))
		return err
	}

	return nil
}

func (m *GrpcRequest) UpdateJSON(ctx context.Context, data []byte) error {
	return m.updateJSON(ctx, &pb.UpdateSlice{
		Data: data,
	})
}

func (m *GrpcRequest) UpdatesBatched(ctx context.Context, data []model.Metrics) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Log.Error("Error marshaling metrics", zap.Error(err))
		return err
	}

	return m.updatesBatched(ctx, &pb.UpdateSlice{
		Data: jsonData,
	})
}
func (m *GrpcRequest) updateJSON(ctx context.Context, update *pb.UpdateSlice) error {
	headers, err := m.preparingReq(ctx, update)
	if err != nil {
		logger.Log.Error("failed create req", zap.Error(err))
		return err
	}

	md := metadata.New(headers)
	ctxMd := metadata.NewOutgoingContext(ctx, md)
	_, err = m.Client.UpdateJson(ctxMd, update)
	if err != nil {
		logger.Log.Error("Error update gauge json", zap.Error(err))
		return err
	}

	return nil
}

func (m *GrpcRequest) updatesBatched(ctx context.Context, update *pb.UpdateSlice) error {
	headers, err := m.preparingReq(ctx, update)
	if err != nil {
		logger.Log.Error("failed create req", zap.Error(err))
		return err
	}

	md := metadata.New(headers)
	ctxMd := metadata.NewOutgoingContext(ctx, md)
	_, err = m.Client.UpdatesBatched(ctxMd, update)
	if err != nil {
		logger.Log.Error("Error update counter json", zap.Error(err))
		return err
	}

	return nil
}

// Preparing json req
func (m *GrpcRequest) preparingReq(ctx context.Context, update *pb.UpdateSlice) (map[string]string, error) {
	headers := make(map[string]string)
	if m.host != "" {
		headers["X-Real-IP"] = m.host
	}

	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	_, err := zb.Write(update.Data)
	if err != nil {
		logger.Log.Error("Failed gzip", zap.Error(err))
		return nil, err
	}
	err = zb.Close()
	if err != nil {
		logger.Log.Error("Failed gzip", zap.Error(err))
		return nil, err
	}

	dataBytes := buf.Bytes()
	// В случае, если ключ не задан
	if m.keyAuth != "" {

		h := hmac.New(sha256.New, []byte(m.keyAuth))

		_, err = h.Write(update.Data)
		if err != nil {
			logger.Log.Error("Failed write", zap.Error(err))
			return nil, err
		}
		dst := h.Sum(nil)
		headers["HashSHA256"] = fmt.Sprintf("%x", dst)
	}

	if m.publicKey != nil {
		dataBytes, err = rsa.EncryptPKCS1v15(rand.Reader, m.publicKey, buf.Bytes())
		if err != nil {
			logger.Log.Error("Failed to encrypt", zap.Error(err))
			return nil, err
		}
	}

	(*update).Data = dataBytes

	return headers, nil
}
