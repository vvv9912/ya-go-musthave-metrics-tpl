package server

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/store"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/logger"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"
	"go.uber.org/zap"
	"log"
	"syscall"
	"time"
)

type PostRequester interface {
	PostReq(ctx context.Context, url string) error
	PostReqJSON(ctx context.Context, url string, data []byte) error
}
type PostRequest struct {
	PostRequester
}

func NewPostRequest() *PostRequest {
	return &PostRequest{}
}

func (p *PostRequest) PostReq(ctx context.Context, url string) error {
	client := resty.New()
	var RetriableTime = []int{1, 3, 5}
	_, err := client.R().SetHeaders(map[string]string{
		"Content-Type": "text/plain",
	}).Post(url)
	if errors.Is(err, syscall.ECONNREFUSED) {
		logger.Log.Info("Failed to post metrics", zap.Error(err))
		for i := range RetriableTime {
			time.Sleep(time.Duration(RetriableTime[i]) * time.Second)
			_, err = client.R().SetHeaders(map[string]string{
				"Content-Type": "text/plain",
			}).Post(url)
			if err != nil {
				if errors.Is(err, syscall.ECONNREFUSED) {
					logger.Log.Info("Failed to post metrics", zap.String("retriable time: ", fmt.Sprint(RetriableTime[i])), zap.Error(err))

				} else {

					logger.Log.Info("Failed to post metrics", zap.String("retriable time: ", fmt.Sprint(RetriableTime[i])), zap.Error(err))
					return err

				}
			} else {
				break
			}
		}
	}
	if err != nil {
		logger.Log.Info("Failed to post metrics", zap.Error(err))
		return err
	}
	return nil
}
func (p *PostRequest) PostReqJSON(ctx context.Context, url string, data []byte) error {
	var RetriableTime = []int{1, 3, 5}

	client := resty.New()

	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	_, err := zb.Write(data)

	if err != nil {
		log.Println(err)
		return err
	}
	zb.Close()

	_, err = client.R().SetHeaders(map[string]string{
		"Content-Type": "application/json", "Content-Encoding": "gzip",
	}).SetBody(buf).Post(url)
	if errors.Is(err, syscall.ECONNREFUSED) {
		logger.Log.Info("Failed to post metrics", zap.Error(err))
		for i := range RetriableTime {
			time.Sleep(time.Duration(RetriableTime[i]) * time.Second)

			_, err = client.R().SetHeaders(map[string]string{
				"Content-Type": "application/json", "Content-Encoding": "gzip",
			}).SetBody(buf).Post(url)

			if err != nil {
				if errors.Is(err, syscall.ECONNREFUSED) {
					logger.Log.Info("Failed to post metrics", zap.String("retriable time: ", fmt.Sprint(store.RetriableTime[i])), zap.Error(err))

				} else {

					logger.Log.Info("Failed to post metrics", zap.String("retriable time: ", fmt.Sprint(store.RetriableTime[i])), zap.Error(err))
					return err

				}
			} else {
				break
			}
		}
	}
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
func (p *PostRequest) PostReqBatched(ctx context.Context, url string, data []model.Metrics) error {
	var RetriableTime = []int{1, 3, 5}

	client := resty.New()
	_, err := client.R().SetBody(data).Post(url)
	if errors.Is(err, syscall.ECONNREFUSED) {
		logger.Log.Info("Failed to post metrics", zap.Error(err))
		for i := range RetriableTime {
			time.Sleep(time.Duration(RetriableTime[i]) * time.Second)

			_, err := client.R().SetBody(data).Post(url)

			if err != nil {
				if errors.Is(err, syscall.ECONNREFUSED) {
					logger.Log.Info("Failed to post metrics", zap.String("retriable time: ", fmt.Sprint(store.RetriableTime[i])), zap.Error(err))

				} else {

					logger.Log.Info("Failed to post metrics", zap.String("retriable time: ", fmt.Sprint(store.RetriableTime[i])), zap.Error(err))
					return err

				}
			} else {
				break
			}
		}
	}
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//на http
//func (p *PostRequest) PostReq(ctx context.Context, url string) error {
//	// Создаем новый запрос url
//	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil) //http://localhost:8080/update/gauge/someMetric/527
//	if err != nil {
//		return err
//	}
//	req.Header.Set("Content-Type", "text/plain")
//	// Выполняем запрос
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//
//	return nil
//}
