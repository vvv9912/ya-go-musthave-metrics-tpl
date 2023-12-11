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
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")
	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

//func (p *PostRequest) PostReq(ctx context.Context, url string) error {
//	client := resty.New()
//	_, err := client.R().SetHeaders(map[string]string{
//		"Content-Type": "text/plain",
//	}).Post(url)
//	if err != nil {
//		fmt.Println(err)
//		return err
//	}
//	return nil
//}
