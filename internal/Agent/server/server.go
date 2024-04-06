package server

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/logger"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"
)

type PostRequester interface {
	PostReq(ctx context.Context, url string) error
	PostReqJSON(ctx context.Context, url string, data []byte) error
	PostReqBatched(ctx context.Context, url string, data []model.Metrics) error
}
type PostRequest struct {
	PostRequester
	keyAuth string
}

func NewPostRequest(keyAuth string) *PostRequest {
	return &PostRequest{keyAuth: keyAuth}
}

func (p *PostRequest) PostReq(ctx context.Context, url string) error {

	client := resty.New()

	_, err := client.R().SetHeaders(map[string]string{
		"Content-Type": "text/plain",
	}).Post(url)
	if err != nil {
		logger.Log.Error("Failed to send metrics", zap.Error(err))
		return err
	}

	return nil
}
func (p *PostRequest) PostReqJSON(ctx context.Context, url string, data []byte) error {
	client := resty.New()

	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	_, err := zb.Write(data)
	if err != nil {
		logger.Log.Error("Failed gzip", zap.Error(err))
		return err
	}
	zb.Close()

	// В случае, если ключ не задан
	if p.keyAuth != "" {

		h := hmac.New(sha256.New, []byte(p.keyAuth))

		h.Write(data)

		dst := h.Sum(nil)

		_, err = client.R().SetHeaders(map[string]string{
			"HashSHA256": fmt.Sprintf("%x", dst), "Content-Type": "application/json", "Content-Encoding": "gzip",
		}).SetBody(buf).Post(url)

		if err != nil {
			logger.Log.Error("Failed to send metrics", zap.Error(err))
			return err
		}
	} else {

		_, err = client.R().SetHeaders(map[string]string{
			"Content-Type": "application/json", "Content-Encoding": "gzip",
		}).SetBody(buf).Post(url)

		if err != nil {
			logger.Log.Error("Failed to send metrics", zap.Error(err))
			return err
		}
	}

	return nil
}

func (p *PostRequest) PostReqBatched(ctx context.Context, url string, data []model.Metrics) error {

	client := resty.New()

	// Если ключ не задан
	if p.keyAuth != "" {

		h := hmac.New(sha256.New, []byte(p.keyAuth))

		jsonData, err := json.Marshal(data)
		if err != nil {
			logger.Log.Error("Error marshaling metrics", zap.Error(err))
			return err
		}

		h.Write(jsonData)

		dst := h.Sum(nil)

		_, err = client.R().SetHeaders(map[string]string{"HashSHA256": fmt.Sprintf("%x", dst), "Content-Type": "application/json"}).SetBody(data).Post(url)
		if err != nil {
			logger.Log.Error("Failed to send metrics batch", zap.Error(err))
			return err
		}

	} else {

		_, err := client.R().SetBody(data).Post(url)
		if err != nil {
			logger.Log.Error("Failed to send metrics batch", zap.Error(err))
			return err
		}

	}
	return nil
}
