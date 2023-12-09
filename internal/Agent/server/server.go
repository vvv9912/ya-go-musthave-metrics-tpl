package server

import (
	"context"
	"net/http"
)

type PostRequester interface {
	PostReq(ctx context.Context, url string) error
}
type PostRequest struct {
	PostRequester
}

func NewPostRequest() *PostRequest {
	return &PostRequest{}
}

func (p *PostRequest) PostReq(ctx context.Context, url string) error {
	// Создаем новый запрос url
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil) //http://localhost:8080/update/gauge/someMetric/527
	req.Header.Set("Content-Type", "text/plain")
	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		//	fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	if err != nil {

		return err
	}
	return nil
}
