package server

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
	"github.com/go-resty/resty/v2"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/logger"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"
	"go.uber.org/zap"
)

type PostRequester interface {
	PostReq(ctx context.Context, url string) error
	PostReqJSON(ctx context.Context, url string, data []byte) error
	PostReqBatched(ctx context.Context, url string, data []model.Metrics) error
}
type PostRequest struct {
	PostRequester
	keyAuth   string
	publicKey *rsa.PublicKey
	host      string
}

func NewPostRequest(keyAuth string, publicKey *rsa.PublicKey, host string) *PostRequest {
	return &PostRequest{keyAuth: keyAuth, publicKey: publicKey, host: host}
}

func (p *PostRequest) PostReq(ctx context.Context, url string) error {

	client := resty.New()
	req := client.R()
	if p.host != "" {
		req.SetHeaders(map[string]string{
			"X-Real-IP": p.host,
		})
	}
	_, err := req.SetHeaders(map[string]string{
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
	req := client.R()

	if p.host != "" {
		req.SetHeaders(map[string]string{
			"X-Real-IP": p.host,
		})
	}

	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	_, err := zb.Write(data)
	if err != nil {
		logger.Log.Error("Failed gzip", zap.Error(err))
		return err
	}
	err = zb.Close()
	if err != nil {
		logger.Log.Error("Failed gzip", zap.Error(err))
		return err
	}
	//
	dataBytes := buf.Bytes()
	// В случае, если ключ не задан
	if p.keyAuth != "" {

		h := hmac.New(sha256.New, []byte(p.keyAuth))

		_, err = h.Write(data)
		if err != nil {
			logger.Log.Error("Failed write", zap.Error(err))
			return err
		}
		dst := h.Sum(nil)

		req.SetHeaders(map[string]string{"HashSHA256": fmt.Sprintf("%x", dst)})

	}

	if p.publicKey != nil {
		dataBytes, err = rsa.EncryptPKCS1v15(rand.Reader, p.publicKey, buf.Bytes())
		if err != nil {
			logger.Log.Error("Failed to encrypt", zap.Error(err))
			return err
		}
	}

	_, err = req.SetHeaders(map[string]string{
		"Content-Type": "application/json", "Content-Encoding": "gzip",
	}).SetBody(dataBytes).Post(url)

	if err != nil {
		logger.Log.Error("Failed to send metrics", zap.Error(err))
		return err
	}

	return nil
}

func (p *PostRequest) PostReqBatched(ctx context.Context, url string, data []model.Metrics) error {

	client := resty.New()

	req := client.R()
	if p.host != "" {
		req.SetHeaders(map[string]string{
			"X-Real-IP": p.host,
		})
	}

	// При передаче слайса в интерфейс client.R, внутри все равно преобраз. в json
	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Log.Error("Error marshaling metrics", zap.Error(err))
		return err
	}

	// Если ключ не задан (подпись)
	if p.keyAuth != "" {

		h := hmac.New(sha256.New, []byte(p.keyAuth))

		_, err = h.Write(jsonData)
		if err != nil {
			logger.Log.Info("Failed to write", zap.Error(err))
			return err
		}

		dst := h.Sum(nil)

		req.SetHeaders(map[string]string{"HashSHA256": fmt.Sprintf("%x", dst), "Content-Type": "application/json"})

	}

	if p.publicKey != nil {
		jsonData, err = rsa.EncryptPKCS1v15(rand.Reader, p.publicKey, jsonData)
		if err != nil {
			logger.Log.Error("Failed to encrypt", zap.Error(err))
			return err
		}
	}

	_, err = req.SetBody(jsonData).Post(url)
	if err != nil {
		logger.Log.Error("Failed to send metrics batch", zap.Error(err))
		return err
	}

	return nil
}
