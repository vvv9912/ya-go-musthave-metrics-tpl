package notifier

import (
	"context"
	"encoding/json"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"
	"log"
	"strconv"
	"sync"
	"time"
)

type EventsMetric interface {
	UpdateMetricsGauge() *map[string]string
	UpdateMetricsCounter() (uint64, error)
}
type PostRequester interface {
	PostReq(ctx context.Context, url string) error
	PostReqJSON(ctx context.Context, url string, data []byte) error
	PostReqBatched(ctx context.Context, url string, data []model.Metrics) error
}
type Notifier struct {
	EventsMetric
	PostRequester

	TimerUpdate time.Duration
	TimerSend   time.Duration
	URL         string //localhost:8080
}

func NewNotifier(eventsMetric EventsMetric, postReq PostRequester, timeupdate time.Duration, timesend time.Duration, url string) *Notifier {
	return &Notifier{EventsMetric: eventsMetric, PostRequester: postReq, TimerUpdate: timeupdate, TimerSend: timesend, URL: url}
}

func (n *Notifier) NotifyPending() (*map[string]string, uint64, error) {
	gauge := n.UpdateMetricsGauge()
	counter, err := n.UpdateMetricsCounter()
	counter++
	if err != nil {
		log.Fatal(err)
		return nil, 0, err
	}
	return gauge, counter, nil
}
func (n *Notifier) SendNotification(ctx context.Context, gauge *map[string]string, counter uint64) error {
	var wg sync.WaitGroup

	//Передаем gauge
	for key, values := range *gauge {
		wg.Add(1)
		go func(key string, values string) {
			defer wg.Done()
			//todo параллельная отправка
			url := "http://" + n.URL + "/update/" + "gauge" + "/" + key + "/" + values

			err := n.PostReq(ctx, url)
			if err != nil {
				log.Println(err)
				return
			}

			val, err := strconv.ParseFloat(values, 64)
			if err != nil {
				log.Println(err)
			}

			m := model.Metrics{
				ID:    key,
				MType: "gauge",
				Delta: nil,
				Value: &val,
			}

			data, err := json.Marshal(m)

			if err != nil {
				log.Println(err)
			}

			url2 := "http://" + n.URL + "/update/"

			err = n.PostReqJSON(ctx, url2, data)
			if err != nil {
				log.Println(err)
			}
		}(key, values)
	}
	//Передаем counter
	wg.Add(1)
	go func() {
		defer wg.Done()
		//todo параллельная отправка
		coun := strconv.FormatUint(counter, 10)

		url := "http://" + n.URL + "/update/" + "counter" + "/" + "PollCount" + "/" + coun

		err := n.PostReq(ctx, url)
		if err != nil {
			log.Println(err)
			return
		}

		counterInt64 := int64(counter)
		m := model.Metrics{
			ID:    "PollCount",
			MType: "counter",
			Delta: &counterInt64,
			Value: nil,
		}

		data, err := json.Marshal(m)
		if err != nil {
			log.Println(err)
		}

		url2 := "http://" + n.URL + "/update/"
		err = n.PostReqJSON(ctx, url2, data)
		if err != nil {
			log.Println(err)
			return
		}
	}()
	wg.Add(1)
	//отправляем множество метрик
	go func() {
		defer wg.Done()

		url := "http://" + n.URL + "/updates/"
		m := make([]model.Metrics, 0, len(*gauge))
		for key, values := range *gauge {
			val, err := strconv.ParseFloat(values, 64)
			if err != nil {
				log.Println(err)
			}

			m = append(m, model.Metrics{
				ID:    key,
				MType: "gauge",
				Delta: nil,
				Value: &val,
			})

		}
		counterInt64 := int64(counter)
		m = append(m, model.Metrics{
			ID:    "PollCount",
			MType: "counter",
			Delta: &counterInt64,
			Value: nil,
		})

		err := n.PostReqBatched(ctx, url, m)
		if err != nil {
			log.Println(err)
			return
		}
	}()
	wg.Wait()

	return nil
}

func (n *Notifier) StartNotifyCron(ctx context.Context) error {
	var gauge *map[string]string
	var couter uint64
	var err error
	go func() {
		ticker := time.NewTicker(n.TimerUpdate)
		for {
			select {
			case <-ctx.Done():
				// Обработка завершения программы
				return
			case <-ticker.C:
				gauge, couter, err = n.NotifyPending()
				if err != nil {
					//log.Println(err)
					return
				}
				continue
			default:
				continue
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(n.TimerSend)
		for {
			select {
			case <-ctx.Done():
				// Обработка завершения программы
				return
			case <-ticker.C:
				if gauge != nil {
					err = n.SendNotification(ctx, gauge, couter)
				}
				if err != nil {
					//fmt.Println(err)
					return
				}
				continue
			default:
				continue
			}
		}
	}()
	return nil
}
