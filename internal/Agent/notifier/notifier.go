package notifier

import (
	"context"
	"encoding/json"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/delaysend"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/logger"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"
	"go.uber.org/zap"
	"log"
	"strconv"
	"sync"
	"syscall"
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
	URL         string
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

func (n *Notifier) SendGaugeReq(ctx context.Context, gauge *map[string]string, PullGoCh chan struct{}, wgWorker *sync.WaitGroup) {
	defer wgWorker.Done()

	wg := sync.WaitGroup{}

	for key, values := range *gauge {
		wg.Add(1)
		PullGoCh <- struct{}{}
		go func(key string, values string) {
			defer func() {
				//освоблождаем горутину
				<-PullGoCh
				wg.Done()
			}()
			url := "http://" + n.URL + "/update/" + "gauge" + "/" + key + "/" + values

			err := delaysend.NewDelaySend().SetDelay([]int{1, 3, 5}).
				AddExpectedError(syscall.ECONNREFUSED).
				SendDelayed(func() error {
					return n.PostReq(ctx, url)
				})
			if err != nil {
				log.Println(err)
				return
			}

			val, err := strconv.ParseFloat(values, 64)
			if err != nil {
				logger.Log.Error("Failed to parse float", zap.Error(err))
				return
			}

			m := model.Metrics{
				ID:    key,
				MType: "gauge",
				Delta: nil,
				Value: &val,
			}

			data, err := json.Marshal(m)
			if err != nil {
				logger.Log.Error("Failed to marshal JSON", zap.Error(err))
				return
			}

			url2 := "http://" + n.URL + "/update/"

			err = delaysend.NewDelaySend().SetDelay([]int{1, 3, 5}).
				AddExpectedError(syscall.ECONNREFUSED).
				SendDelayed(func() error {
					return n.PostReqJSON(ctx, url2, data)
				})
			if err != nil {
				logger.Log.Error("Failed to send JSON", zap.Error(err))
			}
		}(key, values)
	}

	wg.Wait()
}

func (n *Notifier) SendCountReq(ctx context.Context, counter uint64, PullGoCh chan struct{}, wgWorker *sync.WaitGroup) {

	PullGoCh <- struct{}{}
	defer func() {
		//освоблождаем горутину
		<-PullGoCh
		wgWorker.Done()
	}()

	coun := strconv.FormatUint(counter, 10)

	url := "http://" + n.URL + "/update/" + "counter" + "/" + "PollCount" + "/" + coun

	err := delaysend.NewDelaySend().SetDelay([]int{1, 3, 5}).
		AddExpectedError(syscall.ECONNREFUSED).
		SendDelayed(func() error {
			return n.PostReq(ctx, url)
		})
	if err != nil {
		logger.Log.Error("Failed to send counter PostReq", zap.Error(err))
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
		logger.Log.Error("Failed to marshal JSON", zap.Error(err))
		return
	}

	url2 := "http://" + n.URL + "/update/"
	err = delaysend.NewDelaySend().SetDelay([]int{1, 3, 5}).
		AddExpectedError(syscall.ECONNREFUSED).
		SendDelayed(func() error {
			return n.PostReqJSON(ctx, url2, data)
		})
	if err != nil {
		logger.Log.Error("Failed to send JSON", zap.Error(err))
		return
	}

}

func (n *Notifier) SendReqBatched(ctx context.Context, gauge *map[string]string, counter uint64, PullGoCh chan struct{}, wgWorker *sync.WaitGroup) {

	PullGoCh <- struct{}{}
	defer func() {
		//освоблождаем горутину
		<-PullGoCh
		wgWorker.Done()
	}()

	url := "http://" + n.URL + "/updates/"

	m := make([]model.Metrics, 0, len(*gauge))

	for key, values := range *gauge {

		val, err := strconv.ParseFloat(values, 64)
		if err != nil {
			logger.Log.Error("Failed to parse float", zap.Error(err))
			return
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

	err := delaysend.NewDelaySend().SetDelay([]int{1, 3, 5}).
		AddExpectedError(syscall.ECONNREFUSED).
		SendDelayed(func() error {
			return n.PostReqBatched(ctx, url, m)
		})
	if err != nil {
		logger.Log.Error("Failed to send Batched", zap.Error(err))
		return
	}

}
func (n *Notifier) SendNotification(ctx context.Context, ch chan struct{}, gauge *map[string]string, counter uint64) {
	wgWorker := &sync.WaitGroup{}
	wgWorker.Add(3)

	go n.SendGaugeReq(ctx, gauge, ch, wgWorker)
	go n.SendCountReq(ctx, counter, ch, wgWorker)
	go n.SendReqBatched(ctx, gauge, counter, ch, wgWorker)

	wgWorker.Wait()
}

//func (n *Notifier) worker(jobs int, pullCh chan struct{}) {
//	//Создаем воркер, 3 функции => 3 функции
//	// в каждую по пуллу передаем
//	for w:=1; w<=jobs; w++ {
//
//	}
//
//	for {
//		select {
//		case <-pullCh:
//			return
//		default:
//			time.Sleep(1 * time.Second)
//			continue
//		}
//	}
//}

func (n *Notifier) StartNotifyCron(ctx context.Context) error {
	var gauge *map[string]string
	var couter uint64
	var err error
	mu := sync.Mutex{}
	//Обновление метрик
	go func() {
		ticker := time.NewTicker(n.TimerUpdate)
		for {
			select {
			case <-ctx.Done():
				// Обработка завершения программы
				return
			case <-ticker.C:
				mu.Lock()
				gauge, couter, err = n.NotifyPending()
				if err != nil {
					logger.Log.Info("Failed to pending", zap.Error(err))
					return
				}
				mu.Unlock()
				continue
			default:
				continue
			}
		}
	}()
	// отправка
	/*
		Создадим пул горутин
		Передадим в функцию отправки
		где будет распределение по отправке
	*/
	pullCh := make(chan struct{}, 15)

	go func() {
		ticker := time.NewTicker(n.TimerSend)
		for {
			select {
			case <-ctx.Done():
				// Обработка завершения программы
				close(pullCh)
				return
			case <-ticker.C:

				if gauge != nil {
					mu.Lock()
					n.SendNotification(ctx, pullCh, gauge, couter)
					mu.Unlock()
				}
				continue
			default:
				continue
			}
		}
	}()
	return nil
}
