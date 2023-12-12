package notifier

import (
	"context"
	"fmt"
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

func (n *Notifier) NotifyPending(ctx context.Context) (*map[string]string, uint64, error) {
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
			url := "http://" + n.URL + "/update/" + "gauge" + "/" + key + "/" + values
			err := n.PostReq(ctx, url)
			if err != nil {
				log.Println(err)
				return //return err?
			}
		}(key, values)
	}
	//Передаем counter
	wg.Add(1)
	go func() {
		defer wg.Done()
		coun := strconv.FormatUint(counter, 10)
		url := "http://" + n.URL + "/update/" + "counter" + "/" + "PollCount" + "/" + coun
		err := n.PostReq(ctx, url)
		if err != nil {
			log.Println(err)
			return //return err?
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
		for {
			select {
			case <-ctx.Done():
				// Обработка завершения программы
				return
			default:
			}

			gauge, couter, err = n.NotifyPending(ctx)
			if err != nil {

				log.Fatal(err)
				return
			}
			time.Sleep(n.TimerUpdate)

		}
	}()
	time.Sleep(time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				// Обработка завершения программы
				return
			default:
			}
			if gauge != nil {
				err = n.SendNotification(ctx, gauge, couter)
			}
			if err != nil {
				fmt.Println(err)
				return
			}
			time.Sleep(n.TimerSend)

		}
	}()
	return nil
}
